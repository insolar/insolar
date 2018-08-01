package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"io/ioutil"
	"os"

	"plugin"
	"reflect"

	"github.com/2tvenom/cbor"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type GoInsider struct {
	dir string
}

func NewGoInsider(path string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	return &GoInsider{path}
}

func (t *GoInsider) Call(args goplugin.CallReq, reply *goplugin.CallResp) error {
	path := t.dir + "/" + args.Object.Reference
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := ioutil.WriteFile(path, args.Object.Code, 0666)
		check(err)
	} else {
		check(err)
	}

	p, err := plugin.Open(path)
	check(err)

	export, err := p.Lookup("EXP")
	check(err)

	cbor := cbor.NewEncoder(&bytes.Buffer{})
	some, err := cbor.Unmarshal(args.Object.Data, export)

	r2 := reflect.ValueOf(export)
	m2 := r2.MethodByName(args.Method)
	_ = m2.Call([]reflect.Value{})

	return nil
}

var PATH = "/Users/ruz/go/src/github.com/insolar/insolar/tmp"

func main() {
	log.Print("ginsider launched")
	insider := GoInsider{PATH}
	rpc.Register(insider)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7777")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	<-make(chan byte)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
