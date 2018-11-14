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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
)

// Participant describes one consensus participant
type Participant interface {
	GetActiveNode() core.Node
}

// UnsyncHolder
type UnsyncHolder interface {
	// GetUnsync returns list of local unsync nodes. This list is created
	GetUnsync() []core.Node
	// GetPulse returns actual pulse for current consensus process.
	GetPulse() core.PulseNumber
	// SetHash sets hash of unsync lists for each node of consensus.
	SetHash([]*network.NodeUnsyncHash)
	// GetHash get hash of unsync lists for each node of  If hash is not calculated yet, then this call blocks
	// until the hash is calculated with SetHash() call.
	GetHash(blockTimeout time.Duration) ([]*network.NodeUnsyncHash, error)
}

// Consensus interface provides method to make consensus between participants
type Consensus interface {
	// DoConsensus is sync method, it performs all consensus steps and returns list of synced nodes
	// method should be executed in goroutine
	DoConsensus(ctx context.Context, holder UnsyncHolder, self Participant, allParticipants []Participant) ([]core.Node, error)
}

// Communicator interface is used to exchange messages between participants
type Communicator interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, pulse core.PulseNumber, p Participant, data []core.Node) ([]core.Node, error)

	// ExchangeHash used in second consensus step to exchange only hashes of merged data vectors
	ExchangeHash(ctx context.Context, pulse core.PulseNumber, p Participant, data []*network.NodeUnsyncHash) ([]*network.NodeUnsyncHash, error)
}

// NewConsensus creates consensus
func NewConsensus(communicator Communicator, scheme core.PlatformCryptographyScheme) Consensus {
	return &baseConsensus{communicator: communicator, scheme: scheme}
}
