/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulse

import (
	"context"
	"errors"

	"github.com/insolar/insolar/insolar"
)

type StartPulse interface {
	SetStartPulse(context.Context, insolar.Pulse)
	PulseNumber() (insolar.PulseNumber, error)
}

type startPulse struct {
	pulse *insolar.Pulse
}

func NewStartPulse() StartPulse {
	return &startPulse{}
}

func (sp *startPulse) SetStartPulse(ctx context.Context, pulse insolar.Pulse) {
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
