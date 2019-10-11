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
	"fmt"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/longbits"
	"sync"
)

func CreateCatalogC() CatalogC {
	return &catalogC{}
}

type CatalogC = *catalogC

type CustomSharedState struct {
	key     longbits.ByteString
	Mutex   smachine.SyncLink
	Text    string
	Counter int
}

func (p *CustomSharedState) GetKey() longbits.ByteString {
	return p.key
}

type CustomSharedStateAccessor struct {
	link smachine.SharedDataLink
}

func (v CustomSharedStateAccessor) Prepare(fn func(*CustomSharedState)) smachine.SharedDataAccessor {
	return v.link.PrepareAccess(func(data interface{}) {
		fn(data.(*CustomSharedState))
	})
}

type catalogC struct {
	mutex   sync.RWMutex
	entries map[longbits.ByteString]smachine.SharedDataLink
}

func (p *catalogC) Get(key longbits.ByteString) CustomSharedStateAccessor {
	if v, ok := p.TryGet(key); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", key))
}

func (p *catalogC) TryGet(key longbits.ByteString) (CustomSharedStateAccessor, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p._get(key)
}

func (p *catalogC) _get(key longbits.ByteString) (CustomSharedStateAccessor, bool) {
	if v, ok := p.entries[key]; ok && !v.IsZero() {
		return CustomSharedStateAccessor{v}, true
	}
	return CustomSharedStateAccessor{}, false
}

func (p *catalogC) tryPut(key longbits.ByteString, link smachine.SharedDataLink) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if v, ok := p.entries[key]; ok && !v.IsZero() {
		return false
	}
	if p.entries == nil {
		p.entries = make(map[longbits.ByteString]smachine.SharedDataLink)
	}
	p.entries[key] = link
	return true
}

func (p *catalogC) GetOrCreate(ctx smachine.ExecutionContext, key longbits.ByteString) CustomSharedStateAccessor {
	if v, ok := p.TryGet(key); ok {
		return v
	}

	ctx.InitChild(ctx.GetContext(), func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &catalogCSM{catalog: p, sharedState: CustomSharedState{key: key}}
	})

	return p.Get(key)
}

func (p *catalogC) cleanup() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.entries == nil {
		return true
	}
	for k, v := range p.entries {
		if !v.IsValid() {
			delete(p.entries, k)
		}
	}
	return len(p.entries) == 0
}

type catalogCSM struct {
	smachine.StateMachineDeclTemplate
	catalog     *catalogC
	sharedState CustomSharedState
}

func (sm *catalogCSM) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return sm.Init
}

func (sm *catalogCSM) InjectDependencies(smachine.StateMachine, smachine.SlotLink, smachine.DependencyRegistry) bool {
	return true
}

func (sm *catalogCSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *catalogCSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.Migrate)

	sm.sharedState.Mutex = smachine.NewExclusive()

	sdl := ctx.Share(&sm.sharedState, false)

	if sm.catalog == nil || !sm.catalog.tryPut(sm.sharedState.key, sdl) {
		return ctx.Stop()
	}
	return ctx.Jump(sm.State1)
}

func (sm *catalogCSM) State1(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Sleep().ThenRepeat()
}

func (sm *catalogCSM) Migrate(ctx smachine.MigrationContext) smachine.StateUpdate {
	return ctx.Jump(sm.Cleanup)
}

func (sm *catalogCSM) Cleanup(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if sm.catalog.cleanup() {
		return ctx.Stop()
	}
	return ctx.Jump(sm.State1)
}
