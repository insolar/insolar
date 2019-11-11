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

package outgoing

import (
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/s_artifact"
)

type RPCGetCodeStateMachine struct {
	smachine.StateMachineDeclTemplate

	// input
	CodeReference insolar.Reference

	// dependencies
	ArtifactManager *s_artifact.ArtifactClientServiceAdapter

	// pass between stages
	externalError error
	code          []byte
}

/* -------- Declaration ------------- */

func (s *RPCGetCodeStateMachine) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *RPCGetCodeStateMachine) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, di *injector.DependencyInjector) {
	_ = sm.(*RPCGetCodeStateMachine)
}

/* -------- Instance ------------- */

func (s *RPCGetCodeStateMachine) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *RPCGetCodeStateMachine) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepGetCode)
}

func (s *RPCGetCodeStateMachine) stepGetCode(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()

	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		desc, err := svc.GetCode(goCtx, s.CodeReference)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			if err != nil {
				s.code, s.externalError = desc.Code()
			}
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepReturn)
}

func (s *RPCGetCodeStateMachine) stepReturn(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var err error

	if s.externalError != nil {
		// return error
	} else {
		// return ok
	}

	return ctx.Error(err)
}
