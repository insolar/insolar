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

package nodenetwork

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
)

type ConsensusInfo struct {
	lock                       sync.RWMutex
	tempMapR                   map[insolar.Reference]*host.Host
	tempMapS                   map[insolar.ShortNodeID]*host.Host
	nodesJoinedDuringPrevPulse bool
	isJoiner                   bool
}

func (ci *ConsensusInfo) SetIsJoiner(isJoiner bool) {
	ci.lock.Lock()
	defer ci.lock.Unlock()

	ci.isJoiner = isJoiner
}

func (ci *ConsensusInfo) IsJoiner() bool {
	ci.lock.RLock()
	defer ci.lock.RUnlock()

	return ci.isJoiner
}

func (ci *ConsensusInfo) NodesJoinedDuringPreviousPulse() bool {
	ci.lock.RLock()
	defer ci.lock.RUnlock()

	return ci.nodesJoinedDuringPrevPulse
}

func (ci *ConsensusInfo) AddTemporaryMapping(nodeID insolar.Reference, shortID insolar.ShortNodeID, address string) error {
	h, err := host.NewHostNS(address, nodeID, shortID)
	if err != nil {
		return errors.Wrapf(err, "Failed to generate address (%s, %s, %d)", address, nodeID, shortID)
	}
	ci.lock.Lock()
	ci.tempMapR[nodeID] = h
	ci.tempMapS[shortID] = h
	ci.lock.Unlock()
	log.Infof("Added temporary mapping: %s -> (%s, %d)", address, nodeID, shortID)
	return nil
}

func (ci *ConsensusInfo) ResolveConsensus(shortID insolar.ShortNodeID) *host.Host {
	ci.lock.RLock()
	defer ci.lock.RUnlock()

	return ci.tempMapS[shortID]
}

func (ci *ConsensusInfo) ResolveConsensusRef(nodeID insolar.Reference) *host.Host {
	ci.lock.RLock()
	defer ci.lock.RUnlock()

	return ci.tempMapR[nodeID]
}

func (ci *ConsensusInfo) Flush(nodesJoinedDuringPrevPulse bool) {
	ci.lock.Lock()
	defer ci.lock.Unlock()

	ci.tempMapR = make(map[insolar.Reference]*host.Host)
	ci.tempMapS = make(map[insolar.ShortNodeID]*host.Host)
	ci.nodesJoinedDuringPrevPulse = nodesJoinedDuringPrevPulse
}

func newConsensusInfo() *ConsensusInfo {
	return &ConsensusInfo{
		tempMapR: make(map[insolar.Reference]*host.Host),
		tempMapS: make(map[insolar.ShortNodeID]*host.Host),
	}
}
