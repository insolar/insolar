/*
 *    Copyright 2018 Insolar
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
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir                string
	UpstreamRPCAddress string
	UpstreamRPCClient  *rpc.Client
	plugins            map[string]*plugin.Plugin
	pluginsMutex       sync.Mutex
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
	p, err := t.GI.Plugin(args.Code)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("INSMETHOD_" + args.Method)
	if err != nil {
		return errors.Wrapf(err, "Can't find wrapper for %s", args.Method)
	}

	wrapper, ok := symbol.(func(ph proxyctx.ProxyHelper, object []byte,
		data []byte, context *core.LogicCallContext) ([]byte, []byte, error))
	if !ok {
		return errors.New("Wrapper with wrong signature")
	}

	state, result, err := wrapper(t.GI, args.Data, args.Arguments, args.Context) // may be entire args???

	if err != nil {
		return errors.Wrapf(err, "Method call returned error")
	}
	reply.Data = state
	reply.Ret = result
	return nil
}

// CallConstructor is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallConstructor(args rpctypes.DownCallConstructorReq, reply *rpctypes.DownCallConstructorResp) error {
	p, err := t.GI.Plugin(args.Code)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("INSCONSTRUCTOR_" + args.Name)
	if err != nil {
		return errors.Wrapf(err, "Can't find wrapper for %s", args.Name)
	}

	f, ok := symbol.(func(ph proxyctx.ProxyHelper, data []byte) ([]byte, error))
	if !ok {
		return errors.New("Wrapper with wrong signature")
	}

	resValues, err := f(t.GI, args.Arguments)
	if err != nil {
		return errors.Wrapf(err, "Can't call constructor %s", args.Name)
	}

	reply.Ret = resValues

	return nil
}

// Upstream returns RPC client connected to upstream server (goplugin)
func (gi *GoInsider) Upstream() (*rpc.Client, error) {
	if gi.UpstreamRPCClient != nil {
		return gi.UpstreamRPCClient, nil
	}

	client, err := rpc.DialHTTP("tcp", gi.UpstreamRPCAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial '%s'", gi.UpstreamRPCAddress)
	}

	gi.UpstreamRPCClient = client
	return gi.UpstreamRPCClient, nil
}

// ObtainCode returns path on the file system to the plugin, fetches it from a provider
// if it's not in the storage
func (gi *GoInsider) ObtainCode(ref core.RecordRef) (string, error) {
	path := filepath.Join(gi.dir, ref.String())
	_, err := os.Stat(path)

	if err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", errors.Wrap(err, "file !notexists()")
	}

	client, err := gi.Upstream()
	if err != nil {
		return "", err
	}

	log.Debugf("obtaining code %q", ref)
	res := rpctypes.UpGetCodeResp{}
	err = client.Call("RPC.GetCode", rpctypes.UpGetCodeReq{Code: ref}, &res)
	if err != nil {
		return "", errors.Wrap(err, "on calling main API")
	}

	err = ioutil.WriteFile(path, res.Code, 0666)
	if err != nil {
		return "", errors.Wrap(err, "on writing file down")
	}

	return path, nil
}

// Plugin loads Go plugin by reference and returns `*plugin.Plugin`
// ready to lookup symbols
func (gi *GoInsider) Plugin(ref core.RecordRef) (*plugin.Plugin, error) {
	gi.pluginsMutex.Lock()
	defer gi.pluginsMutex.Unlock()
	key := ref.String()
	if gi.plugins[key] != nil {
		return gi.plugins[key], nil
	}

	path, err := gi.ObtainCode(ref)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't obtain code")
	}

	log.Debugf("Opening plugin %q from file %q", ref, path)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open plugin")
	}

	gi.plugins[key] = p
	return p, nil
}

// RouteCall ...
func (gi *GoInsider) RouteCall(ref core.RecordRef, wait bool, method string, args []byte) ([]byte, error) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}

	req := rpctypes.UpRouteReq{
		Wait:      wait,
		Object:    ref,
		Method:    method,
		Arguments: args,
	}

	res := rpctypes.UpRouteResp{}
	err = client.Call("RPC.RouteCall", req, &res)
	if err != nil {
		return nil, errors.Wrap(err, "on calling main API")
	}

	return []byte(res.Result), nil
}

// RouteConstructorCall ...
func (gi *GoInsider) RouteConstructorCall(ref core.RecordRef, name string, args []byte) ([]byte, error) {
	client, err := gi.Upstream()
	if err != nil {
		return []byte{}, err
	}

	req := rpctypes.UpRouteConstructorReq{
		Reference:   ref,
		Constructor: name,
		Arguments:   args,
	}

	res := rpctypes.UpRouteConstructorResp{}
	err = client.Call("RPC.RouteConstructorCall", req, &res)
	if err != nil {
		return []byte{}, errors.Wrap(err, "on calling main API")
	}

	return res.Data, nil
}

// SaveAsChild ...
func (gi *GoInsider) SaveAsChild(parentRef, classRef core.RecordRef, data []byte) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpSaveAsChildReq{
		Parent: parentRef,
		Class:  classRef,
		Data:   data,
	}

	res := rpctypes.UpSaveAsChildResp{}
	err = client.Call("RPC.SaveAsChild", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return res.Reference, nil
}

// GetObjChildren ...
func (gi *GoInsider) GetObjChildren(obj core.RecordRef, class core.RecordRef) ([]core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}

	res := rpctypes.UpGetObjChildrenResp{}
	req := rpctypes.UpGetObjChildrenReq{Obj: obj, Class: class}
	err = client.Call("RPC.GetObjChildren", req, &res)
	if err != nil {
		return nil, errors.Wrap(err, "on calling main API RPC.GetObjChildren")
	}

	return res.Children, nil
}

// SaveAsDelegate ...
func (gi *GoInsider) SaveAsDelegate(intoRef, classRef core.RecordRef, data []byte) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpSaveAsDelegateReq{
		Into:  intoRef,
		Class: classRef,
		Data:  data,
	}

	res := rpctypes.UpSaveAsDelegateResp{}
	err = client.Call("RPC.SaveAsDelegate", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return res.Reference, nil
}

// GetDelegate ...
func (gi *GoInsider) GetDelegate(object, ofType core.RecordRef) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpGetDelegateReq{
		Object: object,
		OfType: ofType,
	}

	res := rpctypes.UpGetDelegateResp{}
	err = client.Call("RPC.GetDelegate", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return res.Object, nil
}

// Serialize - CBOR serializer wrapper: `what` -> `to`
func (gi *GoInsider) Serialize(what interface{}, to *[]byte) error {
	ch := new(codec.CborHandle)
	log.Debugf("serializing %+v", what)
	return codec.NewEncoderBytes(to, ch).Encode(what)
}

// Deserialize - CBOR de-serializer wrapper: `from` -> `into`
func (gi *GoInsider) Deserialize(from []byte, into interface{}) error {
	ch := new(codec.CborHandle)
	log.Debugf("de-serializing %+v", from)
	return codec.NewDecoderBytes(from, ch).Decode(into)
}
