package generator

import (
	"strings"
	"io/ioutil"
	"go/token"
	"go/parser"
	"go/ast"
	"log"
)

type Parser struct {
	sourceFilename string
	sourceCode []byte
	sourceNode *ast.File
	generator *Generator
	module string
	Package string
	StateMachines []*stateMachine
}

func (p *Parser) openFile() error {
	var err error
	p.sourceCode, err = ioutil.ReadFile(p.sourceFilename)
	if err != nil {
		return err
	}
	fSet := token.NewFileSet()
	p.sourceNode, err = parser.ParseFile(fSet, p.sourceFilename, nil, parser.ParseComments)
	p.Package = p.sourceNode.Name.Name
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) findEachStateMachine() {
	for _, decl := range p.sourceNode.Decls {
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

			machine := &stateMachine{Name: currType.Name.Name, Module: p.module}
			p.parseStateMachineInterface(machine, currStruct)
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

func (p *Parser) parseStateMachineInterface(machine *stateMachine, source *ast.InterfaceType) {
	var curPos token.Pos = 0
	for _, methodItem := range source.Methods.List {
		if len(methodItem.Names) == 0 {
			continue
		}
		if methodItem.Pos() <= curPos {
			log.Fatal("Incorrect order of methods")
		}
		curPos = methodItem.Pos()
		methodType := methodItem.Type.(*ast.FuncType)

		currentHandler := &handler{
			Name: methodItem.Names[0].Name,
			Params: remap(p.sourceCode, methodType.Params),
			Results: remap(p.sourceCode, methodType.Results),
		}
		switch {
		case currentHandler.Name == "Init":
			p.parseInit(machine, currentHandler)
		case strings.HasPrefix(currentHandler.Name, "State"):
			p.parseState(machine, currentHandler)
		case strings.HasPrefix(currentHandler.Name, "Transit"):
			p.parseTransit(machine, currentHandler)
		case strings.HasPrefix(currentHandler.Name, "Migrate"):
			p.parseMigrate(machine, currentHandler)
		case strings.HasPrefix(currentHandler.Name, "Error"):
			p.parseError(machine, currentHandler)
		default:
			log.Fatal("Unknokn handler:", currentHandler.Name)
		}
	}
	p.StateMachines = append(p.StateMachines, machine)
	p.generator.stateMachines = append(p.generator.stateMachines, machine)
}

func (*Parser) parseInit(machine *stateMachine, h *handler) {
	if len(h.Params) != 1 {
		log.Fatal("Init must have only one parameter")
	}
	if len(h.Results) != 3 {
		log.Fatal("Init should return three values")
	}
	if !strings.HasPrefix(h.Results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if h.Results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if h.Results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if machine.Init != nil {
		log.Fatal("Only one init handler for state machine")
	}
	machine.Init = &initHandler{
		EventType: h.Params[0],
		payloadType: h.Results[0],
	}
}

func (*Parser) parseState(machine *stateMachine, h *handler) {
	if len(h.Params) != 0 {
		log.Fatal("State must not have any parameters")
	}
	if len(h.Results) != 1 || h.Results[0] != "common.ElState" {
		log.Fatal("State should returns only common.ElState")
	}
	machine.States = append(machine.States, state{Name: h.Name})
}

func (*Parser) parseTransit(machine *stateMachine, h *handler) {
	if len(h.Params) != 2 {
		log.Fatal("Transit must have two parameters")
	}
	if h.Params[0] != machine.Init.EventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(h.Params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if len(h.Results) != 3 {
		log.Fatal("Transit should return three values")
	}
	if !strings.HasPrefix(h.Results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if h.Results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if h.Results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].Transit != nil {
		log.Fatal("Only one transit handler for state")
	}
	machine.States[len(machine.States)-1].Transit = h
}

func (*Parser) parseMigrate(machine *stateMachine, h *handler) {
	if len(h.Params) != 2 {
		log.Fatal("Migrate must have two parameters")
	}
	if h.Params[0] != machine.Init.EventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(h.Params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if len(h.Results) != 3 {
		log.Fatal("Migrate should return three values")
	}
	if !strings.HasPrefix(h.Results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if h.Results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if h.Results[2] != "error" {
		log.Fatal("Returned error must be of type error")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].Migrate != nil {
		log.Fatal("Only one migrate handler for state")
	}
	machine.States[len(machine.States)-1].Migrate = h
}

func (*Parser) parseError(machine *stateMachine, h *handler) {
	if len(h.Params) != 3 {
		log.Fatal("Error handler must have three parameters")
	}
	if h.Params[0] != machine.Init.EventType {
		log.Fatal("Event should be of the same type with the event in the init")
	}
	if !strings.HasPrefix(h.Params[1], "*") {
		log.Fatal("Payload must be a pointer")
	}
	if h.Params[2] != "error" {
		log.Fatal("Third parameter must be of type error")
	}
	if len(h.Results) != 2 {
		log.Fatal("Error should return two values")
	}
	if !strings.HasPrefix(h.Results[0], "*") {
		log.Fatal("Returned payload should be a pointer")
	}
	if h.Results[1] != "common.ElState" {
		log.Fatal("Returned state should be common.ElState")
	}
	if len(machine.States) < 1 {
		log.Fatal("Declare state before handler")
	}
	if machine.States[len(machine.States)-1].Error != nil {
		log.Fatal("Only one error handler for state")
	}
	machine.States[len(machine.States)-1].Error = h
}
