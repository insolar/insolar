package logicrunner

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeBuiltin MachineType = iota
	MachineTypeGoPlugin
)

// LogicRunner is a general interface of contract executor
type LogicRunner interface {
	Start()
	Stop()
	Exec(object Object, method string, args Arguments) (ret Arguments, err error)
}
