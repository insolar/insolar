// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package common

import (
	"context"
	"net"
	"net/rpc"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type LogicRunnerRPCStub interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	server    *rpc.Server
	methods   LogicRunnerRPCStub
	listener  net.Listener
	proto     string
	listen    string
	isStarted bool
}

func NewRPC(lr LogicRunnerRPCStub, cfg *configuration.LogicRunner) *RPC {
	rpcService := &RPC{
		server:  rpc.NewServer(),
		methods: lr,
		proto:   cfg.RPCProtocol,
		listen:  cfg.RPCListen,
	}
	if err := rpcService.server.RegisterName("RPC", rpcService.methods); err != nil {
		panic("Fail to register LogicRunner RPC Service: " + err.Error())
	}

	return rpcService
}

// StartRPC starts RPC server for isolated executors to use
func (rpc *RPC) Start(ctx context.Context) {
	if rpc == nil {
		panic("Calling start on nil")
	}
	var err error
	logger := inslogger.FromContext(ctx)

	rpc.listener, err = net.Listen(rpc.proto, rpc.listen)
	if err != nil {
		logger.Fatalf("couldn't setup listener on %q over %q: %s", rpc.listen, rpc.proto, err)
	}

	logger.Infof("starting LogicRunner RPC service on %q over %q", rpc.listen, rpc.proto)
	rpc.isStarted = true

	go func() {
		rpc.server.Accept(rpc.listener)
		logger.Info("LogicRunner RPC service stopped")
	}()
}

func (rpc *RPC) Stop(_ context.Context) error {
	if rpc == nil {
		return nil
	}
	if rpc.isStarted {
		rpc.isStarted = false
		if err := rpc.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}
