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
	PulseData() pulse.Data
	State() PulseSlotState

	MakePresent(pd pulse.Data)
	MakePast()
}

var _ pulseDataHolder = &futurePulseDataHolder{}

type futurePulseDataHolder struct {
	mutex  sync.RWMutex
	pd     pulse.Data
	isPast bool
}

func (p *futurePulseDataHolder) PulseData() pulse.Data {
	p.mutex.RLock()
	pd := p.pd
	p.mutex.RUnlock()
	return pd
}

func (p *futurePulseDataHolder) State() PulseSlotState {
	p.mutex.RLock()
	switch {
	case p.isPast:
		p.mutex.RUnlock()
		return Past
	case p.pd.IsEmpty():
		p.mutex.RUnlock()
		return 0
	case p.pd.IsExpectedPulse():
		p.mutex.RUnlock()
		return Future
	default:
		p.mutex.RUnlock()
		return Present
	}
}

func (p *futurePulseDataHolder) MakePresent(pd pulse.Data) {
	pd.EnsurePulsarData()

	switch p.State() {
	case Future:
		break
	case Present:
		return
	default:
		panic("illegal state")
	}

	p.mutex.Lock()
	p.pd = pd
	p.mutex.Unlock()
}

func (p *futurePulseDataHolder) MakePast() {
	switch p.State() {
	case Present:
		break
	case Past:
		return
	default:
		panic("illegal state")
	}

	p.mutex.Lock()
	p.isPast = true
	p.mutex.Unlock()
}

var _ pulseDataHolder = &presentPulseDataHolder{}

type presentPulseDataHolder struct {
	pd     pulse.Data
	isPast uint32 //atomic
}

func (p *presentPulseDataHolder) PulseData() pulse.Data {
	return p.pd
}

func (p *presentPulseDataHolder) State() PulseSlotState {
	if atomic.LoadUint32(&p.isPast) == 0 {
		return Present
	}
	return Past
}

func (p *presentPulseDataHolder) MakePresent(pd pulse.Data) {
	if p.State() != Present {
		panic("illegal state")
	}
	if p.pd != pd {
		panic("illegal value")
	}
}

func (p *presentPulseDataHolder) MakePast() {
	atomic.StoreUint32(&p.isPast, 1)
}

var _ pulseDataHolder = &antiquePulseDataHolder{}

type antiquePulseDataHolder struct {
}

func (p antiquePulseDataHolder) PulseData() pulse.Data {
	panic("illegal state")
}

func (p antiquePulseDataHolder) State() PulseSlotState {
	return Antique
}

func (p antiquePulseDataHolder) MakePresent(pd pulse.Data) {
	panic("illegal state")
}

func (p antiquePulseDataHolder) MakePast() {
	panic("illegal state")
}
