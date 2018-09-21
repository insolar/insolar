package pulsarstorage

import (
	"github.com/insolar/insolar/core"
)

type PulsarStorage interface {
	GetLastPulse() (*core.Pulse, error)
	SetLastPulse(pulse *core.Pulse) error
	SavePulse(pulse *core.Pulse) error
}
