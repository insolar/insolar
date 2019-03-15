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
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	insolarRep           = "github.com/insolar/insolar"
	stateMachineTemplate = "conveyor/generator/generator/templates/state_machine.go.tpl"
	matrixTemplate       = "conveyor/generator/generator/templates/matrix.go.tpl"
)

type stateMachine struct {
	Package        string
	Name           string
	InputEventType *string
	PayloadType    *string
	States         []state
}

type Generator struct {
	stateMachines       []*stateMachine
	imports             map[string]interface{}
	fullPathToInsolar   string
	pathToStateMachines string
	pathToMatrixFile    string
}

func NewGenerator(pathToStateMachines string, pathToMatrixFile string) *Generator {
	_, me, _, ok := runtime.Caller(0)
	if !ok {
		exitWithError("couldn't get self full path")
	}
	idx := strings.LastIndex(string(me), insolarRep)
	return &Generator{
		imports:             make(map[string]interface{}),
		fullPathToInsolar:   string(me)[0 : idx+len(insolarRep)],
		pathToStateMachines: pathToStateMachines,
		pathToMatrixFile:    pathToMatrixFile,
	}
}

func (g *Generator) ParseFile(dir string, filename string) {
	g.imports[path.Join(insolarRep, g.pathToStateMachines, dir)] = nil

	file := path.Join(g.fullPathToInsolar, g.pathToStateMachines, dir, filename)
	p := Parser{generator: g, sourceFilename: file}
	p.readStateMachinesInterfaceFile()
	p.findEachStateMachine()
	outFileName := file[0:len(file)-3] + "_generated.go"
	outFile, err := os.Create(outFileName)
	checkErr(err)
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	p.Generate(w)
	err = w.Flush()
	checkErr(err)
}
