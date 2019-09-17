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

import "context"

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	s *Slot
}

func (p *constructionContext) SlotLink() SlotLink {
	return p.s.NewLink()
}

func (p *constructionContext) GetContext() context.Context {
	return p.s.ctx
}

func (p *constructionContext) SetContext(ctx context.Context) {
	if ctx == nil {
		panic("illegal value")
	}
	p.s.ctx = ctx
}

//func (p *constructionContext) GetContainer() SlotMachineState {
//	return p.machine.containerState
//}

func (p *constructionContext) ParentLink() SlotLink {
	return p.s.parent
}

func (p *constructionContext) SetParent(parent SlotLink) {
	p.s.parent = parent
}

func (p *constructionContext) executeCreate(nextCreate CreateFunc) StateMachine {
	p.setMode(updCtxConstruction)
	defer p.setDiscarded()

	return nextCreate(p)
}

/* ========================================================================= */

var _ InitializationContext = &initializationContext{}

type initializationContext struct {
	slotContext
}

func (p *initializationContext) executeInitialization(fn InitFunc) (stateUpdate StateUpdate) {
	p.setMode(updCtxInit)
	defer func() {
		p.discardAndUpdate("initialization", recover(), &stateUpdate)
	}()

	return p.ensureAndPrepare(p.s, fn(p))
}

/* ========================================================================= */

var _ MigrationContext = &migrationContext{}

type migrationContext struct {
	slotContext
}

func (p *migrationContext) executeMigration(fn MigrateFunc) (stateUpdate StateUpdate) {
	p.setMode(updCtxMigrate)
	defer func() {
		p.discardAndUpdate("migration", recover(), &stateUpdate)
	}()

	return p.ensureAndPrepare(p.s, fn(p))
}
