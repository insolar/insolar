package generator

import (
	"text/template"
	"bytes"
	"os"
)

var (
	genTmpl = template.Must(template.New("genTmpl").Parse(`
type _states struct {
	transit func(element common.SlotElementHelper) (interface{}, common.ElState, error)
	migrate func(element common.SlotElementHelper) (interface{}, common.ElState, error)
	error func(element common.SlotElementHelper, err error) (interface{}, common.ElState)
}

var matrix map[string][]_states

func _() {
    matrix["{{.Name}}"] = []_states {
        {{.States}}
    }
}
`))
	genStates = template.Must(template.New("genStates").Parse(`{
                transit: {{.TransitHandler}},
                migrate: {{.MigrateHandler}},
                error: {{.ErrorHandler}},
        },`))
)

func (g *Generator) GenMatrix () {

	states := new(bytes.Buffer)
	for _, state := range g.stateMachines[0].States {
		genStates.Execute(states, struct {
			TransitHandler string
			MigrateHandler string
			ErrorHandler string
		}{
			TransitHandler: state.transit.name,
			MigrateHandler: state.migrate.name,
			ErrorHandler: state.error.name,
		})
	}
	genTmpl.Execute(os.Stdout, struct{
		Name string
		States string
	}{
		Name: g.stateMachines[0].Name,
		States: states.String(),
	})
}


