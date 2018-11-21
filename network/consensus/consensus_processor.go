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
	"github.com/insolar/insolar/network"
)

// CommunicatorReceiver
type CommunicatorReceiver interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, pulse core.PulseNumber, from core.RecordRef, data []core.Node) ([]core.Node, error)
	// ExchangeHash used in second consensus step to exchange only hashes of merged data vectors
	ExchangeHash(ctx context.Context, pulse core.PulseNumber, from core.RecordRef, data []*network.NodeUnsyncHash) ([]*network.NodeUnsyncHash, error)
}

// Processor is an interface to bind all functionality related to consensus with the network layer
type Processor interface {
	// ProcessPulse is called when we get new pulse from pulsar. Should be called in goroutine
	ProcessPulse(ctx context.Context, pulse core.Pulse)
	// IsPartOfConsensus returns whether we should perform all consensus interactions or not
	IsPartOfConsensus() bool
	// ReceiverHandler return handler that is responsible to handle consensus network requests
	ReceiverHandler() CommunicatorReceiver
	// SetNodeKeeper set NodeKeeper for the processor to integrate Processor with unsync -> sync -> active pipeline
	SetNodeKeeper(keeper network.NodeKeeper)
}
