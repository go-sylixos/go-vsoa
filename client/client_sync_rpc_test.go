package client

import (
	"encoding/json"
	"flag"
	"go-vsoa/protocol"
	"testing"
)

var (
	rpc_addr = flag.String("rpc_addr", "localhost:3002", "server address")
)

type RpcTestParam struct {
	Num int `json:"Test Num"`
}

func TestRPC(t *testing.T) {
	flag.Parse()

	clientOption := Option{
		Password: "123456",
	}

	c := NewClient(clientOption)
	_, err := c.Connect("tcp", *rpc_addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	req := protocol.NewMessage()
	reply := protocol.NewMessage()

	reply, err = c.Call("/a/b/c", protocol.TypeRPC, protocol.RpcMethodGet, req)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("Seq:", reply.SeqNo(), "Param:", (reply.Param))
	}

	req.Param, _ = json.RawMessage(`{"Test Num":123}`).MarshalJSON()

	reply, err = c.Call("/a/b/c", protocol.TypeRPC, protocol.RpcMethodGet, req)
	if err != nil {
		t.Fatal(err)
	} else {
		DstParam := new(RpcTestParam)
		json.Unmarshal(reply.Param, DstParam)
		t.Log("Seq:", reply.SeqNo(), "Param:", DstParam, "Unmarshaled data:", DstParam.Num)
	}
}
