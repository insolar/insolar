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
	stateMachineTpl = template.Must(template.New("stateMachineTpl").Parse(`
package {{.Package}}

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"errors"
)

{{range $i, $machine := .StateMachines}}type SMFID{{$machine.Name}} struct {
}
{{range $i, $state := .States}}
func (*SMFID{{$machine.Name}}) {{$state.Name}}() common.ElState {
    return {{$i}}
}{{end}}

type SMRH{{$machine.Name}} struct {
	cleanHandlers {{$machine.Name}}
}

func SMRH{{$machine.Name}}Export() []common.State {
    m := SMRH{{$machine.Name}}{
        cleanHandlers: &{{$machine.Name}}Implementation{},
    }
    var x []common.State
    return append(x,{{range $i, $state := .States}}
        common.State{
	        Transit: m.{{$state.Transit.Name}},
	        Migrate: m.{{$state.Migrate.Name}},
	        Error: m.{{$state.Error.Name}},
        },
    {{end}})
}

func (s *SMRH{{$machine.Name}}) Init(element common.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    payload, state, err := s.cleanHandlers.Init(aInput)
    if err != nil {
        return payload, state, err
    }
    return s.cleanHandlers.{{(index $machine.States 0).Transit.Name}}(aInput, payload)
}
{{range $i, $state := $machine.States}}
func (s *SMRH{{$machine.Name}}) {{$state.Transit.Name}}(element common.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().({{index $state.Transit.Params 1}})
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.{{$state.Transit.Name}}(aInput, aPayload)
}

func (s *SMRH{{$machine.Name}}) {{$state.Migrate.Name}}(element common.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().({{index $state.Migrate.Params 1}})
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.{{$state.Migrate.Name}}(aInput, aPayload)
}

func (s *SMRH{{$machine.Name}}) {{$state.Error.Name}}(element common.SlotElementHelper, err error) (interface{}, common.ElState) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong input event type")
        return nil, 0
    }
    aPayload, ok := element.GetPayload().({{index $state.Error.Params 1}})
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong payload type")
        return nil, 0
    }
    return s.cleanHandlers.{{$state.Error.Name}}(aInput, aPayload, err)
}
{{end}}
{{end}}
`))
)

func (p *Parser) Generate(w io.Writer) {
	stateMachineTpl.Execute(w, p)
}
