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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/insolar/insolar/conveyor/adapter/adapterhelper"
)

const (
	insolarRep           = "github.com/insolar/insolar"
	stateMachineTemplate = "conveyor/generator/generator/templates/state_machine.go.tpl"
	matrixTemplate       = "conveyor/generator/generator/templates/matrix.go.tpl"
	generatedMatrix      = "conveyor/generator/matrix/matrix.go"
	generatedSuffix      = "_generated.go"
)

// todo remove this type
type TemporaryCustomAdapterHelper struct{}

type Generator struct {
	stateMachines        []*StateMachine
	fullPathToInsolar    string
	adapterHelperCatalog map[string]string
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func exitWithError(errMsg string, a ...interface{}) {
	panic(fmt.Sprintf(errMsg, a...))
}

func NewGenerator() *Generator {
	_, me, _, ok := runtime.Caller(0)
	if !ok {
		exitWithError("couldn't get self full path")
	}
	idx := strings.LastIndex(string(me), insolarRep)
	return &Generator{
		fullPathToInsolar:    string(me)[0 : idx+len(insolarRep)],
		adapterHelperCatalog: make(map[string]string),
	}
}

func (g *Generator) CheckAllMachines() {
	for _, machine := range g.stateMachines {

		checkHasInitHandlers(machine)

		for stateIndex, state := range machine.States {
			for _, pulseState := range []PulseState{Future, Present, Past} {
				for _, handlerType := range []handlerType{Transition, Migration, AdapterResponse} {
					currentHandler := state.handlers[pulseState][handlerType]

					if currentHandler == nil {
						continue
					}

					if len(currentHandler.params) < 3 || currentHandler.params[2] != *machine.InputEventType {
						exitWithError("[%s %s] Third parameter should be %s\n", machine.Name, currentHandler.Name, *machine.InputEventType)
					}

					// check Init handlers
					if stateIndex == 0 && handlerType == Transition {

						checkInitHandlersSignature(currentHandler, machine)

					} else {
						if currentHandler.params[3] != *machine.PayloadType {
							exitWithError("[%s %s] Fourth parameter should be %s not %s\n", machine.Name, currentHandler.Name, *machine.PayloadType, currentHandler.params[3])
						}
						if len(currentHandler.results) != 1 {
							exitWithError("[%s %s] Handlers should return only fsm.ElementState\n", machine.Name, currentHandler.Name)
						}
					}

					checkTransitionHandlerParams(stateIndex, handlerType, currentHandler, machine)

					checkAdapterRespHandler(handlerType, currentHandler, machine)
				}
			}
		}
	}
}

func checkTransitionHandlerParams(stateIndex int, handlerType handlerType, currentHandler *handler, machine *StateMachine) {
	if stateIndex != 0 && handlerType == Transition {
		if len(currentHandler.params) != 4 && len(currentHandler.params) != 5 {
			exitWithError("[%s %s] Transition handlers should have 4 or 5 (with adapter helper) parameters\n", machine.Name, currentHandler.Name)
		}
	}
}

func checkAdapterRespHandler(handlerType handlerType, currentHandler *handler, machine *StateMachine) {
	if handlerType == AdapterResponse {
		if len(currentHandler.params) != 5 {
			exitWithError("[%s %s] AdapterResponse handlers should have 5 parameters\n", machine.Name, currentHandler.Name)
		}
	}
}

func checkInitHandlersSignature(currentHandler *handler, machine *StateMachine) {
	if currentHandler.params[3] != "interface {}" {
		exitWithError("[%s %s] Init handlers should have interface{} as payload parameter\n", machine.Name, currentHandler.Name)
	}
	if currentHandler.results[1] != *machine.PayloadType {
		exitWithError("[%s %s] Init handlers should return payload as %s\n", machine.Name, currentHandler.Name, *machine.PayloadType)
	}
}

func checkHasInitHandlers(machine *StateMachine) {
	if machine.States[0].handlers[Present][Transition] == nil {
		exitWithError("[%s] Present Init handler should be defined", machine.Name)
	}
	if machine.States[0].handlers[Future][Transition] == nil {
		exitWithError("[%s] Future Init handler should be defined", machine.Name)
	}
}

func (g *Generator) ParseAdapterHelpers() {
	t := reflect.TypeOf(adapterhelper.Catalog{})
	for i := 0; i < t.NumField(); i++ {
		g.adapterHelperCatalog[t.Field(i).Type.Name()] = t.Field(i).Name
	}
}

type stateMachineWithID struct {
	StateMachine
	ID int
}

func (g *Generator) GenerateStateMachines() {
	for i, machine := range g.stateMachines {
		tplBody, err := ioutil.ReadFile(path.Join(g.fullPathToInsolar, stateMachineTemplate))
		checkErr(err)

		file, err := os.Create(machine.File[:len(machine.File)-3] + generatedSuffix)
		checkErr(err)

		out := bufio.NewWriter(file)

		err = template.Must(template.New("smTmpl").Funcs(templateFuncs).
			Parse(string(tplBody))).
			Execute(out, stateMachineWithID{StateMachine: *machine, ID: i + 1})
		checkErr(err)

		err = out.Flush()
		checkErr(err)
		err = file.Close()
		checkErr(err)
	}

}

type matrixParams struct {
	Imports  []string
	Machines []*StateMachine
}

func getDirFromPath(p string) string {
	dir, _ := path.Split(p)
	if strings.HasSuffix(dir, "/") {
		return dir[:len(dir)-1]
	}
	return dir
}

func fileToImport(f string) string {
	if idx := strings.Index(f, "github.com/insolar/insolar"); idx >= 0 {
		return getDirFromPath(f[idx:])
	}
	return getDirFromPath(f)
}

func (g *Generator) sortedImports() []string {
	importsMap := make(map[string]struct{})
	for _, machine := range g.stateMachines {
		importsMap[fileToImport(machine.File)] = struct{}{}
	}

	var imports []string
	for key := range importsMap {
		imports = append(imports, key)
	}
	sort.Strings(imports)
	return imports
}

func (g *Generator) GenerateMatrix() {
	params := matrixParams{
		Machines: g.stateMachines,
		Imports:  g.sortedImports(),
	}

	tplBody, err := ioutil.ReadFile(path.Join(g.fullPathToInsolar, matrixTemplate))
	checkErr(err)

	file, err := os.Create(path.Join(g.fullPathToInsolar, generatedMatrix))
	checkErr(err)

	defer file.Close()
	out := bufio.NewWriter(file)

	err = template.Must(template.New("MtTmpl").Funcs(templateFuncs).
		Parse(string(tplBody))).
		Execute(out, params)
	checkErr(err)

	err = out.Flush()
	checkErr(err)
}
