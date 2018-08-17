/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"plugin"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin/girpc"
)

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir        string
	RPCAddress string
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path string, address string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	return &GoInsider{dir: path, RPCAddress: address}
}

// RPC struct with methods representing RPC interface of this code runner
type RPC struct {
	gi *GoInsider
}

// Call is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) Call(args girpc.CallReq, reply *girpc.CallResp) error {
	path, err := t.gi.ObtainCode(args.Reference)
	if err != nil {
		return errors.Wrap(err, "couldn't obtain code")
	}

	p, err := plugin.Open(path)
	if err != nil {
		return errors.Wrap(err, "couldn't open plugin")
	}

	export, err := p.Lookup("INSEXPORT")
	if err != nil {
		return errors.Wrap(err, "couldn't lookup 'INSEXPORT' in '"+path+"'")
	}

	ch := new(codec.CborHandle)

	err = codec.NewDecoderBytes(args.Data, ch).Decode(export)
	if err != nil {
		return errors.Wrapf(err, "couldn't decode data into %T", export)
	}

	method := reflect.ValueOf(export).MethodByName(args.Method)
	if !method.IsValid() {
		return errors.New("wtf, no method " + args.Method + "in the plugin")
	}

	inLen := method.Type().NumIn()

	mask := make([]interface{}, inLen)
	for i := 0; i < inLen; i++ {
		argType := method.Type().In(i)
		mask[i] = reflect.Zero(argType).Interface()
	}

	err = codec.NewDecoderBytes(args.Arguments, ch).Decode(&mask)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal CBOR for arguments of the method")
	}

	in := make([]reflect.Value, inLen)
	for i := 0; i < inLen; i++ {
		in[i] = reflect.ValueOf(mask[i])
	}

	resValues := method.Call(in)

	err = codec.NewEncoderBytes(&reply.Data, ch).Encode(export)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal new object data into cbor")
	}

	res := make([]interface{}, len(resValues))
	for i, v := range resValues {
		res[i] = v.Interface()
	}

	var resSerialized []byte
	err = codec.NewEncoderBytes(&resSerialized, ch).Encode(res)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal returned values into cbor")
	}

	reply.Ret = resSerialized

	return nil
}

// ObtainCode returns path on the file system to the plugin, fetches it from a provider
// if it's not in the storage
func (t *GoInsider) ObtainCode(ref logicrunner.Reference) (string, error) {
	path := t.dir + "/" + string(ref)
	_, err := os.Stat(path)

	if err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", errors.Wrap(err, "file !notexists()")
	}

	client, err := rpc.DialHTTP("tcp", t.RPCAddress)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't dial '%s'", t.RPCAddress)
	}

	res := logicrunner.Object{}
	err = client.Call("RPC.GetObject", ref, &res)
	if err != nil {
		return "", errors.Wrap(err, "on calling main API")
	}

	err = ioutil.WriteFile(path, res.Data, 0666)
	if err != nil {
		return "", errors.Wrap(err, "on writing file down")
	}

	return path, nil
}

func main() {
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	pflag.Parse()

	insider := NewGoInsider(*path, *rpcAddress)
	err := rpc.Register(&RPC{insider})
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
