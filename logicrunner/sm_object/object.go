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

package sm_object

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_sender"
)

type ObjectInfo struct {
	ObjectReference insolar.Reference
	IsReadyToWork   bool

	artifactClient *s_artifact.ArtifactClientServiceAdapter
	sender         *s_sender.SenderServiceAdapter
	pulseSlot      *conveyor.PulseSlot
	externalError  error

	ObjectLatestDescriptor artifacts.ObjectDescriptor

	ImmutableExecute smachine.SyncLink
	MutableExecute   smachine.SyncLink

	PreviousExecutorState payload.PreviousExecutorState
}

type SharedObjectState struct {
	SemaphoreReadyToWork              smachine.SyncLink
	SemaphorePreviousExecutorFinished smachine.SyncLink
	SemaphorePreviousResultSaved      smachine.SyncLink
	ObjectInfo
}

func NewObjectSM(objectReference insolar.Reference, exists bool) *ObjectSM {
	return &ObjectSM{
		SharedObjectState: SharedObjectState{
			ObjectInfo: ObjectInfo{ObjectReference: objectReference},
		},
		oldObject: exists,
	}
}

type ObjectSM struct {
	smachine.StateMachineDeclTemplate

	SharedObjectState
	readyToWorkCtl           smachine.BoolConditionalLink
	previousExecutorFinished smachine.BoolConditionalLink
	previousResultSaved      smachine.BoolConditionalLink
	oldObject                bool
}

/* -------- Declaration ------------- */

func (sm *ObjectSM) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&sm.artifactClient)
	injector.MustInject(&sm.sender)
	injector.MustInject(&sm.pulseSlot)
}

func (sm *ObjectSM) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return sm.Init
}

/* -------- Instance ------------- */

func (sm *ObjectSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *ObjectSM) sendPayloadToVirtual(ctx smachine.ExecutionContext, pl payload.Payload) {
	goCtx := ctx.GetContext()

	resultsMessage, err := payload.NewMessage(pl)
	if err == nil {
		objectReference := sm.ObjectReference

		sm.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
			_, done := svc.SendRole(goCtx, resultsMessage, insolar.DynamicRoleVirtualExecutor, objectReference)
			done()
		}).DelayedSend()
	} else {
		logger := inslogger.FromContext(goCtx)
		logger.Error("Failed to serialize message: ", err.Error())
	}
}

func (sm *ObjectSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.migrateSendStateBeforeExecution)

	sm.readyToWorkCtl = smachine.NewConditionalBool(false, "readyToWork")
	sm.SemaphoreReadyToWork = sm.readyToWorkCtl.SyncLink()

	sm.previousExecutorFinished = smachine.NewConditionalBool(false, "previousExecutorFinished")
	sm.SemaphorePreviousExecutorFinished = sm.readyToWorkCtl.SyncLink()

	sm.previousResultSaved = smachine.NewConditionalBool(false, "previousResultSaved")
	sm.SemaphorePreviousResultSaved = sm.previousResultSaved.SyncLink()

	sm.ImmutableExecute = smachine.NewFixedSemaphore(5, "immutable calls")
	sm.MutableExecute = smachine.NewFixedSemaphore(1, "mutable calls") // TODO here we need an ORDERED queue

	sdl := ctx.Share(&sm.SharedObjectState, 0)
	if !ctx.Publish(sm.ObjectReference, sdl) {
		return ctx.Stop()
	}
	return ctx.Jump(sm.stepCheckPreviousExecutor)
}

func (sm *ObjectSM) stepCheckPreviousExecutor(ctx smachine.ExecutionContext) smachine.StateUpdate {
	switch sm.PreviousExecutorState {
	case payload.PreviousExecutorUnknown:
		return ctx.Jump(sm.stepGetPendingsInformation)
	case payload.PreviousExecutorProbablyExecutes, payload.PreviousExecutorExecutes:
		// we should wait here till PendingFinished/ExecutorResults came, retry and then change state to PreviousExecutorFinished
		if ctx.AcquireForThisStep(sm.SemaphorePreviousExecutorFinished).IsNotPassed() {
			return ctx.Sleep().ThenRepeat()
		}

		// we shouldn't be here
		// if we came to that place - means MutableRequestsAreReady, but PreviousExecutor still executes)
		panic("unreachable")
	case payload.PreviousExecutorFinished:
		return ctx.Jump(sm.stepGetLatestValidatedState)
	default:
		panic("unreachable")
	}
}

func (sm *ObjectSM) stepGetPendingsInformation(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()

	objectReference := sm.ObjectReference

	return sm.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		logger := inslogger.FromContext(goCtx)

		hasAbandonedRequests, err := svc.HasPendings(goCtx, objectReference)
		if err != nil {
			logger.Error("couldn't check pending state: ", err.Error())
		}

		var newState payload.PreviousExecutorState
		if hasAbandonedRequests {
			logger.Debug("ledger has requests older than one pulse")
			newState = payload.PreviousExecutorProbablyExecutes
		} else {
			logger.Debug("no requests on ledger older than one pulse")
			newState = payload.PreviousExecutorFinished
		}

		return func(ctx smachine.AsyncResultContext) {
			if sm.PreviousExecutorState == payload.PreviousExecutorUnknown {
				sm.PreviousExecutorState = newState
			} else {
				logger.Info("state already changed, ignoring check")
			}
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(sm.stepCheckPreviousExecutor)
}

// we should check here only if not creation request here
func (sm *ObjectSM) stepGetLatestValidatedState(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.migrateStop)

	goCtx := ctx.GetContext()
	objectReference := sm.ObjectReference

	if !sm.oldObject {
		sm.oldObject = true
		sm.IsReadyToWork = true
		return ctx.Jump(sm.stateGotLatestValidatedStatePrototypeAndCode)
	}

	return sm.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		var err error

		failCallback := func(ctx smachine.AsyncResultContext) {
			inslogger.FromContext(goCtx).Error("Failed to obtain objects: ", err)
			sm.externalError = err
		}

		objectDescriptor, err := svc.GetObject(goCtx, objectReference, nil)
		if err != nil {
			err = errors.Wrap(err, "Failed to obtain object descriptor")
			return failCallback
		}

		return func(ctx smachine.AsyncResultContext) {
			sm.ObjectLatestDescriptor = objectDescriptor
			sm.IsReadyToWork = true
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(sm.stateGotLatestValidatedStatePrototypeAndCode)
}

func (sm *ObjectSM) stateGotLatestValidatedStatePrototypeAndCode(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if sm.externalError != nil {
		ctx.Error(sm.externalError)
	} else if sm.IsReadyToWork != true {
		return ctx.Sleep().ThenJump(sm.stateGotLatestValidatedStatePrototypeAndCode)
	}

	ctx.ApplyAdjustment(sm.readyToWorkCtl.NewValue(true))

	return ctx.JumpExt(smachine.SlotStep{
		Transition: sm.waitForMigration,
		Migration:  sm.migrateSendStateAfterExecution,
	})
}

func (sm *ObjectSM) waitForMigration(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Sleep().ThenRepeat()
}

// //////////////////////////////////////

func (sm *ObjectSM) migrateSendStateBeforeExecution(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(sm.stateSendStateBeforeExecution)
}

func (sm *ObjectSM) stateSendStateBeforeExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	_, immutableLeft := sm.ImmutableExecute.GetCounts()
	_, mutableLeft := sm.MutableExecute.GetCounts()

	ledgerHasMoreRequests := immutableLeft+mutableLeft > 0

	var newState payload.PreviousExecutorState

	switch sm.PreviousExecutorState {
	case payload.PreviousExecutorUnknown:
		newState = payload.PreviousExecutorFinished
	case payload.PreviousExecutorProbablyExecutes:
		newState = payload.PreviousExecutorFinished
	case payload.PreviousExecutorExecutes:
		newState = payload.PreviousExecutorUnknown
	case payload.PreviousExecutorFinished:
		newState = payload.PreviousExecutorFinished
	default:
		panic("unreachable")
	}

	sm.sendPayloadToVirtual(ctx, &payload.ExecutorResults{
		ObjectReference:       sm.ObjectReference,
		LedgerHasMoreRequests: ledgerHasMoreRequests,
		State:                 newState,
	})

	return ctx.Stop()
}

// //////////////////////////////////////

func (sm *ObjectSM) migrateSendStateAfterExecution(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(sm.stateSendStateAfterExecution)
}

func (sm *ObjectSM) stateSendStateAfterExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	immutableInProgress, immutableLeft := sm.ImmutableExecute.GetCounts()
	mutableInProgress, mutableLeft := sm.MutableExecute.GetCounts()

	ledgerHasMoreRequests := immutableLeft+mutableLeft > 0
	pendingCount := uint32(immutableInProgress + mutableInProgress)

	var newState payload.PreviousExecutorState

	switch sm.PreviousExecutorState {
	case payload.PreviousExecutorFinished:
		if pendingCount > 0 {
			newState = payload.PreviousExecutorExecutes
		} else {
			newState = payload.PreviousExecutorFinished
		}
	default:
		panic("unreachable")
	}

	if pendingCount > 0 || ledgerHasMoreRequests {
		sm.sendPayloadToVirtual(ctx, &payload.ExecutorResults{
			ObjectReference:       sm.ObjectReference,
			LedgerHasMoreRequests: ledgerHasMoreRequests,
			State:                 newState,
		})
	}

	return ctx.Jump(sm.stateWaitFinishExecutionAfterMigration)
}

func (sm *ObjectSM) stateWaitFinishExecutionAfterMigration(ctx smachine.ExecutionContext) smachine.StateUpdate {
	mc, _ := sm.MutableExecute.GetCounts()
	ic, _ := sm.ImmutableExecute.GetCounts()
	if mc > 0 || ic > 0 {
		return ctx.Poll().ThenRepeat()
	}

	sm.sendPayloadToVirtual(ctx, &payload.PendingFinished{
		ObjectRef: sm.ObjectReference,
	})

	return ctx.Stop()
}

func (sm *ObjectSM) migrateStop(ctx smachine.MigrationContext) smachine.StateUpdate {
	return ctx.Stop()
}
