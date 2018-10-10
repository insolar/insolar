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

package consensus

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/pkg/errors"
)

type participantWrapper struct {
	node *core.ActiveNode
}

// GetActiveNode implements Participant interface for ActiveNode wrapper.
func (an *participantWrapper) GetActiveNode() *core.ActiveNode {
	return an.node
}

type selfWrapper struct {
	keeper nodekeeper.NodeKeeper
}

// GetActiveNode implements Participant interface for NodeKeeper wrapper.
func (s *selfWrapper) GetActiveNode() *core.ActiveNode {
	return s.keeper.GetSelf()
}

// NetworkConsensus binds all functionality related to consensus with the network layer
type NetworkConsensus struct {
	consensus       consensus.Consensus
	communicatorSnd consensus.Communicator
	communicatorRcv consensus.Communicator
	keeper          nodekeeper.NodeKeeper
	self            *selfWrapper
}

// ProcessPulse is called when we get new pulse from pulsar. Should be called in goroutine
func (ic *NetworkConsensus) ProcessPulse(ctx context.Context, pulse core.Pulse) {
	activeNodes := ic.keeper.GetActiveNodes()
	if len(activeNodes) == 0 {
		return
	}
	participants := make([]consensus.Participant, len(activeNodes))
	for i, activeNode := range activeNodes {
		participants[i] = &participantWrapper{activeNode}
	}
	success, unsyncList := ic.keeper.SetPulse(pulse.PulseNumber)
	holder := nodekeeper.NewUnsyncHolder(pulse.PulseNumber, unsyncList)
	if !success {
		log.Error("InsolarConsensus: could not set new pulse to NodeKeeper, aborting")
		return
	}
	unsyncCandidates, err := ic.consensus.DoConsensus(holder, ctx, ic.self, participants)
	if err != nil {
		log.Errorf("InsolarConsensus: error performing consensus steps: %s", err.Error())
	}
	// We have to keep in mind a scenario when DoConsensus takes too long time and a new ProcessPulse is called
	// simultaneously with the current call. It will happen if DoConsensus takes more time than the delay between two
	// consecutive pulses.
	// In this scenario ic.keeper.SetPulse(pulse + 1) will happen earlier than ic.keeper.Sync(syncCandidates, pulse).
	// ic.keeper.SetPulse(pulse + 1) will internally call Sync(nil) to update NodeKeeper's unsync, sync and active lists.
	// That's why we have to pass PulseNumber to ic.keeper.Sync to check relevance of the pulse and to ignore the call
	// if we detect this kind of race condition.
	ic.keeper.Sync(unsyncCandidates, pulse.PulseNumber)
}

// IsPartOfConsensus returns whether we should perform all consensus interactions or not
func (ic *NetworkConsensus) IsPartOfConsensus() bool {
	return ic.keeper.GetSelf() != nil
}

// ReceiverHandler return handler that is responsible to handle consensus network requests
func (ic *NetworkConsensus) ReceiverHandler() consensus.Communicator {
	return ic.communicatorRcv
}

// NewInsolarConsensus creates new object to handle all consensus events
func NewInsolarConsensus(keeper nodekeeper.NodeKeeper, handler hosthandler.HostHandler) (consensus.InsolarConsensus, error) {
	communicatorSnd := &communicatorSender{handler}
	communicatorRcv := &communicatorReceiver{keeper}
	consensus, err := consensus.NewConsensus(communicatorSnd)
	if err != nil {
		return nil, errors.Wrap(err, "error creating insolar consensus")
	}
	return &NetworkConsensus{
		consensus:       consensus,
		communicatorSnd: communicatorSnd,
		communicatorRcv: communicatorRcv,
		keeper:          keeper,
		self:            &selfWrapper{keeper},
	}, nil
}
