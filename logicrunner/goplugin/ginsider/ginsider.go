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
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"plugin"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir                string
	UpstreamRPCAddress string
	UpstreamRPCClient  *rpc.Client
	plugins            map[string]*plugin.Plugin
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path string, address string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	res := GoInsider{dir: path, UpstreamRPCAddress: address}
	res.plugins = make(map[string]*plugin.Plugin)
	return &res
}

// RPC struct with methods representing RPC interface of this code runner
type RPC struct {
	GI *GoInsider
}

// CallMethod is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallMethod(args rpctypes.DownCallMethodReq, reply *rpctypes.DownCallMethodResp) error {
	p, err := t.GI.Plugin(args.Reference)
	if err != nil {
		return err
	}

	export, err := p.Lookup("INSEXPORT")
	if err != nil {
		return errors.Wrap(err, "couldn't lookup 'INSEXPORT' in plugin")
	}

	ch := new(codec.CborHandle)

	err = codec.NewDecoderBytes(args.Data, ch).Decode(export)
	if err != nil {
		return errors.Wrapf(err, "couldn't decode data into %T", export)
	}

	setContext := reflect.ValueOf(export).MethodByName("SetContext")
	if !setContext.IsValid() {
		return errors.New("this is not a contract, it not supports SetContext method")
	}
	cc := foundation.CallContext{
		Me: foundation.Reference("contract address"),
		// fill me
	}
	setContext.Call([]reflect.Value{reflect.ValueOf(&cc)})

	method := reflect.ValueOf(export).MethodByName(args.Method)
	if !method.IsValid() {
		return errors.New("no method " + args.Method + " in the plugin")
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

	log.Debugf(
		"Calling method %q in contract %q with %d arguments",
		args.Method, args.Reference, inLen,
	)
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

// CallConstructor is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallConstructor(args rpctypes.DownCallConstructorReq, reply *rpctypes.DownCallConstructorResp) error {
	p, err := t.GI.Plugin(args.Reference)
	if err != nil {
		return err
	}

	export, err := p.Lookup(args.Name)
	if err != nil {
		return errors.Wrapf(err, "couldn't lookup symbol %q in plugin", args.Name)
	}

	method := reflect.ValueOf(export)
	if !method.IsValid() {
		return fmt.Errorf("%q is not valid symbol", args.Name)
	}
	if method.Kind() != reflect.Func {
		return fmt.Errorf("%q is not a function", args.Name)
	}

	inLen := method.Type().NumIn()

	mask := make([]interface{}, inLen)
	for i := 0; i < inLen; i++ {
		argType := method.Type().In(i)
		mask[i] = reflect.Zero(argType).Interface()
	}

	ch := new(codec.CborHandle)

	err = codec.NewDecoderBytes(args.Arguments, ch).Decode(&mask)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal CBOR for arguments of the constructor")
	}

	in := make([]reflect.Value, inLen)
	for i := 0; i < inLen; i++ {
		in[i] = reflect.ValueOf(mask[i])
	}

	log.Debugf(
		"Calling constructor %q in contract %q with %d arguments",
		args.Name, args.Reference, inLen,
	)
	resValues := method.Call(in)

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
func (t *GoInsider) ObtainCode(ref core.RecordRef) (string, error) {
	path := t.dir + "/" + ref.String()
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

	log.Debugf("obtaining plugin %q", ref)
	res := rpctypes.Object{}
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

// Plugin loads Go plugin by reference and returns `*plugin.Plugin`
// ready to lookup symbols
func (t *GoInsider) Plugin(ref core.RecordRef) (*plugin.Plugin, error) {
	key := ref.String()
	if t.plugins[key] != nil {
		return t.plugins[key], nil
	}

	path, err := t.ObtainCode(ref)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't obtain code")
	}

	log.Debugf("Opening plugin %q from file %q", ref, path)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open plugin")
	}

	t.plugins[key] = p
	return p, nil
}

// RouteCall ...
func (t *GoInsider) RouteCall(ref string, method string, args []byte) ([]byte, error) {
	client, err := t.Upstream()
	if err != nil {
		return nil, err
	}

	req := rpctypes.UpRouteReq{
		Reference: core.String2Ref(ref),
		Method:    method,
		Arguments: args,
	}

	res := rpctypes.UpRouteResp{}
	err = client.Call("RPC.RouteCall", req, &res)
	if err != nil {
		return nil, errors.Wrap(err, "on calling main API")
	}

	return []byte(res.Result), res.Err
}

// RouteConstructorCall ...
func (t *GoInsider) RouteConstructorCall(ref string, name string, args []byte) (string, error) {
	client, err := t.Upstream()
	if err != nil {
		return "", err
	}

	req := rpctypes.UpRouteConstructorReq{
		Reference:   core.String2Ref(ref),
		Constructor: name,
		Arguments:   args,
	}

	res := rpctypes.UpRouteConstructorResp{}
	err = client.Call("RPC.RouteConstructorCall", req, &res)
	if err != nil {
		return "", errors.Wrap(err, "on calling main API")
	}

	return res.Reference.String(), res.Err
}

// Serialize - CBOR serializer wrapper: `what` -> `to`
func (t *GoInsider) Serialize(what interface{}, to *[]byte) error {
	ch := new(codec.CborHandle)
	log.Printf("serializing %+v", what)
	return codec.NewEncoderBytes(to, ch).Encode(what)
}

// Deserialize - CBOR de-serializer wrapper: `from` -> `into`
func (t *GoInsider) Deserialize(from []byte, into interface{}) error {
	ch := new(codec.CborHandle)
	log.Printf("de-serializing %+v", from)
	return codec.NewDecoderBytes(from, ch).Decode(into)
}
