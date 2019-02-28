package generator

import (
	"text/template"
	"os"
	"bufio"
)

var (
	matrixTmpl = template.Must(template.New("matrixTmpl").
		Parse(`package matrix

import (
    "github.com/insolar/insolar/conveyor/generator/common"
    {{range .Imports}}"{{.}}"
    {{end}}
)

type matrix struct {
    matrix  []common.StateMachine
    indexes map[string]int
}

var Matrix matrix

func init() {
    Matrix := matrix{
        indexes: make(map[string]int),
    }
    Matrix.matrix = append(Matrix.matrix,
        {{range .Machines}}{{.Module}}.SMRH{{.Name}}Export(),
        {{end}}
	)

    {{range $i, $machine := .Machines}}Matrix.indexes["{{.Module}}.{{$machine.Name}}"] = {{$i}}
    {{end}}
}

func (m *matrix) GetHandlers(machine int, state int) *common.State {
    return &m.matrix[machine].States[state]
}

func (m *matrix) GetIdx(machine string) int {
    return m.indexes[machine]
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


