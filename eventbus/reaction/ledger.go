package reaction

import (
	"io"

	"github.com/insolar/insolar/core"
)

type Code struct {
	code        []byte
	machineType core.MachineType
}

func (e *Code) Serialize() (io.Reader, error) {
	return serialize(e, TypeCode)
}
