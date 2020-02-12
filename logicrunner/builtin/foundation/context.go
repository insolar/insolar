// Copyright 2020 Insolar Network Ltd.
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

package foundation

import (
	"github.com/tylerb/gls"

	"github.com/insolar/insolar/insolar"
)

const glsCallContextKey = "callCtx"

// GetLogicalContext returns current calling context.
func GetLogicalContext() *insolar.LogicCallContext {
	ctx := gls.Get(glsCallContextKey)
	if ctx == nil {
		panic("object has no context")
	}

	if ctx, ok := ctx.(*insolar.LogicCallContext); ok {
		return ctx
	}

	panic("wrong type of context")
}

// SetLogicalContext saves current calling context
func SetLogicalContext(ctx *insolar.LogicCallContext) {
	gls.Set(glsCallContextKey, ctx)
}

// ClearContext clears underlying gls context
func ClearContext() {
	gls.Cleanup()
}
