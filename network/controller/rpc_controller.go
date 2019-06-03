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
	"strings"
	"time"

	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

type RPCController interface {
	component.Initer

	// hack for DI, else we receive ServiceNetwork injection in RPCController instead of rpcController that leads to stack overflow
	IAmRPCController()

	SendMessage(nodeID insolar.Reference, name string, msg insolar.Parcel) ([]byte, error)
	SendBytes(ctx context.Context, nodeID insolar.Reference, name string, msgBytes []byte) ([]byte, error)
	SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error
	RemoteProcedureRegister(name string, method insolar.RemoteProcedure)
}

type rpcController struct {
	Scheme      insolar.PlatformCryptographyScheme `inject:""`
	Network     network.HostNetwork                `inject:""`
	NodeNetwork insolar.NodeNetwork                `inject:""`

	options     *common.Options
	methodTable map[string]insolar.RemoteProcedure
}

func (rpc *rpcController) IAmRPCController() {
	// hack for DI, else we receive ServiceNetwork injection in RPCController instead of rpcController that leads to stack overflow
}

func (rpc *rpcController) RemoteProcedureRegister(name string, method insolar.RemoteProcedure) {
	rpc.methodTable[name] = method
}

func (rpc *rpcController) invoke(ctx context.Context, name string, data []byte) ([]byte, error) {
	method, exists := rpc.methodTable[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("RPC with name %s is not registered", name))
	}
	return method(ctx, data)
}

func (rpc *rpcController) SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	ctx, span := instracer.StartSpan(context.Background(), "RPCController.SendCascadeMessage")
	span.AddAttributes(
		trace.StringAttribute("method", method),
		trace.StringAttribute("msg.Type", msg.Type().String()),
		trace.StringAttribute("msg.DefaultTarget", msg.DefaultTarget().String()),
	)
	defer span.End()
	ctx = msg.Context(ctx)
	return rpc.initCascadeSendMessage(ctx, data, false, method, message.ParcelToBytes(msg))
}

func (rpc *rpcController) initCascadeSendMessage(ctx context.Context, data insolar.Cascade,
	findCurrentNode bool, method string, args []byte) error {

	_, span := instracer.StartSpan(context.Background(), "RPCController.initCascadeSendMessage")
	span.AddAttributes(
		trace.StringAttribute("method", method),
	)
	defer span.End()
	if len(data.NodeIds) == 0 {
		return errors.New("node IDs list should not be empty")
	}
	if data.ReplicationFactor == 0 {
		return errors.New("replication factor should not be zero")
	}

	var nextNodes []insolar.Reference
	var err error

	if findCurrentNode {
		nodeID := rpc.NodeNetwork.GetOrigin().ID()
		nextNodes, err = cascade.CalculateNextNodes(rpc.Scheme, data, &nodeID)
	} else {
		nextNodes, err = cascade.CalculateNextNodes(rpc.Scheme, data, nil)
	}
	if err != nil {
		return errors.Wrap(err, "Failed to CalculateNextNodes")
	}
	if len(nextNodes) == 0 {
		return nil
	}

	var failedNodes []string
	for _, nextNode := range nextNodes {
		err = rpc.requestCascadeSendMessage(ctx, data, nextNode, method, args)
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Failed to send cascade message to node %s: %s", nextNode, err.Error())
			failedNodes = append(failedNodes, nextNode.String())
		}
	}

	if len(failedNodes) > 0 {
		return errors.New("Failed to send cascade message to nodes: " + strings.Join(failedNodes, ", "))
	}
	inslogger.FromContext(ctx).Debug("Cascade message successfully sent to all nodes of the next layer")
	return nil
}

func (rpc *rpcController) requestCascadeSendMessage(ctx context.Context, data insolar.Cascade, nodeID insolar.Reference,
	method string, args []byte) error {

	_, span := instracer.StartSpan(context.Background(), "RPCController.requestCascadeSendMessage")
	defer span.End()
	request := &packet.CascadeRequest{
		TraceID: inslogger.TraceID(ctx),
		RPC: &packet.RPCRequest{
			Method: method,
			Data:   args,
		},
		Cascade: &packet.Cascade{
			NodeIds:           data.NodeIds,
			Entropy:           data.Entropy,
			ReplicationFactor: uint32(data.ReplicationFactor),
		},
	}

	future, err := rpc.Network.SendRequest(ctx, types.Cascade, request, nodeID)
	if err != nil {
		return err
	}

	go func(ctx context.Context, receiver insolar.Reference, f network.Future, duration time.Duration) {
		response, err := f.WaitResponse(duration)
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Failed to get response to cascade message request from node %s: %s",
				receiver, err.Error())
			return
		}
		if response.GetResponse() == nil || response.GetResponse().GetBasic() == nil {
			inslogger.FromContext(ctx).Warnf("Failed to get response to cascade message request from node %s: "+
				"got invalid response protobuf message: %s", receiver, response)
		}
		data := response.GetResponse().GetBasic()
		if !data.Success {
			inslogger.FromContext(ctx).Warnf("Error response to cascade message request from node %s: %s",
				response.GetSender(), data.Error)
			return
		}
	}(ctx, nodeID, future, rpc.options.PacketTimeout)

	return nil
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
	response, err := future.WaitResponse(rpc.options.PacketTimeout)
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

func (rpc *rpcController) SendMessage(nodeID insolar.Reference, name string, msg insolar.Parcel) ([]byte, error) {
	msgBytes := message.ParcelToBytes(msg)
	ctx := context.Background() // TODO: ctx as argument
	ctx = insmetrics.InsertTag(ctx, tagMessageType, msg.Type().String())
	stats.Record(ctx, statParcelsSentSizeBytes.M(int64(len(msgBytes))))
	request := &packet.RPCRequest{
		Method: name,
		Data:   msgBytes,
	}

	start := time.Now()
	ctx = msg.Context(ctx)
	logger := inslogger.FromContext(ctx)
	// TODO: change sendrequest signature to have request as argument
	logger.Debugf("Before SendParcel with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		name, msg.DefaultTarget().String())
	future, err := rpc.Network.SendRequest(ctx, types.RPC, request, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending RPC request to node %s", nodeID.String())
	}
	logger.Debugf("SendParcel with nodeID = %s method = %s, message reference = %s, RequestID = %d", nodeID.String(),
		name, msg.DefaultTarget().String(), future.Request().GetRequestID())
	response, err := future.WaitResponse(rpc.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting RPC response from node %s", nodeID.String())
	}
	if response.GetResponse() == nil || response.GetResponse().GetRPC() == nil {
		inslogger.FromContext(ctx).Warnf("Error getting RPC response from node %s: "+
			"got invalid response protobuf message: %s", nodeID, response)
	}
	data := response.GetResponse().GetRPC()
	logger.Debugf("Inside SendParcel: type - '%s', target - %s, caller - %s, targetRole - %s, time - %s",
		msg.Type(), msg.DefaultTarget(), msg.GetCaller(), msg.DefaultRole(), time.Since(start))
	if data.Result == nil {
		return nil, errors.New("RPC call returned error: " + data.Error)
	}
	stats.Record(ctx, statParcelsReplySizeBytes.M(int64(len(data.Result))))
	return data.Result, nil
}

func (rpc *rpcController) processMessage(ctx context.Context, request network.Packet) (network.Packet, error) {
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

func (rpc *rpcController) processCascade(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetCascade() == nil {
		inslogger.FromContext(ctx).Warnf("process cascade: got invalid request protobuf message: %s", request)
	}

	payload := request.GetRequest().GetCascade()
	ctx, logger := inslogger.WithTraceField(ctx, payload.TraceID)

	generalError := ""
	_, invokeErr := rpc.invoke(ctx, payload.RPC.Method, payload.RPC.Data)
	if invokeErr != nil {
		logger.Debugf("failed to invoke RPC: %s", invokeErr.Error())
		generalError += invokeErr.Error() + "; "
	}
	cascade := insolar.Cascade{
		NodeIds:           payload.Cascade.NodeIds,
		Entropy:           payload.Cascade.Entropy,
		ReplicationFactor: uint(payload.Cascade.ReplicationFactor),
	}
	sendErr := rpc.initCascadeSendMessage(ctx, cascade, true, payload.RPC.Method, payload.RPC.Data)
	if sendErr != nil {
		logger.Debugf("failed to send message to next cascade layer: %s", sendErr.Error())
		generalError += sendErr.Error()
	}

	if generalError != "" {
		return rpc.Network.BuildResponse(ctx, request, &packet.BasicResponse{Success: false, Error: generalError}), nil
	}
	return rpc.Network.BuildResponse(ctx, request, &packet.BasicResponse{Success: true}), nil
}

func (rpc *rpcController) Init(ctx context.Context) error {
	rpc.Network.RegisterRequestHandler(types.RPC, rpc.processMessage)
	rpc.Network.RegisterRequestHandler(types.Cascade, rpc.processCascade)
	return nil
}

func NewRPCController(options *common.Options) RPCController {
	return &rpcController{options: options, methodTable: make(map[string]insolar.RemoteProcedure)}
}
