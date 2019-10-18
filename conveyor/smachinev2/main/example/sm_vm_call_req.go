//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package example

import (
	"github.com/insolar/insolar/conveyor/injector"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/longbits"
)

type StateMachineCallRequest struct {
	catalogObj LocalObjectCatalog

	callerObj, calleeObj longbits.ByteString
	callMethod           string

	sharedStateLink SharedObjectStateAccessor

	objInfo    ObjectInfo
	callType   ContractCallType
	callResult CallResult
}

/* -------- Declaration ------------- */

var declCallRequest smachine.StateMachineDeclaration = declarationCallRequest{}

type declarationCallRequest struct{}

func (declarationCallRequest) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
}

func (declarationCallRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

func (declarationCallRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (declarationCallRequest) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineCallRequest)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineCallRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declCallRequest
}

func (s *StateMachineCallRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.callerObj = longbits.NewByteStringOf("testObjectA")
	s.calleeObj = longbits.NewByteStringOf("testObjectB")
	s.callMethod = "someMethod"

	ctx.SetDefaultMigration(s.cancelOnMigrate)
	ctx.SetDefaultErrorHandler(s.sendReplyOnError)

	return ctx.Jump(s.stateGetSharedReadyToWork)
}

func (s *StateMachineCallRequest) cancelOnMigrate(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)
	return ctx.Jump(s.sendReplyOnCancel)
}

func (s *StateMachineCallRequest) stateGetSharedReadyToWork(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.sharedStateLink.IsZero() {
		s.sharedStateLink = s.catalogObj.GetOrCreate(ctx, s.calleeObj)
	}

	var readyToWork smachine.SyncLink

	switch s.sharedStateLink.Prepare(
		func(state *SharedObjectState) {
			readyToWork = state.SemaReadyToWork
			s.objInfo = state.ObjectInfo // it may need to be re-fetched
		}).TryUse(ctx).GetDecision() {
	case smachine.NotPassed:
		return ctx.WaitShared(s.sharedStateLink.SharedDataLink).ThenRepeat()
	case smachine.Impossible:
		// the holder of the sharedState is stopped
		return ctx.Stop()
	}

	if !s.objInfo.IsReadyToWork && ctx.AcquireForThisStep(readyToWork).IsNotPassed() {
		return ctx.Sleep().ThenRepeat()
	}

	objCode := s.objInfo.ObjectLatestValidCode
	callMethod := s.callMethod

	s.objInfo.ContractRunner.PrepareAsync(ctx, func(svc ContractRunnerService) smachine.AsyncResultFunc {
		callType := svc.ClassifyCall(objCode, callMethod)
		return func(ctx smachine.AsyncResultContext) {
			s.callType = callType
			ctx.WakeUp()
		}
	}).Start()

	return ctx.Sleep().ThenJump(s.stateSharedReadyToWork)
}

func (s *StateMachineCallRequest) stateSharedReadyToWork(ctx smachine.ExecutionContext) smachine.StateUpdate {
	switch s.callType {
	case ContractCallMutable:
		return ctx.Jump(s.statePrepareMutableCall)
	case ContractCallImmutable:
		return ctx.Jump(s.stateStartImmutableCall)
	case ContractCallSaga:
		return ctx.Jump(s.stateRegisterSagaCall)
	default:
		panic("illegal state")
	}
}

/* ================ Immutable call scenario ============== */

func (s *StateMachineCallRequest) stateStartImmutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if !ctx.AcquireForThisStep(s.objInfo.ImmutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	objCode := s.objInfo.ObjectLatestValidCode
	objState := s.objInfo.ObjectLatestValidState
	callMethod := s.callMethod

	s.objInfo.ContractRunner.PrepareAsync(ctx, func(svc ContractRunnerService) smachine.AsyncResultFunc {
		result := svc.CallImmutableMethod(objCode, callMethod, objState)
		return func(ctx smachine.AsyncResultContext) {
			s.callResult = result
			ctx.WakeUp()
		}
	}).Start()

	return ctx.Sleep().ThenJump(s.stateDoneImmutableCall)
}

func (s *StateMachineCallRequest) stateDoneImmutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
			Steps:
		    1. register resulting immutable state on ledger (and request if a normal call)
		    2. report resulting immutable state to SharedObjectStateAccessor
			3. send results back
	*/
	panic("unimplemented")
}

/* ================ Saga call scenario ============== */

func (s *StateMachineCallRequest) stateRegisterSagaCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
		Steps:
		1. register Saga call on ledger
		2. send confirmation to caller
		3. continue as mutable call
	*/
	panic("unimplemented")
}

/* ================ Mutable call scenario ============== */

func (s *StateMachineCallRequest) statePrepareMutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
			Steps:
			1. lock on limiter for mutable calls
			2. check ordering (reenter queue on mutable limiter if ordering is wrong)
			3. get last (unverified) mutable state from SharedObjectStateAccessor
			4. start mutable call in VM
		    5. register resulting mutable state on ledger (and request if a normal call)
		    6. report resulting mutable state to SharedObjectStateAccessor
			7. unlock on limiter for mutable calls
			8. send results back
	*/
	panic("unimplemented")
}

func (s *StateMachineCallRequest) sendReplyOnError(ctx smachine.FailureContext) {
	// TODO send reply
}

func (s *StateMachineCallRequest) sendReplyOnCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO send reply
	return ctx.Stop()
}
