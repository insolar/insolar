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
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"

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
	"github.com/insolar/insolar/platformpolicy"
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

type Permission struct {
	JoinerPublicKey []byte
	Signature       []byte
	UTC             []byte
	ReconnectTo     string
	DiscoveryRef    insolar.Reference
}

type Bootstrapper interface {
	component.Initer

	Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error)
	BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error)
	ZeroBootstrap(ctx context.Context) (*network.BootstrapResult, error)
	SetLastPulse(number insolar.PulseNumber)
	GetLastPulse() insolar.PulseNumber
}

type bootstrapper struct {
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

	genesisRequestsReceived map[insolar.Reference]*GenesisRequest
	genesisLock             sync.Mutex

	firstPulseTime time.Time

	reconnectToNewNetwork func(ctx context.Context, address string)
}

func (p *Permission) RawBytes() []byte {
	res := make([]byte, 0)
	res = append(res, p.JoinerPublicKey...)
	res = append(res, p.DiscoveryRef.Bytes()...)
	res = append(res, []byte(p.ReconnectTo)...)
	res = append(res, p.UTC...)
	return res
}

func (bc *bootstrapper) GetFirstFakePulseTime() time.Time {
	return bc.firstPulseTime
}

func (bc *bootstrapper) getRequest(ref insolar.Reference) *GenesisRequest {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	return bc.genesisRequestsReceived[ref]
}

func (bc *bootstrapper) setRequest(ref insolar.Reference, req *GenesisRequest) {
	bc.genesisLock.Lock()
	defer bc.genesisLock.Unlock()

	bc.genesisRequestsReceived[ref] = req
}

type NodeBootstrapRequest struct {
	// TODO: change to mandate cuz cert not registered for gob
	// Certificate   insolar.Certificate
	JoinClaim packets.NodeJoinClaim
	// LastNodePulse is a last received pulse number.
	LastNodePulse insolar.PulseNumber
	// Permission is a information for reconnect to another discovery node.
	Permission Permission
}

type NodeBootstrapResponse struct {
	Code         Code
	RejectReason string
	// ETA - promise to accept joiner node to the network (in seconds).
	ETA int
	// AssignShortID is an demand to use this short id.
	AssignShortID insolar.ShortNodeID
	// UpdateSincePulse is a pulse number from which origin have to update storage.
	UpdateSincePulse insolar.PulseNumber
	// NetworkSize is a size of the network from bootstrap node.
	NetworkSize int
	// Permission is a information for reconnect to another discovery node.
	Permission Permission
}

type GenesisRequest struct {
	LastPulse insolar.PulseNumber
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
	ID      insolar.Reference
	SID     insolar.ShortNodeID
	Role    insolar.StaticRole
	PK      []byte
	Address string
	Version string
}

func newNode(n *NodeStruct) (insolar.NetworkNode, error) {
	pk, err := platformpolicy.NewKeyProcessor().ImportPublicKeyBinary(n.PK)
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing node public key")
	}

	result := node.NewNode(n.ID, n.Role, pk, n.Address, n.Version)
	mNode := result.(node.MutableNode)
	mNode.SetShortID(n.SID)
	return mNode, nil
}

func newNodeStruct(node insolar.NetworkNode) (*NodeStruct, error) {
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
	UpdateSchedule
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

func (bc *bootstrapper) SetLastPulse(number insolar.PulseNumber) {
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

func (bc *bootstrapper) forceSetLastPulse(number insolar.PulseNumber) {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.forceSetLastPulse wait lastPulseLock")
	bc.lastPulseLock.Lock()
	span.End()
	defer bc.lastPulseLock.Unlock()

	log.Infof("Network will start from pulse %d + delta", number)
	bc.lastPulse = number
}

func (bc *bootstrapper) GetLastPulse() insolar.PulseNumber {
	_, span := instracer.StartSpan(context.Background(), "Bootstrapper.GetLastPulse wait lastPulseLock")
	bc.lastPulseLock.RLock()
	span.End()
	defer bc.lastPulseLock.RUnlock()

	return bc.lastPulse
}

func (bc *bootstrapper) checkActiveNode(node insolar.NetworkNode) error {
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

func (bc *bootstrapper) ZeroBootstrap(ctx context.Context) (*network.BootstrapResult, error) {
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

func (bc *bootstrapper) BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
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

	return parseBotstrapResults(bootstrapResults), nil
}

func (bc *bootstrapper) calculateLastIgnoredPulse(ctx context.Context, lastPulses []insolar.PulseNumber) insolar.PulseNumber {
	maxLastPulse := bc.GetLastPulse()
	inslogger.FromContext(ctx).Debugf("NetworkNode %s (origin) LastIgnoredPulse: %d", bc.NodeKeeper.GetOrigin().ID(), maxLastPulse)
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
	request := bc.Network.NewRequestBuilder().Type(types.Genesis).Data(&GenesisRequest{
		LastPulse: bc.GetLastPulse(),
		Discovery: discovery,
	}).Build()
	future, err := bc.Network.SendRequestToHost(ctx, request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send genesis request to address %s", h)
	}
	response, err := future.WaitResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to genesis request from address %s", h)
	}
	data := response.GetData().(*GenesisResponse)
	if data.Response.Discovery == nil {
		return nil, errors.New("Error genesis response from discovery node: " + data.Error)
	}
	return data, nil
}

func (bc *bootstrapper) getDiscoveryNodesChannel(ctx context.Context, discoveryNodes []insolar.DiscoveryNode, needResponses int) <-chan *network.BootstrapResult {
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
			bootstrapResult, err := bootstrap(ctx, address, bc.options, bc.startBootstrap, nil)
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

func (bc *bootstrapper) waitResultFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult) *network.BootstrapResult {
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

func (bc *bootstrapper) waitGenesisResults(ctx context.Context, ch <-chan *GenesisResponse, count int) ([]insolar.NetworkNode, []insolar.PulseNumber, error) {
	result := make([]insolar.NetworkNode, 0)
	lastPulses := make([]insolar.PulseNumber, 0)
	for {
		select {
		case res := <-ch:
			discovery, err := newNode(res.Response.Discovery)
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

func bootstrap(ctx context.Context, address string, options *common.Options, bootstrapF func(context.Context, string, *Permission) (*network.BootstrapResult, error), perm *Permission) (*network.BootstrapResult, error) {
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

func (bc *bootstrapper) startBootstrap(ctx context.Context, address string, perm *Permission) (*network.BootstrapResult, error) {
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

	bootstrapReq := &NodeBootstrapRequest{
		JoinClaim:     *claim,
		LastNodePulse: lastPulse.PulseNumber,
	}

	if perm == nil {
		proc := platformpolicy.NewKeyProcessor()
		key, err := proc.ExportPublicKeyBinary(bc.Certificate.GetPublicKey())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to export an origin pub key")
		}
		bootstrapReq.Permission.JoinerPublicKey = key
	} else {
		bootstrapReq.Permission = *perm
	}

	request := bc.Network.NewRequestBuilder().Type(types.Bootstrap).Data(bootstrapReq).Build()
	future, err := bc.Network.SendRequestToHost(ctx, request, bootstrapHost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send bootstrap request to address %s", address)
	}

	response, err := future.WaitResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to bootstrap request from address %s", address)
	}

	data := response.GetData().(*NodeBootstrapResponse)
	logger := inslogger.FromContext(ctx)

	switch data.Code {
	case Rejected:
		return nil, errors.New("Rejected: " + data.RejectReason)
	case Redirected:
		logger.Infof("bootstrap redirected from %s to %s", bc.NodeKeeper.GetOrigin().Address(), data.Permission.ReconnectTo)
		return bootstrap(ctx, data.Permission.ReconnectTo, bc.options, bc.startBootstrap, &data.Permission)
	case UpdateSchedule:
		// TODO: INS-1960
	}

	return &network.BootstrapResult{
		Host:              response.GetSenderHost(),
		ReconnectRequired: data.Code == ReconnectRequired,
		NetworkSize:       data.NetworkSize,
	}, nil
}

func (bc *bootstrapper) startCyclicBootstrap(ctx context.Context) {
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
			index := bc.getLagerNetorkIndex(ctx, results)
			if index >= 0 {
				bc.reconnectToNewNetwork(ctx, nodes[index].GetHost())
			}
		}
		time.Sleep(time.Second * bootstrapTimeout)
	}
}

func (bc *bootstrapper) getLagerNetorkIndex(ctx context.Context, results []*network.BootstrapResult) int {
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

func (bc *bootstrapper) StopCyclicBootstrap() {
	atomic.StoreInt32(&bc.cyclicBootstrapStop, 1)
}

func (bc *bootstrapper) processBootstrap(ctx context.Context, request network.Request) (network.Response, error) {
	var code Code
	bootstrapRequest := request.GetData().(*NodeBootstrapRequest)
	if bootstrapRequest == nil {
		return nil, errors.New("received broken bootstrap request")
	}
	var shortID insolar.ShortNodeID
	if CheckShortIDCollision(bc.NodeKeeper, bootstrapRequest.JoinClaim.ShortNodeID) {
		shortID = GenerateShortID(bc.NodeKeeper, bootstrapRequest.JoinClaim.GetNodeID())
	} else {
		shortID = bootstrapRequest.JoinClaim.ShortNodeID
	}
	lastPulse, err := bc.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	if permissionIsEmpty(bootstrapRequest.Permission) {
		code = Redirected
		err := bc.updatePermissionsOnRequest(bootstrapRequest)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update a permission in request")
		}
	} else {
		code, err = bc.getCodeFromPermission(bootstrapRequest.Permission)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get a code from permission")
		}
	}

	if bc.Gatewayer.Gateway().GetState() == insolar.CompleteNetworkState {
		code = ReconnectRequired
	}

	return bc.Network.BuildResponse(ctx, request,
		&NodeBootstrapResponse{
			Code:         code,
			RejectReason: "",
			// TODO: calculate an ETA
			AssignShortID:    shortID,
			UpdateSincePulse: lastPulse.PulseNumber,
			NetworkSize:      len(bc.NodeKeeper.GetAccessor().GetActiveNodes()),
			Permission:       bootstrapRequest.Permission,
		}), nil
}

func (bc *bootstrapper) getCodeFromPermission(permission Permission) (Code, error) {
	verified, err := bc.checkPermissionSign(permission)
	if err != nil {
		return Rejected, errors.Wrap(err, "failed to check a permission sign")
	}
	if !verified {
		return Rejected, errors.New("failed to verify a permission sign")
	}

	// TODO: INS-1960
	// etaDiff := time.Since(permission.UTC)
	// if etaDiff > updateScheduleETA {
	// 	return UpdateSchedule, nil
	// }

	return Accepted, nil
}

func (bc *bootstrapper) updatePermissionsOnRequest(request *NodeBootstrapRequest) error {
	request.Permission.DiscoveryRef = bc.NodeKeeper.GetOrigin().ID()
	t, err := time.Now().GobEncode()
	if err != nil {
		return errors.Wrap(err, "failed to encode a time")
	}
	request.Permission.UTC = t
	request.Permission.ReconnectTo = bc.getRandActiveDiscoveryAddress()

	sign, err := bc.getPermissionSign(request.Permission)
	if err != nil {
		return errors.Wrap(err, "failed to get a permission sign")
	}

	request.Permission.Signature = sign
	return nil
}

func (bc *bootstrapper) processGenesis(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*GenesisRequest)
	discovery, err := newNodeStruct(bc.NodeKeeper.GetOrigin())
	if err != nil {
		return bc.Network.BuildResponse(ctx, request, &GenesisResponse{Error: err.Error()}), nil
	}
	bc.SetLastPulse(data.LastPulse)
	bc.setRequest(request.GetSender(), data)
	return bc.Network.BuildResponse(ctx, request, &GenesisResponse{
		Response: GenesisRequest{Discovery: discovery, LastPulse: bc.GetLastPulse()},
	}), nil
}

func (bc *bootstrapper) Init(ctx context.Context) error {
	bc.firstPulseTime = time.Now()
	bc.pinger = pinger.NewPinger(bc.Network)
	bc.Network.RegisterRequestHandler(types.Bootstrap, bc.processBootstrap)
	bc.Network.RegisterRequestHandler(types.Genesis, bc.processGenesis)
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

func (bc *bootstrapper) getInactivenodes() []insolar.DiscoveryNode {
	res := make([]insolar.DiscoveryNode, 0)
	for _, node := range bc.Certificate.GetDiscoveryNodes() {
		if bc.NodeKeeper.GetAccessor().GetActiveNode(*node.GetNodeRef()) != nil {
			res = append(res, node)
		}
	}
	return res
}

func (bc *bootstrapper) checkPermissionSign(permission Permission) (bool, error) {
	nodes := bc.Certificate.GetDiscoveryNodes()
	var discoveryPubKey crypto.PublicKey
	found := false
	for _, node := range nodes {
		if node.GetNodeRef().Equal(permission.DiscoveryRef) {
			discoveryPubKey = node.GetPublicKey()
			found = true
		}
	}
	if !found {
		return false, errors.New("failed to find a discovery node from reference in permission")
	}
	verified := bc.Cryptography.Verify(discoveryPubKey, insolar.SignatureFromBytes(permission.Signature), permission.RawBytes())
	return verified, nil
}

func (bc *bootstrapper) getRandActiveDiscoveryAddress() string {
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

func (bc *bootstrapper) getPermissionSign(perm Permission) ([]byte, error) {
	sign, err := bc.Cryptography.Sign(perm.RawBytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a permission")
	}
	return sign.Bytes(), nil
}

func permissionIsEmpty(perm Permission) bool {
	empty := false
	if len(perm.ReconnectTo) == 0 {
		empty = true
	}
	if perm.DiscoveryRef.IsEmpty() {
		empty = true
	}
	if len(perm.Signature) == 0 {
		empty = true
	}
	return empty
}

func NewBootstrapper(options *common.Options, reconnectToNewNetwork func(ctx context.Context, address string)) Bootstrapper {
	return &bootstrapper{
		options:                 options,
		bootstrapLock:           make(chan struct{}),
		genesisRequestsReceived: make(map[insolar.Reference]*GenesisRequest),
		reconnectToNewNetwork:   reconnectToNewNetwork,
	}
}
