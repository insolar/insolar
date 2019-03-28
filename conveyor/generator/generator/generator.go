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
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

const (
	insolarRep           = "github.com/insolar/insolar"
	stateMachineTemplate = "conveyor/generator/generator/templates/state_machine.go.tpl"
	matrixTemplate       = "conveyor/generator/generator/templates/matrix.go.tpl"
	generatedMatrix      = "conveyor/generator/matrix/matrix.go"
)

// todo remove this type
type TemporaryCustomAdapterHelper struct{}

type Generator struct {
	stateMachines     []*StateMachine
	fullPathToInsolar string
}

var gen *Generator

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exitWithError(err string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(err, a...))
	os.Exit(1)
}

func init() {
	_, me, _, ok := runtime.Caller(0)
	if !ok {
		exitWithError("couldn't get self full path")
	}
	idx := strings.LastIndex(string(me), insolarRep)
	gen = &Generator{
		fullPathToInsolar: string(me)[0 : idx+len(insolarRep)],
	}
}

func CheckAllMachines() {
	for _, machine := range gen.stateMachines {
		if machine.States[0].handlers[Present][Transition] == nil {
			log.Fatal("Present Init handler should be defined")
		}
		if machine.States[0].handlers[Future][Transition] == nil {
			log.Fatal("Future Init handler should be defined")
		}

		for i, s := range machine.States {
			for _, ps := range []PulseState{Future, Present, Past} {
				for _, ht := range []handlerType{Transition, Migration, AdapterResponse} {
					if s.handlers[ps][ht] == nil {
						continue
					}

					if len(s.handlers[ps][ht].params) < 3 || s.handlers[ps][ht].params[2] != *machine.InputEventType {
						log.Fatal("[", s.handlers[ps][ht].funcName, "] Forth parameter should be ", *machine.InputEventType)
					}

					if i == 0 && ht == Transition {
						if s.handlers[ps][ht].params[3] != "interface {}" {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] Init handlers should have interface{} as payload parameter")
						}
						if s.handlers[ps][ht].results[1] != *machine.PayloadType {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] Init handlers should return payload as ", *machine.PayloadType)
						}
					} else {
						if s.handlers[ps][ht].params[3] != *machine.PayloadType {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] Handlers payload should be ", *machine.PayloadType, " current ", s.handlers[ps][ht].params[3])
						}
						if len(s.handlers[ps][ht].results) != 1 {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] Handlers should return only fsm.ElementState")
						}
					}

					if i != 0 && ht == Transition {
						if len(s.handlers[ps][ht].params) != 4 && len(s.handlers[ps][ht].params) != 5 {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] Transition handlers should have 4 or 5 (with adapher helper) parameters")
						}
					}

					if ht == AdapterResponse {
						if len(s.handlers[ps][ht].params) != 5 {
							log.Fatal("[", s.handlers[ps][ht].funcName, "] AdapterResponse handlers should have 5 parameters")
						}
					}
				}
			}
		}
	}
}

type stateMachineWithId struct {
	StateMachine
	ID int
}

func GenerateStateMachines() {
	for i, machine := range gen.stateMachines {
		tplBody, err := ioutil.ReadFile(path.Join(gen.fullPathToInsolar, stateMachineTemplate))
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(machine.File[:len(machine.File)-3] + "_generated.go")
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		out := bufio.NewWriter(file)

		err = template.Must(template.New("newSmTmpl").Funcs(templateFuncs).
			Parse(string(tplBody))).
			Execute(out, stateMachineWithId{StateMachine: *machine, ID: i + 1})
		if err != nil {
			log.Fatal(err)
		}
		err = out.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}

}
func GenerateMatrix() {
	tplBody, err := ioutil.ReadFile(path.Join(gen.fullPathToInsolar, matrixTemplate))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(path.Join(gen.fullPathToInsolar, generatedMatrix))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	out := bufio.NewWriter(file)

	err = template.Must(template.New("newMtrxTmpl").Funcs(templateFuncs).
		Parse(string(tplBody))).
		Execute(out, gen.stateMachines)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
