package server

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-sylixos/go-vsoa/protocol"
)

// publisher is a method of the VsoaServer struct that sends a publish message to all active clients subscribed to a specific service path at the specified time interval.
//
// Parameters:
// - servicePath: a string representing the service path to publish.
// - timeDriction: a time.Duration value representing the time interval between each publish message.
// - pubs: a function that takes two parameters: a pointer to a protocol.Message and a pointer to another protocol.Message. It is called to initialize the request message before publishing.
func (s *Server) publisher(servicePath string, timeOrTrigger any, pubs func(*protocol.Message, *protocol.Message)) {
	req := protocol.NewMessage()

	var ticker *time.Ticker
	isTrigger := false

	switch v := timeOrTrigger.(type) {
	case time.Duration:
		ticker = time.NewTicker(v)
		defer ticker.Stop()
	case chan struct{}:
		s.triggerChan[servicePath] = v
		isTrigger = true
	default:
		panic("Invalid type for timeOrTrigger")
	}

	for {
		var wg sync.WaitGroup
		var ctx context.Context
		var cancel context.CancelFunc
		var timeout time.Duration

		if isTrigger {
			<-s.triggerChan[servicePath]
			timeout = time.Duration(len(s.clients)) * time.Millisecond
			ctx, cancel = context.WithTimeout(context.Background(), timeout)
		} else {
			<-ticker.C
			timeout = 4 * time.Duration(timeOrTrigger.(time.Duration)) / 5
			ctx, cancel = context.WithTimeout(context.Background(), timeout)
		}

		pubs(req, nil)

		for _, c := range s.clients {
			if c.Subscribes[servicePath] && c.Authed {
				wg.Add(1)
				go func(c *client) {
					defer wg.Done()
					reqCopy := *req // Aviod change req object at the same time.
					reqCopy.URL = []byte(servicePath)
					s.sendMessageWithContext(ctx, &reqCopy, c.Conn, timeout)
				}(c)
			}
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All sends completed within the period
		case <-ctx.Done():
			// Timeout, 4/5 of the period elapsed
		}

		cancel()
	}
}

// Normal channel Publish Message
func (s *Server) sendMessageWithContext(ctx context.Context, req *protocol.Message, conn net.Conn, timeout time.Duration) {
	select {
	case <-ctx.Done():
		// Context cancelled or timed out
		return
	default:
		// Send the message
		req.SetMessageType(protocol.TypePublish)

		req.SetReply(false)

		tmp, err := req.Encode(protocol.ChannelNormal)
		if err != nil {
			log.Panicln(err)
			return
		}

		conn.SetWriteDeadline(time.Now().Add(timeout))

		_, err = conn.Write(tmp)
		protocol.PutData(&tmp)
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}

		return
	}
}
