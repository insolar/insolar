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
	"os"
	"github.com/pkg/errors"
	"bufio"
)

type handler struct {
	Name string
	Params []string
	Results []string
}

type state struct {
	Name string
	Transit *handler
	Error *handler
	Migrate *handler
}

type initHandler struct {
	EventType string
	payloadType string
}

type initErrorHandler struct {
	Name string
	EventType string
	payloadType string
}

type stateMachine struct {
	Module string
	Name string
	Init *initHandler
	InitError *initErrorHandler
	States []state
}

type Generator struct {
	stateMachines []*stateMachine
	imports map[string]interface{}
	matrix string
	base string
	path string
}

func NewGenerator(base string, path string, matrix string) *Generator{
	return &Generator{
		imports: make(map[string]interface{}),
		matrix: matrix,
		base: base,
		path: path,
	}
}

func (g *Generator) ParseFile(dir string, filename string) error {
	g.imports[g.importPath(dir)] = nil

	file := g.sourceFile(dir, filename)
	p := Parser{generator: g, module: g.modulePath(dir), sourceFilename: file}
	err := p.openFile()
	if err != nil {
		return err
	}
	p.findEachStateMachine()
	outFile, err := os.Create(g.generatedFile(file))
	if err != nil {
		return errors.Wrap(err, "Couldn't create file")
	}
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	p.Generate(w)
	w.Flush()
	return nil
}

