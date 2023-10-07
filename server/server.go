package server

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"go-vsoa/protocol"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ErrServerClosed is returned by the Server's Serve, ListenAndServe after a call to Shutdown or Close.
var (
	ErrServerClosed      = errors.New("VSOA: Server closed")
	ErrReqReachLimit     = errors.New("request reached rate limit")
	ErrAlreadyRegistered = errors.New("URL has been Registered")
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024

	DefaultTimeout = 5 * time.Minute
)

// VSOA Server need to konw the client infos
// until TCP connection closed
type client struct {
	Uid  uint32 // Server send this to client, helps client send quick massage
	Conn net.Conn
	// Quick channel Conn is UDP conn so we can't close real client conn
	// But we have to know the QAddr to delete quickChannel map
	QAddr          *net.UDPAddr
	beforeServInfo bool
	Authed         bool
	Subscribes     map[string]bool // key: URL, value: If Subs
}

// Handler declares the signature of a function that can be bound to a Route.
type Handler func(req *protocol.Message, resp *protocol.Message)

// VsoaServer is the VSOA server that use TCP with UDP.
type VsoaServer struct {
	Name         string //Used for ServInfo
	address      string
	option       Option
	ln           net.Listener
	readTimeout  time.Duration
	writeTimeout time.Duration

	routerMapMu sync.RWMutex
	routeMap    map[string]Handler

	mu            sync.RWMutex
	activeClients map[uint32]*client
	// When QuickChannel get RemoteAddr we need to use it to check if we have the activeClient
	quickChannel map[string]uint32
	clientsCount atomic.Uint32
	doneChan     chan struct{}

	inShutdown int32
	// onShutdown []func(s *VsoaServer)
	// onRestart  []func(s *VsoaServer)

	// TLSConfig for creating tls tcp connection.
	tlsConfig *tls.Config

	handlerMsgNum int32

	// HandleServiceError is used to get all service errors. You can use it write logs or others.
	HandleServiceError func(error)

	// ServerErrorFunc is a customized error handlers and you can use it to return customized error strings to clients.
	// If not set, it use err.Error()
	ServerErrorFunc func(res *protocol.Message, err error) string
}

// NewServer returns a server.
func NewServer(name string, so Option) *VsoaServer {
	s := &VsoaServer{
		Name:   name,
		option: so,
		// this can cause server close connection
		readTimeout:   DefaultTimeout,
		writeTimeout:  DefaultTimeout,
		quickChannel:  make(map[string]uint32),
		activeClients: make(map[uint32]*client),
		doneChan:      make(chan struct{}),
		routeMap:      make(map[string]Handler),
	}

	return s
}

// Serve starts and listens VSOA normal channel requests.
// TODO: we need to start listen Quick channel too!
// It is blocked until receiving connections from clients.
func (s *VsoaServer) Serve(address string) (err error) {
	var ln net.Listener
	ln, err = s.makeListener("tcp", address)
	if err != nil {
		return err
	}

	s.address = address

	// Go quick channel listener
	go s.serveQuickListener(address)

	return s.serveListener(ln)
}

// serveListener accepts incoming connections on the Listener ln,
// creating a new service goroutine for each.
// The service goroutines read requests and then call services to reply to them.
func (s *VsoaServer) serveListener(ln net.Listener) error {
	var tempDelay time.Duration

	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if s.isShutdown() {
				<-s.doneChan
				return ErrServerClosed
			}

			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
				continue
			}

			return e
		}
		tempDelay = 0

		CUid := s.clientsCount.Add(1)

		s.mu.Lock()
		// We don't know the QuickChannel info yet.
		s.activeClients[CUid] = &client{
			Conn: conn,
			Uid:  CUid,
			// until we got ServInfo
			QAddr:          nil,
			beforeServInfo: true,
			Authed:         false,
		}
		s.mu.Unlock()

		go s.serveConn(conn, CUid)
	}
}

// sendResponse sends a response to the client.
//
// It takes in a res *protocol.Message object, representing the response to be sent,
// and a conn net.Conn object, representing the network connection to the client.
// The function encodes the response message and writes it to the connection.
// If a write timeout is set, it sets the write deadline before writing the response.
// The function returns no values.
func (s *VsoaServer) sendResponse(res *protocol.Message, conn net.Conn) {
	// Do service method
	tmp, err := res.Encode(protocol.ChannelNormal)
	if err != nil {
		log.Panicln(err)
		return
	}

	if s.writeTimeout != 0 {
		conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	}
	conn.Write(tmp)
	protocol.PutData(&tmp)
}

// serveConn serves a connection and handles incoming requests for the VsoaServer.
// This is a session for one client, so we need to keep auth logic inside it
//
// Parameters:
// - conn: the net.Conn representing the connection.
// - ClientUid: the uint32 representing the client's UID.
//
// Returns: none.
func (s *VsoaServer) serveConn(conn net.Conn, ClientUid uint32) {
	if s.isShutdown() {
		s.closeConn(ClientUid)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Printf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
		}

		// make sure all inflight requests are handled and all drained
		if s.isShutdown() {
			<-s.doneChan
		}

		// close TCP&UDPconn
		s.closeConn(ClientUid)
	}()

	if tlsConn, ok := conn.(*tls.Conn); ok {
		if d := s.readTimeout; d != 0 {
			conn.SetReadDeadline(time.Now().Add(d))
		}
		if d := s.writeTimeout; d != 0 {
			conn.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			return
		}
	}

	r := bufio.NewReaderSize(conn, ReaderBuffsize)

	// read requests and handle it
	for {
		if s.isShutdown() {
			return
		}

		t0 := time.Now()
		// If client send nothing during readTimeout to server, server will kill the connection!
		if s.readTimeout != 0 {
			conn.SetReadDeadline(t0.Add(s.readTimeout))
		}

		// read a request from the underlying connection
		req := protocol.NewMessage()
		// If client send nothing during readTimeout to server, can cause error
		err := req.Decode(r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Vsoa client has closed this connection: %s", conn.RemoteAddr().String())
			}

			if s.HandleServiceError != nil {
				s.HandleServiceError(err)
			}
			return
		}

		if !req.IsPingEcho() && !req.IsServInfo() {
			if !s.activeClients[ClientUid].Authed {
				// Close unauthed client
				log.Printf("auth failed for conn %s: %v", conn.RemoteAddr().String(), protocol.StatusText(protocol.StatusPassword))
				return
			}
		}

		go s.processOneRequest(req, conn, ClientUid)
	}
}

// processOneRequest is a method of the VsoaServer struct that processes a single request.
//
// It takes the following parameters:
// - req: a pointer to a protocol.Message struct, representing the request message
// - conn: a net.Conn object, representing the network connection
// - ClientUid: an unsigned 32-bit integer, representing the client UID
//
// There is no return value.
func (s *VsoaServer) processOneRequest(req *protocol.Message, conn net.Conn, ClientUid uint32) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, true)]

			log.Printf("failed to handle the request: %v, stacks: %s", r, buf)
		}
	}()

	atomic.AddInt32(&s.handlerMsgNum, 1)
	defer atomic.AddInt32(&s.handlerMsgNum, -1)

	res := req.CloneHeader()

	if req.IsServInfo() {
		err := s.servInfoHandler(req, res, ClientUid)
		if err != nil {
			log.Printf("Failed to Auth Client: %d, err: %s", ClientUid, err)
		}
		s.sendResponse(res, conn)
		return
	}

	res.SetReply(true)

	if req.IsPingEcho() {
		s.sendResponse(res, conn)
		return
	}

	if !req.IsOneway() {
		// TODO: father URL logic
		if handle, ok := s.routeMap["RPC."+string(req.URL)+"."+req.MessageRpcMethodText()]; ok {
			handle(req, res)
		} else if _, ok := s.routeMap["SUBS/UNSUBS."+string(req.URL)]; ok {
			s.subs(req, ClientUid)
			res.SetStatusType(protocol.StatusSuccess)
		} else {
			res.SetStatusType(protocol.StatusInvalidUrl)
		}

		s.sendResponse(res, conn)
	} else {
		if handle, ok := s.routeMap["DATAGRAME."+string(req.URL)]; ok {
			handle(req, res)
		} else if handle, ok := s.routeMap["DATAGRAME.DEFAULT"]; ok {
			handle(req, res)
		}
	}
}

// servInfoHandler handles the server information request from a client.
//
// It takes in the request message, the response message, and the client UID as parameters.
// It returns an error if any error occurs during the process.
func (s *VsoaServer) servInfoHandler(req *protocol.Message, resp *protocol.Message, ClientUid uint32) error {
	r := new(protocol.ServInfoResParam)
	r.Info = s.Name
	s.mu.Lock()
	s.activeClients[ClientUid].beforeServInfo = false
	s.mu.Unlock()
	if s.option.Password == "" {
		s.mu.Lock()
		s.activeClients[ClientUid].Authed = true
		// quick channel register
		if req.TunID() != 0 {
			qAddr := (*net.UDPAddr)(s.activeClients[ClientUid].Conn.RemoteAddr().(*net.TCPAddr))
			qAddr.Port = int(req.TunID())
			qString := qAddr.String()
			s.activeClients[ClientUid].QAddr = (qAddr)
			s.quickChannel[(qString)] = ClientUid
		}
		s.activeClients[ClientUid].Subscribes = make(map[string]bool)
		s.mu.Unlock()
		r.NewGoodMessage(protocol.ServInfoResAsString, resp, ClientUid)
	} else {
		infoParam := new(protocol.ServInfoReqParam)
		err := json.Unmarshal(req.Param, infoParam)
		if err != nil {
			return err
		} else {
			if infoParam.Password != s.option.Password {
				r.NewErrMessage(resp)
				// After this call server to close normal channel conn
				return protocol.ErrMessagePasswd
			} else {
				s.mu.Lock()
				s.activeClients[ClientUid].Authed = true
				// quick channel register
				if req.TunID() != 0 {
					qAddr := (*net.UDPAddr)(s.activeClients[ClientUid].Conn.RemoteAddr().(*net.TCPAddr))
					qAddr.Port = int(req.TunID())
					qString := qAddr.String()
					s.activeClients[ClientUid].QAddr = (qAddr)
					s.quickChannel[(qString)] = ClientUid
				}
				s.activeClients[ClientUid].Subscribes = make(map[string]bool)
				s.mu.Unlock()
				r.NewGoodMessage(protocol.ServInfoResAsString, resp, ClientUid)
			}
		}
	}

	// TODO: handle other client options like ping echo seting logic
	return nil
}

// TODO: add both GET&SET with same handler!
// AddRpcHandler adds an RPC handler to the VsoaServer.
//
// It takes in the following parameters:
// - servicePath: the path of the service
// - serviceMethod: the type of the RPC message
// - handler: the function to handle the RPC message
//
// It returns an error.
func (s *VsoaServer) AddRpcHandler(servicePath string, serviceMethod protocol.RpcMessageType, handler func(*protocol.Message, *protocol.Message)) (err error) {
	s.routerMapMu.Lock()
	defer s.routerMapMu.Unlock()
	if _, ok := s.routeMap["RPC."+servicePath+"."+protocol.RpcMethodText(serviceMethod)]; !ok {
		s.routeMap["RPC."+servicePath+"."+protocol.RpcMethodText(serviceMethod)] = handler
	} else {
		return ErrAlreadyRegistered
	}
	return nil
}

// AddOneWayHandler adds a DATAGRAME handler to the VsoaServer.
//
// It takes in the servicePath string and the handler function, and returns an error.
func (s *VsoaServer) AddOneWayHandler(servicePath string, handler func(*protocol.Message, *protocol.Message)) (err error) {
	s.routerMapMu.Lock()
	defer s.routerMapMu.Unlock()
	if _, ok := s.routeMap["DATAGRAME."+servicePath]; !ok {
		s.routeMap["DATAGRAME."+servicePath] = handler
	} else {
		return ErrAlreadyRegistered
	}
	return nil
}

// AddDefaultOndataHandler adds a default DATAGRAME handler to the VsoaServer.
//
// The handler parameter is a function that takes two parameters: a pointer to a protocol.Message
// and a pointer to another protocol.Message. It is responsible for handling the ondata event.
// This function does not return anything.
func (s *VsoaServer) AddDefaultOndataHandler(handler func(*protocol.Message, *protocol.Message)) {
	s.routerMapMu.Lock()
	defer s.routerMapMu.Unlock()
	s.routeMap["DATAGRAME.DEFAULT"] = handler
}

// AddPublisher adds a publisher to the VsoaServer.
// Use this function to register a publish handler after Subs calls.
//
// The function takes the following parameter(s):
// - servicePath: a string representing the service path
// - timeDriction: a time duration representing the time duration
// - pubs: a function that takes two pointers to protocol.Message and returns nothing
//
// It returns an error.
func (s *VsoaServer) AddPublisher(servicePath string, timeDriction time.Duration, pubs func(*protocol.Message, *protocol.Message)) (err error) {
	s.routerMapMu.Lock()
	defer s.routerMapMu.Unlock()

	if _, ok := s.routeMap["SUBS/UNSUBS."+servicePath]; !ok {
		// No need to have handler save in the routeMap
		s.routeMap["SUBS/UNSUBS."+servicePath] = pubs
		// Maybe it's bad to run a Publisher for each pub
		go s.publisher(servicePath, timeDriction, pubs)
	} else {
		return ErrAlreadyRegistered
	}
	return nil
}

// subs updates the subscription status of a client.
//
// It takes a request message and a client UID as parameters and updates the
// subscription status of the client accordingly. If the request is a subscribe
// request, it sets the subscription status of the client to true for the given
// URL. Otherwise, it sets the subscription status to false.
func (s *VsoaServer) subs(req *protocol.Message, ClientUid uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if req.IsSubscribe() {
		s.activeClients[ClientUid].Subscribes[string(req.URL)] = true
	} else {
		s.activeClients[ClientUid].Subscribes[string(req.URL)] = false
	}
}

// isShutdown checks if the VsoaServer is in the shutdown state.
//
// It returns a boolean value indicating whether the VsoaServer is in the
// shutdown state or not.
func (s *VsoaServer) isShutdown() bool {
	return atomic.LoadInt32(&s.inShutdown) == 1
}

// closeConn closes the connection for a given client in the VsoaServer.
//
// It takes the ClientUid as a parameter, which is the unique identifier of the client.
// There is no return type for this function.
func (s *VsoaServer) closeConn(ClientUid uint32) {
	c := s.activeClients[ClientUid]
	// Quick Channel connection needs to close by client
	c.Conn.Close()

	s.mu.Lock()
	defer s.mu.Unlock()
	// Clear Client Subscribes map
	for k := range s.activeClients[ClientUid].Subscribes {
		delete(s.activeClients[ClientUid].Subscribes, k)
	}
	delete(s.quickChannel, s.activeClients[ClientUid].QAddr.String())
	delete(s.activeClients, ClientUid)
}

// Option contains all options for creating server.
type Option struct {
	Password string
	// TLSConfig for tcp and quic
	TLSConfig *tls.Config
}
