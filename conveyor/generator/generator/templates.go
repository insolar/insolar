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

package generator

import (
	"text/template"
	"io"
)

var (
	stateMachineIdsTpl = template.Must(template.New("stateMachineIdsTpl").Parse(`
package sample

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"errors"
)

type SMFID{{.Name}} struct { // STFID = State Machine Flow IDs
}
`))

	stateTpl = template.Must(template.New("stateTpl").Parse(`
func (*SMFID{{.Machine}}) {{.Name}}() common.ElState {
    return {{.Value}}
}
`))

	stateMachineRawTpl = template.Must(template.New("stateMachineRawTpl").Parse(`
type SMRH{{.Name}} struct { // SMRH = State Machine Raw Handlers
	cleanHandlers {{.Name}}
}

func NewSMRH{{.Name}}() SMRH{{.Name}} {
	return SMRH{{.Name}}{
		// cleanHandlers: &{{.Name}}Implementation{},
	}
}
`))

	initTpl = template.Must(template.New("initTpl").Parse(`
func (s *SMRH{{.Machine}}) Init(element common.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().({{.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    payload, state, err := s.cleanHandlers.Init(aInput)
    if err != nil {
        return payload, state, err
    }
    return s.cleanHandlers.{{.FirstHandlerName}}(aInput, payload)
}
`))

	transitMigrateTpl = template.Must(template.New("transitMigrateTpl").Parse(`
func (s *SMRH{{.Machine}}) {{.Name}}(element common.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().({{.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().({{.PayloadType}})
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.{{.Name}}(aInput, aPayload)
}
`))

	errorTpl = template.Must(template.New("errorTpl").Parse(`
func (s *SMRH{{.Machine}}) {{.Name}}(element common.SlotElementHelper, err error) (interface{}, common.ElState) {
    aInput, ok := element.GetInputEvent().({{.EventType}})
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong input event type")
        return nil, 0
    }
    aPayload, ok := element.GetPayload().({{.PayloadType}})
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong payload type")
        return nil, 0
    }
    return s.cleanHandlers.{{.Name}}(aInput, aPayload, err)
}
`))

)

type stateParams struct {
	Machine string
	Name string
	Value int
}

func (g *Generator) GenerateStateMachine(w io.Writer, idx int) {
	stateMachineIdsTpl.Execute(w, g.stateMachines[idx])
	for i, state := range g.stateMachines[idx].States {
		stateTpl.Execute(w, stateParams{
			Machine: g.stateMachines[idx].Name,
			Name: state.name,
			Value: i + 1,
		})
	}
}

type initParams struct {
	Machine string
	EventType string
	FirstHandlerName string
}

type handlerParams struct {
	Machine string
	Name string
	EventType string
	PayloadType string
}

func (g *Generator) GenerateRawHandlers(w io.Writer, idx int) {
	stateMachineRawTpl.Execute(w, g.stateMachines[idx])
	initTpl.Execute(w, initParams{
		Machine: g.stateMachines[idx].Name,
		EventType: g.stateMachines[idx].Init.eventType,
		FirstHandlerName: g.stateMachines[idx].States[0].transit.name,
	})
	for _, state := range g.stateMachines[idx].States {
		transitMigrateTpl.Execute(w, handlerParams{
			Machine: g.stateMachines[idx].Name,
			Name: state.transit.name,
			EventType: g.stateMachines[idx].Init.eventType,
			PayloadType: state.transit.params[1],
		})
		transitMigrateTpl.Execute(w, handlerParams{
			Machine: g.stateMachines[idx].Name,
			Name: state.migrate.name,
			EventType: g.stateMachines[idx].Init.eventType,
			PayloadType: state.migrate.params[1],
		})
		errorTpl.Execute(w, handlerParams{
			Machine: g.stateMachines[idx].Name,
			Name: state.error.name,
			EventType: g.stateMachines[idx].Init.eventType,
			PayloadType: state.error.params[1],
		})
	}
}
