///
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
///

package smachine

var _ StateMachine = &sharedStateMachine{}

type sharedStateMachine struct {
	state      SharedState
	sharedLink SharedDataLink
}

func (p *sharedStateMachine) GetStateMachineDeclaration() StateMachineDeclaration {
	return sharedStateMachineD
}

func (p *sharedStateMachine) initState(ctx InitializationContext) StateUpdate {
	ctx.SetDefaultFlags(StepWeak)
	ctx.SetDefaultMigration(p.migrateState)
	p.sharedLink = ctx.Share(p.state, true)
	return ctx.Jump(p.mainState)
}

func (p *sharedStateMachine) mainState(ctx ExecutionContext) StateUpdate {
	if p.state.CanBeDisposed(ctx) {
		return ctx.Jump(p.stopState)
	}
	return ctx.Sleep().ThenRepeat()
}

func (p *sharedStateMachine) migrateState(ctx MigrationContext) StateUpdate {
	if p.state.CanBeMigrated(ctx) {
		return ctx.Stay()
	}
	return ctx.Jump(p.stopState)
}

func (p *sharedStateMachine) stopState(ctx ExecutionContext) StateUpdate {
	//p.state.OnDisposed()
	return ctx.Stop()
}

var sharedStateMachineD StateMachineDeclaration = &sharedStateMachineDecl{}

type sharedStateMachineDecl struct {
}

func (*sharedStateMachineDecl) IsConsecutive(cur, next StateFunc) bool {
	return true
}

func (*sharedStateMachineDecl) GetInitStateFor(sm StateMachine) InitFunc {
	return sm.(*sharedStateMachine).initState
}

func (*sharedStateMachineDecl) GetMigrateFn(StateFunc) MigrateFunc {
	return nil
}
