package client

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/go-sylixos/go-vsoa/protocol"
)

var (
	publish_server_quick_addr = flag.String("publish_server_quick_addr", "localhost:3003", "server address")
)

type QPublishTestParam struct {
	Publish string `json:"publish"`
}

// TestSub is a test function for testing the Subscribe and UnSubscribe methods of the Client struct.
//
// This function starts a server, initializes a callback, parses flags, creates a client with a password option, and connects to the server.
// It then subscribes to a channel and checks for any errors. If there is an error, it checks if the error is an invalid URL and logs a pass message.
// After a delay, it unsubscribes from the channel and then subscribes again.
// Finally, it waits for another delay.
func TestSubQuickChannel(t *testing.T) {
	startServer()

	// Do this to make sure the server is ready on slow machine
	time.Sleep(50 * time.Millisecond)

	cb := new(callback)
	cb.T = t

	flag.Parse()

	clientOption := Option{
		Password: "123456",
	}

	c := NewClient(clientOption)
	_, err := c.Connect("vsoa", *publish_server_quick_addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	//client don't know if it's quick channel or not
	err = c.Subscribe("/p/q", cb.getQPublishParam)
	if err != nil {
		if err == strErr(protocol.StatusText(protocol.StatusInvalidUrl)) {
			t.Log("Pass: Invalid URL")
		} else {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)

	err = c.UnSubscribe("/p/q")
	if err != nil {
		if err == strErr(protocol.StatusText(protocol.StatusInvalidUrl)) {
			t.Log("Pass: Invalid URL")
		} else {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)

	err = c.Subscribe("/p/q", cb.getQPublishParam)
	if err != nil {
		if err == strErr(protocol.StatusText(protocol.StatusInvalidUrl)) {
			t.Log("Pass: Invalid URL")
		} else {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)
}

// User can create callback struct to put/get more info into callback func
func (c callback) getQPublishParam(m *protocol.Message) {
	DstParam := new(QPublishTestParam)
	json.Unmarshal(m.Param, DstParam)
	c.T.Log("Param:", DstParam.Publish)
}
