package logicrunner

type MachineType int

const (
	MachineTypeBuiltin MachineType = iota
	MachineTypeGoPlugin
)

type LogicRunner interface {
	Start()
	Stop()
	Exec(object Object, method string, args Arguments) (ret Arguments, err error)
}
