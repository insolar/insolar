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

package ginsider

import (
	"io/ioutil"
	"net/rpc"
	"os"
	"plugin"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir                string
	UpstreamRPCAddress string
	UpstreamRPCClient  *rpc.Client
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path string, address string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	return &GoInsider{dir: path, UpstreamRPCAddress: address}
}

// RPC struct with methods representing RPC interface of this code runner
type RPC struct {
	GI *GoInsider
}

// Call is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) Call(args rpctypes.DownCallReq, reply *rpctypes.DownCallResp) error {
	path, err := t.GI.ObtainCode(args.Reference)
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

// Upstream returns RPC client connected to upstream server (goplugin)
func (t *GoInsider) Upstream() (*rpc.Client, error) {
	if t.UpstreamRPCClient != nil {
		return t.UpstreamRPCClient, nil
	}

	client, err := rpc.DialHTTP("tcp", t.UpstreamRPCAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial '%s'", t.UpstreamRPCAddress)
	}

	t.UpstreamRPCClient = client
	return t.UpstreamRPCClient, nil
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

	client, err := t.Upstream()
	if err != nil {
		return "", err
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

// Exec
func (t *GoInsider) Exec(ref string, method string, args []byte) (data []byte, res []byte, err error) {
	return data, res, err
}

// CurrentGoInsider - hackish way to give proxies access to the current environment
var CurrentGoInsider *GoInsider
