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

package sm_execute_request

import (
	"context"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

type ExecuteIncomingRequest struct {
	RequestReference       insolar.Reference
	RequestObjectReference insolar.Reference
	RequestDeduplicated    bool
	Request                *record.IncomingRequest
	Result                 *record.Result
}

/* -------- Declaration ------------- */

func (s ExecuteIncomingRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	panic("implement me")
}

func (s ExecuteIncomingRequest) InjectDependencies(smachine.StateMachine, smachine.SlotLink, *injector.DependencyInjector) {
	panic("implement me")
}

func (s ExecuteIncomingRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	panic("implement me")
}

func (s ExecuteIncomingRequest) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	panic("implement me")
}

func (s ExecuteIncomingRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	panic("implement me")
}

/* -------- Instance ------------- */

func (s ExecuteIncomingRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return &s
}
