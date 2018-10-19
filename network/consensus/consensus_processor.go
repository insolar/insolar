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
)

// NodeKeeper manages unsync, sync and active lists
type NodeKeeper interface {
	// GetID get current node ID
	GetID() core.RecordRef
	// GetSelf get active node for the current insolard. Returns nil if the current insolard is not an active node.
	GetSelf() *core.ActiveNode
	// GetActiveNode get active node by its reference. Returns nil if node is not found.
	GetActiveNode(ref core.RecordRef) *core.ActiveNode
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// GetActiveNodesByRole get active nodes by role
	GetActiveNodesByRole(role core.JetRole) []core.RecordRef
	// AddActiveNodes add active nodes.
	AddActiveNodes([]*core.ActiveNode)
	// SetPulse sets internal PulseNumber to number. Returns true if set was successful, false if number is less
	// or equal to internal PulseNumber. If set is successful, returns collected unsync list and starts collecting new unsync list
	SetPulse(number core.PulseNumber) (bool, UnsyncList)
	// Sync initiates transferring syncCandidates -> sync, sync -> active.
	// If number is less than internal PulseNumber then ignore Sync.
	Sync(syncCandidates []*core.ActiveNode, number core.PulseNumber)
	// AddUnsync add unsync node to the unsync list. Returns channel that receives active node on successful sync.
	// Channel will return nil node if added node has not passed the consensus.
	// Returns error if current node is not active and cannot participate in consensus.
	AddUnsync(nodeID core.RecordRef, roles []core.NodeRole, address string /*, publicKey *ecdsa.PublicKey*/) (chan *core.ActiveNode, error)
	// GetUnsyncHolder get unsync list executed in consensus for specific pulse.
	// 1. If pulse is less than internal NodeKeeper pulse, returns error.
	// 2. If pulse is equal to internal NodeKeeper pulse, returns unsync list holder for currently executed consensus.
	// 3. If pulse is more than internal NodeKeeper pulse, blocks till next SetPulse or duration timeout and then acts like in par. 2
	GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (UnsyncList, error)
}

type UnsyncList interface {
	// GetUnsync returns list of local unsync nodes. This list is created
	GetUnsync() []*core.ActiveNode
	// GetPulse returns actual pulse for current consensus process.
	GetPulse() core.PulseNumber
	// SetHash sets hash of unsync lists for each node of consensus.
	SetHash([]*NodeUnsyncHash)
	// GetHash get hash of unsync lists for each node of consensus. If hash is not calculated yet, then this call blocks
	// until the hash is calculated with SetHash() call
	GetHash(blockTimeout time.Duration) ([]*NodeUnsyncHash, error)
	// AddUnsyncList add unsync list for remote ref
	AddUnsyncList(ref core.RecordRef, unsync []*core.ActiveNode)
	// AddUnsyncHash add unsync hash for remote ref
	AddUnsyncHash(ref core.RecordRef, hash []*NodeUnsyncHash)
	// GetUnsyncList get unsync list for remote ref
	GetUnsyncList(ref core.RecordRef) ([]*core.ActiveNode, bool)
	// GetUnsyncHash get unsync hash for remote ref
	GetUnsyncHash(ref core.RecordRef) ([]*NodeUnsyncHash, bool)
}

// CommunicatorReceiver
type CommunicatorReceiver interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, pulse core.PulseNumber, from core.RecordRef, data []*core.ActiveNode) ([]*core.ActiveNode, error)
	// ExchangeHash used in second consensus step to exchange only hashes of merged data vectors
	ExchangeHash(ctx context.Context, pulse core.PulseNumber, from core.RecordRef, data []*NodeUnsyncHash) ([]*NodeUnsyncHash, error)
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
	SetNodeKeeper(keeper NodeKeeper)
}
