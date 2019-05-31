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

package bootstrap

import (
	"context"
	"crypto"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/version"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/controller/pinger"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

const bootstrapTimeout time.Duration = 2 // seconds
// const updateScheduleETA time.Duration = 60 // seconds

var (
	ErrReconnectRequired = errors.New("NetworkNode should connect via consensus bootstrap")
)

type DiscoveryNode struct {
	Host *host.Host
	Node insolar.DiscoveryNode
}

type Bootstrapper interface {
	component.Initer

	Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error)
	ZeroBootstrap(ctx context.Context) (*network.BootstrapResult, error)
}

type DiscoveryBootstrapper interface {
	BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error)
	SetLastPulse(number insolar.PulseNumber)
	GetLastPulse() insolar.PulseNumber
}

type Bootstrap struct {
	Certificate   insolar.Certificate         `inject:""`
	NodeKeeper    network.NodeKeeper          `inject:""`
	Network       network.HostNetwork         `inject:""`
	Gatewayer     network.Gatewayer           `inject:""`
	PulseAccessor pulse.Accessor              `inject:""`
	Cryptography  insolar.CryptographyService `inject:""`

	options *common.Options
	pinger  *pinger.Pinger

	lastPulse      insolar.PulseNumber
	lastPulseLock  sync.RWMutex
	pulsePersisted bool

	bootstrapLock       chan struct{}
	cyclicBootstrapStop int32

	genesisRequestsReceived map[insolar.Reference]*packet.GenesisRequest
	genesisLock             sync.Mutex

	firstPulseTime time.Time

	reconnectToNewNetwork func(ctx context.Context, address string)
}

func (bc *Bootstrap) GetFirstFakePulseTime() time.Time {
	return bc.firstPulseTime
}

func (bc *Bootstrap) getRequest(ref insolar.Reference) *packet.GenesisRequest {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	return bc.genesisRequestsReceived[ref]
}

func (bc *Bootstrap) setRequest(ref insolar.Reference, req *packet.GenesisRequest) {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	bc.genesisRequestsReceived[ref] = req
}

// Bootstrap on the discovery node (step 1 of the bootstrap process)
func (bc *Bootstrap) Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error) {
	log.Info("Bootstrapping to discovery node")
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.Bootstrap")
	defer span.End()
	discoveryNodes := bc.Certificate.GetDiscoveryNodes()
	if utils.OriginIsDiscovery(bc.Certificate) {
		discoveryNodes = RemoveOrigin(discoveryNodes, bc.NodeKeeper.GetOrigin().ID())
	}
	if len(discoveryNodes) == 0 {
		return nil, nil, errors.New("There are 0 discovery nodes to connect to")
	}
	ch := bc.getDiscoveryNodesChannel(ctx, discoveryNodes, 1)
	result := bc.waitResultFromChannel(ctx, ch)
	if result == nil {
		return nil, nil, errors.New("Failed to bootstrap to any of discovery nodes")
	}
	discovery := FindDiscovery(bc.Certificate, result.Host.NodeID)
	return result, &DiscoveryNode{result.Host, discovery}, nil
}

func (bc *Bootstrap) SetLastPulse(number insolar.PulseNumber) {
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

func (bc *Bootstrap) forceSetLastPulse(number insolar.PulseNumber) {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.forceSetLastPulse wait lastPulseLock")
	bc.lastPulseLock.Lock()
	span.End()
	defer bc.lastPulseLock.Unlock()

	log.Infof("Network will start from pulse %d + delta", number)
	bc.lastPulse = number
}

func (bc *Bootstrap) GetLastPulse() insolar.PulseNumber {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.GetLastPulse wait lastPulseLock")
	bc.lastPulseLock.RLock()
	span.End()
	defer bc.lastPulseLock.RUnlock()

	return bc.lastPulse
}

func (bc *Bootstrap) checkActiveNode(node insolar.NetworkNode) error {
	n := bc.NodeKeeper.GetAccessor().GetActiveNode(node.ID())
	if n != nil {
		return errors.Errorf("NetworkNode ID collision: %s", n.ID())
	}
	n = bc.NodeKeeper.GetAccessor().GetActiveNodeByShortID(node.ShortID())
	if n != nil {
		return errors.Errorf("Short ID collision: %d", n.ShortID())
	}
	if node.Version() != bc.NodeKeeper.GetOrigin().Version() {
		return errors.Errorf("NetworkNode %s version %s does not match origin version %s",
			node.ID(), node.Version(), bc.NodeKeeper.GetOrigin().Version())
	}
	return nil
}

func (bc *Bootstrap) ZeroBootstrap(ctx context.Context) (*network.BootstrapResult, error) {
	host, err := host.NewHostN(bc.NodeKeeper.GetOrigin().Address(), bc.NodeKeeper.GetOrigin().ID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a host")
	}
	inslogger.FromContext(ctx).Info("[ Bootstrap ] Zero bootstrap")
	bc.NodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{bc.NodeKeeper.GetOrigin()})
	return &network.BootstrapResult{
		Host: host,
		// FirstPulseTime: nb.Bootstrapper.GetFirstFakePulseTime(),
	}, nil
}

// BootstrapDiscovery bootstrapping as discovery node
func (bc *Bootstrap) BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ BootstrapDiscovery ] Network bootstrap between discovery nodes")
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.BootstrapDiscovery")
	defer span.End()
	discoveryNodes := RemoveOrigin(bc.Certificate.GetDiscoveryNodes(), *bc.Certificate.GetNodeRef())
	discoveryCount := len(discoveryNodes)
	if discoveryCount == 0 {
		return bc.ZeroBootstrap(ctx)
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
			logger.Infof("[ BootstrapDiscovery ] Connected to %d/%d discovery nodes", len(hosts), discoveryCount)
		}
	}
	reconnectRequests := 0
	for _, bootstrapResult := range bootstrapResults {
		if bootstrapResult.ReconnectRequired {
			reconnectRequests++
		}
	}
	minRequests := int(math.Floor(0.5*float64(discoveryCount))) + 1
	if reconnectRequests >= minRequests {
		logger.Infof("[ BootstrapDiscovery ] Need to reconnect as joiner (requested by %d/%d discovery nodes)",
			reconnectRequests, discoveryCount)
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
		activeNode.(node.MutableNode).SetState(insolar.NodeUndefined)
		activeNodesStr = append(activeNodesStr, activeNode.ID().String())
	}
	bc.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodeUndefined)
	activeNodes = append(activeNodes, bc.NodeKeeper.GetOrigin())
	bc.NodeKeeper.SetInitialSnapshot(activeNodes)
	logger.Infof("[ BootstrapDiscovery ] Added active nodes: %s", strings.Join(activeNodesStr, ", "))

	if bc.options.CyclicBootstrapEnabled {
		go bc.startCyclicBootstrap(ctx)
	}

	return parseBootstrapResults(bootstrapResults), nil
}

func (bc *Bootstrap) calculateLastIgnoredPulse(ctx context.Context, lastPulses []insolar.PulseNumber) insolar.PulseNumber {
	maxLastPulse := bc.GetLastPulse()
	inslogger.FromContext(ctx).Debugf("NetworkNode %s (origin) LastIgnoredPulse: %d", bc.NodeKeeper.GetOrigin().ID(), maxLastPulse)
	for _, pulse := range lastPulses {
		if pulse > maxLastPulse {
			maxLastPulse = pulse
		}
	}
	return maxLastPulse
}

func (bc *Bootstrap) sendGenesisRequest(ctx context.Context, h *host.Host) (*packet.GenesisResponse, error) {
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.sendGenesisRequest")
	defer span.End()
	discovery, err := bc.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to prepare genesis request to address %s", h)
	}
	request := &packet.GenesisRequest{
		LastPulse: bc.GetLastPulse(),
		Discovery: discovery,
	}
	future, err := bc.Network.SendRequestToHost(ctx, types.Genesis, request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send genesis request to address %s", h)
	}
	response, err := future.WaitResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to genesis request from address %s", h)
	}
	if response.GetResponse() == nil || response.GetResponse().GetGenesis() == nil {
		return nil, errors.Errorf("Failed to get response to genesis request from address %s: "+
			"got incorrect response: %s", h, response)
	}
	data := response.GetResponse().GetGenesis()
	if data.Response.Discovery == nil {
		return nil, errors.New("Error genesis response from discovery node: " + data.Error)
	}
	return data, nil
}

func (bc *Bootstrap) getDiscoveryNodesChannel(ctx context.Context, discoveryNodes []insolar.DiscoveryNode, needResponses int) <-chan *network.BootstrapResult {
	// we need only one host to bootstrap
	bootstrapResults := make(chan *network.BootstrapResult, needResponses)
	for _, discoveryNode := range discoveryNodes {
		go func(ctx context.Context, address string) {
			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
			ctx, span := instracer.StartSpan(ctx, "Bootstrapper.getDiscoveryNodesChannel")
			defer span.End()
			span.AddAttributes(
				trace.StringAttribute("Bootstrap node", address),
			)
			bootstrapResult, err := bootstrap(ctx, address, bc.options, bc.startBootstrap, nil)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
				return
			}
			bootstrapResults <- bootstrapResult
		}(ctx, discoveryNode.GetHost())
	}

	return bootstrapResults
}

func (bc *Bootstrap) getGenesisRequestsChannel(ctx context.Context, discoveryHosts []*host.Host) chan *packet.GenesisResponse {
	result := make(chan *packet.GenesisResponse)
	for _, discoveryHost := range discoveryHosts {
		go func(ctx context.Context, address *host.Host, ch chan<- *packet.GenesisResponse) {
			logger := inslogger.FromContext(ctx)
			ctx, span := instracer.StartSpan(ctx, "Bootstrapper.getGenesisRequestChannel")
			span.AddAttributes(
				trace.StringAttribute("genesis request to", address.String()),
			)
			defer span.End()
			cachedReq := bc.getRequest(address.NodeID)
			if cachedReq != nil {
				logger.Infof("Got genesis info of node %s from cache", address)
				ch <- &packet.GenesisResponse{Response: cachedReq}
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

func (bc *Bootstrap) waitResultFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult) *network.BootstrapResult {
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

func (bc *Bootstrap) waitResultsFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult, count int) ([]*network.BootstrapResult, []*host.Host) {
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

func (bc *Bootstrap) waitGenesisResults(ctx context.Context, ch <-chan *packet.GenesisResponse,
	count int) ([]insolar.NetworkNode, []insolar.PulseNumber, error) {

	result := make([]insolar.NetworkNode, 0)
	lastPulses := make([]insolar.PulseNumber, 0)
	for {
		select {
		case res := <-ch:
			discovery, err := node.ClaimToNode(version.Version, res.Response.Discovery)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Error deserializing node from discovery node")
			}
			result = append(result, discovery)
			lastPulses = append(lastPulses, res.Response.LastPulse)
			inslogger.FromContext(ctx).Debugf("NetworkNode %s LastIgnoredPulse: %d", discovery.ID(), res.Response.LastPulse)
			if len(result) == count {
				return result, lastPulses, nil
			}
		case <-time.After(bc.options.BootstrapTimeout):
			return nil, nil, errors.New(fmt.Sprintf("Genesis bootstrap timeout, successful genesis requests: %d/%d", len(result), count))
		}
	}
}

type bootstrapFunc func(context.Context, string, *packet.Permission) (*network.BootstrapResult, error)

func bootstrap(ctx context.Context, address string, options *common.Options, bootstrapF bootstrapFunc,
	perm *packet.Permission) (*network.BootstrapResult, error) {

	minTO := options.MinTimeout
	if !options.InfinityBootstrap {
		return bootstrapF(ctx, address, perm)
	}
	for {
		result, err := bootstrapF(ctx, address, perm)
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

func (bc *Bootstrap) startBootstrap(ctx context.Context, address string,
	perm *packet.Permission) (*network.BootstrapResult, error) {

	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.startBootstrap")
	defer span.End()
	bootstrapHost, err := bc.pinger.Ping(ctx, address, bc.options.PingTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to ping address %s", address)
	}
	claim, err := bc.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a join claim")
	}
	lastPulse, err := bc.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	request := &packet.BootstrapRequest{
		JoinClaim:     claim,
		LastNodePulse: lastPulse.PulseNumber,
		Permission:    perm,
	}
	future, err := bc.Network.SendRequestToHost(ctx, types.Bootstrap, request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send bootstrap request to address %s", address)
	}

	response, err := future.WaitResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to bootstrap request from address %s", address)
	}
	if response.GetResponse() == nil || response.GetResponse().GetBootstrap() == nil {
		return nil, errors.Errorf("Failed to get response to bootstrap request from address %s: "+
			"got incorrect response: %s", address, response)
	}
	data := response.GetResponse().GetBootstrap()
	logger := inslogger.FromContext(ctx)
	switch data.Code {
	case packet.Rejected:
		return nil, errors.New("Rejected: " + data.RejectReason)
	case packet.Redirected:
		// TODO: handle this case somehow
		if data.Permission == nil {
			return nil, errors.Errorf("discovery node %s returned empty permission for redirect", address)
		}
		logger.Infof("bootstrap redirected from %s to %s", address, data.Permission.Payload.ReconnectTo)
		return bootstrap(ctx, data.Permission.Payload.ReconnectTo, bc.options, bc.startBootstrap, data.Permission)
	case packet.UpdateSchedule:
		// TODO: INS-1960
	}

	return &network.BootstrapResult{
		Host:              response.GetSenderHost(),
		ReconnectRequired: data.Code == packet.ReconnectRequired,
		NetworkSize:       int(data.NetworkSize),
	}, nil
}

func (bc *Bootstrap) startCyclicBootstrap(ctx context.Context) {
	for atomic.LoadInt32(&bc.cyclicBootstrapStop) == 0 {
		results := make([]*network.BootstrapResult, 0)
		nodes := bc.getInactivenodes()
		for _, node := range nodes {
			res, err := bc.startBootstrap(ctx, node.GetHost(), nil)
			if err != nil {
				logger := inslogger.FromContext(ctx)
				logger.Errorf("[ StartCyclicBootstrap ] ", err)
				continue
			}
			results = append(results, res)
		}
		if len(results) != 0 {
			index := bc.getLargerNetworkIndex(results)
			if index >= 0 {
				bc.reconnectToNewNetwork(ctx, nodes[index].GetHost())
			}
		}
		time.Sleep(time.Second * bootstrapTimeout)
	}
}

func (bc *Bootstrap) getLargerNetworkIndex(results []*network.BootstrapResult) int {
	networkSize := results[0].NetworkSize
	index := 0
	for i := 1; i < len(results); i++ {
		if results[i].NetworkSize > networkSize {
			networkSize = results[i].NetworkSize
			index = i
		}
	}
	if networkSize > len(bc.NodeKeeper.GetAccessor().GetActiveNodes()) {
		return index
	}
	return -1
}

func (bc *Bootstrap) StopCyclicBootstrap() {
	atomic.StoreInt32(&bc.cyclicBootstrapStop, 1)
}

func (bc *Bootstrap) nodeShouldReconnectAsJoiner(nodeID insolar.Reference) bool {
	return bc.Gatewayer.Gateway().GetState() == insolar.CompleteNetworkState &&
		utils.IsDiscovery(nodeID, bc.Certificate)
}

func (bc *Bootstrap) processBootstrap(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetBootstrap() == nil {
		return nil, errors.Errorf("process bootstrap: got invalid protobuf request message: %s", request)
	}

	code := packet.Accepted
	data := request.GetRequest().GetBootstrap()
	var shortID insolar.ShortNodeID
	if CheckShortIDCollision(bc.NodeKeeper, data.JoinClaim.ShortNodeID) {
		shortID = GenerateShortID(bc.NodeKeeper, data.JoinClaim.GetNodeID())
	} else {
		shortID = data.JoinClaim.ShortNodeID
	}
	lastPulse, err := bc.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	var perm *packet.Permission
	if bc.nodeShouldReconnectAsJoiner(data.JoinClaim.NodeRef) { //nolint
		code = packet.ReconnectRequired
	} else if data.Permission == nil {
		code = packet.Redirected
		perm, err = bc.generatePermission(data.JoinClaim.NodePK[:])
		if err != nil {
			err = errors.Wrapf(err, "failed to generate permission")
			return bc.rejectBootstrapRequest(ctx, request, err.Error()), nil
		}
	} else {
		err = bc.checkPermission(data.Permission)
		if err != nil {
			err = errors.Wrapf(err, "failed to check permission")
			return bc.rejectBootstrapRequest(ctx, request, err.Error()), nil
		}
	}

	networkSize := uint32(len(bc.NodeKeeper.GetAccessor().GetActiveNodes()))
	return bc.Network.BuildResponse(ctx, request,
		&packet.BootstrapResponse{
			Code: code,
			// TODO: calculate ETA
			AssignShortID:    uint32(shortID),
			UpdateSincePulse: lastPulse.PulseNumber,
			NetworkSize:      networkSize,
			Permission:       perm,
		}), nil
}

func (bc *Bootstrap) rejectBootstrapRequest(ctx context.Context, request network.Packet, reason string) network.Packet {
	inslogger.FromContext(ctx).Errorf("Rejected bootstrap request from node %s: %s", request.GetSender(), reason)
	return bc.Network.BuildResponse(ctx, request, &packet.BootstrapResponse{Code: packet.Rejected, RejectReason: reason})
}

func (bc *Bootstrap) generatePermission(joinerPublicKey []byte) (*packet.Permission, error) {
	result := packet.Permission{
		Payload: packet.PermissionPayload{
			DiscoveryRef:    bc.NodeKeeper.GetOrigin().ID(),
			UTC:             time.Now().Unix(),
			ReconnectTo:     bc.getRandActiveDiscoveryAddress(),
			JoinerPublicKey: joinerPublicKey,
		}}

	sign, err := bc.signPermission(&result.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign permission")
	}

	result.Signature = sign
	return &result, nil
}

func (bc *Bootstrap) processGenesis(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetGenesis() == nil {
		return nil, errors.Errorf("process genesis: got invalid protobuf request message: %s", request)
	}

	data := request.GetRequest().GetGenesis()
	discovery, err := bc.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return bc.Network.BuildResponse(ctx, request, &packet.GenesisResponse{Error: err.Error()}), nil
	}
	bc.SetLastPulse(data.LastPulse)
	bc.setRequest(request.GetSender(), data)
	return bc.Network.BuildResponse(ctx, request, &packet.GenesisResponse{
		Response: &packet.GenesisRequest{Discovery: discovery, LastPulse: bc.GetLastPulse()},
	}), nil
}

func (bc *Bootstrap) Init(ctx context.Context) error {
	bc.firstPulseTime = time.Now()
	bc.pinger = pinger.NewPinger(bc.Network)
	bc.Network.RegisterRequestHandler(types.Bootstrap, bc.processBootstrap)
	bc.Network.RegisterRequestHandler(types.Genesis, bc.processGenesis)
	return nil
}

func parseBootstrapResults(results []*network.BootstrapResult) *network.BootstrapResult {
	minIDIndex := 0
	minID := results[0].Host.NodeID
	for i, result := range results {
		if minID.Compare(result.Host.NodeID) > 0 {
			minIDIndex = i
		}
	}
	return results[minIDIndex]
}

func (bc *Bootstrap) getInactivenodes() []insolar.DiscoveryNode {
	res := make([]insolar.DiscoveryNode, 0)
	for _, node := range bc.Certificate.GetDiscoveryNodes() {
		if bc.NodeKeeper.GetAccessor().GetActiveNode(*node.GetNodeRef()) != nil {
			res = append(res, node)
		}
	}
	return res
}

func (bc *Bootstrap) checkPermission(permission *packet.Permission) error {
	nodes := bc.Certificate.GetDiscoveryNodes()
	var discoveryPubKey crypto.PublicKey
	found := false
	for _, node := range nodes {
		if node.GetNodeRef().Equal(permission.Payload.DiscoveryRef) {
			discoveryPubKey = node.GetPublicKey()
			found = true
		}
	}
	if !found {
		return errors.New("failed to find a discovery node from reference in permission")
	}
	payload, err := permission.Payload.Marshal()
	if err != nil {
		return errors.New("failed to marshal bootstrap permission payload part")
	}
	verified := bc.Cryptography.Verify(discoveryPubKey, insolar.SignatureFromBytes(permission.Signature), payload)
	if !verified {
		return errors.New("bootstrap permission payload verification failed")
	}
	return nil
}

func (bc *Bootstrap) getRandActiveDiscoveryAddress() string {
	if len(bc.NodeKeeper.GetAccessor().GetActiveNodes()) <= 1 {
		return bc.NodeKeeper.GetOrigin().Address()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(bc.Certificate.GetDiscoveryNodes()))
	node := bc.NodeKeeper.GetAccessor().GetActiveNode(*bc.Certificate.GetDiscoveryNodes()[index].GetNodeRef())
	if (node != nil) && (node.GetState() == insolar.NodeReady) {
		return bc.Certificate.GetDiscoveryNodes()[index].GetHost()
	}

	return bc.NodeKeeper.GetOrigin().Address()
}

func (bc *Bootstrap) signPermission(perm *packet.PermissionPayload) ([]byte, error) {
	data, err := perm.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal bootstrap permission")
	}
	sign, err := bc.Cryptography.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign bootstrap permission")
	}
	return sign.Bytes(), nil
}

func NewBootstrapper(options *common.Options, reconnectToNewNetwork func(ctx context.Context, address string)) *Bootstrap {
	return &Bootstrap{
		options:                 options,
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[insolar.Reference]*packet.GenesisRequest),
		reconnectToNewNetwork:   reconnectToNewNetwork,
	}
}
