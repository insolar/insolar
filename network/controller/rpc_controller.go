//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
	response, err := future.WaitResponse(rpc.options.AckPacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting RPC response from node %s", nodeID.String())
	}
	if response.GetResponse() == nil || response.GetResponse().GetRPC() == nil {
		inslogger.FromContext(ctx).Warnf("Error getting RPC response from node %s: "+
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
