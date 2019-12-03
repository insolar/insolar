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

package sm_execute_request

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
)

type RequestCatalog struct{}

func (p RequestCatalog) Get(ctx smachine.ExecutionContext, requestReference insolar.Reference) SharedRequestStateAccessor {
	if v, ok := p.TryGet(ctx, requestReference); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", requestReference.String()))
}

func (p RequestCatalog) TryGet(ctx smachine.ExecutionContext, requestReference insolar.Reference) (SharedRequestStateAccessor, bool) {
	if v := ctx.GetPublishedLink(requestReference); v.IsAssignableTo((*SharedRequestState)(nil)) {
		return SharedRequestStateAccessor{v}, true
	}
	return SharedRequestStateAccessor{}, false
}

// //////////////////////////////////////

type SharedRequestStateAccessor struct {
	smachine.SharedDataLink
}

func (v SharedRequestStateAccessor) Prepare(fn func(*SharedRequestState)) smachine.SharedDataAccessor {
	return v.PrepareAccess(func(data interface{}) bool {
		fn(data.(*SharedRequestState))
		return false
	})
}
