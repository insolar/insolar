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

	"github.com/insolar/insolar/conveyor/adapter/adapterhelper"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/statemachine"
    {{range .Imports}}{{.}}
    {{end}}
)

func Raw{{.Name}}PresentFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: {{.ID}},
		States: []statemachine.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetMigration)}}Migration: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
            		aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong payload type") }
            		ctx := context.TODO()
            		state := {{.GetMigration.Name}}(ctx, element, aInput, aPayload)
            		return aPayload, state, nil
            	},{{end}}
				{{if (handlerExists $state.GetTransition)}}Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransition.Name}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
					// todo here must be real adapter helper
					state := {{.GetTransition.Name}}(ctx, element, aInput, aPayload{{getAdapterHelper $ .GetTransition.GetAdapterHelperType}})
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponse)}}AdapterResponse: func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.({{unPackage .GetAdapterResponse.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponse.Name}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}

func Raw{{.Name}}PastFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: {{.ID}},
		States: []statemachine.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetTransitionPast)}}Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransitionPast.Name}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
                    // todo here must be real adapter helper
					state := {{.GetTransitionPast.Name}}(ctx, element, aInput, aPayload{{getAdapterHelper $ .GetTransitionPast.GetAdapterHelperType}})
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponsePast)}}AdapterResponse: func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.({{unPackage .GetAdapterResponsePast.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponsePast.Name}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}

func Raw{{.Name}}FutureFactory(helpers *adapterhelper.Catalog) *statemachine.StateMachine {
	return &statemachine.StateMachine{
		ID: {{.ID}},
		States: []statemachine.State{
			{{range $i, $state := .States}}{
				{{if (handlerExists $state.GetMigrationFuturePresent)}}Migration: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
            		aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong payload type") }
            		ctx := context.TODO()
            		state := {{.GetMigrationFuturePresent.Name}}(ctx, element, aInput, aPayload)
            		return aPayload, state, nil
            	},{{end}}
				{{if (handlerExists $state.GetTransitionFuture)}}Transition: func(element fsm.SlotElementHelper) (interface{}, fsm.ElementState, error) {
    		        aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
            		if !ok { return nil, 0, errors.New("wrong input event type") }
            		ctx := context.TODO()
				    {{if (isNull $i)}}state, payload := {{.GetTransitionFuture.Name}}(ctx, element, aInput, element.GetPayload())
                    return payload, state, nil
					{{else}}aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
                    if !ok { return nil, 0, errors.New("wrong payload type") }
                    // todo here must be real adapter helper
					helper := CustomAdapterHelper{}
					state := {{.GetTransitionFuture.Name}}(ctx, element, aInput, aPayload{{getAdapterHelper $ .GetTransitionFuture.GetAdapterHelperType}})
                    return aPayload, state, nil
					{{end}}
				},{{end}}
				{{if (handlerExists $state.GetAdapterResponseFuture)}}AdapterResponse: func(element fsm.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error) {
					aInput, ok := element.GetInputEvent().({{unPackage $.InputEventType $.Package}})
					if !ok { return nil, 0, errors.New("wrong input event type") }
					aPayload, ok := element.GetPayload().({{unPackage $.PayloadType $.Package}})
					if !ok { return nil, 0, errors.New("wrong payload type") }
					aResponse, ok := response.({{unPackage .GetAdapterResponseFuture.GetResponseAdapterType $.Package}})
					if !ok { return nil, 0, errors.New("wrong response type") }
					ctx := context.TODO()
					state := {{.GetAdapterResponseFuture.Name}}(ctx, element, aInput, aPayload, aResponse)
					return aPayload, state, nil
                },{{end}}
			},{{end}}
		},
	}
}
