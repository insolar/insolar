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

func NewVMObjectSM(objKey longbits.ByteString) *vmObjectSM {
	return &vmObjectSM{SharedObjectState: SharedObjectState{ObjectInfo: ObjectInfo{ObjKey: objKey}}}
}

type vmObjectSM struct {
	smachine.StateMachineDeclTemplate

	SharedObjectState
	readyToWorkCtl smachine.BoolConditionalLink
}

type ObjectInfo struct {
	ObjKey        longbits.ByteString
	IsReadyToWork bool

	ArtifactClient *ArtifactClientServiceAdapter
	ContractRunner *ContractRunnerServiceAdapter

	ObjectLatestValidState ArtifactBinary
	ObjectLatestValidCode  ArtifactBinary

	ImmutableExecute smachine.SyncLink
	MutableExecute   smachine.SyncLink
}

type SharedObjectState struct {
	SemaReadyToWork smachine.SyncLink
	ObjectInfo
}

//////////////////////////

func (sm *vmObjectSM) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&sm.ArtifactClient)
	injector.MustInject(&sm.ContractRunner)
}

func (sm *vmObjectSM) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return sm.Init
}

func (sm *vmObjectSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *vmObjectSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.migrateStop)

	sm.readyToWorkCtl = smachine.NewConditionalBool(false, "readyToWork")
	sm.SemaReadyToWork = sm.readyToWorkCtl.SyncLink()
	sm.ImmutableExecute = smachine.NewFixedSemaphore(5, "immutable calls")
	sm.MutableExecute = smachine.NewFixedSemaphore(1, "mutable calls") // TODO here we need an ORDERED queue

	sdl := ctx.Share(&sm.SharedObjectState, 0)
	if !ctx.Publish(sm.ObjKey, sdl) {
		return ctx.Stop()
	}
	return ctx.Jump(sm.stateGetLatestValidatedState)
}

func (sm *vmObjectSM) stateGetLatestValidatedState(ctx smachine.ExecutionContext) smachine.StateUpdate {
	sm.ArtifactClient.PrepareAsync(ctx, func(svc ArtifactClientService) smachine.AsyncResultFunc {
		stateObj, codeObj := svc.GetLatestValidatedStateAndCode()

		return func(ctx smachine.AsyncResultContext) {
			sm.ObjectLatestValidState = stateObj
			sm.ObjectLatestValidCode = codeObj
			sm.IsReadyToWork = true
			ctx.WakeUp()
		}
	})

	return ctx.Sleep().ThenJump(sm.stateGotLatestValidatedState)
}

func (sm *vmObjectSM) stateGotLatestValidatedState(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if sm.ObjectLatestValidState == nil {
		return ctx.Sleep().ThenRepeat()
	}
	ctx.ApplyAdjustment(sm.readyToWorkCtl.NewValue(true))

	return ctx.JumpExt(smachine.SlotStep{Transition: sm.waitForMigration, Migration: sm.migrateSendState})
}

func (sm *vmObjectSM) waitForMigration(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Sleep().ThenRepeat()
}

func (sm *vmObjectSM) migrateStop(ctx smachine.MigrationContext) smachine.StateUpdate {
	return ctx.Stop()
}

func (sm *vmObjectSM) migrateSendState(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)
	return ctx.Jump(sm.stateCompleteExecution)
}

func (sm *vmObjectSM) stateCompleteExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*  TODO
	Steps:
	1. Send last state to next executor
	2. Send transcript to validators
	*/
	return ctx.Stop()
}
