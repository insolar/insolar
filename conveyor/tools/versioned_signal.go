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
	"unsafe"
)

func NewVersionedSignal() VersionedSignal {
	return VersionedSignal{}
}

type VersionedSignal struct {
	signalVersion unsafe.Pointer // atomic *SignalVersion
}

func (p *VersionedSignal) NextBroadcast() {
	sv := (*SignalVersion)(atomic.SwapPointer(&p.signalVersion, nil))
	if sv != nil {
		sv.signal()
	}
}

func (p *VersionedSignal) BroadcastAndMark() *SignalVersion {
	nsv := newSignalVersion()
	sv := (*SignalVersion)(atomic.SwapPointer(&p.signalVersion, (unsafe.Pointer)(nsv)))
	if sv != nil {
		sv.signal()
	}
	return nsv
}

func (p *VersionedSignal) Mark() *SignalVersion {
	var nsv *SignalVersion
	for {
		sv := (*SignalVersion)(atomic.LoadPointer(&p.signalVersion))
		switch {
		case sv != nil:
			return sv
		case nsv == nil: // avoid repetitive new
			nsv = newSignalVersion()
		}
		if atomic.CompareAndSwapPointer(&p.signalVersion, nil, (unsafe.Pointer)(nsv)) {
			return nsv
		}
	}
}

func newSignalVersion() *SignalVersion {
	sv := SignalVersion{}
	sv.wg.Add(1)
	return &sv
}

type signalChannel = chan struct{}

type SignalVersion struct {
	next *SignalVersion
	wg   sync.WaitGroup
	c    unsafe.Pointer // atomic *signalChannel
}

func (p *SignalVersion) signal() {
	if p.next != nil {
		p.next.signal() // older signals must fire first
	}

	var closedSignal *signalChannel // explicit type decl to avoid passing of something wrong into unsafe.Pointer conversion
	closedSignal = &closedChan

	atomic.CompareAndSwapPointer(&p.c, nil, (unsafe.Pointer)(closedSignal))
	p.wg.Done()
}

func (p *SignalVersion) Wait() {
	if p == nil {
		return
	}

	p.wg.Wait()
}

func (p *SignalVersion) ChannelIf(choice bool, def <-chan struct{}) <-chan struct{} {
	if choice {
		return p.Channel()
	}
	return def
}

func (p *SignalVersion) Channel() <-chan struct{} {
	if p == nil {
		return ClosedChannel()
	}

	var wcp *signalChannel
	for {
		sc := (*signalChannel)(atomic.LoadPointer(&p.c))
		switch {
		case sc != nil:
			return *sc
		case wcp == nil:
			wcp = new(signalChannel)
		}

		if atomic.CompareAndSwapPointer(&p.c, nil, (unsafe.Pointer)(wcp)) {
			go func() {
				p.wg.Wait()
				close(*wcp)
			}()
			return *wcp
		}
	}
}

func (p *SignalVersion) HasSignal() bool {
	if p == nil {
		return true
	}

	sc := (*signalChannel)(atomic.LoadPointer(&p.c))
	if sc == nil {
		return false
	}
	select {
	case <-*sc:
		return true
	default:
		return false
	}
}
