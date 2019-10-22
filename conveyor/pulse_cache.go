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
	"fmt"
	"sync"

	"github.com/insolar/insolar/pulse"
)

//func NewPulseDataCache(accessRotations int) PulseDataCache {
//	if accessRotations <= 0 {
//		panic("illegal value")
//	}
//	return PulseDataCache{
//		access:make([]map[pulse.Number]struct{}, accessRotations),
//	}
//}

type PulseDataCache struct {
	mutex     sync.RWMutex
	minRange  uint32
	cache     map[pulse.Number]pulse.Data
	access    []map[pulse.Number]struct{}
	accessIdx int
}

func (p *PulseDataCache) Init(minRange uint32, accessRotations int) {
	if p.access != nil {
		panic("illegal state")
	}
	if accessRotations < 0 {
		panic("illegal value")
	}
	p.minRange = minRange
	p.access = make([]map[pulse.Number]struct{}, accessRotations)
	p.access[0] = make(map[pulse.Number]struct{})
}

func (p *PulseDataCache) GetMinRange() uint32 {
	return p.minRange
}

func (p *PulseDataCache) EvictAndRotate(currentPN pulse.Number) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p._evict(currentPN)
	p._rotate()
	p._touch(currentPN) // to retain current PD at corner cases
}

func (p *PulseDataCache) EvictNoRotate(currentPN pulse.Number) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p._evict(currentPN)
}

func (p *PulseDataCache) Rotate() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p._rotate()
}

func (p *PulseDataCache) _evict(currentPN pulse.Number) {
	cpn := currentPN.AsUint32()
	if uint32(pulse.MinTimePulse)+p.minRange >= cpn {
		// must keep all
		return
	}
	minPN := pulse.OfUint32(cpn - p.minRange)

outer:
	for pn := range p.cache {
		if pn >= minPN {
			continue
		}
		for _, am := range p.access {
			if _, ok := am[pn]; ok {
				continue outer
			}
		}

		delete(p.cache, pn)
	}
}

type accessState uint8

const (
	miss accessState = iota
	hit
	hitNoTouch
)

func (p *PulseDataCache) getRO(pn pulse.Number) (pulse.Data, accessState) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.cache != nil {
		if pd, ok := p.cache[pn]; ok {
			if p._wasTouched(pn) {
				return pd, hit
			}
			return pd, hitNoTouch
		}
	}
	return pulse.Data{}, miss
}

func (p *PulseDataCache) Get(pn pulse.Number) (pulse.Data, bool) {
	pd, m := p.getRO(pn)
	if m != hitNoTouch {
		return pd, m != miss
	}

	p.mutex.Lock()
	p._touch(pn)
	p.mutex.Unlock()

	return pd, true
}

func (p *PulseDataCache) Check(pn pulse.Number) (pulse.Data, bool) {
	pd, m := p.getRO(pn)
	return pd, m != miss
}

func (p *PulseDataCache) Contains(pn pulse.Number) bool {
	_, m := p.getRO(pn)
	return m != miss
}

func (p *PulseDataCache) Touch(pn pulse.Number) bool {
	switch _, m := p.getRO(pn); m {
	case miss:
		return false
	case hit:
		return true
	}
	p.mutex.Lock()
	p._touch(pn)
	p.mutex.Unlock()
	return true
}

func (p *PulseDataCache) Put(pd pulse.Data) {
	switch epd, m := p.getRO(pd.PulseNumber); {
	case m == miss:
		//break
	case pd != epd:
		panic(fmt.Errorf("duplicate pulseData: before=%v after=%v", epd, pd))
	case m == hitNoTouch:
		p.mutex.Lock()
		p._touch(pd.PulseNumber)
		p.mutex.Unlock()
		return
	default: //m == hit:
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.cache == nil {
		p.cache = make(map[pulse.Number]pulse.Data)
		p.cache[pd.PulseNumber] = pd
	} else {
		switch epd, ok := p.cache[pd.PulseNumber]; {
		case !ok:
			p.cache[pd.PulseNumber] = pd
		case pd != epd:
			panic(fmt.Errorf("duplicate pulseData: before=%v after=%v", epd, pd))
		}
	}
	p._touch(pd.PulseNumber)
}

func (p *PulseDataCache) _wasTouched(pn pulse.Number) bool {
	_, ok := p.access[p.accessIdx][pn]
	return ok
}

func (p *PulseDataCache) _touch(pn pulse.Number) {
	p.access[p.accessIdx][pn] = struct{}{}
}

func (p *PulseDataCache) _rotate() {
	p.accessIdx++
	if p.accessIdx >= len(p.access) {
		p.accessIdx = 0
	}
	switch m := p.access[p.accessIdx]; {
	case m == nil:
		p.access[p.accessIdx] = make(map[pulse.Number]struct{})
	case len(m) == 0:
		// reuse
	default:
		p.access[p.accessIdx] = make(map[pulse.Number]struct{}, len(m))
	}
}
