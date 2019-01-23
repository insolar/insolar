/*
 *    Copyright 2019 Insolar
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

package bootstrap

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/controller/pinger"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type Bootstrapper struct {
	options   *common.Options
	transport network.InternalTransport
	pinger    *pinger.Pinger
	cert      core.Certificate
	keeper    network.NodeKeeper

	lastPulse      core.PulseNumber
	lastPulseLock  sync.RWMutex
	pulsePersisted bool

	bootstrapLock chan struct{}

	genesisRequestsReceived map[core.RecordRef]*GenesisRequest
	genesisLock             sync.Mutex
}

func (bc *Bootstrapper) getRequest(ref core.RecordRef) *GenesisRequest {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	return bc.genesisRequestsReceived[ref]
}

func (bc *Bootstrapper) setRequest(ref core.RecordRef, req *GenesisRequest) {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	bc.genesisRequestsReceived[ref] = req
}

type NodeBootstrapRequest struct{}

type NodeBootstrapResponse struct {
	Code         Code
	RedirectHost string
	RejectReason string
}

type GenesisRequest struct {
	LastPulse core.PulseNumber
	Discovery *NodeStruct
}

type GenesisResponse struct {
	Response GenesisRequest
	Error    string
}

type StartSessionRequest struct{}

type StartSessionResponse struct {
	SessionID SessionID
}

type NodeStruct struct {
	ID      core.RecordRef
	SID     core.ShortNodeID
	Role    core.StaticRole
	PK      []byte
	Address string
	Version string
}

func newNode(n *NodeStruct) (core.Node, error) {
	pk, err := platformpolicy.NewKeyProcessor().ImportPublicKeyPEM(n.PK)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing node public key")
	}

	result := nodenetwork.NewNode(n.ID, n.Role, pk, n.Address, n.Version)
	mNode := result.(nodenetwork.MutableNode)
	mNode.SetShortID(n.SID)
	return mNode, nil
}

func newNodeStruct(node core.Node) (*NodeStruct, error) {
	pk, err := platformpolicy.NewKeyProcessor().ExportPublicKeyPEM(node.PublicKey())
	if err != nil {
		return nil, errors.Wrap(err, "error serializing node public key")
	}

	return &NodeStruct{
		ID:      node.ID(),
		SID:     node.ShortID(),
		Role:    node.Role(),
		PK:      pk,
		Address: node.PhysicalAddress(),
		Version: node.Version(),
	}, nil
}

type Code uint8

const (
	Accepted = Code(iota + 1)
	Rejected
	Redirected
)

func init() {
	gob.Register(&NodeBootstrapRequest{})
	gob.Register(&NodeBootstrapResponse{})
	gob.Register(&StartSessionRequest{})
	gob.Register(&StartSessionResponse{})
	gob.Register(&GenesisRequest{})
	gob.Register(&GenesisResponse{})
}

// Bootstrap on the discovery node (step 1 of the bootstrap process)
func (bc *Bootstrapper) Bootstrap(ctx context.Context) (*DiscoveryNode, error) {
	log.Info("Bootstrapping to discovery node")
	ch := bc.getDiscoveryNodesChannel(ctx, bc.cert.GetDiscoveryNodes(), 1)
	host := bc.waitResultFromChannel(ctx, ch)
	if host == nil {
		return nil, errors.New("Failed to bootstrap to any of discovery nodes")
	}
	discovery := FindDiscovery(bc.cert, host.NodeID)
	return &DiscoveryNode{Host: host, Node: discovery}, nil
}

func (bc *Bootstrapper) SetLastPulse(number core.PulseNumber) {
	bc.lastPulseLock.Lock()
	defer bc.lastPulseLock.Unlock()

	if !bc.pulsePersisted {
		bc.lastPulse = number
		close(bc.bootstrapLock)
		bc.pulsePersisted = true
	}
}

func (bc *Bootstrapper) forceSetLastPulse(number core.PulseNumber) {
	bc.lastPulseLock.Lock()
	defer bc.lastPulseLock.Unlock()

	log.Debugf("Network will start from pulse %d", number)
	bc.lastPulse = number
}

func (bc *Bootstrapper) GetLastPulse() core.PulseNumber {
	bc.lastPulseLock.RLock()
	defer bc.lastPulseLock.RUnlock()

	return bc.lastPulse
}

func (bc *Bootstrapper) checkActiveNode(node core.Node) error {
	n := bc.keeper.GetActiveNode(node.ID())
	if n != nil {
		return errors.New(fmt.Sprintf("Node ID collision: %s", n.ID()))
	}
	n = bc.keeper.GetActiveNodeByShortID(node.ShortID())
	if n != nil {
		return errors.New(fmt.Sprintf("Short ID collision: %d", n.ShortID()))
	}
	return nil
}

func (bc *Bootstrapper) BootstrapDiscovery(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("Network bootstrap between discovery nodes")
	discoveryNodes := bc.cert.GetDiscoveryNodes()
	var err error
	discoveryNodes, err = RemoveOrigin(discoveryNodes, *bc.cert.GetNodeRef())
	if err != nil {
		return errors.Wrapf(err, "Discovery bootstrap failed")
	}
	discoveryCount := len(discoveryNodes)
	if discoveryCount == 0 {
		return nil
	}

	var hosts []*host.Host
	for {
		ch := bc.getDiscoveryNodesChannel(ctx, discoveryNodes, discoveryCount)
		hosts = bc.waitResultsFromChannel(ctx, ch, discoveryCount)
		if len(hosts) == discoveryCount {
			// we connected to all discovery nodes
			break
		}
	}
	activeNodes := make([]core.Node, 0)
	activeNodesStr := make([]string, 0)

	<-bc.bootstrapLock
	logger.Debugf("After bootstrap lock")

	ch := bc.getGenesisRequestsChannel(ctx, hosts)
	activeNodes, lastPulses, err := bc.waitGenesisResults(ctx, ch, len(hosts))
	if err != nil {
		return err
	}
	bc.forceSetLastPulse(bc.calculateLastIgnoredPulse(ctx, lastPulses))
	for _, activeNode := range activeNodes {
		err = bc.checkActiveNode(activeNode)
		if err != nil {
			return errors.Wrapf(err, "Discovery check of node %s failed", activeNode.ID())
		}
		activeNodesStr = append(activeNodesStr, activeNode.ID().String())
	}
	bc.keeper.AddActiveNodes(activeNodes)
	logger.Infof("Added active nodes: %s", strings.Join(activeNodesStr, ", "))
	return nil
}

func (bc *Bootstrapper) calculateLastIgnoredPulse(ctx context.Context, lastPulses []core.PulseNumber) core.PulseNumber {
	maxLastPulse := bc.GetLastPulse()
	inslogger.FromContext(ctx).Debugf("Node %s (origin) LastIgnoredPulse: %d", bc.keeper.GetOrigin().ID(), maxLastPulse)
	for _, pulse := range lastPulses {
		if pulse > maxLastPulse {
			maxLastPulse = pulse
		}
	}
	return maxLastPulse
}

func (bc *Bootstrapper) sendGenesisRequest(ctx context.Context, h *host.Host) (*GenesisResponse, error) {
	discovery, err := newNodeStruct(bc.keeper.GetOrigin())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to prepare genesis request to address %s", h)
	}
	request := bc.transport.NewRequestBuilder().Type(types.Genesis).Data(&GenesisRequest{
		LastPulse: bc.GetLastPulse(),
		Discovery: discovery,
	}).Build()
	future, err := bc.transport.SendRequestPacket(ctx, request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send genesis request to address %s", h)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to genesis request from address %s", h)
	}
	data := response.GetData().(*GenesisResponse)
	if data.Response.Discovery == nil {
		return nil, errors.New("Error genesis response from discovery node: " + data.Error)
	}
	return data, nil
}

func (bc *Bootstrapper) getDiscoveryNodesChannel(ctx context.Context, discoveryNodes []core.DiscoveryNode, needResponses int) <-chan *host.Host {
	// we need only one host to bootstrap
	bootstrapHosts := make(chan *host.Host, needResponses)
	for _, discoveryNode := range discoveryNodes {
		go func(ctx context.Context, address string, ch chan<- *host.Host) {
			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
			bootstrapHost, err := bootstrap(ctx, address, bc.options, bc.startBootstrap)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapHosts <- bootstrapHost
		}(ctx, discoveryNode.GetHost(), bootstrapHosts)
	}
	return bootstrapHosts
}

func (bc *Bootstrapper) getGenesisRequestsChannel(ctx context.Context, discoveryHosts []*host.Host) chan *GenesisResponse {
	result := make(chan *GenesisResponse)
	for _, discoveryHost := range discoveryHosts {
		go func(ctx context.Context, address *host.Host, ch chan<- *GenesisResponse) {
			logger := inslogger.FromContext(ctx)
			cachedReq := bc.getRequest(address.NodeID)
			if cachedReq != nil {
				logger.Infof("Got genesis info of node %s from cache", address)
				ch <- &GenesisResponse{Response: *cachedReq}
				return
			}

			logger.Infof("Sending genesis bootstrap request to address %s", address)
			response, err := bc.sendGenesisRequest(ctx, address)
			if err != nil {
				logger.Warnf("Discovery bootstrap to host %s failed: %s", address, err)
				return
			}
			result <- response
		}(ctx, discoveryHost, result)
	}
	return result
}

func (bc *Bootstrapper) waitResultFromChannel(ctx context.Context, ch <-chan *host.Host) *host.Host {
	for {
		select {
		case bootstrapHost := <-ch:
			return bootstrapHost
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warn("Bootstrap timeout")
			return nil
		}
	}
}

func (bc *Bootstrapper) waitResultsFromChannel(ctx context.Context, ch <-chan *host.Host, count int) []*host.Host {
	result := make([]*host.Host, 0)
	for {
		select {
		case bootstrapHost := <-ch:
			result = append(result, bootstrapHost)
			if len(result) == count {
				return result
			}
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warnf("Bootstrap timeout, successful bootstraps: %d/%d", len(result), count)
			return result
		}
	}
}

func (bc *Bootstrapper) waitGenesisResults(ctx context.Context, ch <-chan *GenesisResponse, count int) ([]core.Node, []core.PulseNumber, error) {
	result := make([]core.Node, 0)
	lastPulses := make([]core.PulseNumber, 0)
	for {
		select {
		case res := <-ch:
			discovery, err := newNode(res.Response.Discovery)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Error deserializing node from discovery node")
			}
			result = append(result, discovery)
			lastPulses = append(lastPulses, res.Response.LastPulse)
			inslogger.FromContext(ctx).Debugf("Node %s LastIgnoredPulse: %d", discovery.ID(), res.Response.LastPulse)
			if len(result) == count {
				return result, lastPulses, nil
			}
		case <-time.After(bc.options.BootstrapTimeout):
			return nil, nil, errors.New(fmt.Sprintf("Genesis bootstrap timeout, successful genesis requests: %d/%d", len(result), count))
		}
	}
}

func bootstrap(ctx context.Context, address string, options *common.Options, bootstrapF func(context.Context, string) (*host.Host, error)) (*host.Host, error) {
	minTO := options.MinTimeout
	if !options.InfinityBootstrap {
		return bootstrapF(ctx, address)
	}
	for {
		result, err := bootstrapF(ctx, address)
		if err == nil {
			return result, nil
		}
		time.Sleep(minTO)
		minTO *= options.TimeoutMult
		if minTO > options.MaxTimeout {
			minTO = options.MaxTimeout
		}
	}
}

func (bc *Bootstrapper) startBootstrap(ctx context.Context, address string) (*host.Host, error) {
	bootstrapHost, err := bc.pinger.Ping(ctx, address, bc.options.PingTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	request := bc.transport.NewRequestBuilder().Type(types.Bootstrap).Data(&NodeBootstrapRequest{}).Build()
	future, err := bc.transport.SendRequestPacket(ctx, request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send bootstrap request to address %s", address)
	}
	response, err := future.GetResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to bootstrap request from address %s", address)
	}
	data := response.GetData().(*NodeBootstrapResponse)
	if data.Code == Rejected {
		return nil, errors.New("Rejected: " + data.RejectReason)
	}
	if data.Code == Redirected {
		return bootstrap(ctx, data.RedirectHost, bc.options, bc.startBootstrap)
	}
	return response.GetSenderHost(), nil
}

func (bc *Bootstrapper) processBootstrap(ctx context.Context, request network.Request) (network.Response, error) {
	// TODO: redirect logic
	return bc.transport.BuildResponse(ctx, request, &NodeBootstrapResponse{Code: Accepted}), nil
}

func (bc *Bootstrapper) processGenesis(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*GenesisRequest)
	discovery, err := newNodeStruct(bc.keeper.GetOrigin())
	if err != nil {
		return bc.transport.BuildResponse(ctx, request, &GenesisResponse{Error: err.Error()}), nil
	}
	bc.SetLastPulse(data.LastPulse)
	bc.setRequest(request.GetSender(), data)
	return bc.transport.BuildResponse(ctx, request, &GenesisResponse{
		Response: GenesisRequest{Discovery: discovery, LastPulse: bc.GetLastPulse()},
	}), nil
}

func (bc *Bootstrapper) Start(keeper network.NodeKeeper) {
	bc.keeper = keeper
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
	bc.transport.RegisterPacketHandler(types.Genesis, bc.processGenesis)
}

func NewBootstrapper(options *common.Options, certificate core.Certificate, transport network.InternalTransport) *Bootstrapper {
	return &Bootstrapper{
		options:                 options,
		cert:                    certificate,
		transport:               transport,
		pinger:                  pinger.NewPinger(transport),
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[core.RecordRef]*GenesisRequest),
	}
}
