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

package phases

import (
	"fmt"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type StateMatrix struct {
	data   [][]packets.TriState
	mapper packets.BitSetMapper
}

func NewStateMatrix(mapper packets.BitSetMapper) *StateMatrix {
	data := make([][]packets.TriState, mapper.Length())
	for i := 0; i < mapper.Length(); i++ {
		data[i] = make([]packets.TriState, mapper.Length())
		for j := 0; j < mapper.Length(); j++ {
			data[i][j] = packets.TimedOut
		}
	}
	return &StateMatrix{data: data, mapper: mapper}
}

func (sm *StateMatrix) ApplyBitSet(sender core.RecordRef, set packets.BitSet) error {
	cells, err := set.GetCells(sm.mapper)
	if err != nil {
		return errors.Wrap(err, "Can't get cells from bitset")
	}
	i, err := sm.mapper.RefToIndex(sender)
	if err != nil {
		return errors.Wrap(err, "Can't map sender reference to matrix index")
	}
	for _, cell := range cells {
		j, err := sm.mapper.RefToIndex(cell.NodeID)
		if err != nil {
			return errors.Wrap(err, "Can't map cell reference to matrix index")
		}
		sm.data[i][j] = cell.State
	}
	return nil
}

type AdditionalRequest struct {
	RequestIndex int
	Candidates   []core.RecordRef
}

type Phase2MatrixState struct {
	// wether any nodes in current consensus need to advance to phase 2.1
	NeedPhase21 bool

	Active                   []core.RecordRef
	AdditionalRequestsPhase2 []*AdditionalRequest
}

func newPhase2MatrixState() *Phase2MatrixState {
	return &Phase2MatrixState{
		NeedPhase21:              false,
		Active:                   make([]core.RecordRef, 0),
		AdditionalRequestsPhase2: make([]*AdditionalRequest, 0),
	}
}

func (sm *StateMatrix) CalculatePhase2(origin core.RecordRef) (*Phase2MatrixState, error) {
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
		active := consensusReachedMajority(timedOuts, count)
		if !active {
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
		sm.fillState2Result(result, j)
	}
	return result, nil
}

func (sm *StateMatrix) calculateAdditionalRequest(timedOutNodeIndex int) (*AdditionalRequest, error) {
	candidates := make([]core.RecordRef, 0)
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
		return sm.appendAdditionalRequest(result, index)
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
