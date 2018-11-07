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

package dhtnetwork

import (
	"strings"
	"time"

	consensus2 "github.com/insolar/insolar/network/dhtnetwork/consensus"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/id"
	"github.com/insolar/insolar/network/transport/packet/types"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/insolar/insolar/network/dhtnetwork/resolver"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/pkg/errors"
)

type Wrapper struct {
	HostNetwork hosthandler.HostHandler
}

// initCascadeSendMessage initiates the RPC call on target host and sends messages to next cascade layers
func (w *Wrapper) initCascadeSendMessage(data core.Cascade, findCurrentNode bool, method string, args [][]byte) error {
	if len(data.NodeIds) == 0 {
		return errors.New("node IDs list should not be empty")
	}
	if data.ReplicationFactor == 0 {
		return errors.New("replication factor should not be zero")
	}

	var nextNodes []core.RecordRef
	var err error

	if findCurrentNode {
		nodeID := w.HostNetwork.GetNodeID()
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
		hostID := resolver.ResolveHostID(nextNode)
		err = w.HostNetwork.CascadeSendMessage(data, hostID, method, args)
		if err != nil {
			log.Debugln("failed to send cascade message: ", err)
			failedNodes = append(failedNodes, nextNode.String())
		}
	}

	if len(failedNodes) > 0 {
		return errors.New("failed to send cascade message to nodes: " + strings.Join(failedNodes, ", "))
	}

	return nil
}

// SendMessage send message to nodeID
func (w *Wrapper) SendMessage(nodeID core.RecordRef, method string, msg core.SignedMessage) ([]byte, error) {
	start := time.Now()
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	hostID := resolver.ResolveHostID(nodeID)
	buff, err := message.SignedToBytes(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize event")
	}

	log.Debugf("SendMessage with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		method,  message.ExtractTarget(msg).String())

	metrics.NetworkMessageSentTotal.Inc()
	res, err := w.HostNetwork.RemoteProcedureCall(CreateDHTContext(w.HostNetwork), hostID, method, [][]byte{buff})
	log.Debugf("Inside SendMessage: type - '%s', target - %s, caller - %s, targetRole - %s, time - %s",
		msg.Type(), message.ExtractTarget(msg).String(), msg.GetCaller(), message.ExtractRole(msg), time.Since(start))
	return res, err
}

func (w *Wrapper) SendCascadeMessage(data core.Cascade, method string, msg core.SignedMessage) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	buff, err := message.SignedToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize event")
	}

	return w.initCascadeSendMessage(data, false, method, [][]byte{buff})
}

// RemoteProcedureRegister register remote procedure that will be executed when message is received
func (w *Wrapper) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	w.HostNetwork.RemoteProcedureRegister(name, method)
}

// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
func (w *Wrapper) Bootstrap() error {
	err := w.HostNetwork.Bootstrap()
	if err != nil {
		return errors.Wrap(err, "error bootstraping DHT network")
	}
	w.HostNetwork.GetHostsFromBootstrap()
	return nil
}

// AnalyzeNetwork legacy method for old DHT network (should be removed in
func (w *Wrapper) AnalyzeNetwork() error {
	ctx := CreateDHTContext(w.HostNetwork)
	err := w.HostNetwork.ObtainIP()
	if err != nil {
		return errors.Wrap(err, "Failed to ObtainIP")
	}

	err = w.HostNetwork.AnalyzeNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to AnalyzeNetwork")
	}
	return nil
}

// Authorize start authorization process on discovery node.
func (w *Wrapper) Authorize() error {
	return w.HostNetwork.StartAuthorize()
}

// GetNodeID get self node id (should be removed in far future)
func (w *Wrapper) GetNodeID() core.RecordRef {
	return w.HostNetwork.GetNodeID()
}

// Start listening to network requests.
func (w *Wrapper) Start() {
	go func() {
		err := w.HostNetwork.Listen()
		if err != nil {
			log.Error(err)
		}
	}()
}

// Disconnect stop listening to network requests.
func (w *Wrapper) Stop() {
	log.Infoln("Stop network")
	w.HostNetwork.Disconnect()
}

// PublicAddress returns public address that can be published for all nodes.
func (w *Wrapper) PublicAddress() string {
	return w.HostNetwork.GetOriginHost().Address.String()
}

// SendRequest send request to a remote node.
func (w *Wrapper) SendRequest(network.Request, core.RecordRef) (network.Future, error) {
	panic("not used in DHT implementation")
}

// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
func (w *Wrapper) RegisterRequestHandler(t types.PacketType, handler network.RequestHandler) {
	panic("not used in DHT implementation")
}

// NewRequestBuilder create packet builder for an outgoing request with sender set to current node.
func (w *Wrapper) NewRequestBuilder() network.RequestBuilder {
	panic("not used in DHT implementation")
}

func (w *Wrapper) BuildResponse(request network.Request, responseData interface{}) network.Response {
	panic("not used in DHT implementation")
}

// ResendPulseToKnownHosts resend pulse when we receive pulse from pulsar daemon
func (w *Wrapper) ResendPulseToKnownHosts(pulse core.Pulse) {
	p := &packet.RequestPulse{Pulse: pulse}
	activeNodes := w.HostNetwork.GetActiveNodesList()
	hosts := make([]host.Host, 0)
	for _, node := range activeNodes {
		address, err := host.NewAddress(node.PhysicalAddress())
		if err != nil {
			log.Error("error resolving address while resending pulse: " + node.PhysicalAddress())
			continue
		}
		id := id.FromBase58(resolver.ResolveHostID(node.ID()))
		hosts = append(hosts, host.Host{ID: id, Address: address})
	}
	ResendPulseToKnownHosts(w.HostNetwork, hosts, p)
}

// Inject inject components
func (w *Wrapper) Inject(components core.Components) {
	if components.NodeNetwork == nil {
		log.Error("active node component is nil")
	} else {
		nodeKeeper := components.NodeNetwork.(network.NodeKeeper)
		w.HostNetwork.SetNodeKeeper(nodeKeeper)
	}
	if components.NetworkCoordinator == nil {
		log.Error("network coordinator is nil")
	} else {
		w.HostNetwork.GetNetworkCommonFacade().SetNetworkCoordinator(components.NetworkCoordinator)
	}
}

func (w *Wrapper) GetConsensus() consensus.Processor {
	return w.HostNetwork.GetNetworkCommonFacade().GetConsensus()
}

func NewDhtHostNetwork(conf configuration.Configuration, certificate core.Certificate, pulseCallback network.OnPulse) (network.HostNetwork, error) {
	cascade1 := &cascade.Cascade{}
	dht, err := NewHostNetwork(conf, cascade1, certificate, pulseCallback)
	if err != nil {
		return nil, err
	}

	w := &Wrapper{HostNetwork: dht}
	f := func(data core.Cascade, method string, args [][]byte) error {
		return w.initCascadeSendMessage(data, true, method, args)
	}
	cascade1.SendMessage = f
	return w, nil
}

func NewDhtNetworkController(network network.HostNetwork) (network.Controller, error) {
	// hack for new interface
	w := network.(*Wrapper)
	return w, nil
}

func CreateDHTContext(handler hosthandler.HostHandler) hosthandler.Context {
	ctx, err := NewContextBuilder(handler).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}

func NewNetworkConsensus(network network.HostNetwork) consensus.Processor {
	handler := network.(*Wrapper).HostNetwork
	return consensus2.NewNetworkConsensus(handler)
}
