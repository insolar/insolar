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

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir        string
	RPCAddress string
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	return &GoInsider{dir: path}
}

// Call is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *GoInsider) Call(args goplugin.CallReq, reply *goplugin.CallResp) error {
	path, err := t.ObtainCode(args.Object)
	check(err)

	p, err := plugin.Open(path)
	check(err)

	export, err := p.Lookup("INSEXPORT")
	check(err)

	var dataBuf bytes.Buffer
	cbor := cbor.NewEncoder(&dataBuf)
	_, err = cbor.Unmarshal(args.Object.Data, export)
	check(err)

	method := reflect.ValueOf(export).MethodByName("INSMETHOD__" + args.Method)
	if !method.IsValid() {
		panic("wtf, no method " + args.Method + "in the plugin")
	}

	res := method.Call([]reflect.Value{})

	_, err = cbor.Marshal(export)
	check(err)

	reply.Data = dataBuf.Bytes()

	log.Printf("res: %+v\n", res)

	return nil
}

// ObtainCode returns path on the file system to the plugin, fetches it from a provider
// if it's not in the storage
func (t *GoInsider) ObtainCode(obj logicrunner.Object) (string, error) {
	path := t.dir + "/" + string(obj.Reference)
	_, err := os.Stat(path)

	if err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	client, err := rpc.DialHTTP("tcp", t.RPCAddress)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't dial '%s'", t.RPCAddress)
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
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	pflag.Parse()

	insider := GoInsider{dir: *path, RPCAddress: *rpcAddress}
	err := rpc.Register(&insider)
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
	go http.Serve(listener, nil)
	<-make(chan byte)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
