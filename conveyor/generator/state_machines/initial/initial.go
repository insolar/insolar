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

package initial

import (
	"context"

	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const (
	InitState fsm.ElementState = iota
)

func Register(g *generator.Generator) {
	g.AddMachine("Initial").
		InitFuture(ParseInputEvent).
		Init(ParseInputEvent)
}

func ParseInputEvent(ctx context.Context, helper fsm.SlotElementHelper, input interface{}, payload interface{}) (fsm.ElementState, interface{}) {
	parcel, ok := helper.GetInputEvent().(insolar.Parcel)
	if !ok {
		inslogger.FromContext(ctx).Warnf("[ ParseInputEvent ] InputEvent must be insolar.Parcel. Actual: %+v", helper.GetInputEvent())
		return 0, nil
	}
	switch parcel.Type() {
	case insolar.TypeGetCode:
		return fsm.NewElementState(0, 0), nil
	default:
		inslogger.FromContext(ctx).Warnf("[ ParseInputEvent ] Unknown parcel type: %s", parcel.Type().String())
		return 0, nil
	}
	return 0, nil
}
