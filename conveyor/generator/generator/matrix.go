package generator

import (
	"text/template"
	"os"
	"bufio"
)

var (
	matrixFuncMap = template.FuncMap{
		"isNull": func(i int) bool {
			return i == 0
		},
	}
	matrixTmpl = template.Must(template.New("matrixTmpl").Funcs(matrixFuncMap).
		Parse(`package matrix

import (
    "github.com/insolar/insolar/conveyor/generator/common"
    {{range .Imports}}"{{.}}"
    {{end}}
)

type matrix struct {
    matrix  []*common.StateMachine
}

type MachineType int
var Matrix matrix

const (
    {{range $i, $m := .Machines}}{{$m.Name}}{{if (isNull $i)}} MachineType = iota + 1{{end}}{{end}}
)

func init() {
    Matrix := matrix{}
    Matrix.matrix = append(Matrix.matrix, nil,
        {{range .Machines}}{{.Module}}.SMRH{{.Name}}Factory(),
        {{end}})
}

func (m *matrix) GetStateMachineByType(mType MachineType) *common.StateMachine {
    return m.matrix[int(mType)]
}
`))
)

func (g *Generator) getImports() []string {
	keys := make([]string, len(g.imports))
	i := 0
	for k := range g.imports {
		keys[i] = k
		i++
	}
	return keys
}

func (g *Generator) GenMatrix () error {
	file, err := os.Create(g.matrix)
	if err != nil {
		return err
	}
	defer file.Close()
	out := bufio.NewWriter(file)
	matrixTmpl.Execute(out, struct{
		Imports []string
		Machines []*stateMachine
	}{
		Imports: g.getImports(),
		Machines: g.stateMachines,

	})
	out.Flush()
	return nil
}


