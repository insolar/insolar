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

func (lr *LogicRunner) GetExecution(ref Ref) *ExecutionState {
	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()
	res, ok := lr.execution[ref]
	if !ok {
		return nil
	}
	return res
}

func (lr *LogicRunner) UpsertExecution(ref Ref) *ExecutionState {
	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()
	if _, ok := lr.execution[ref]; !ok {
		lr.execution[ref] = &ExecutionState{
			Ref: &ref,
			caseBind: core.CaseBind{
				Requests: make([]core.CaseRequest, 0),
			},
		}
	}
	return lr.execution[ref]
}

func (lr *LogicRunner) MustExecutionState(ref Ref) *ExecutionState {
	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()
	res, ok := lr.execution[ref]
	if !ok {
		panic("No requested execution state")
	}
	return res
}

func (lr *LogicRunner) addObjectCaseRecord(ref Ref, cr core.CaseRecord) {
	lr.MustExecutionState(ref).AddCaseRecord(cr)
}

func (lr *LogicRunner) nextValidationStep(ref Ref) (*core.CaseRecord, int) {
	lr.caseBindReplaysMutex.Lock()
	defer lr.caseBindReplaysMutex.Unlock()

	r, ok := lr.caseBindReplays[ref]
	if !ok {
		return nil, -1
	}
	record, step := r.NextStep()
	lr.caseBindReplays[ref] = r
	return record, step
}

func (lr *LogicRunner) pulse(ctx context.Context) *core.Pulse {
	pulse, err := lr.PulseManager.Current(ctx)
	if err != nil {
		panic(err)
	}
	return pulse
}

func (lr *LogicRunner) GetConsensus(ctx context.Context, r Ref) (*Consensus, bool) {
	lr.consensusMutex.Lock()
	defer lr.consensusMutex.Unlock()
	c, ok := lr.consensus[r]
	if !ok {
		validators, err := lr.JetCoordinator.QueryRole(
			ctx,
			core.DynamicRoleVirtualValidator,
			&r,
			lr.pulse(ctx).PulseNumber,
		)
		if err != nil {
			panic("cannot QueryRole")
		}
		// TODO INS-732 check pulse of message and ensure we deal with right validator
		c = newConsensus(lr, validators)
		lr.consensus[r] = c
	}
	return c, ok
}

func (lr *LogicRunner) RefreshConsensus() {
	lr.consensusMutex.Lock()
	defer lr.consensusMutex.Unlock()
	if lr.consensus == nil {
		lr.consensus = make(map[Ref]*Consensus)
		return
	}
	for k, c := range lr.consensus {
		if c.ready {
			delete(lr.consensus, k)
		} else {
			c.ready = true
		}
	}
}
