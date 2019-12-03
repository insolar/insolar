//
// Copyright 2019 Insolar Technologies GmbHlf
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
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ObjectAccessor struct {
	ObjectReference insolar.Reference

	objectCatalog   sm_object.LocalObjectCatalog // dependency
	sharedStateLink sm_object.SharedObjectStateAccessor
}

func (a *ObjectAccessor) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&a.objectCatalog)
}

func (a *ObjectAccessor) prepare(ctx smachine.ExecutionContext) {
	if a.ObjectReference.IsEmpty() {
		panic("ObjectReference in object accessor is empty")
	}
	if a.sharedStateLink.IsZero() {
		a.sharedStateLink = a.objectCatalog.GetOrCreate(ctx, a.ObjectReference)
	}
}

func (a *ObjectAccessor) prepareAccess(ctx smachine.ExecutionContext, cb func(objectState *sm_object.SharedObjectState)) smachine.StateUpdate {
	a.prepare(ctx)

	switch a.sharedStateLink.Prepare(cb).TryUse(ctx).GetDecision() {
	case smachine.NotPassed:
		return ctx.WaitShared(a.sharedStateLink.SharedDataLink).ThenRepeat()
	case smachine.Impossible:
		// the holder of the sharedState is stopped
		return ctx.Stop()
	case smachine.Passed:
	default:
		panic("unknown state from TryUse")
	}

	return smachine.StateUpdate{}
}

func (a *ObjectAccessor) setPreviousExecutorState(ctx smachine.ExecutionContext, state payload.PreviousExecutorState) smachine.StateUpdate {
	a.prepare(ctx)

	return a.prepareAccess(ctx, func(objectState *sm_object.SharedObjectState) {
		switch {
		case objectState.PreviousExecutorState > state:
			// do nothing here, since our previous information is more "reliable" at the moment
		case state == payload.PreviousExecutorFinished:
			// wake up and notify about new state
			// all other checks should be inside of PreviousExecutorNotification function
			if !ctx.ApplyAdjustment(objectState.PreviousExecutorFinished.NewValue(true)) {
				panic("failed to apply adjustement")
			}

			fallthrough
		default:
			// change state and forget about it (??)
			objectState.PreviousExecutorState = state
		}
	})
}

type StateMachineStillExecuting struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.StillExecuting

	objectAccessor ObjectAccessor
}

/* -------- Declaration ------------- */

var declStillExecuting smachine.StateMachineDeclaration = &declarationStillExecuting{}

type declarationStillExecuting struct {
	smachine.StateMachineDeclTemplate
}

func (declarationStillExecuting) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineStillExecuting)
	return s.Init
}

func (declarationStillExecuting) InjectDependencies(sm smachine.StateMachine, sl smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachineStillExecuting)
	s.objectAccessor.InjectDependencies(sm, sl, injector)
}

/* -------- Instance ------------- */

func (s *StateMachineStillExecuting) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declStillExecuting
}

func (s *StateMachineStillExecuting) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.objectAccessor.ObjectReference = s.Payload.ObjectRef
	// TODO[bigbes]: we should do there something with executed RequestReferences
	return ctx.Jump(s.stepSetStillExecuting)
}

func (s *StateMachineStillExecuting) stepSetStillExecuting(ctx smachine.ExecutionContext) smachine.StateUpdate {
	stateUpdate := s.objectAccessor.setPreviousExecutorState(ctx, payload.PreviousExecutorFinished)

	if !stateUpdate.IsZero() {
		return stateUpdate
	}

	return ctx.Stop()
}
