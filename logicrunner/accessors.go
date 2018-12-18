/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// RegisterExecutor registers an executor for particular `MachineType`
func (lr *LogicRunner) RegisterExecutor(t core.MachineType, e core.MachineLogicExecutor) error {
	lr.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (lr *LogicRunner) GetExecutor(t core.MachineType) (core.MachineLogicExecutor, error) {
	if res := lr.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.Errorf("No executor registered for machine %d", int(t))
}

func (lr *LogicRunner) GetObjectState(ref Ref) *ObjectState {
	lr.stateMutex.Lock()
	defer lr.stateMutex.Unlock()
	res, ok := lr.state[ref]
	if !ok {
		return nil
	}
	return res
}

func (lr *LogicRunner) UpsertObjectState(ref Ref) *ObjectState {
	lr.stateMutex.Lock()
	defer lr.stateMutex.Unlock()
	if _, ok := lr.state[ref]; !ok {
		lr.state[ref] = &ObjectState{Ref: &ref}
	}
	return lr.state[ref]
}

func (lr *LogicRunner) MustObjectState(ref Ref) *ObjectState {
	res := lr.GetObjectState(ref)
	if res == nil {
		panic("No requested object state")
	}
	return res
}

func (lr *LogicRunner) pulse(ctx context.Context) *core.Pulse {
	pulse, err := lr.PulseStorage.Current(ctx)
	if err != nil {
		panic(err)
	}
	return pulse
}

func (lr *LogicRunner) GetConsensus(ctx context.Context, ref Ref) *Consensus {
	state := lr.UpsertObjectState(ref)

	state.Lock()
	defer state.Unlock()

	if state.Consensus == nil {
		validators, err := lr.JetCoordinator.QueryRole(
			ctx,
			core.DynamicRoleVirtualValidator,
			ref.Record(),
			lr.pulse(ctx).PulseNumber,
		)
		if err != nil {
			panic("cannot QueryRole")
		}
		// TODO INS-732 check pulse of message and ensure we deal with right validator
		state.Consensus = newConsensus(lr, validators)
	}
	return state.Consensus
}

func (st *ObjectState) RefreshConsensus() {
	if st.Consensus == nil {
		return
	}

	st.Consensus.ready = true
	st.Consensus = nil
}

func (st *ObjectState) StartValidation() *ExecutionState {
	st.Lock()
	defer st.Unlock()

	if st.Validation != nil {
		panic("Unexpected. Validation already in progress")
	}
	st.Validation = &ExecutionState{}
	return st.Validation
}
