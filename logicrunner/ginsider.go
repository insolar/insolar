package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
)

type Args struct {
	code      []byte
	data      []byte
	method    string
	arguments []byte
}

type CallReply struct {
	data   []byte
	result []byte
}

type GoInsider int

func (t *GoInsider) Call(args *Args, reply *CallReply) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	insider := new(GoInsider)
	rpc.Register(insider)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
