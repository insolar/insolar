package reaction

import (
	"io"

	"github.com/insolar/insolar/core"
)

type Code struct {
	code        []byte           // nolint
	machineType core.MachineType // nolint
}

func (e *Code) Serialize() (io.Reader, error) {
	return serialize(e, TypeCode)
}
