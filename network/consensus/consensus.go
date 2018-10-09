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
	"errors"

	"github.com/insolar/insolar/core"
)

// Participant describes one consensus participant
type Participant interface {
	GetActiveNode() *core.ActiveNode
}

// DataProvider for data manipulation
type DataProvider interface {
	// GetDataList get active nodes to exchange with other parties of the consensus
	GetDataList() []*core.ActiveNode
	// MergeDataList merge active nodes from other parties of the consensus
	MergeDataList([]*core.ActiveNode) error
	// GetHash get hash of merged data vectors
	GetHash() (hash []byte, err error)
}

// Consensus interface provides method to make consensus between participants
type Consensus interface {
	// DoConsensus is sync method, it make all consensus steps and returns boolean result
	// method should be executed in goroutine
	DoConsensus(ctx context.Context, self Participant, allParticipants []Participant) (bool, error)
}

// Communicator interface is used to exchange messages between participants
type Communicator interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, p Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error)

	// ExchangeHash used in second consensus step to exchange only hashes of merged data vectors
	ExchangeHash(ctx context.Context, p Participant, data []byte) ([]byte, error)
}

// NewConsensus creates consensus
func NewConsensus(provider DataProvider, communicator Communicator) (Consensus, error) {
	return nil, errors.New("not implemented")
}
