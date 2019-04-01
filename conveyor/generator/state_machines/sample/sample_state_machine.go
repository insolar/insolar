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

package sample

import (
	"context"

	a "github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/insolar"
)

// custom types
type CustomEvent struct{}
type CustomPayload struct{}

type TestResult struct{}

func (e *TestResult) Type() insolar.ReplyType {
	return insolar.ReplyType(42)
}

const (
	InitState fsm.ElementState = iota
	StateFirst
	StateSecond
	StateThird
)

func Register(g *generator.Generator) {
	g.AddMachine("SampleStateMachine").
		InitFuture(initFutureHandler).
		Init(initPresentHandler, StateFirst).
		Transition(StateFirst, transitPresentFirst, StateSecond).
		AdapterResponse(StateFirst, responseAdapterHelper, StateThird).
		Transition(StateThird, transitPresentThird, 0)
}

func initPresentHandler(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload interface{}) (fsm.ElementState, *CustomPayload) {
	return StateFirst, nil
}

func initFutureHandler(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload interface{}) (fsm.ElementState, *CustomPayload) {
	panic("implement me")
}

func transitPresentFirst(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, adapterHelper a.SendResponseHelper) fsm.ElementState {
	helper.DeactivateTill(fsm.Response)
	err := adapterHelper.SendResponse(helper, &TestResult{}, uint32(StateFirst))
	if err != nil {
		panic(err)
	}
	return StateSecond
}

func responseAdapterHelper(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload, respPayload *TestResult) fsm.ElementState {
	return StateThird
}

func transitPresentThird(ctx context.Context, helper fsm.SlotElementHelper, input CustomEvent, payload *CustomPayload) fsm.ElementState {
	return 0
}
