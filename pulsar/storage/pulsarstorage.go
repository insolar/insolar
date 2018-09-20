package pulsarstorage

import (
	"github.com/insolar/insolar/core"
)

type PulsarStorage interface {
	GetLastPulse() (*core.Pulse, error)
	UpdatePulse(*core.Pulse) error
}
