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

package {{.Package}}

import (
	"context"
	"errors"

	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
)

func Raw{{.Name}}PresentFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: {{.ID}},
		States: []common.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetMigration)}}Migration: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
            		aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong payload type") }
            		ctx := context.TODO()
            		state := {{.GetMigration.GetName}}(ctx, element, aInput, aPayload)
            		return aPayload, state, nil
            	},{{end}}
				{{if (handlerExists $state.GetTransition)}}Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransition.GetName}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := {{.GetTransition.GetName}}(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponse)}}AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().({{unPackage .GetAdapterResponse.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponse.GetName}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}

func Raw{{.Name}}PastFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: {{.ID}},
		States: []common.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetTransitionPast)}}Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransitionPast.GetName}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := {{.GetTransitionPast.GetName}}(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponsePast)}}AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().({{unPackage .GetAdapterResponsePast.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponsePast.GetName}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}

func Raw{{.Name}}FutureFactory() statemachine.StateMachine {
	return &common.StateMachine{
		ID: {{.ID}},
		States: []common.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetMigrationFuturePresent)}}Migration: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
            		aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong payload type") }
            		ctx := context.TODO()
            		state := {{.GetMigrationFuturePresent.GetName}}(ctx, element, aInput, aPayload)
            		return aPayload, state, nil
            	},{{end}}
				{{if (handlerExists $state.GetTransitionFuture)}}Transition: func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransitionFuture.GetName}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo change to real adapter helper
					helper := CustomAdapterHelper{}
					state := {{.GetTransitionFuture.GetName}}(ctx, element, aInput, aPayload, helper)
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponseFuture)}}AdapterResponse: func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.GetRespPayload().({{unPackage .GetAdapterResponseFuture.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponseFuture.GetName}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}
