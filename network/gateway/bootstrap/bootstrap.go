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
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/version"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

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

// TODO: remove
type DiscoveryNode struct {
	Host *host.Host
	Node insolar.DiscoveryNode
}

// type Bootstrapper interface {
// 	component.Initer
//
// 	Bootstrap(ctx context.Context) (*network.BootstrapResult, *DiscoveryNode, error)
// }

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

	firstPulseTime time.Time

	// reconnectToNewNetwork func(ctx context.Context, address string)
}

// func (bc *Bootstrap) GetFirstFakePulseTime() time.Time {
// 	return bc.firstPulseTime
// }

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

// BootstrapDiscovery bootstrapping as discovery node
func (bc *Bootstrap) BootstrapDiscovery(ctx context.Context) (*network.BootstrapResult, error) {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ BootstrapDiscovery ] Network bootstrap between discovery nodes")
	ctx, span := instracer.StartSpan(ctx, "Bootstrapper.BootstrapDiscovery")
	defer span.End()
	discoveryNodes := network.ExcludeOrigin(bc.Certificate.GetDiscoveryNodes(), *bc.Certificate.GetNodeRef())
	discoveryCount := len(discoveryNodes)
	if discoveryCount == 0 {
		panic("Zero bootstrap not allowed")
	}

	// var bootstrapResults []*network.BootstrapResult
	// var hosts []*host.Host
	// for {
	// 	ch := bc.GetDiscoveryNodesChannel(ctx, discoveryNodes, discoveryCount)
	// 	bootstrapResults, hosts = waitResultsFromChannel(ctx, ch, discoveryCount, bc.options.BootstrapTimeout)
	// 	if len(hosts) == discoveryCount {
	// 		// we connected to all discovery nodes
	// 		break
	// 	} else {
	// 		logger.Infof("[ BootstrapDiscovery ] Connected to %d/%d discovery nodes", len(hosts), discoveryCount)
	// 	}
	// }
	// reconnectRequests := 0
	// for _, bootstrapResult := range bootstrapResults {
	// 	if bootstrapResult.ReconnectRequired {
	// 		reconnectRequests++
	// 	}
	// }
	// minRequests := int(math.Floor(0.5*float64(discoveryCount))) + 1
	// if reconnectRequests >= minRequests {
	// 	logger.Infof("[ BootstrapDiscovery ] Need to reconnect as joiner (requested by %d/%d discovery nodes)",
	// 		reconnectRequests, discoveryCount)
	// 	return nil, ErrReconnectRequired
	// }
	// activeNodesStr := make([]string, 0)
	//
	// <-bc.bootstrapLock
	// logger.Debugf("[ BootstrapDiscovery ] After bootstrap lock")
	//
	// ch := bc.getGenesisRequestsChannel(ctx, hosts)
	// activeNodes, lastPulses, err := bc.waitGenesisResults(ctx, ch, len(hosts))
	// if err != nil {
	// 	return nil, err
	// }
	// bc.forceSetLastPulse(bc.calculateLastIgnoredPulse(ctx, lastPulses))
	// for _, activeNode := range activeNodes {
	// 	err = bc.checkActiveNode(activeNode)
	// 	if err != nil {
	// 		return nil, errors.Wrapf(err, "Discovery check of node %s failed", activeNode.ID())
	// 	}
	// 	activeNode.(node.MutableNode).SetState(insolar.NodeUndefined)
	// 	activeNodesStr = append(activeNodesStr, activeNode.ID().String())
	// }
	activeNodes := make([]insolar.NetworkNode, 0)
	bc.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodeUndefined)
	activeNodes = append(activeNodes, bc.NodeKeeper.GetOrigin())
	bc.NodeKeeper.SetInitialSnapshot(activeNodes)
	// logger.Infof("[ BootstrapDiscovery ] Added active nodes: %s", strings.Join(activeNodesStr, ", "))
	//
	// if bc.options.CyclicBootstrapEnabled {
	// 	go bc.startCyclicBootstrap(ctx)
	// }

	return nil, nil
	// return parseBootstrapResults(bootstrapResults), nil
}

// func (bc *Bootstrap) calculateLastIgnoredPulse(ctx context.Context, lastPulses []insolar.PulseNumber) insolar.PulseNumber {
// 	maxLastPulse := bc.GetLastPulse()
// 	inslogger.FromContext(ctx).Debugf("NetworkNode %s (origin) LastIgnoredPulse: %d", bc.NodeKeeper.GetOrigin().ID(), maxLastPulse)
// 	for _, pulse := range lastPulses {
// 		if pulse > maxLastPulse {
// 			maxLastPulse = pulse
// 		}
// 	}
// 	return maxLastPulse
// }

// func (bc *Bootstrap) GetDiscoveryNodesChannel(ctx context.Context, discoveryNodes []insolar.DiscoveryNode, needResponses int) <-chan *network.BootstrapResult {
// 	// we need only one host to bootstrap
// 	bootstrapResults := make(chan *network.BootstrapResult, needResponses)
// 	for _, discoveryNode := range discoveryNodes {
// 		go func(ctx context.Context, address string) {
// 			inslogger.FromContext(ctx).Infof("Starting bootstrap to address %s", address)
// 			ctx, span := instracer.StartSpan(ctx, "Bootstrapper.GetDiscoveryNodesChannel")
// 			defer span.End()
// 			span.AddAttributes(
// 				trace.StringAttribute("Bootstrap node", address),
// 			)
// 			bootstrapResult, err := bc.startBootstrap(ctx, address, nil) //bootstrap(ctx, address, bc.options, bc.startBootstrap, nil)
// 			if err != nil {
// 				inslogger.FromContext(ctx).Errorf("Error bootstrapping to address %s: %s", address, err.Error())
// 				return
// 			}
// 			bootstrapResults <- bootstrapResult
// 		}(ctx, discoveryNode.GetHost())
// 	}
//
// 	return bootstrapResults
// }

// timeout - bc.options.BootstrapTimeout
func WaitResultFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult, timeout time.Duration) *network.BootstrapResult {
	for {
		select {
		case bootstrapHost := <-ch:
			return bootstrapHost
		case <-time.After(timeout):
			inslogger.FromContext(ctx).Warn("Bootstrap timeout")
			return nil
		}
	}
}

func waitResultsFromChannel(ctx context.Context, ch <-chan *network.BootstrapResult, count int, timeout time.Duration) ([]*network.BootstrapResult, []*host.Host) {
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
		case <-time.After(timeout):
			inslogger.FromContext(ctx).Warnf("Bootstrap timeout, successful bootstraps: %d/%d", len(result), count)
			return result, hosts
		}
	}
}

func (bc *Bootstrap) startBootstrap(ctx context.Context, address string,
	perm *packet.Permission) (*network.BootstrapResult, error) {

	// PING request --------
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

	// BOOTSTRAP request --------
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

	// Process result
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
		// return bootstrap(ctx, data.Permission.Payload.ReconnectTo, bc.options, bc.startBootstrap, data.Permission)
	case packet.UpdateSchedule:
		// TODO: INS-1960
	}

	return &network.BootstrapResult{
		Host:              response.GetSenderHost(),
		ReconnectRequired: data.Code == packet.ReconnectRequired,
		NetworkSize:       int(data.NetworkSize),
	}, nil
}

// func (bc *Bootstrap) startCyclicBootstrap(ctx context.Context) {
// 	for atomic.LoadInt32(&bc.cyclicBootstrapStop) == 0 {
// 		results := make([]*network.BootstrapResult, 0)
// 		nodes := bc.getInactivenodes()
// 		for _, node := range nodes {
// 			res, err := bc.startBootstrap(ctx, node.GetHost(), nil)
// 			if err != nil {
// 				logger := inslogger.FromContext(ctx)
// 				logger.Errorf("[ StartCyclicBootstrap ] ", err)
// 				continue
// 			}
// 			results = append(results, res)
// 		}
// 		if len(results) != 0 {
// 			index := bc.getLargerNetworkIndex(results)
// 			if index >= 0 {
// 				bc.reconnectToNewNetwork(ctx, nodes[index].GetHost())
// 			}
// 		}
// 		time.Sleep(time.Second * bootstrapTimeout)
// 	}
// }

// func (bc *Bootstrap) getLargerNetworkIndex(results []*network.BootstrapResult) int {
// 	networkSize := results[0].NetworkSize
// 	index := 0
// 	for i := 1; i < len(results); i++ {
// 		if results[i].NetworkSize > networkSize {
// 			networkSize = results[i].NetworkSize
// 			index = i
// 		}
// 	}
// 	if networkSize > len(bc.NodeKeeper.GetAccessor().GetActiveNodes()) {
// 		return index
// 	}
// 	return -1
// }

// func (bc *Bootstrap) StopCyclicBootstrap() {
// 	atomic.StoreInt32(&bc.cyclicBootstrapStop, 1)
// }

// func (bc *Bootstrap) Init(ctx context.Context) error {
// 	bc.firstPulseTime = time.Now()
// 	bc.pinger = pinger.NewPinger(bc.Network)
// 	bc.Network.RegisterRequestHandler(types.Bootstrap, bc.processBootstrap)
// 	bc.Network.RegisterRequestHandler(types.Genesis, bc.processGenesis)
// 	return nil
// }

// func parseBootstrapResults(results []*network.BootstrapResult) *network.BootstrapResult {
// 	minIDIndex := 0
// 	minID := results[0].Host.NodeID
// 	for i, result := range results {
// 		if minID.Compare(result.Host.NodeID) > 0 {
// 			minIDIndex = i
// 		}
// 	}
// 	return results[minIDIndex]
// }

// func (bc *Bootstrap) getInactivenodes() []insolar.DiscoveryNode {
// 	res := make([]insolar.DiscoveryNode, 0)
// 	for _, node := range bc.Certificate.GetDiscoveryNodes() {
// 		if bc.NodeKeeper.GetAccessor().GetActiveNode(*node.GetNodeRef()) != nil {
// 			res = append(res, node)
// 		}
// 	}
// 	return res
// }

// func NewBootstrapper(options *common.Options /*, reconnectToNewNetwork func(ctx context.Context, address string)*/) *Bootstrap {
// 	return &Bootstrap{
// 		options:                 options,
// 		bootstrapLock:           make(chan struct{}),
// 		// reconnectToNewNetwork:   reconnectToNewNetwork,
// 	}
// }
