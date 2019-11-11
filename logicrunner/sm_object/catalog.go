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
	"fmt"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
)

type LocalObjectCatalog struct{}

func (p LocalObjectCatalog) Get(ctx smachine.ExecutionContext, objectReference insolar.Reference) SharedObjectStateAccessor {
	if v, ok := p.TryGet(ctx, objectReference); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", objectReference.String()))
}

func (p LocalObjectCatalog) TryGet(ctx smachine.ExecutionContext, objectReference insolar.Reference) (SharedObjectStateAccessor, bool) {
	if v := ctx.GetPublishedLink(objectReference); v.IsAssignableTo((*SharedObjectState)(nil)) {
		return SharedObjectStateAccessor{v}, true
	}
	return SharedObjectStateAccessor{}, false
}

func (p LocalObjectCatalog) Create(ctx smachine.ExecutionContext, objectReference insolar.Reference) SharedObjectStateAccessor {
	if _, ok := p.TryGet(ctx, objectReference); ok {
		panic(fmt.Sprintf("already exists: %s", objectReference.String()))
	}

	ctx.InitChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		ctx.SetTracerId(fmt.Sprintf("object-%s", objectReference.String()))
		return NewObjectSM(objectReference, false)
	})

	return p.Get(ctx, objectReference)
}

func (p LocalObjectCatalog) GetOrCreate(ctx smachine.ExecutionContext, objectReference insolar.Reference) SharedObjectStateAccessor {
	if v, ok := p.TryGet(ctx, objectReference); ok {
		return v
	}

	ctx.InitChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		ctx.SetTracerId(fmt.Sprintf("object-%s", objectReference.String()))
		return NewObjectSM(objectReference, true)
	})

	return p.Get(ctx, objectReference)
}

// //////////////////////////////////////

type SharedObjectStateAccessor struct {
	smachine.SharedDataLink
}

func (v SharedObjectStateAccessor) Prepare(fn func(*SharedObjectState)) smachine.SharedDataAccessor {
	return v.PrepareAccess(func(data interface{}) bool {
		fn(data.(*SharedObjectState))
		return false
	})
}
