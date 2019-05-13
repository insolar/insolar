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

package phases

import (
	"fmt"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type StateMatrix struct {
	data   [][]packets.BitSetState
	mapper packets.BitSetMapper
}

func NewStateMatrix(mapper packets.BitSetMapper) *StateMatrix {
	data := make([][]packets.BitSetState, mapper.Length())
	for i := 0; i < mapper.Length(); i++ {
		data[i] = make([]packets.BitSetState, mapper.Length())
		for j := 0; j < mapper.Length(); j++ {
			if i == j {
				data[i][j] = packets.Legit
			} else {
				data[i][j] = packets.TimedOut
			}
		}
	}
	return &StateMatrix{data: data, mapper: mapper}
}

func (sm *StateMatrix) ApplyBitSet(sender insolar.Reference, set packets.BitSet) error {
	array, err := set.GetTristateArray()
	if err != nil {
		return errors.Wrap(err, "Can't get tristate array from bitset")
	}
	if len(array) != sm.mapper.Length() {
		return errors.Wrapf(err, "Incorrect bitset length: %d != %d", len(array), sm.mapper.Length())
	}
	i, err := sm.mapper.RefToIndex(sender)
	if err != nil {
		return errors.Wrap(err, "Can't map sender reference to matrix index")
	}
	sm.data[i] = array
	return nil
}

type AdditionalRequest struct {
	RequestIndex int
	Candidates   []insolar.Reference
}

type Phase2MatrixState struct {
	// wether any nodes in current consensus need to advance to phase 2.1
	NeedPhase21 bool

	Active                   []insolar.Reference
	TimedOut                 []insolar.Reference
	AdditionalRequestsPhase2 []*AdditionalRequest
}

func newPhase2MatrixState() *Phase2MatrixState {
	return &Phase2MatrixState{
		NeedPhase21:              false,
		Active:                   make([]insolar.Reference, 0),
		TimedOut:                 make([]insolar.Reference, 0),
		AdditionalRequestsPhase2: make([]*AdditionalRequest, 0),
	}
}

func (sm *StateMatrix) CalculatePhase2(origin insolar.Reference) (*Phase2MatrixState, error) {
	originIndex, err := sm.mapper.RefToIndex(origin)
	if err != nil {
		return nil, errors.Wrap(err, "Can't map origin reference to matrix index")
	}
	result := newPhase2MatrixState()
	count := len(sm.data)
	// iterate through columns
	for j := 0; j < count; j++ {
		timedOuts := 0
		currentNeedsPhase21 := false
		for i := 0; i < count; i++ {
			if sm.data[i][j] == packets.TimedOut {
				timedOuts++

				if i == originIndex {
					// this flag is only actual if the node at index `j` is supposed to be active by majority of nodes
					currentNeedsPhase21 = true
				}
			}
		}
		active := consensusReachedMajority(count-timedOuts, count)
		if !active {
			timedOutRef, err := sm.mapper.IndexToRef(j)
			if err == nil {
				result.TimedOut = append(result.TimedOut, timedOutRef)
			}
			// else we can ignore error, we can do nothing with unknown node on such index
			continue
		}
		if timedOuts > 0 {
			result.NeedPhase21 = true
		}
		if currentNeedsPhase21 {
			err := sm.appendAdditionalRequest(result, j)
			if err != nil {
				return nil, err
			}
		}
		err = sm.fillState2Result(result, j)
		if err != nil {
			return nil, errors.Wrapf(err, "Can't set active matrix state for index %d", j)
		}
	}
	return result, nil
}

func (sm *StateMatrix) ReceivedProofFromNode(origin, nodeID insolar.Reference) error {
	originIndex, err := sm.mapper.RefToIndex(origin)
	if err != nil {
		return errors.Wrap(err, "Can't map origin reference to matrix index")
	}
	nodeIndex, err := sm.mapper.RefToIndex(nodeID)
	if err != nil {
		return errors.Wrap(err, "Can't map node reference to matrix index")
	}
	sm.data[originIndex][nodeIndex] = packets.Legit
	return nil
}

func (sm *StateMatrix) calculateAdditionalRequest(timedOutNodeIndex int) (*AdditionalRequest, error) {
	candidates := make([]insolar.Reference, 0)
	for i := 0; i < len(sm.data); i++ {
		if sm.data[i][timedOutNodeIndex] != packets.Legit {
			continue
		}
		nodeID, err := sm.mapper.IndexToRef(i)
		if err == packets.ErrBitSetNodeIsMissing {
			log.Warnf("Error mapping matrix index %d to node reference ID", i)
			continue
		}
		if err == packets.ErrBitSetIncorrectNode {
			return nil, errors.Wrapf(err, "Error mapping matrix index %d to node reference ID", i)
		}
		candidates = append(candidates, nodeID)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("Could not get candidates for matrix index %d", timedOutNodeIndex))
	}
	return &AdditionalRequest{RequestIndex: timedOutNodeIndex, Candidates: candidates}, nil
}

func (sm *StateMatrix) fillState2Result(result *Phase2MatrixState, index int) error {
	nodeID, err := sm.mapper.IndexToRef(index)
	if err == packets.ErrBitSetIncorrectNode {
		return errors.Wrapf(err, "Error mapping matrix index %d to node reference ID", index)
	}
	if err == packets.ErrBitSetNodeIsMissing {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "fillState2Result unknown error")
	}
	result.Active = append(result.Active, nodeID)
	return nil
}

func (sm *StateMatrix) appendAdditionalRequest(result *Phase2MatrixState, index int) error {
	req, err := sm.calculateAdditionalRequest(index)
	if err != nil {
		return errors.Wrapf(err, "Could not generate additional phase 2.1 request for index %d", index)
	}
	result.AdditionalRequestsPhase2 = append(result.AdditionalRequestsPhase2, req)
	return nil
}
