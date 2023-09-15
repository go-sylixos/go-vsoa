package server

import (
	"go-vsoa/protocol"
	"log"
	"net"
	"time"
)

func (s *VsoaServer) publisher(servicePath string, timeDriction time.Duration, pubs func(*protocol.Message, *protocol.Message)) {
	req := protocol.NewMessage()
	pubs(req, nil)

	ticker := time.NewTicker(timeDriction)
	defer ticker.Stop()

	for range ticker.C {
		for _, client := range s.activeClients {
			if client.Subscribes[servicePath] {
				req.URL = []byte(servicePath)
				s.sendMessage(req, client.Conn)
			}
		}
	}
}

// Normal channel Publish Message
func (s *VsoaServer) sendMessage(req *protocol.Message, conn net.Conn) error {
	req.SetMessageType(protocol.TypePublish)

	//	req.SeqNo(seq)
	req.SetReply(false)

	tmp, err := req.Encode(protocol.ChannelNormal)
	if err != nil {
		log.Panicln(err)
		return err
	}

	if s.writeTimeout != 0 {
		conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	}

	_, err = conn.Write(tmp)
	protocol.PutData(&tmp)

	return err
}
