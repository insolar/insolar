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

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/longbits"
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
	return v.link.PrepareAccess(func(data interface{}) bool {
		fn(data.(*CustomSharedState))
		return false
	})
}

type catalogC struct {
}

func (p *catalogC) Get(ctx smachine.ExecutionContext, key longbits.ByteString) CustomSharedStateAccessor {
	if v, ok := p.TryGet(ctx, key); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", key))
}

func (p *catalogC) TryGet(ctx smachine.ExecutionContext, key longbits.ByteString) (CustomSharedStateAccessor, bool) {

	if v := ctx.GetPublishedLink(key); v.IsAssignableTo((*CustomSharedState)(nil)) {
		return CustomSharedStateAccessor{v}, true
	}
	return CustomSharedStateAccessor{}, false
}

func (p *catalogC) GetOrCreate(ctx smachine.ExecutionContext, key longbits.ByteString) CustomSharedStateAccessor {
	if v, ok := p.TryGet(ctx, key); ok {
		return v
	}

	ctx.InitChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &catalogEntryCSM{sharedState: CustomSharedState{
			key: key,
			//Mutex: smachine.NewExclusiveWithFlags("", 0), //smachine.QueueAllowsPriority),
			Mutex: smachine.NewSemaphoreWithFlags(2, "", smachine.QueueAllowsPriority).SyncLink(),
		}}
	})

	return p.Get(ctx, key)
}

type catalogEntryCSM struct {
	smachine.StateMachineDeclTemplate
	sharedState CustomSharedState
}

func (sm *catalogEntryCSM) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return sm.Init
}

func (sm *catalogEntryCSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *catalogEntryCSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	sdl := ctx.Share(&sm.sharedState, 0)
	if !ctx.Publish(sm.sharedState.key, sdl) {
		return ctx.Stop()
	}
	return ctx.JumpExt(smachine.SlotStep{Transition: sm.State1, Flags: smachine.StepWeak})
}

func (sm *catalogEntryCSM) State1(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Sleep().ThenRepeat()
}
