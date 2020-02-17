// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"context"
	"errors"
	"sync"

	"github.com/insolar/insolar/insolar"
)

type StartPulse interface {
	SetStartPulse(context.Context, insolar.Pulse)
	PulseNumber() (insolar.PulseNumber, error)
}

type startPulse struct {
	sync.RWMutex
	pulse *insolar.Pulse
}

func NewStartPulse() StartPulse {
	return &startPulse{}
}

func (sp *startPulse) SetStartPulse(ctx context.Context, pulse insolar.Pulse) {
	sp.Lock()
	defer sp.Unlock()

	if sp.pulse == nil {
		sp.pulse = &pulse
	}
}

func (sp *startPulse) PulseNumber() (insolar.PulseNumber, error) {
	sp.RLock()
	defer sp.RUnlock()

	if sp.pulse == nil {
		return 0, errors.New("start pulse in nil")
	}
	return sp.pulse.PulseNumber, nil
}
