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
	"github.com/insolar/insolar/core"
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

// NewParticipant create new participant for consensus interfaces.
func NewParticipant(node *core.ActiveNode) Participant {
	return &participantWrapper{node: node}
}

type dataProviderWrapper struct {
	nodekeeper nodekeeper.NodeKeeper
}

func (dt *dataProviderWrapper) GetDataList() []*core.ActiveNode {
	return dt.nodekeeper.GetUnsync()
}

func (dt *dataProviderWrapper) MergeDataList(data []*core.ActiveNode) error {
	dt.nodekeeper.AddUnsyncGossip(data)
	return nil
}

// InsolarConsensus binds all functionality related to consensus with the network layer
type InsolarConsensus struct {
	consensus       Consensus
	communicatorSnd Communicator
	communicatorRcv Communicator
	nodekeeper      nodekeeper.NodeKeeper
}

// ProcessPulse is called when we get new pulse from pulsar. Should be called in goroutine
func (ic *InsolarConsensus) ProcessPulse(pulse core.Pulse) {

}

// NewInsolarConsensus creates new object to handle all consensus events
func NewInsolarConsensus(keeper nodekeeper.NodeKeeper) (*InsolarConsensus, error) {
	communicatorSnd := &communicatorSender{}
	communicatorRcv := &communicatorReceiver{}
	consensus, err := NewConsensus(&dataProviderWrapper{keeper}, communicatorSnd)
	if err != nil {
		return nil, errors.Wrap(err, "error creating insolar consensus")
	}
	return &InsolarConsensus{
		consensus:       consensus,
		communicatorSnd: communicatorSnd,
		communicatorRcv: communicatorRcv,
		nodekeeper:      keeper,
	}, nil
}
