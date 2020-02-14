// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package ginsider

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/metrics"
)

type pluginRec struct {
	sync.Mutex
	plugin *plugin.Plugin
}

// GoInsider is an RPC interface to run code of plugins
type GoInsider struct {
	dir              string
	upstreamProtocol string
	upstreamAddress  string

	upstreamMutex  sync.Mutex // lock UpstreamClient change
	UpstreamClient *rpc.Client

	plugins      map[insolar.Reference]*pluginRec
	pluginsMutex sync.Mutex

	lrCommon.Serializer
	lrCommon.SystemError
}

// NewGoInsider creates a new GoInsider instance validating arguments
func NewGoInsider(path, network, address string) *GoInsider {
	//TODO: check that path exist, it's a directory and writable
	res := GoInsider{dir: path, upstreamProtocol: network, upstreamAddress: address}
	res.plugins = make(map[insolar.Reference]*pluginRec)
	lrCommon.CurrentProxyCtx = &res
	res.Serializer = lrCommon.NewCBORSerializer()
	res.SystemError = lrCommon.NewSystemError()
	return &res
}

// RPC struct with methods representing RPC interface of this code runner
type RPC struct {
	GI *GoInsider
}

func recoverRPC(ctx context.Context, err *error) {
	if r := recover(); r != nil {
		if err != nil {
			if *err == nil {
				*err = errors.New(fmt.Sprint(r))
			} else {
				*err = errors.New(fmt.Sprint(*err, r))
			}
		}
		inslogger.FromContext(ctx).Error("panic: ", r, string(debug.Stack()))
	}
}

// CallMethod is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallMethod(args rpctypes.DownCallMethodReq, reply *rpctypes.DownCallMethodResp) (err error) {
	start := time.Now()
	metrics.InsgorundCallsTotal.Inc()
	ctx := inslogger.ContextWithTrace(context.Background(), args.Context.TraceID)
	inslogger.FromContext(ctx).Debugf("Calling method %q on object %q", args.Method, args.Context.Callee)
	defer recoverRPC(ctx, &err)

	foundation.SetLogicalContext(args.Context)
	defer foundation.ClearContext()

	p, err := t.GI.Plugin(ctx, args.Code)
	if err != nil {
		return errors.Wrapf(err, "Couldn't get plugin by code reference %s", args.Code.String())
	}

	if args.Context.Caller.IsEmpty() {
		attr, err := p.Lookup("INSATTR_" + args.Method + "_API")
		if err != nil {
			return errors.Wrapf(
				err, "Calling non INSATTRAPI method %s (code ref: %s)",
				args.Method, args.Code.String(),
			)
		}
		api, ok := attr.(*bool)
		if !ok {
			return errors.Errorf("INSATTRAPI attribute for method %s is not boolean", args.Method)
		}
		if !*api {
			return errors.Errorf("Calling non INSATTRAPI method ")
		}
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

	state, result, err := wrapper(args.Data, args.Arguments) // may be entire args???
	if err != nil {
		return errors.Wrapf(err, "Method call returned error")
	}

	reply.Data = state
	reply.Ret = result

	metrics.InsgorundContractExecutionTime.WithLabelValues(args.Method).Observe(time.Since(start).Seconds())

	return nil
}

// CallConstructor is an RPC that runs a method on an object and
// returns a new state of the object and result of the method
func (t *RPC) CallConstructor(args rpctypes.DownCallConstructorReq, reply *rpctypes.DownCallConstructorResp) (err error) {
	metrics.InsgorundCallsTotal.Inc()
	ctx := inslogger.ContextWithTrace(context.Background(), args.Context.TraceID)
	inslogger.FromContext(ctx).Debugf("Calling constructor %q in code %q", args.Name, args.Code)
	defer recoverRPC(ctx, &err)

	foundation.SetLogicalContext(args.Context)
	defer foundation.ClearContext()

	p, err := t.GI.Plugin(ctx, args.Code)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("INSCONSTRUCTOR_" + args.Name)
	if err != nil {
		return errors.Wrapf(err, "Can't find wrapper for %s", args.Name)
	}

	f, ok := symbol.(func(ref insolar.Reference, data []byte) ([]byte, []byte, error))
	if !ok {
		return errors.New("Wrapper with wrong signature")
	}

	objRef := insolar.NewReference(*args.Context.Request.GetLocal())
	state, result, err := f(*objRef, args.Arguments)
	if err != nil {
		return errors.Wrapf(err, "Can't call constructor %s", args.Name)
	}

	reply.Data = state
	reply.Ret = result

	return nil
}

// Upstream returns RPC client connected to upstream server (goplugin)
func (gi *GoInsider) Upstream() (*rpc.Client, error) {
	gi.upstreamMutex.Lock()
	defer gi.upstreamMutex.Unlock()
	if gi.UpstreamClient != nil {
		return gi.UpstreamClient, nil
	}

	client, err := rpc.Dial(gi.upstreamProtocol, gi.upstreamAddress)
	if err != nil {
		log.Fatalf("can't connect to upstream, protocol: %s, address: %s", gi.upstreamProtocol, gi.upstreamAddress)
		os.Exit(0)
	}

	gi.UpstreamClient = client
	return gi.UpstreamClient, nil
}

// ObtainCode returns path on the file system to the plugin, fetches it from a provider
// if it's not in the storage
func (gi *GoInsider) ObtainCode(ctx context.Context, ref insolar.Reference) (string, error) {
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

	inslogger.FromContext(ctx).Debugf("obtaining code %q", ref)
	req := rpctypes.UpGetCodeReq{
		UpBaseReq: MakeUpBaseReq(),
		Code:      ref,
		MType:     insolar.MachineTypeGoPlugin,
	}
	res := rpctypes.UpGetCodeResp{}
	err = client.Call("RPC.GetCode", req, &res)
	if err != nil {
		if err == rpc.ErrShutdown {
			log.Error("Insgorund can't connect to Insolard")
			os.Exit(0)
		}
		return "", errors.Wrap(err, "[ ObtainCode ] on calling main API")
	}

	err = ioutil.WriteFile(path, res.Code, 0666)
	if err != nil {
		return "", errors.Wrap(err, "[ ObtainCode ] on writing file down")
	}

	return path, nil
}

// Plugin loads Go plugin by reference and returns `*plugin.Plugin`
// ready to lookup symbols
func (gi *GoInsider) Plugin(ctx context.Context, ref insolar.Reference) (*plugin.Plugin, error) {
	rec := gi.getPluginRec(ref)

	rec.Lock()
	defer rec.Unlock()

	if rec.plugin != nil {
		return rec.plugin, nil
	}

	path, err := gi.ObtainCode(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't obtain code")
	}

	inslogger.FromContext(ctx).Debugf("Opening plugin %q from file %q", ref, path)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open plugin")
	}

	rec.plugin = p
	return p, nil
}

// getPluginRec return existed gi.plugins[ref] or create a new one
// also set gi.plugins[ref].Lock()
func (gi *GoInsider) getPluginRec(ref insolar.Reference) *pluginRec {
	gi.pluginsMutex.Lock()
	defer gi.pluginsMutex.Unlock()

	if gi.plugins[ref] == nil {
		gi.plugins[ref] = &pluginRec{}
	}
	res := gi.plugins[ref]
	return res
}

// MakeUpBaseReq makes base of request from current CallContext
func MakeUpBaseReq() rpctypes.UpBaseReq {
	callCtx := foundation.GetLogicalContext()

	return rpctypes.UpBaseReq{
		Mode:            callCtx.Mode,
		Callee:          *callCtx.Callee,
		CalleePrototype: *callCtx.Prototype,
		Request:         *callCtx.Request,
	}
}

// RouteCall ...
func (gi *GoInsider) RouteCall(ref insolar.Reference, immutable bool, saga bool, method string, args []byte, proxyPrototype insolar.Reference) ([]byte, error) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}
	if gi.GetSystemError() != nil {
		return nil, gi.GetSystemError()
	}

	req := rpctypes.UpRouteReq{
		UpBaseReq: MakeUpBaseReq(),
		Immutable: immutable,
		Saga:      saga,
		Object:    ref,
		Method:    method,
		Arguments: args,
		Prototype: proxyPrototype,
	}

	res := rpctypes.UpRouteResp{}
	err = client.Call("RPC.RouteCall", req, &res)
	if err != nil {
		gi.SetSystemError(err)
		if err == rpc.ErrShutdown {
			log.Error("Insgorund can't connect to Insolard")
			os.Exit(0)
		}
		return nil, errors.Wrap(err, "[ RouteCall ] on calling main API")
	}

	return []byte(res.Result), nil
}

// SaveAsChild ...
func (gi *GoInsider) SaveAsChild(
	parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte,
) (
	[]byte, error,
) {
	client, err := gi.Upstream()
	if err != nil {
		return nil, err
	}
	if gi.GetSystemError() != nil {
		return nil, gi.GetSystemError()
	}

	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq:       MakeUpBaseReq(),
		Parent:          parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	res := rpctypes.UpSaveAsChildResp{}
	err = client.Call("RPC.SaveAsChild", req, &res)
	if err != nil {
		gi.SetSystemError(err)
		if err == rpc.ErrShutdown {
			log.Error("Insgorund can't connect to Insolard")
			os.Exit(0)
		}
		return nil, errors.Wrap(err, "[ SaveAsChild ] on calling main API")
	}

	return res.Result, nil
}

// DeactivateObject ...
func (gi *GoInsider) DeactivateObject(object insolar.Reference) error {
	client, err := gi.Upstream()
	if err != nil {
		return err
	}
	if gi.GetSystemError() != nil {
		return gi.GetSystemError()
	}

	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: MakeUpBaseReq(),
	}

	res := rpctypes.UpDeactivateObjectResp{}
	err = client.Call("RPC.DeactivateObject", req, &res)
	if err != nil {
		gi.SetSystemError(err)
		if err == rpc.ErrShutdown {
			log.Error("Insgorund can't connect to Insolard")
			os.Exit(0)
		}
		return errors.Wrap(err, "[ DeactivateObject ] on calling main API")
	}

	return nil
}

// Serialize - CBOR serializer wrapper: `what` -> `to`
func (gi *GoInsider) Serialize(what interface{}, to *[]byte) (err error) {
	*to, err = insolar.Serialize(what)
	return err
}

// Deserialize - CBOR de-serializer wrapper: `from` -> `into`
func (gi *GoInsider) Deserialize(from []byte, into interface{}) error {
	return insolar.Deserialize(from, into)
}

// MakeErrorSerializable converts errors satisfying error interface to foundation.Error
func (gi *GoInsider) MakeErrorSerializable(e error) error {
	if e == nil || e == (*foundation.Error)(nil) || reflect.ValueOf(e).IsNil() {
		return nil
	}
	return &foundation.Error{S: e.Error()}
}

// AddPlugin inject plugin by ref in gi memory
func (gi *GoInsider) AddPlugin(ref insolar.Reference, path string) error {
	rec := gi.getPluginRec(ref)

	rec.Lock()
	defer rec.Unlock()

	if rec.plugin != nil {
		return errors.New("ref already in use")
	}

	p, err := plugin.Open(path)
	if err != nil {
		return errors.Wrap(err, "[ AddPlugin ] couldn't open plugin")
	}

	inslogger.FromContext(context.TODO()).Debugf("AddPlugin plugins %+v", gi.plugins)
	rec.plugin = p
	return nil
}
