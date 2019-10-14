///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package conveyor

import (
	"github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/pulse"
)

type PulseService struct {
	preparePulse    func(chan<- PreparedState) bool
	cancelPulse     func() bool
	preparePulseOut chan<- PreparedState
}

func (p *PulseService) subscribe(preparePulse func(chan<- PreparedState) bool, cancelPulse func() bool) {
	if p.preparePulse != nil || p.cancelPulse != nil {
		panic("illegal state")
	}
	p.preparePulse = preparePulse
	p.cancelPulse = cancelPulse
}

func (p PulseService) onPreparePulseChange(out chan<- PreparedState) {
	if p.preparePulse != nil {
		p.preparePulse(out)
	}
}

func (p PulseService) onCancelPulseChange() {
	p.preparePulseOut = nil
	if p.cancelPulse != nil {
		p.cancelPulse()
	}
}

func (p PulseService) onCommitPulseChange(pulse.Data) {
	p.preparePulseOut = nil
	p.preparePulse = nil
	p.cancelPulse = nil
}

type PulseServiceAdapter struct {
	svc  PulseService
	exec smachine.ExecutionAdapter
}

func (a *PulseServiceAdapter) Call(ctx smachine.ExecutionContext, fn func(svc *PulseService)) {
	a.exec.PrepareSync(ctx, func() smachine.AsyncResultFunc {
		fn(&a.svc)
		return nil
	})
}
