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

package sm_request

import (
	"context"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/payload"
)

type StateMachineStillExecuting struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.StillExecuting
}

var declStillExecuting smachine.StateMachineDeclaration = declarationStillExecuting{}

type declarationStillExecuting struct{}

func (declarationStillExecuting) GetStepLogger(context.Context, smachine.StateMachine) smachine.StateMachineStepLoggerFunc {
	return nil
}

func (declarationStillExecuting) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	_ = sm.(*StateMachineStillExecuting)
}

func (declarationStillExecuting) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

func (declarationStillExecuting) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (declarationStillExecuting) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineStillExecuting)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineStillExecuting) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declStillExecuting
}

func (s *StateMachineStillExecuting) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}
