package pulse

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"errors"
)

type StartPulse interface {
	OnPulse(context.Context, insolar.Pulse)
	PulseNumber() (insolar.PulseNumber, error)
}

type startPulse struct {
	pulse *insolar.Pulse
}

func NewStartPulse() StartPulse {
	return &startPulse{}
}

func (sp *startPulse) OnPulse(ctx context.Context, pulse insolar.Pulse) {
	if sp.pulse == nil {
		sp.pulse = &pulse
	}
}

func (sp *startPulse) PulseNumber() (insolar.PulseNumber, error) {
	if sp.pulse == nil {
		return 0, errors.New("start pulse in nil")
	}
	return sp.pulse.PulseNumber, nil
}
