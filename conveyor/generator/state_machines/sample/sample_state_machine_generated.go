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
	"errors"

	"github.com/insolar/insolar/conveyor/adapter/adapterhelper"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/statemachine"
    "github.com/insolar/insolar/conveyor/generator/state_machines/sample/custom"
    
)

func RawSampleStateMachinePresentFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 2,
		States: []statemachine.State{
			{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := initPresentHandler(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*custom.Payload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo here must be real adapter helper
					state := transitPresentFirst(ctx, element, aInput, aPayload, helpers.SendResponseHelper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(custom.Event)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*custom.Payload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.(*TestResult)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := responseAdapterHelper(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*custom.Payload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo here must be real adapter helper
					state := transitPresentThird(ctx, element, aInput, aPayload)
                    return aPayload, state, nil
					
				},
				
			},
		},
	}
}

func RawSampleStateMachinePastFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 2,
		States: []statemachine.State{
			{
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := initPresentHandler(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*custom.Payload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
                    // todo here must be real adapter helper
					state := transitPresentFirst(ctx, element, aInput, aPayload, helpers.SendResponseHelper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(custom.Event)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*custom.Payload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.(*TestResult)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := responseAdapterHelper(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*custom.Payload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
                    // todo here must be real adapter helper
					state := transitPresentThird(ctx, element, aInput, aPayload)
                    return aPayload, state, nil
					
				},
				
			},
		},
	}
}

func RawSampleStateMachineFutureFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: 2,
		States: []statemachine.State{
			{
				
				Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(custom.Event)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := initFutureHandler(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},
		},
	}
}
