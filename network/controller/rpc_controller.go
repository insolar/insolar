// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package controller

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

// RemoteProcedure is remote procedure call function.
type RemoteProcedure func(ctx context.Context, args []byte) ([]byte, error)

//go:generate minimock -i github.com/insolar/insolar/network/controller.RPCController -o . -s _mock.go -g

type RPCController interface {
	SendBytes(ctx context.Context, nodeID insolar.Reference, name string, msgBytes []byte) ([]byte, error)
	RemoteProcedureRegister(name string, method RemoteProcedure)
}

type rpcController struct {
	Network network.HostNetwork `inject:""`

	options     *network.Options
	methodTable map[string]RemoteProcedure
}

func (rpc *rpcController) RemoteProcedureRegister(name string, method RemoteProcedure) {
	rpc.methodTable[name] = method
}

func (rpc *rpcController) invoke(ctx context.Context, name string, data []byte) ([]byte, error) {
	method, exists := rpc.methodTable[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("RPC with name %s is not registered", name))
	}
	return method(ctx, data)
}

func (rpc *rpcController) SendBytes(ctx context.Context, nodeID insolar.Reference, name string, msgBytes []byte) ([]byte, error) {
	request := &packet.RPCRequest{
		Method: name,
		Data:   msgBytes,
	}

	future, err := rpc.Network.SendRequest(ctx, types.RPC, request, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending RPC request to node %s", nodeID.String())
	}
	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"request_id":     future.Request().GetRequestID(),
		"target_node_id": nodeID.String(),
	})
	logger.Debug("sent RPC request")
	response, err := future.WaitResponse(rpc.options.AckPacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting RPC response from node %s", nodeID.String())
	}
	logger.Debug("received RPC response")
	if response.GetResponse() == nil || response.GetResponse().GetRPC() == nil {
		logger.Warnf("Error getting RPC response from node %s: "+
			"got invalid response protobuf message: %s", nodeID, response)
	}
	data := response.GetResponse().GetRPC()
	if data.Result == nil {
		return nil, errors.New("RPC call returned error: " + data.Error)
	}
	stats.Record(ctx, statParcelsReplySizeBytes.M(int64(len(data.Result))))
	return data.Result, nil
}

func (rpc *rpcController) processMessage(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetRPC() == nil {
		inslogger.FromContext(ctx).Warnf("process RPC: got invalid request protobuf message: %s", request)
	}

	ctx = insmetrics.InsertTag(ctx, tagPacketType, request.GetType().String())
	stats.Record(ctx, statPacketsReceived.M(1))

	payload := request.GetRequest().GetRPC()
	result, err := rpc.invoke(ctx, payload.Method, payload.Data)
	if err != nil {
		return rpc.Network.BuildResponse(ctx, request, &packet.RPCResponse{Error: err.Error()}), nil
	}
	return rpc.Network.BuildResponse(ctx, request, &packet.RPCResponse{Result: result}), nil
}

func (rpc *rpcController) Init(ctx context.Context) error {
	rpc.Network.RegisterRequestHandler(types.RPC, rpc.processMessage)
	return nil
}

func NewRPCController(options *network.Options) RPCController {
	return &rpcController{options: options, methodTable: make(map[string]RemoteProcedure)}
}
