package logicrunner

import (
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type MachineType int

const (
	MachineTypeBuiltin MachineType = iota
	MachineTypeGoPlugin
)

type LogicRunner interface {
	Exec(object Object, method string, args Arguments) (ret Arguments, err error)
}

func NewLogicRunner(t MachineType, b API) LogicRunner {
	switch t {
	case MachineTypeGoPlugin:
		return goplugin.New(b)
	}
	return nil
}
