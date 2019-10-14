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

package tools

import (
	"sync"
	"sync/atomic"
)

func NewVersionedSignal() VersionedSignal {
	return VersionedSignal{}
}

type VersionedSignal struct {
	mutex   sync.Mutex
	cond    *sync.Cond
	version uint32 //read:atomic; write:lock+atomic
	verCh   chan struct{}
}

func (p *VersionedSignal) getVersion() uint32 {
	return atomic.LoadUint32(&p.version)
}

func (p *VersionedSignal) Mark() SignalVersion {
	return SignalVersion{p, p.getVersion(), nil}
}

/* Increment version and send a signal */
func (p *VersionedSignal) NextBroadcast() {
	p.NextBroadcastAndMark()
}

/* Increment version and send a signal */
func (p *VersionedSignal) NextBroadcastAndMark() SignalVersion {
	p.mutex.Lock()
	p.verCh = nil
	v := p.getVersion() + 1
	if !atomic.CompareAndSwapUint32(&p.version, v-1, v) {
		p.mutex.Unlock()
		panic("illegal state")
	}
	if p.cond != nil {
		p.cond.Broadcast()
	}
	p.mutex.Unlock()
	return SignalVersion{p, v, nil}
}

/* Broadcasts on every signal */
func (p *VersionedSignal) GetCond() *sync.Cond {
	p.mutex.Lock()
	c := p.getCond()
	p.mutex.Unlock()
	return c
}

func (p *VersionedSignal) getCond() *sync.Cond {
	c := p.cond
	if c == nil {
		c = sync.NewCond(&p.mutex)
		p.cond = c
	}
	return c
}

func (p *VersionedSignal) getChannel(v uint32) <-chan struct{} {
	if p.getVersion() != v {
		return ClosedChannel()
	}

	p.mutex.Lock()
	if p.getVersion() != v {
		return ClosedChannel()
	}

	ch := p.verCh
	if ch == nil {
		ch = make(chan struct{})
		p.verCh = ch

		go p.workerCloser(v, p.getCond(), ch)
	}
	p.mutex.Unlock()
	return ch
}

func (p *VersionedSignal) workerCloser(v uint32, cd *sync.Cond, ch chan struct{}) {
	cd.L.Lock()
	for p.getVersion() == v {
		cd.Wait()
	}
	cd.L.Unlock()
	close(ch)
}

type SignalVersion struct {
	s *VersionedSignal
	v uint32
	c <-chan struct{}
}

func (p SignalVersion) HasSignal() bool {
	return p.v != p.s.getVersion()
}

func (p SignalVersion) GetCond() *sync.Cond {
	return p.s.GetCond()
}

func (p SignalVersion) Wait() {
	cd := p.GetCond()

	cd.L.Lock()
	for p.s.getVersion() == p.v {
		cd.Wait()
	}
	cd.L.Unlock()
}

func (p *SignalVersion) Channel() <-chan struct{} {
	if p.c == nil {
		p.c = p.s.getChannel(p.v)
	}
	return p.c
}
