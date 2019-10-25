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

package sm_request

import (
	"fmt"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/pulse"
)

func DefaultHandlersFactory(_ pulse.Number, input conveyor.InputEvent) smachine.CreateFunc {
	switch inputConverted := input.(type) {
	case *payload.Meta:
		return metaHandlerFactory(inputConverted)
	default:
		panic(fmt.Sprintf("unknoen event type, got %T", input))
	}
}

func metaHandlerFactory(messageMeta *payload.Meta) smachine.CreateFunc {
	payloadBytes := messageMeta.Payload
	payloadType, err := payload.UnmarshalType(payloadBytes)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal payload type: %s", err.Error()))
	}

	switch payloadType {
	case payload.TypeCallMethod:
		pl := payload.CallMethod{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.CallMethod: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineCallMethod{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypeSagaCallAcceptNotification:
		pl := payload.SagaCallAcceptNotification{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.SagaCallAcceptNotification: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineSagaAccept{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypeUpdateJet:
		pl := payload.UpdateJet{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.UpdateJet: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineUpdateJet{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypePendingFinished:
		pl := payload.PendingFinished{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.PendingFinished: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachinePendingFinished{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypeExecutorResults:
		pl := payload.ExecutorResults{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.ExecutorResults: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineExecutorResults{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypeStillExecuting:
		pl := payload.StillExecuting{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.StillExecuting: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineStillExecuting{Meta: messageMeta, Payload: &pl}
		}
	case payload.TypeAdditionalCallFromPreviousExecutor:
		pl := payload.AdditionalCallFromPreviousExecutor{}
		if err := pl.Unmarshal(payloadBytes); err != nil {
			panic(fmt.Sprintf("failed to unmarshal payload.AdditionalCallFromPreviousExecutor: %s", err.Error()))
		}
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			return &StateMachineAdditionalCall{Meta: messageMeta, Payload: &pl}
		}
	default:
		panic(fmt.Sprintf(" no handler for message type %s", payloadType.String()))
	}
}
