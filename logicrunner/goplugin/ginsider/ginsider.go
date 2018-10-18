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
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/tylerb/gls"
)

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir              string
	UpstreamProtocol string
	UpstreamAddress  string
	UpstreamClient   *rpc.Client
	plugins          map[string]*plugin.Plugin
	pluginsMutex     sync.Mutex
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path, network, address string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	res := GoInsider{dir: path, UpstreamProtocol: network, UpstreamAddress: address}
	res.plugins = make(map[string]*plugin.Plugin)
	return &res
}

// RPC struct with methods representing RPC interface of this code runner
type RPC struct {
	GI *GoInsider
}

func recoverRPC(err *error) {
	if r := recover(); r != nil {
		if err != nil {
			if *err == nil {
				*err = errors.New(fmt.Sprint(r))
			} else {
				*err = errors.New(fmt.Sprint(*err, r))
			}
		}
		log.Errorln("panic: ", r)
	}
}

// CallMethod is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallMethod(args rpctypes.DownCallMethodReq, reply *rpctypes.DownCallMethodResp) (err error) {
	log.Debugf("Calling method %q on object %q", args.Method, args.Context.Callee)
	defer recoverRPC(&err)

	p, err := t.GI.Plugin(args.Code)
	if err != nil {
		return errors.Wrapf(err, "Couldn't get plugin by code reference %s", args.Code.String())
	}

	symbol, err := p.Lookup("INSMETHOD_" + args.Method)
	if err != nil {
		return errors.Wrapf(
			err, "Can't find wrapper for %s (code ref: %s)",
			args.Method, args.Code.String(),
		)
	}

	wrapper, ok := symbol.(func(object []byte, data []byte) ([]byte, []byte, error))
	if !ok {
		return errors.New("Wrapper with wrong signature")
	}

	gls.Set("ctx", args.Context)
	state, result, err := wrapper(args.Data, args.Arguments) // may be entire args???
	gls.Cleanup()

	if err != nil {
		return errors.Wrapf(err, "Method call returned error")
	}
	reply.Data = state
	reply.Ret = result
	return nil
}

// CallConstructor is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallConstructor(args rpctypes.DownCallConstructorReq, reply *rpctypes.DownCallConstructorResp) (err error) {
	log.Debugf("Calling constructor %q in code %q", args.Name, args.Code)
	defer recoverRPC(&err)

	p, err := t.GI.Plugin(args.Code)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("INSCONSTRUCTOR_" + args.Name)
	if err != nil {
		return errors.Wrapf(err, "Can't find wrapper for %s", args.Name)
	}

	f, ok := symbol.(func(data []byte) ([]byte, error))
	if !ok {
		return errors.New("Wrapper with wrong signature")
	}

	resValues, err := f(args.Arguments)
	if err != nil {
		return errors.Wrapf(err, "Can't call constructor %s", args.Name)
	}

	reply.Ret = resValues

	return nil
}

// Upstream returns RPC client connected to upstream server (goplugin)
func (gi *GoInsider) Upstream() (*rpc.Client, error) {
	if gi.UpstreamClient != nil {
		return gi.UpstreamClient, nil
	}

	client, err := rpc.Dial(gi.UpstreamProtocol, gi.UpstreamAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial '%s' over %s", gi.UpstreamAddress, gi.UpstreamProtocol)
	}

	gi.UpstreamClient = client
	return gi.UpstreamClient, nil
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
	err = client.Call("RPC.GetCode", rpctypes.UpGetCodeReq{Code: ref, MType: core.MachineTypeGoPlugin}, &res)
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

// MakeUpBaseReq makes base of request from current CallContext
func MakeUpBaseReq() rpctypes.UpBaseReq {
	if ctx, ok := gls.Get("ctx").(*core.LogicCallContext); ok {
		return rpctypes.UpBaseReq{
			Callee:  *ctx.Callee,
			Request: *ctx.Request,
		}
	}
	panic("Wrong or unexistent context")
}

// RouteCall ...
func (gi *GoInsider) RouteCall(ref core.RecordRef, wait bool, method string, args []byte) ([]byte, error) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}
	req := rpctypes.UpRouteReq{
		UpBaseReq: MakeUpBaseReq(),
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

// SaveAsChild ...
func (gi *GoInsider) SaveAsChild(parentRef, classRef core.RecordRef, constructorName string, argsSerialized []byte) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq:       MakeUpBaseReq(),
		Parent:          parentRef,
		Class:           classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	res := rpctypes.UpSaveAsChildResp{}
	err = client.Call("RPC.SaveAsChild", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return *res.Reference, nil
}

// GetObjChildren ...
func (gi *GoInsider) GetObjChildren(obj core.RecordRef, class core.RecordRef) ([]core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}

	res := rpctypes.UpGetObjChildrenResp{}
	req := rpctypes.UpGetObjChildrenReq{
		UpBaseReq: MakeUpBaseReq(),
		Obj:       obj,
		Class:     class,
	}
	err = client.Call("RPC.GetObjChildren", req, &res)
	if err != nil {
		return nil, errors.Wrap(err, "on calling main API RPC.GetObjChildren")
	}

	return res.Children, nil
}

// SaveAsDelegate ...
func (gi *GoInsider) SaveAsDelegate(intoRef, classRef core.RecordRef, constructorName string, argsSerialized []byte) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpSaveAsDelegateReq{
		UpBaseReq:       MakeUpBaseReq(),
		Into:            intoRef,
		Class:           classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	res := rpctypes.UpSaveAsDelegateResp{}
	err = client.Call("RPC.SaveAsDelegate", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return *res.Reference, nil
}

// GetDelegate ...
func (gi *GoInsider) GetDelegate(object, ofType core.RecordRef) (core.RecordRef, error) {
	client, err := gi.Upstream()
	if err != nil {
		return core.NewRefFromBase58(""), err
	}

	req := rpctypes.UpGetDelegateReq{
		UpBaseReq: MakeUpBaseReq(),
		Object:    object,
		OfType:    ofType,
	}

	res := rpctypes.UpGetDelegateResp{}
	err = client.Call("RPC.GetDelegate", req, &res)
	if err != nil {
		return core.NewRefFromBase58(""), errors.Wrap(err, "on calling main API")
	}

	return res.Object, nil
}

// DeactivateObject ...
func (gi *GoInsider) DeactivateObject(object core.RecordRef) error {
	client, err := gi.Upstream()
	if err != nil {
		return err
	}

	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: MakeUpBaseReq(),
		Object:    object,
	}

	res := rpctypes.UpDeactivateObjectResp{}
	err = client.Call("RPC.DeactivateObject", req, &res)
	if err != nil {
		return errors.Wrap(err, "on calling main API")
	}

	return nil
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

// MakeErrorSerializable converts errors satisfying error interface to foundation.Error
func (gi *GoInsider) MakeErrorSerializable(e error) error {
	if e == nil || e == (*foundation.Error)(nil) || reflect.ValueOf(e).IsNil() {
		return nil
	}
	return &foundation.Error{S: e.Error()}
}
