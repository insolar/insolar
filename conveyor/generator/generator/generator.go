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

const insolarRep = "github.com/insolar/insolar"

type stateMachine struct {
	Package string
	Name string
	InputEventType *string
	PayloadType *string
	States []state
}

type Generator struct {
	stateMachines []*stateMachine
	imports map[string]interface{}
	fullPathToInsolar string
	PathToStateMachines string
	pathToMatrixFile string
}

func NewGenerator(pathToStateMachines string, pathToMatrixFile string) *Generator{
	_, me, _, ok := runtime.Caller(0)
	if ok == false {
		exitWithError("couldn't get self full path")
	}
	idx := strings.LastIndex(string(me), insolarRep)
	return &Generator{
		imports: make(map[string]interface{}),
		fullPathToInsolar: string(me)[0:idx + len(insolarRep)],
		PathToStateMachines: pathToStateMachines,
		pathToMatrixFile: pathToMatrixFile,
	}
}

func (g *Generator) ParseFile(dir string, filename string) {
	g.imports[path.Join(insolarRep, g.PathToStateMachines, dir)] = nil

	file := path.Join(g.fullPathToInsolar, g.PathToStateMachines, dir, filename)
	p := Parser{generator: g, sourceFilename: file}
	p.openFile()
	p.findEachStateMachine()
	outFileName := file[0:len(file)-3] + "_generated.go"
	outFile, err := os.Create(outFileName)
	checkErr(err)
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	p.Generate(w)
	w.Flush()
}

