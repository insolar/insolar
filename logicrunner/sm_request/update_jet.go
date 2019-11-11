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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/s_jet_storage"
	"github.com/insolar/insolar/logicrunner/s_sender"
)

type StateMachineUpdateJet struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.UpdateJet

	sender     *s_sender.SenderServiceAdapter
	jetStorage *s_jet_storage.JetStorageServiceAdapter

	externalError error
}

var declUpdateJet smachine.StateMachineDeclaration = &declarationUpdateJet{}

type declarationUpdateJet struct {
	smachine.StateMachineDeclTemplate
}

/* -------- Declaration ------------- */

func (declarationUpdateJet) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineUpdateJet)
	return s.Init
}

func (declarationUpdateJet) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachineUpdateJet)

	injector.MustInject(&s.sender)
	injector.MustInject(&s.jetStorage)
}

/* -------- Instance ------------- */

func (s *StateMachineUpdateJet) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declUpdateJet
}

func (s *StateMachineUpdateJet) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepUpdateJet)
}

func (s *StateMachineUpdateJet) stepUpdateJet(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()
	pl := s.Payload

	return s.jetStorage.PrepareAsync(ctx, func(svc s_jet_storage.JetStorageService) smachine.AsyncResultFunc {
		err := svc.Update(goCtx, pl.Pulse, true, pl.JetID)
		return func(ctx smachine.AsyncResultContext) {
			s.externalError = errors.Wrap(err, "failed to update jets")
		}
	}).DelayedStart().Sleep().ThenJump(s.stepStop)
}

func (s *StateMachineUpdateJet) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()

	if s.externalError != nil {
		s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
			bus.ReplyError(goCtx, svc, *s.Meta, s.externalError)
		}).DelayedSend()

		return ctx.Error(s.externalError)
	}

	return ctx.Stop()
}
