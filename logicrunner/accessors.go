//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

// RegisterExecutor registers an executor for particular `MachineType`
func (lr *LogicRunner) RegisterExecutor(t insolar.MachineType, e insolar.MachineLogicExecutor) error {
	lr.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (lr *LogicRunner) GetExecutor(t insolar.MachineType) (insolar.MachineLogicExecutor, error) {
	if res := lr.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.Errorf("No executor registered for machine %d", int(t))
}

func (lr *LogicRunner) GetObjectState(ref Ref) *ObjectState {
	lr.stateMutex.RLock()
	res, ok := lr.state[ref]
	lr.stateMutex.RUnlock()
	if !ok {
		return nil
	}
	return res
}

func (lr *LogicRunner) UpsertObjectState(ref Ref) *ObjectState {
	lr.stateMutex.RLock()
	if res, ok := lr.state[ref]; ok {
		lr.stateMutex.RUnlock()
		return res
	}
	lr.stateMutex.RUnlock()

	lr.stateMutex.Lock()
	defer lr.stateMutex.Unlock()
	if _, ok := lr.state[ref]; !ok {
		lr.state[ref] = &ObjectState{}
	}
	return lr.state[ref]
}

func (lr *LogicRunner) MustObjectState(ref Ref) *ObjectState {
	res := lr.GetObjectState(ref)
	if res == nil {
		panic("No requested object state. ref: " + ref.String())
	}
	return res
}

func (lr *LogicRunner) GetExecutionState(ref Ref) *ExecutionState {
	os := lr.GetObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()
	return os.ExecutionState
}

func (lr *LogicRunner) pulse(ctx context.Context) *insolar.Pulse {
	pulse, err := lr.PulseAccessor.Latest(ctx)
	if err != nil {
		panic(err)
	}
	return &pulse
}

func (lr *LogicRunner) GetConsensus(ctx context.Context, ref Ref) *Consensus {
	state := lr.UpsertObjectState(ref)

	state.Lock()
	defer state.Unlock()

	if state.Consensus == nil {
		validators, err := lr.JetCoordinator.QueryRole(
			ctx,
			insolar.DynamicRoleVirtualValidator,
			*ref.Record(),
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

func (st *ObjectState) StartValidation(ref Ref) *ExecutionState {
	st.Lock()
	defer st.Unlock()

	if st.Validation != nil {
		panic("Unexpected. Validation already in progress")
	}
	st.Validation = NewExecutionState(ref)
	return st.Validation
}
