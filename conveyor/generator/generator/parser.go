package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
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

func (p *Parser) openFile() {
	var err error
	p.sourceCode, err = ioutil.ReadFile(p.sourceFilename)
	checkErr(err)
	fSet := token.NewFileSet()
	p.sourceNode, err = parser.ParseFile(fSet, p.sourceFilename, nil, parser.ParseComments)
	checkErr(err)
	p.Package = p.sourceNode.Name.Name
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

			machine := &stateMachine{
				Name: currType.Name.Name,
				Module: p.module,
				States: []state{{Name: "Init"}},
			}
			p.parseStateMachineInterface(machine, currStruct)
		}
	}
}

func isStateMachineTag(group *ast.CommentGroup) bool {
	for _, comment := range group.List {
		if strings.Contains(comment.Text, "conveyor: state_machine") {
			return true
		}
	}
	return false
}

func getFieldTypes(code []byte, fieldList *ast.FieldList) []string {
	if fieldList == nil {
		return make([]string, 0)
	}
	result := make([]string, len(fieldList.List))
	for i, field := range fieldList.List {
		fieldType := field.Type
		result[i] = string(code[fieldType.Pos()-1 : fieldType.End()-1])
	}
	return result
}

func (p *Parser) parseStateMachineInterface(machine *stateMachine, source *ast.InterfaceType) {
	curPos := token.Pos(0)
	for _, methodItem := range source.Methods.List {
		if len(methodItem.Names) == 0 {
			continue
		}
		if methodItem.Pos() <= curPos {
			exitWithError("Incorrect order of methods")
		}
		curPos = methodItem.Pos()
		methodType := methodItem.Type.(*ast.FuncType)

		currentHandler := &handler{
			machine: machine,
			state: len(machine.States) -1,
			Name: methodItem.Names[0].Name,
			Params: getFieldTypes(p.sourceCode, methodType.Params),
			Results: getFieldTypes(p.sourceCode, methodType.Results),
		}

		switch {
		case currentHandler.Name == "TID":
		case strings.HasPrefix(currentHandler.Name, "s_"):
			currentHandler.setAsState()
		case strings.HasPrefix(currentHandler.Name, "i_"):
			currentHandler.setAsInit()
		case strings.HasPrefix(currentHandler.Name, "if_"):
			currentHandler.setAsInitFuture()
		case strings.HasPrefix(currentHandler.Name, "ip_"):
			currentHandler.setAsInitPast()
		case strings.HasPrefix(currentHandler.Name, "es_"):
			currentHandler.setAsErrorState()
		case strings.HasPrefix(currentHandler.Name, "esf_"):
			currentHandler.setAsErrorStateFuture()
		case strings.HasPrefix(currentHandler.Name, "esp_"):
			currentHandler.setAsErrorStatePast()
		case strings.HasPrefix(currentHandler.Name, "m_"):
			currentHandler.setAsMigration()
		case strings.HasPrefix(currentHandler.Name, "mfp_"):
			currentHandler.setAsMigrationFuturePresent()
		case strings.HasPrefix(currentHandler.Name, "t_"):
			currentHandler.setAsTransition()
		case strings.HasPrefix(currentHandler.Name, "tf_"):
			currentHandler.setAsTransitionFuture()
		case strings.HasPrefix(currentHandler.Name, "tp_"):
			currentHandler.setAsTransitionPast()
		/*case strings.HasPrefix(currentHandler.Name, "f_"):
			currentHandler.setAsFinalization()
		case strings.HasPrefix(currentHandler.Name, "ff_"):
			currentHandler.setAsFinalizationFuture()
		case strings.HasPrefix(currentHandler.Name, "fp_"):
			currentHandler.setAsFinalizationPast()*/
		case strings.HasPrefix(currentHandler.Name, "a_"):
			currentHandler.setAsAdapterResponse()
		case strings.HasPrefix(currentHandler.Name, "af_"):
			currentHandler.setAsAdapterResponseFuture()
		case strings.HasPrefix(currentHandler.Name, "ap_"):
			currentHandler.setAsAdapterResponsePast()
		case strings.HasPrefix(currentHandler.Name, "ea_"):
			currentHandler.setAsAdapterResponseError()
		case strings.HasPrefix(currentHandler.Name, "eaf_"):
			currentHandler.setAsAdapterResponseErrorFuture()
		case strings.HasPrefix(currentHandler.Name, "eap_"):
			currentHandler.setAsAdapterResponseErrorPast()
		default:
			exitWithError("Unknown handler: %s", currentHandler.Name)
		}
	}
	p.StateMachines = append(p.StateMachines, machine)
	p.generator.stateMachines = append(p.generator.stateMachines, machine)
}
