package goplugin

import "github.com/insolar/insolar/logicrunner"

type GoPlugin struct {
}

func (gp *GoPlugin) Exec(logicrunner.Object) (interface{}, error) {
	panic("implement me")
}

func New() GoPlugin {
	return GoPlugin{} // TODO
}
