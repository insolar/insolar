//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package population

import "github.com/insolar/insolar/network/consensus/gcpv2/api/member"

func NewPanicDispatcher(panicMsg string) EventDispatcher {
	return NewSinkDispatcher(func(EventClosureFunc) {
		panic(panicMsg)
	})
}

func NewSinkDispatcher(sink EventDispatchFunc) EventDispatcher {
	return NewInterceptorAndSink(nil, sink)
}

func NewIgnoringDispatcher() EventDispatcher {
	return NewInterceptorAndSink(nil, nil)
}

func NewInterceptor(intercept EventDispatcher) EventDispatcher {
	return NewInterceptorAndSink(intercept, nil)
}

func NewInterceptorAndSink(intercept EventDispatcher, sink EventDispatchFunc) EventDispatcher {
	return &EventWrapper{GetEventDispatchFn(intercept, sink)}
}

func GetEventDispatchFn(intercept EventDispatcher, sink EventDispatchFunc) EventDispatchFunc {
	if intercept == nil {
		if sink == nil {
			return func(fn EventClosureFunc) {}
		}
		return sink
	}
	if sink == nil {
		return func(fn EventClosureFunc) {
			fn(intercept)
		}
	}
	return func(fn EventClosureFunc) {
		fn(intercept)
		sink(fn)
	}
}

var _ EventDispatcher = &EventWrapper{}

type EventWrapper struct {
	sink EventDispatchFunc
}

func (p *EventWrapper) OnTrustUpdated(populationVersion uint32, n *NodeAppearance,
	before member.TrustLevel, after member.TrustLevel, hasFullProfile bool) {
	p.sink(func(d EventDispatcher) {
		d.OnTrustUpdated(populationVersion, n, before, after, hasFullProfile)
	})
}

func (p *EventWrapper) OnNodeStateAssigned(populationVersion uint32, n *NodeAppearance) {
	p.sink(func(d EventDispatcher) {
		d.OnNodeStateAssigned(populationVersion, n)
	})
}

func (p *EventWrapper) OnDynamicNodeUpdate(populationVersion uint32, n *NodeAppearance, flags UpdateFlags) {
	p.sink(func(d EventDispatcher) {
		d.OnDynamicNodeUpdate(populationVersion, n, flags)
	})
}

func (p *EventWrapper) OnPurgatoryNodeUpdate(populationVersion uint32, n MemberPacketSender, flags UpdateFlags) {
	p.sink(func(d EventDispatcher) {
		d.OnPurgatoryNodeUpdate(populationVersion, n, flags)
	})
}

func (p *EventWrapper) OnCustomEvent(populationVersion uint32, n *NodeAppearance, event interface{}) {
	p.sink(func(d EventDispatcher) {
		d.OnCustomEvent(populationVersion, n, event)
	})
}

func (p *EventWrapper) OnDynamicPopulationCompleted(populationVersion uint32, indexedCount int) {
	p.sink(func(d EventDispatcher) {
		d.OnDynamicPopulationCompleted(populationVersion, indexedCount)
	})
}
