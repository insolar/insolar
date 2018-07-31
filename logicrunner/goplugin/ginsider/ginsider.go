package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"io/ioutil"
	"os"

	"plugin"
	"reflect"

	"github.com/insolar/insolar/logicrunner/goplugin"
)

type GoInsider struct {
}

var PATH = "/Users/ruz/go/src/github.com/insolar/insolar/tmp"

func (t *GoInsider) Call(args goplugin.CallReq, reply *goplugin.CallResp) error {
	path := PATH + "/" + args.Object.Reference
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := ioutil.WriteFile(path, args.Object.Code, 0666)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	p, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}

	export, err := p.Lookup("EXP")
	if err != nil {
		panic(err)
	}

	r2 := reflect.ValueOf(export)
	m2 := r2.MethodByName(args.Method)
	_ = m2.Call([]reflect.Value{})

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
