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

package getobject

import (
	"context"
	"errors"

	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
)

func RawGetObjectStateMachinePresentFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: 1,
		States: []common.State{
			{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := Init(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				Migration: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
            		aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		aPayload, ok := element.GetPayload().(*CustomPayload)
            		if !ok { return nil, 0, errors.New("wrong payload type") }
            		ctx := context.TODO()
            		state := MigrateToPresent(ctx, element, aInput, aPayload)
            		return aPayload, state, nil
            	},
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := GetJet(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := GetJetResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchJet(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := FetchJetResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := WaitHotData(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitHotDataResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckIndex(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckIndex(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchIndex(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchIndex(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckJetForState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckJetForState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchJetForState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchJetForState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},
		},
	}
}

func RawGetObjectStateMachinePastFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: 1,
		States: []common.State{
			{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := Init(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := GetJet(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := GetJetResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchJet(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := FetchJetResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := WaitHotData(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitHotDataResponse(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckIndex(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckIndex(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchIndex(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchIndex(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := CheckJetForState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitCheckJetForState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchJetForState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchJetForState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    aPayload, ok := element.GetPayload().(*CustomPayload)
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := FetchState(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					
				},
				AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().(CustomEvent)
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().(*CustomPayload)
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().(CustomAdapterResponsePayload)
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := WaitFetchState(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },
			},{
				
				
			},{
				
				
			},{
				
				
			},
		},
	}
}

func RawGetObjectStateMachineFutureFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: 1,
		States: []common.State{
			{
				
				Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().(CustomEvent)
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    state, payload := InitFuture(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					
				},
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},{
				
				
				
			},
		},
	}
}
