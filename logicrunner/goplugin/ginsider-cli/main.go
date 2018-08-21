package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/spf13/pflag"

	"github.com/insolar/insolar/logicrunner/goplugin/ginsider"
)

func main() {
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	pflag.Parse()

	insider := ginsider.NewGoInsider(*path, *rpcAddress)
	ginsider.CurrentGoInsider = insider

	err := rpc.Register(&ginsider.RPC{insider})
	if err != nil {
		log.Fatal("Couldn't register RPC interface: ", err)
		os.Exit(1)
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", *listen)

	log.Print("ginsider launched, listens " + *listen)
	if err != nil {
		log.Fatal("listen error:", err)
		os.Exit(1)
	}
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("couldn't start server: ", err)
		os.Exit(1)
	}
	log.Print("bye\n")
}
