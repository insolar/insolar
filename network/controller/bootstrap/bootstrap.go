/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package bootstrap

import (
	"context"
	"encoding/gob"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	coreutils "github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/controller/pinger"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

var (
	ErrReconnectRequired = errors.New("Node should connect via consensus bootstrap")
)

type DiscoveryNode struct {
	Host *host.Host
	Node core.DiscoveryNode
}

type Bootstrapper interface {
	component.Initer

	Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error)
	BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error)
	SetLastPulse(number core.PulseNumber)
	GetLastPulse() core.PulseNumber
	// GetFirstFakePulseTime() time.Time
}

type bootstrapper struct {
	Certificate     core.Certificate     `inject:""`
	NodeKeeper      network.NodeKeeper   `inject:""`
	NetworkSwitcher core.NetworkSwitcher `inject:""`
	Rules           network.Rules        `inject:""`

	options   *common.Options
	transport network.InternalTransport
	pinger    *pinger.Pinger

	lastPulse      core.PulseNumber
	lastPulseLock  sync.RWMutex
	pulsePersisted bool

	bootstrapLock chan struct{}

	genesisRequestsReceived map[core.RecordRef]*GenesisRequest
	genesisLock             sync.Mutex

	firstPulseTime time.Time
}

func (bc *bootstrapper) GetFirstFakePulseTime() time.Time {
	return bc.firstPulseTime
}

func (bc *bootstrapper) getRequest(ref core.RecordRef) *GenesisRequest {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	return bc.genesisRequestsReceived[ref]
}

func (bc *bootstrapper) setRequest(ref core.RecordRef, req *GenesisRequest) {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	bc.genesisRequestsReceived[ref] = req
}

type NodeBootstrapRequest struct{}

type NodeBootstrapResponse struct {
	Code           Code
	RedirectHost   string
	RejectReason   string
	DiscoveryCount int
	// FirstPulseTimeUnix int64
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
	pk, err := platformpolicy.NewKeyProcessor().ImportPublicKeyBinary(n.PK)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing node public key")
	}

	result := nodenetwork.NewNode(n.ID, n.Role, pk, n.Address, n.Version)
	mNode := result.(nodenetwork.MutableNode)
	mNode.SetShortID(n.SID)
	return mNode, nil
}

func newNodeStruct(node core.Node) (*NodeStruct, error) {
	pk, err := platformpolicy.NewKeyProcessor().ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		return nil, errors.Wrap(err, "error serializing node public key")
	}

	return &NodeStruct{
		ID:      node.ID(),
		SID:     node.ShortID(),
		Role:    node.Role(),
		PK:      pk,
		Address: node.Address(),
		Version: node.Version(),
	}, nil
}

type Code uint8

const (
	Accepted = Code(iota + 1)
	Rejected
	Redirected
	ReconnectRequired
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
func (bc *bootstrapper) Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Bootstrapping to discovery node")
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.Bootstrap")
	defer span.End()

	discoveryCount := len(bc.Certificate.GetDiscoveryNodes())
	ch := bc.getDiscoveryNodesChannel(ctx, bc.Certificate.GetDiscoveryNodes(), discoveryCount)

	bootstrapResults, hosts := bc.waitResultsFromChannel(ctx, ch, discoveryCount)
	logger.Infof("[ Bootstrap ] Connected to %d/%d discovery nodes", len(hosts), discoveryCount)

	if len(bootstrapResults) == 0 {
		return nil, nil, errors.New("[ Bootstrap ] Failed to get bootstrap results")
	}

	majorityRule := bc.Certificate.GetMajorityRule()
	b, isMajority := getDiscoveryFromBootstrapResults(bootstrapResults, majorityRule)
	if utils.OriginIsDiscovery(bc.Certificate) || isMajority {
		return b, &DiscoveryNode{b.Host, findDiscovery(bc.Certificate, b.Host.NodeID)}, nil
	}

	return nil, nil, errors.New("majority rule failed")
}

func getDiscoveryFromBootstrapResults(bootstrapResults []*network.BootstrapResult, majorityRule int) (*network.BootstrapResult, bool) {
	sort.Slice(bootstrapResults, func(i, j int) bool {
		return bootstrapResults[i].DiscoveryCount > bootstrapResults[j].DiscoveryCount
	})

	maxBootstrapResults := make([]*network.BootstrapResult, 0, len(bootstrapResults))
	for i, result := range bootstrapResults {
		if i == 0 || maxBootstrapResults[0].DiscoveryCount == result.DiscoveryCount {
			maxBootstrapResults = append(maxBootstrapResults, result)
		} else {
			break
		}
	}

	i := coreutils.RandomInt(len(maxBootstrapResults))
	randomMaxResult := bootstrapResults[i]
	return randomMaxResult, randomMaxResult.DiscoveryCount >= majorityRule
}

func (bc *bootstrapper) SetLastPulse(number core.PulseNumber) {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.SetLastPulse wait lastPulseLock")
	bc.lastPulseLock.Lock()
	span.End()
	defer bc.lastPulseLock.Unlock()

	if !bc.pulsePersisted {
		bc.lastPulse = number
		close(bc.bootstrapLock)
		bc.pulsePersisted = true
	}
}

func (bc *bootstrapper) forceSetLastPulse(number core.PulseNumber) {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.forceSetLastPulse wait lastPulseLock")
	bc.lastPulseLock.Lock()
	span.End()
	defer bc.lastPulseLock.Unlock()

	log.Infof("Network will start from pulse %d + delta", number)
	bc.lastPulse = number
}

func (bc *bootstrapper) GetLastPulse() core.PulseNumber {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.GetLastPulse wait lastPulseLock")
	bc.lastPulseLock.RLock()
	span.End()
	defer bc.lastPulseLock.RUnlock()

	return bc.lastPulse
}

func (bc *bootstrapper) checkActiveNode(node core.Node) error {
	n := bc.NodeKeeper.GetActiveNode(node.ID())
	if n != nil {
		return errors.New(fmt.Sprintf("Node ID collision: %s", n.ID()))
	}
	n = bc.NodeKeeper.GetActiveNodeByShortID(node.ShortID())
	if n != nil {
		return errors.New(fmt.Sprintf("Short ID collision: %d", n.ShortID()))
	}
	return nil
}

func (bc *bootstrapper) BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ BootstrapDiscovery ] Network bootstrap between discovery nodes")
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.BootstrapDiscovery")
	defer span.End()

	discoveryNodes := bc.Certificate.GetDiscoveryNodes()
	var err error
	discoveryNodes, err = removeOrigin(discoveryNodes, *bc.Certificate.GetNodeRef())
	if err != nil {
		return nil, errors.Wrapf(err, "Discovery bootstrap failed")
	}
	discoveryCount := len(discoveryNodes)
	if discoveryCount == 0 {
		host, err := host.NewHostN(bc.NodeKeeper.GetOrigin().Address(), bc.NodeKeeper.GetOrigin().ID())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a host")
		}
		return &network.BootstrapResult{
			Host: host,
			// FirstPulseTime: bc.firstPulseTime,
		}, nil
	}

	var bootstrapResults []*network.BootstrapResult
	var hosts []*host.Host
	for {
		ch := bc.getDiscoveryNodesChannel(ctx, discoveryNodes, discoveryCount)
		bootstrapResults, hosts = bc.waitResultsFromChannel(ctx, ch, discoveryCount)
		if len(hosts) == discoveryCount {
			// we connected to all discovery nodes
			break
		} else {
			logger.WithFields(map[string]interface{}{
				"connected":      len(hosts),
				"discoveryCount": discoveryCount,
			}).Info("[ BootstrapDiscovery ] Connected to discovery nodes")

			reconnectRequests := getReconnectCount(bootstrapResults)
			// Few discovery nodes are down and network was in complete state
			if reconnectRequests == len(hosts) {
				logger.WithFields(map[string]interface{}{
					"connected":         len(hosts),
					"reconnectRequests": reconnectRequests,
				}).Info("[ BootstrapDiscovery ] Need to reconnect as joiner (all connected discoveries require reconnect)")
				return nil, ErrReconnectRequired
			}
		}
	}

	reconnectRequests := getReconnectCount(bootstrapResults)
	minRequests := int(math.Floor(0.5*float64(discoveryCount))) + 1

	if reconnectRequests >= minRequests {
		logger.WithFields(map[string]interface{}{
			"reconnectRequests": reconnectRequests,
			"discoveryCount":    discoveryCount,
		}).Info("[ BootstrapDiscovery ] Need to reconnect as joiner (requested by discovery nodes)")
		return nil, ErrReconnectRequired
	}

	activeNodesStr := make([]string, 0)

	<-bc.bootstrapLock
	logger.Debugf("[ BootstrapDiscovery ] After bootstrap lock")

	ch := bc.getGenesisRequestsChannel(ctx, hosts)
	activeNodes, lastPulses, err := bc.waitGenesisResults(ctx, ch, len(hosts))
	if err != nil {
		return nil, err
	}
	bc.forceSetLastPulse(bc.calculateLastIgnoredPulse(ctx, lastPulses))
	for _, activeNode := range activeNodes {
		err = bc.checkActiveNode(activeNode)
		if err != nil {
			return nil, errors.Wrapf(err, "Discovery check of node %s failed", activeNode.ID())
		}
		activeNode.(nodenetwork.MutableNode).SetState(core.NodeDiscovery)
		activeNodesStr = append(activeNodesStr, activeNode.ID().String())
	}

	bc.NodeKeeper.AddActiveNodes(activeNodes)
	bc.NodeKeeper.GetOrigin().(nodenetwork.MutableNode).SetState(core.NodeDiscovery)
	logger.Infof("[ BootstrapDiscovery ] Added active nodes: %s", strings.Join(activeNodesStr, ", "))
	return parseBotstrapResults(bootstrapResults), nil
}

func (bc *bootstrapper) calculateLastIgnoredPulse(ctx context.Context, lastPulses []core.PulseNumber) core.PulseNumber {
	maxLastPulse := bc.GetLastPulse()
	inslogger.FromContext(ctx).Debugf("Node %s (origin) LastIgnoredPulse: %d", bc.NodeKeeper.GetOrigin().ID(), maxLastPulse)
	for _, pulse := range lastPulses {
		if pulse > maxLastPulse {
			maxLastPulse = pulse
		}
	}
	return maxLastPulse
}

func (bc *bootstrapper) sendGenesisRequest(ctx context.Context, h *host.Host) (*GenesisResponse, error) {
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.sendGenesisRequest")
	defer span.End()
	discovery, err := newNodeStruct(bc.NodeKeeper.GetOrigin())
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

func (bc *bootstrapper) getDiscoveryNodesChannel(ctx context.Context, discoveryNodes []core.DiscoveryNode, needResponses int) <-chan *network.BootstrapResult {
	// we need only one host to bootstrap
	bootstrapResults := make(chan *network.BootstrapResult, needResponses)
	for _, discoveryNode := range discoveryNodes {
		go func(ctx context.Context, address string, ch chan<- *network.BootstrapResult) {
			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
			ctx, span := instracer.StartSpan(ctx, "Bootstrapper.getDiscoveryNodesChannel")
			defer span.End()
			span.AddAttributes(
				trace.StringAttribute("Bootstrap node", address),
			)
			bootstrapResult, err := bootstrap(ctx, address, bc.options, bc.startBootstrap)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapResults <- bootstrapResult
		}(ctx, discoveryNode.GetHost(), bootstrapResults)
	}

	return bootstrapResults
}

func (bc *bootstrapper) getGenesisRequestsChannel(ctx context.Context, discoveryHosts []*host.Host) chan *GenesisResponse {
	result := make(chan *GenesisResponse)
	for _, discoveryHost := range discoveryHosts {
		go func(ctx context.Context, address *host.Host, ch chan<- *GenesisResponse) {
			logger := inslogger.FromContext(ctx)
			ctx, span := instracer.StartSpan(ctx, "Bootsytrapper.getGenesisRequestChannel")
			span.AddAttributes(
				trace.StringAttribute("genesis request to", address.String()),
			)
			defer span.End()
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

func (bc *bootstrapper) waitResultsFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult, count int) ([]*network.BootstrapResult, []*host.Host) {
	result := make([]*network.BootstrapResult, 0)
	hosts := make([]*host.Host, 0)
	for {
		select {
		case bootstrapResult := <-ch:
			result = append(result, bootstrapResult)
			hosts = append(hosts, bootstrapResult.Host)
			if len(result) == count {
				return result, hosts
			}
		case <-time.After(bc.options.BootstrapTimeout):
			inslogger.FromContext(ctx).Warnf("Bootstrap timeout, successful bootstraps: %d/%d", len(result), count)
			return result, hosts
		}
	}
}

func (bc *bootstrapper) waitGenesisResults(ctx context.Context, ch <-chan *GenesisResponse, count int) ([]core.Node, []core.PulseNumber, error) {
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

func bootstrap(ctx context.Context, address string, options *common.Options, bootstrapF func(context.Context, string) (*network.BootstrapResult, error)) (*network.BootstrapResult, error) {
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

func (bc *bootstrapper) startBootstrap(ctx context.Context, address string) (*network.BootstrapResult, error) {
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.startBootstrap")
	defer span.End()
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

	switch data.Code {
	case Rejected:
		return nil, errors.New("Rejected: " + data.RejectReason)
	case Redirected:
		return bootstrap(ctx, data.RedirectHost, bc.options, bc.startBootstrap)
	}
	return &network.BootstrapResult{
		// FirstPulseTime:    time.Unix(data.FirstPulseTimeUnix, 0),
		Host:              response.GetSenderHost(),
		ReconnectRequired: data.Code == ReconnectRequired,
		DiscoveryCount:    data.DiscoveryCount,
	}, nil
}

func (bc *bootstrapper) processBootstrap(ctx context.Context, request network.Request) (network.Response, error) {
	// TODO: redirect logic
	var code Code
	if bc.NetworkSwitcher.WasInCompleteState() {
		code = ReconnectRequired
	} else {
		code = Accepted
	}

	_, activeDiscoveryNodesLen := bc.Rules.CheckMajorityRule()

	return bc.transport.BuildResponse(ctx, request,
		&NodeBootstrapResponse{
			Code:           code,
			DiscoveryCount: activeDiscoveryNodesLen,
			// FirstPulseTimeUnix: bc.firstPulseTime.Unix(),
		}), nil
}

func (bc *bootstrapper) processGenesis(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*GenesisRequest)
	discovery, err := newNodeStruct(bc.NodeKeeper.GetOrigin())
	if err != nil {
		return bc.transport.BuildResponse(ctx, request, &GenesisResponse{Error: err.Error()}), nil
	}
	bc.SetLastPulse(data.LastPulse)
	bc.setRequest(request.GetSender(), data)
	return bc.transport.BuildResponse(ctx, request, &GenesisResponse{
		Response: GenesisRequest{Discovery: discovery, LastPulse: bc.GetLastPulse()},
	}), nil
}

func (bc *bootstrapper) Init(ctx context.Context) error {
	bc.firstPulseTime = time.Now()
	bc.transport.RegisterPacketHandler(types.Bootstrap, bc.processBootstrap)
	bc.transport.RegisterPacketHandler(types.Genesis, bc.processGenesis)
	return nil
}

func parseBotstrapResults(results []*network.BootstrapResult) *network.BootstrapResult {
	minIDIndex := 0
	minID := results[0].Host.NodeID
	for i, result := range results {
		if minID.Compare(result.Host.NodeID) > 0 {
			minIDIndex = i
		}
	}
	return results[minIDIndex]
}

func getReconnectCount(results []*network.BootstrapResult) int {
	reconnectRequests := 0
	for _, bootstrapResult := range results {
		if bootstrapResult.ReconnectRequired {
			reconnectRequests++
		}
	}

	return reconnectRequests
}

func NewBootstrapper(options *common.Options, transport network.InternalTransport) Bootstrapper {
	return &bootstrapper{
		options:       options,
		transport:     transport,
		pinger:        pinger.NewPinger(transport),
		bootstrapLock: make(chan struct{}),

		genesisRequestsReceived: make(map[core.RecordRef]*GenesisRequest),
	}
}
