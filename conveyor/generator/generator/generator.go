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
	"go/ast"
	"io/ioutil"
	"go/token"
	"go/parser"
	"strings"
	"log"
)

type method struct {
	name string
	params []string
	results []string
}
type state struct {
	name string
	transit *method
	error *method
	migrate *method
}

type initHandler struct {
	eventType string
	payloadType string
}

type stateMachine struct {
	Name string
	Init *initHandler
	States []state
}

type Generator struct {
	sourceFilename string
	sourceCode []byte
	sourceNode *ast.File
	stateMachines []*stateMachine
}

func (g *Generator) ParseFile(filename string) error {
	g.sourceFilename = filename
	err := g.openFile()
	if err != nil {
		return err
	}
	g.findEachStateMachine()
	return nil
}

func (g *Generator) openFile() error {
	var err error
	g.sourceCode, err = ioutil.ReadFile(g.sourceFilename)
	if err != nil {
		return err
	}
	fSet := token.NewFileSet()
	g.sourceNode, err = parser.ParseFile(fSet, g.sourceFilename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) findEachStateMachine() {
	for _, decl := range g.sourceNode.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			currStruct, ok := currType.Type.(*ast.InterfaceType)
			if !ok || !isStateMachineTag(genDecl.Doc) {
				continue
			}

			machine := &stateMachine{Name: currType.Name.Name}
			g.parseStateMachineInterface(machine, currStruct)
		}
	}
}

func isStateMachineTag(group *ast.CommentGroup) bool {
	for _, comment := range group.List {
		if strings.Contains(comment.Text, "conveyer: state_machine") {
			return true
		}
	}
	return false
}

func remap(code []byte, m *ast.FieldList) []string {
	if m == nil {
		return make([]string, 0)
	}
	result := make([]string, len(m.List))
	for i, z := range m.List {
		t := z.Type
		result[i] = string(code[t.Pos()-1 : t.End()-1])
	}
	return result
}

func (g *Generator) parseStateMachineInterface(machine *stateMachine, source *ast.InterfaceType) {
	for _, methodItem := range source.Methods.List {
		if len(methodItem.Names) == 0 {
			continue
		}
		methodType := methodItem.Type.(*ast.FuncType)

		currentMethod := &method{
			name: methodItem.Names[0].Name,
			params: remap(g.sourceCode, methodType.Params),
			results: remap(g.sourceCode, methodType.Results),
		}
		switch {
		case currentMethod.name == "Init":
			g.parseInit(machine, currentMethod)
		case strings.HasPrefix(currentMethod.name, "State"):
			g.parseState(machine, currentMethod)
		case strings.HasPrefix(currentMethod.name, "Transit"):
			g.parseTransit(machine, currentMethod)
		case strings.HasPrefix(currentMethod.name, "Migrate"):
			g.parseMigrate(machine, currentMethod)
		case strings.HasPrefix(currentMethod.name, "Error"):
			g.parseError(machine, currentMethod)
		default:
			log.Fatal("Unknokn handler:", currentMethod.name)
		}
	}
	g.stateMachines = append(g.stateMachines, machine)
}

func (g *Generator) parseInit(machine *stateMachine, m *method) {
	if len(m.params) != 1 {
		log.Fatal("Init must have only one parameter")
	}
	if len(m.results) != 3 {
		log.Fatal("Init should return three values")
	}
	if !strings.HasPrefix(m.results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if m.results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if m.results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if machine.Init != nil {
		log.Fatal("Only one init handler for state machine")
	}
	machine.Init = &initHandler{
		eventType: m.params[0],
		payloadType: m.results[0],
	}
}

func (g *Generator) parseState(machine *stateMachine, m *method) {
	if len(m.params) != 0 {
		log.Fatal("State must not have any parameters")
	}
	if len(m.results) != 1 || m.results[0] != "common.ElState" {
		log.Fatal("State should returns only common.ElState")
	}
	machine.States = append(machine.States, state{name: m.name})
}

func (g *Generator) parseTransit(machine *stateMachine, m *method) {
	if len(m.params) != 2 {
		log.Fatal("Transit must have two parameters")
	}
	if m.params[0] != machine.Init.eventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(m.params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if len(m.results) != 3 {
		log.Fatal("Transit should return three values")
	}
	if !strings.HasPrefix(m.results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if m.results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if m.results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].transit != nil {
		log.Fatal("Only one transit handler for state")
	}
	machine.States[len(machine.States)-1].transit = m
}

func (g *Generator) parseMigrate(machine *stateMachine, m *method) {
	if len(m.params) != 2 {
		log.Fatal("Migrate must have two parameters")
	}
	if m.params[0] != machine.Init.eventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(m.params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if len(m.results) != 3 {
		log.Fatal("Migrate should return three values")
	}
	if !strings.HasPrefix(m.results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if m.results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if m.results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].migrate != nil {
		log.Fatal("Only one migrate handler for state")
	}
	machine.States[len(machine.States)-1].migrate = m
}

func (g *Generator) parseError(machine *stateMachine, m *method) {
	if len(m.params) != 3 {
		log.Fatal("Error handler must have three parameters")
	}
	if m.params[0] != machine.Init.eventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(m.params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if m.params[2] != "error" {
		log.Fatal("Third parameter must be of type error")
	}
	if len(m.results) != 2 {
		log.Fatal("Error should return two values")
	}
	if !strings.HasPrefix(m.results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if m.results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].error != nil {
		log.Fatal("Only one error handler for state")
	}
	machine.States[len(machine.States)-1].error = m
}