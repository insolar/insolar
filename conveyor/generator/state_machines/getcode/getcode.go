/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package getcode

import (
	"context"

	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/insolar"
)

const (
	InitState fsm.ElementState = iota
	GetingCode
)

type GetCodePayload struct {
	reply insolar.Reply
	err   error
}

// Register adds init state machine to generator
func Register(g *generator.Generator) {
	g.AddMachine("GetCode").
		InitFuture(ParseInputEvent).
		Init(ParseInputEvent)
}

// ParseInputEvent parse input ivent
func ParseInputEvent(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload interface{}) (fsm.ElementState, interface{}) {
	return GetingCode, payload
}
