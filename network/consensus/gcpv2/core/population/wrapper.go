// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
