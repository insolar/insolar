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

package init

import (
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

const (
	InitState gen.ElState = iota
)

func Register() {
	gen.AddMachine("Init").
		InitFuture(ParseInputEvent).
		Init(IncorrectAction, IncorrectAction)
}

func ParseInputEvent(helper slot.SlotElementHelper, input interface{}, payload interface{}) (interface{}, fsm.ElementState) {
	return nil, fsm.NewElementState(4, 0)
}

func IncorrectAction() {
	panic("[ IncorrectAction ] We shouldn't be here")
}
