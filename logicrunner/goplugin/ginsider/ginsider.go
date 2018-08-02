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
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type GoInsider struct {
	dir        string
	RpcAddress string
}

func NewGoInsider(path string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	return &GoInsider{dir: path}
}

func (t *GoInsider) Call(args goplugin.CallReq, reply *goplugin.CallResp) error {
	path, err := t.ObtainCode(args.Object)
	check(err)

	p, err := plugin.Open(path)
	check(err)

	export, err := p.Lookup("INSEXPORT")
	check(err)

	var data_buf bytes.Buffer
	cbor := cbor.NewEncoder(&data_buf)
	_, err = cbor.Unmarshal(args.Object.Data, export)
	check(err)

	method := reflect.ValueOf(export).MethodByName("INSMETHOD__" + args.Method)
	if !method.IsValid() {
		panic("wtf, no method " + args.Method + "in the plugin")
	}

	res := method.Call([]reflect.Value{})

	cbor.Marshal(export)
	reply.Data = data_buf.Bytes()

	log.Printf("res: %+v\n", res)

	return nil
}

func (t *GoInsider) ObtainCode(obj logicrunner.Object) (string, error) {
	path := t.dir + "/" + string(obj.Reference)
	_, err := os.Stat(path)

	if err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	client, err := rpc.DialHTTP("tcp", t.RpcAddress)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't dial '%s'", t.RpcAddress)
	}

	res := logicrunner.Object{}
	err = client.Call("GoPluginRPC.GetObject", obj.Reference, &res)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path, res.Data, 0666)
	if err != nil {
		return "", err
	}

	return path, nil
}

func main() {
	listen := pflag.StringP("listen", "l", ":7778", "address and port to listen")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpc_address := pflag.String("rpc", "localhost:7777", "address and port of RPC API")
	pflag.Parse()

	insider := GoInsider{dir: *path, RpcAddress: *rpc_address}
	rpc.Register(&insider)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", *listen)

	log.Print("ginsider launched, listens " + *listen)
	if err != nil {
		log.Fatal("listen error:", err)
		os.Exit(1)
	}
	go http.Serve(listener, nil)
	<-make(chan byte)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
