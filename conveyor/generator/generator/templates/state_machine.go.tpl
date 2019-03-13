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
	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"

	"errors"
)

{{range $i, $machine := .StateMachines}}type Base{{$machine.Name}} struct {}

func (*Base{{$machine.Name}}) GetTypeID() statemachine.ID {
    return {{inc $i}}
}
{{range $i, $state := .States}}{{if (gtNull $i)}}
func (*Base{{$machine.Name}}) {{$state.Name}}() statemachine.StateID {
    return {{$i}}
}{{end}}{{end}}

type Raw{{$machine.Name}} struct {
    cleanStateMachine {{$machine.Name}}
}

func Raw{{$machine.Name}}Factory() [3]common.StateMachine {
    m := Raw{{$machine.Name}}{
        cleanStateMachine: &Clean{{$machine.Name}}{},
    }

    var x = [3][]common.State{}
    // future state machine
    x[0] = append(x[0], common.State{
        Transition: m.{{(index .States 0).GetTransitionFutureName}},
        ErrorState: m.{{(index .States 0).GetErrorStateFutureName}},
    },{{range $i, $state := .States}}{{if (gtNull $i)}}
    common.State{
        {{if (handlerExists $state.TransitionFuture)}}Transition: m.{{$state.GetTransitionFutureName}},{{end}}
        {{if (handlerExists $state.AdapterResponseFuture)}}AdapterResponse: m.{{$state.GetAdapterResponseFutureName}},{{end}}
        ErrorState: m.{{$state.GetErrorStateFutureName}},
        {{if (handlerExists $state.AdapterResponseErrorFuture)}}AdapterResponseError: m.{{$state.GetAdapterResponseErrorFutureName}},{{end}}
    },{{end}}{{end}})

    // present state machine
    x[1] = append(x[1], common.State{
        Transition: m.{{(index .States 0).GetTransitionName}},
        ErrorState: m.{{(index .States 0).GetErrorStateName}},
    },{{range $i, $state := .States}}{{if (gtNull $i)}}
    common.State{
        Migration: m.{{$state.GetMigrationFuturePresentName}},
        Transition: m.{{$state.GetTransitionName}},
        {{if (handlerExists $state.AdapterResponse)}}AdapterResponse: m.{{$state.GetAdapterResponseName}},{{end}}
        ErrorState: m.{{$state.GetErrorStateName}},
        {{if (handlerExists $state.AdapterResponseError)}}AdapterResponseError: m.{{$state.GetAdapterResponseErrorName}},{{end}}
    },{{end}}{{end}})

    // past state machine
    x[2] = append(x[2], common.State{
        Transition: m.{{(index .States 0).GetTransitionPastName}},
        ErrorState: m.{{(index .States 0).GetErrorStatePastName}},
    },{{range $i, $state := .States}}{{if (gtNull $i)}}
    common.State{
        Transition: m.{{$state.GetTransitionPastName}},
        {{if (handlerExists $state.AdapterResponsePast)}}AdapterResponse: m.{{$state.GetAdapterResponsePastName}},{{end}}
        ErrorState: m.{{$state.GetErrorStatePastName}},
        {{if (handlerExists $state.AdapterResponseErrorPast)}}AdapterResponseError: m.{{$state.GetAdapterResponseErrorPastName}},{{end}}
    },{{end}}{{end}})

    return [3]common.StateMachine{
        {
            ID: m.cleanStateMachine.({{$machine.Name}}).GetTypeID(),
            States: x[0],
        },
        {
            ID: m.cleanStateMachine.({{$machine.Name}}).GetTypeID(),
            States: x[0],
        },
        {
            ID: m.cleanStateMachine.({{$machine.Name}}).GetTypeID(),
            States: x[0],
        },
    }
}

{{template "initHandler" (params $machine (index .States 0).Transition)}}
{{template "initHandler" (params $machine (index .States 0).TransitionFuture)}}
{{template "initHandler" (params $machine (index .States 0).TransitionPast)}}
{{template "errorStateHandler" (params $machine (index .States 0).ErrorState)}}
{{template "errorStateHandler" (params $machine (index .States 0).ErrorStateFuture)}}
{{template "errorStateHandler" (params $machine (index .States 0).ErrorStatePast)}}

{{range $i, $state := $machine.States}}{{if (gtNull $i)}}
{{template "transitionHandler" (params $machine $state.Transition)}}
{{template "transitionHandler" (params $machine $state.TransitionFuture)}}
{{template "transitionHandler" (params $machine $state.TransitionPast)}}
{{template "migrationHandler" (params $machine $state.Migration)}}
{{template "migrationHandler" (params $machine $state.MigrationFuturePresent)}}
{{template "errorStateHandler" (params $machine $state.ErrorState)}}
{{template "errorStateHandler" (params $machine $state.ErrorStateFuture)}}
{{template "errorStateHandler" (params $machine $state.ErrorStatePast)}}
{{template "adapterResponseHandler" (params $machine $state.AdapterResponse)}}
{{template "adapterResponseHandler" (params $machine $state.AdapterResponseFuture)}}
{{template "adapterResponseHandler" (params $machine $state.AdapterResponsePast)}}
{{template "adapterResponseErrorHandler" (params $machine $state.AdapterResponseError)}}
{{template "adapterResponseErrorHandler" (params $machine $state.AdapterResponseErrorFuture)}}
{{template "adapterResponseErrorHandler" (params $machine $state.AdapterResponseErrorPast)}}
{{end}}{{end}}
{{end}}

{{define "initHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().({{.Machine.InputEventType}})
    if !ok { return nil, 0, errors.New("wrong input event type") }
    payload, state, err := s.cleanStateMachine.{{.Handler.Name}}(aInput, element.GetPayload())
    return payload, state, err
}{{end}}{{end}}
{{define "errorStateHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.{{.Handler.Name}}(element.GetInputEvent(), element.GetPayload(), err)
    return payload, state
}{{end}}{{end}}
{{define "transitionHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().({{.Machine.InputEventType}})
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().({{.Machine.PayloadType}})
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.{{.Handler.Name}}(aInput, aPayload)
    return payload, state, err
}{{end}}{{end}}
{{define "migrationHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().({{.Machine.InputEventType}})
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().({{.Machine.PayloadType}})
    if !ok { return nil, 0, errors.New("wrong payload type") }
    payload, state, err := s.cleanStateMachine.{{.Handler.Name}}(aInput, aPayload)
    return payload, state, err
}{{end}}{{end}}
{{define "adapterResponseHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper, ar adapter.IAdapterResponse) (interface{}, statemachine.ElementState, error) {
    aInput, ok := element.GetInputEvent().({{.Machine.InputEventType}})
    if !ok { return nil, 0, errors.New("wrong input event type") }
    aPayload, ok := element.GetPayload().({{.Machine.PayloadType}})
    if !ok { return nil, 0, errors.New("wrong payload type") }
    aResponse, ok := ar.GetRespPayload().({{index .Handler.Params 2}})
    if !ok { return nil, 0, errors.New("wrong response type") }
    payload, state, err := s.cleanStateMachine.{{.Handler.Name}}(aInput, aPayload, aResponse)
    return payload, state, err
}{{end}}{{end}}
{{define "adapterResponseErrorHandler"}}{{if (handlerExists .Handler)}}func (s *Raw{{.Machine.Name}}) {{.Handler.Name}}(element slot.SlotElementHelper, ar adapter.IAdapterResponse, err error) (interface{}, statemachine.ElementState) {
    payload, state := s.cleanStateMachine.{{.Handler.Name}}(element.GetInputEvent(), element.GetPayload(), ar, err)
    return payload, state
}{{end}}{{end}}
