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

package smachine

import (
	"sync/atomic"
	"time"
)

type boostPermit struct {
	active uint32 //atomic
}

func (p *boostPermit) _get() uint32 {
	return atomic.LoadUint32(&p.active)
}

// nil-safe
func (p *boostPermit) isActive() bool {
	return p != nil && p._get() != 0
}

func (p *boostPermit) discard() {
	atomic.StoreUint32(&p.active, 0)
}

var activeBoost = &boostPermit{1}
var inactiveBoost = &boostPermit{0}

type chainedBoostPermit struct {
	boostPermit
	timeMark int64
	next     *chainedBoostPermit
}

func (p *chainedBoostPermit) discardOlderThan(t time.Time) *chainedBoostPermit {
	tn := t.UnixNano()
	n := p
	for ; n != nil && tn >= n.timeMark; n = n.next {
		n.discard()
	}
	return n
}

func (p *chainedBoostPermit) canReuse() bool {
	return p._get() == 1
}

func (p *chainedBoostPermit) use() {
	if atomic.CompareAndSwapUint32(&p.active, 1, 2) {
		return
	}
	if p._get() != 2 {
		panic("illegal state")
	}
}

// nil-safe
func (p *chainedBoostPermit) reuseOrNew(t time.Time) *chainedBoostPermit {
	switch {
	case p == nil:
		//
	case p.canReuse():
		p.timeMark = t.UnixNano()
		return p
	case p.next != nil:
		panic("illegal state")
	}
	n := &chainedBoostPermit{boostPermit{1}, t.UnixNano(), nil}
	if p != nil {
		p.next = n
	}
	return n
}
