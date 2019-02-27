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

var matrix [][]common.State
var indexes = make(map[string]int)

func init() {
    matrix = append(matrix,
        {{range .Machines}}{{.Module}}.SMRH{{.Name}}Export(),
        {{end}}
	)

    {{range $i, $machine := .Machines}}indexes["{{.Module}}.{{$machine.Name}}"] = {{$i}}
    {{end}}
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
	file, err := os.Create("conveyor/generator/matrix/matrix.go")
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


