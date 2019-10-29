//
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
//

package conveyor

import (
	"github.com/insolar/insolar/pulse"
	"sync"
	"sync/atomic"
)

type pulseDataHolder interface {
	// range is nil for Future and Antique
	PulseRange() (pulse.Range, PulseSlotState)
	// data is expected for Future, and empty for Antique
	PulseData() (pulse.Data, PulseSlotState)

	MakePresent(pr pulse.Range)
	MakePast()
}

var _ pulseDataHolder = &futurePulseDataHolder{}

type futurePulseDataHolder struct {
	mutex    sync.RWMutex
	expected pulse.Data
	pr       pulse.Range
	isPast   bool
}

func (p *futurePulseDataHolder) PulseData() (pulse.Data, PulseSlotState) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	switch {
	case p.pr == nil:
		switch {
		case p.expected.IsEmpty():
			return pulse.Data{}, 0
		case p.isPast:
			panic("illegal state")
		}
		return p.expected, Future
	case p.isPast:
		return p.pr.RightBoundData(), Past
	default:
		return p.pr.RightBoundData(), Present
	}
}

func (p *futurePulseDataHolder) PulseRange() (pulse.Range, PulseSlotState) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	switch {
	case p.pr == nil:
		switch {
		case p.expected.IsEmpty():
			return nil, 0
		case p.isPast:
			panic("illegal state")
		}
		return p.pr, Future
	case p.isPast:
		return p.pr, Past
	default:
		return p.pr, Present
	}
}

func (p *futurePulseDataHolder) MakePresent(pr pulse.Range) {
	pr.RightBoundData().EnsurePulsarData()

	if _, ps := p.PulseRange(); ps != Future {
		panic("illegal state")
	}

	p.mutex.Lock()
	p.pr = pr
	p.mutex.Unlock()
}

func (p *futurePulseDataHolder) MakePast() {
	if _, ps := p.PulseRange(); ps != Present {
		panic("illegal state")
	}
	p.mutex.Lock()
	p.isPast = true
	p.mutex.Unlock()
}

var _ pulseDataHolder = &presentPulseDataHolder{}

type presentPulseDataHolder struct {
	pr     pulse.Range
	isPast uint32 //atomic
}

func (p *presentPulseDataHolder) PulseData() (pulse.Data, PulseSlotState) {
	return p.pr.RightBoundData(), p.State()
}

func (p *presentPulseDataHolder) PulseRange() (pulse.Range, PulseSlotState) {
	return p.pr, p.State()
}

func (p *presentPulseDataHolder) State() PulseSlotState {
	if atomic.LoadUint32(&p.isPast) == 0 {
		return Present
	}
	return Past
}

func (p *presentPulseDataHolder) MakePresent(pulse.Range) {
	panic("illegal state")
}

func (p *presentPulseDataHolder) MakePast() {
	atomic.StoreUint32(&p.isPast, 1)
}

var _ pulseDataHolder = &antiqueNoPulseDataHolder{}

type antiqueNoPulseDataHolder struct {
}

func (p antiqueNoPulseDataHolder) PulseData() (pulse.Data, PulseSlotState) {
	return pulse.Data{}, Antique
}

func (p antiqueNoPulseDataHolder) PulseRange() (pulse.Range, PulseSlotState) {
	return nil, Antique
}

func (p antiqueNoPulseDataHolder) MakePresent(pulse.Range) {
	panic("illegal state")
}

func (p antiqueNoPulseDataHolder) MakePast() {
	panic("illegal state")
}
