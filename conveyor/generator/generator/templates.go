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
	"io"
	"text/template"
)

var (
	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"notNil": func(x interface{}) bool {
			return x == nil
		},
	}
	stateMachineTpl = template.Must(template.New("stateMachineTpl").Funcs(funcMap).Parse(`
package {{.Package}}

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"errors"
)

{{range $i, $machine := .StateMachines}}type SMFID{{$machine.Name}} struct {
}
func (*SMFID{{$machine.Name}}) TID() common.ElType {
    return {{$i}}
}
{{range $i, $state := .States}}
func (*SMFID{{$machine.Name}}) {{$state.Name}}() common.ElState {
    return {{inc $i}}
}{{end}}

type SMRH{{$machine.Name}} struct {
	cleanHandlers {{$machine.Name}}
}

func SMRH{{$machine.Name}}Export() common.StateMachine {
    m := SMRH{{$machine.Name}}{
        cleanHandlers: &{{$machine.Name}}Implementation{},
    }
    var x []common.State
    x = append(x, common.State{
	        Transit: m.Init,
	        Error: m.{{$machine.InitError.Name}},
        }, {{range $i, $state := .States}}
        common.State{
	        Transit: m.{{$state.Transit.Name}},
	        Migrate: m.{{$state.Migrate.Name}},
	        Error: m.{{$state.Error.Name}},
        },
    {{end}})
    return common.StateMachine{
        Id: int(m.cleanHandlers.({{$machine.Name}}).TID()),
        States: x,
    }
}

func (s *SMRH{{$machine.Name}}) Init(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    payload, state, err := s.cleanHandlers.Init(aInput)
    if err != nil {
        return payload, state.ToInt(), err
    }
    return payload, state.ToInt(), err
}
func (s *SMRH{{$machine.Name}}) {{$machine.InitError.Name}}(element slot.SlotElementHelper, err error) (interface{}, uint32) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0
    }
    payload, state := s.cleanHandlers.{{$machine.InitError.Name}}(aInput, err)
    return payload, state.ToInt()
}
{{range $i, $state := $machine.States}}
func (s *SMRH{{$machine.Name}}) {{$state.Transit.Name}}(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().({{index $state.Transit.Params 1}})
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    payload, state, err := s.cleanHandlers.{{$state.Transit.Name}}(aInput, aPayload)
    return payload, state.ToInt(), err
}

func (s *SMRH{{$machine.Name}}) {{$state.Migrate.Name}}(element slot.SlotElementHelper) (interface{}, uint32, error) {
    aInput, ok := element.GetInputEvent().({{$machine.Init.EventType}})
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().({{index $state.Migrate.Params 1}})
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    payload, state, err := s.cleanHandlers.{{$state.Migrate.Name}}(aInput, aPayload)
    return payload, state.ToInt(), err
}

func (s *SMRH{{$machine.Name}}) {{$state.Error.Name}}(element slot.SlotElementHelper, err error) (interface{}, uint32) {
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
    payload, state := s.cleanHandlers.{{$state.Error.Name}}(aInput, aPayload, err)
    return payload, state.ToInt()
}
{{end}}
{{end}}
`))
)

func (p *Parser) Generate(w io.Writer) {
	stateMachineTpl.Execute(w, p)
}
