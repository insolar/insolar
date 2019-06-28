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
)

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
