package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/insolar/insolar/logicrunner/goplugin"
)

type GoInsider struct {
}

func (t *GoInsider) Call(args goplugin.CallReq, reply *goplugin.CallResp) error {
	*reply = goplugin.CallResp{
		Ret: []byte{1, 2, 3, 4, 5},
		Err: nil,
	}
	return nil
}

func main() {
	log.Print("ginsider launched")
	insider := new(GoInsider)
	rpc.Register(insider)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7777")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	<-make(chan byte)
}
