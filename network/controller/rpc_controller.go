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

package controller

import (
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type RPCController struct {
	options     *common.Options
	hostNetwork network.HostNetwork
	methodTable map[string]core.RemoteProcedure
}

type RequestRPC struct {
	Method string
	Data   [][]byte
}

type ResponseRPC struct {
	Success bool
	Result  []byte
	Error   string
}

type RequestCascade struct {
	RPC     RequestRPC
	Cascade core.Cascade
}

type ResponseCascade struct {
	Success bool
	Error   string
}

func init() {
	gob.Register(&RequestRPC{})
	gob.Register(&ResponseRPC{})
	gob.Register(&RequestCascade{})
	gob.Register(&ResponseCascade{})
}

func (rpc *RPCController) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	rpc.methodTable[name] = method
}

func (rpc *RPCController) invoke(name string, data [][]byte) ([]byte, error) {
	method, exists := rpc.methodTable[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("RPC with name %s is not registered", name))
	}
	return method(data)
}

func (rpc *RPCController) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	return rpc.initCascadeSendMessage(data, false, method, [][]byte{message.ParcelToBytes(msg)})
}

func (rpc *RPCController) initCascadeSendMessage(data core.Cascade, findCurrentNode bool, method string, args [][]byte) error {
	if len(data.NodeIds) == 0 {
		return errors.New("node IDs list should not be empty")
	}
	if data.ReplicationFactor == 0 {
		return errors.New("replication factor should not be zero")
	}

	var nextNodes []core.RecordRef
	var err error

	if findCurrentNode {
		nodeID := rpc.hostNetwork.GetNodeID()
		nextNodes, err = cascade.CalculateNextNodes(data, &nodeID)
	} else {
		nextNodes, err = cascade.CalculateNextNodes(data, nil)
	}
	if err != nil {
		return errors.Wrap(err, "Failed to CalculateNextNodes")
	}
	if len(nextNodes) == 0 {
		return nil
	}

	var failedNodes []string
	for _, nextNode := range nextNodes {
		err = rpc.requestCascadeSendMessage(data, nextNode, method, args)
		if err != nil {
			log.Warnf("failed to send cascade message to node %s: %s", nextNode, err.Error())
			failedNodes = append(failedNodes, nextNode.String())
		}
	}

	if len(failedNodes) > 0 {
		return errors.New("failed to send cascade message to nodes: " + strings.Join(failedNodes, ", "))
	}

	return nil
}

func (rpc *RPCController) requestCascadeSendMessage(data core.Cascade, nodeID core.RecordRef,
	method string, args [][]byte) error {

	request := rpc.hostNetwork.NewRequestBuilder().Type(types.Cascade).Data(&RequestCascade{
		RPC: RequestRPC{
			Method: method,
			Data:   args,
		},
		Cascade: data,
	}).Build()

	future, err := rpc.hostNetwork.SendRequest(request, nodeID)
	if err != nil {
		return err
	}

	go func(f network.Future, duration time.Duration) {
		response, err := f.GetResponse(duration)
		if err != nil {
			log.Warnf("failed to get response to cascade message request from node %s: %s",
				future.GetRequest().GetSender(), err.Error())
			return
		}
		data := response.GetData().(*ResponseCascade)
		if !data.Success {
			log.Warnf("error response to cascade message request from node %s: %s",
				response.GetSender(), data.Error)
			return
		}
	}(future, rpc.options.PacketTimeout)

	return nil
}

func (rpc *RPCController) SendMessage(nodeID core.RecordRef, name string, msg core.Parcel) ([]byte, error) {
	start := time.Now()
	log.Debugf("SendParcel with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		name, message.ExtractTarget(msg).String())
	request := rpc.hostNetwork.NewRequestBuilder().Type(types.RPC).Data(&RequestRPC{
		Method: name,
		Data:   [][]byte{message.ParcelToBytes(msg)},
	}).Build()
	future, err := rpc.hostNetwork.SendRequest(request, nodeID)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending RPC request to node %s", nodeID.String())
	}
	response, err := future.GetResponse(rpc.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting RPC response from node %s", nodeID.String())
	}
	data := response.GetData().(*ResponseRPC)
	log.Debugf("Inside SendParcel: type - '%s', target - %s, caller - %s, targetRole - %s, time - %s",
		msg.Type(), message.ExtractTarget(msg), msg.GetCaller(), message.ExtractRole(msg), time.Since(start))
	if !data.Success {
		return nil, errors.New("RPC call returned error: " + data.Error)
	}
	return data.Result, nil
}

func (rpc *RPCController) processMessage(request network.Request) (network.Response, error) {
	payload := request.GetData().(*RequestRPC)
	result, err := rpc.invoke(payload.Method, payload.Data)
	if err != nil {
		return rpc.hostNetwork.BuildResponse(request, &ResponseRPC{Success: false, Error: err.Error()}), nil
	}
	return rpc.hostNetwork.BuildResponse(request, &ResponseRPC{Success: true, Result: result}), nil
}

func (rpc *RPCController) processCascade(request network.Request) (network.Response, error) {
	payload := request.GetData().(*RequestCascade)

	generalError := ""
	_, invokeErr := rpc.invoke(payload.RPC.Method, payload.RPC.Data)
	if invokeErr != nil {
		log.Debugf("failed to invoke RPC: %s", invokeErr.Error())
		generalError += invokeErr.Error() + "; "
	}
	sendErr := rpc.initCascadeSendMessage(payload.Cascade, true, payload.RPC.Method, payload.RPC.Data)
	if sendErr != nil {
		log.Debugf("failed to send message to next cascade layer: %s", sendErr.Error())
		generalError += sendErr.Error()
	}

	if generalError != "" {
		return rpc.hostNetwork.BuildResponse(request, &ResponseCascade{Success: false, Error: generalError}), nil
	}
	return rpc.hostNetwork.BuildResponse(request, &ResponseCascade{Success: true}), nil
}

func (rpc *RPCController) Start() {
	rpc.hostNetwork.RegisterRequestHandler(types.RPC, rpc.processMessage)
	rpc.hostNetwork.RegisterRequestHandler(types.Cascade, rpc.processCascade)
}

func NewRPCController(options *common.Options, hostNetwork network.HostNetwork) *RPCController {
	return &RPCController{options: options, hostNetwork: hostNetwork, methodTable: make(map[string]core.RemoteProcedure)}
}
