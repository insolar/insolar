// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package phasebundle

import (
	"fmt"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"
)

func newPopulationEventHandler(nodeCount int) *populationEventHandler {

	if nodeCount == 0 {
		panic("illegal value")
	}

	/* Ensure sufficient sizes of queues to avoid lockups */
	nodeCount = 1 + nodeCount*10

	return &populationEventHandler{
		make(chan population.MemberPacketSender, nodeCount),
		make(chan *population.NodeAppearance, nodeCount),
		make(chan ph2ctl.UpdateSignal, nodeCount*3),
	}
}

type populationEventHandler struct {
	qForPhase1 chan population.MemberPacketSender // can handle duplicate events
	qForPhase2 chan *population.NodeAppearance    // can NOT handle duplicate events
	qForPhase3 chan ph2ctl.UpdateSignal           // can NOT handle duplicate events
}

func (p *populationEventHandler) queueToPhase1(n population.MemberPacketSender) {
	select {
	case p.qForPhase1 <- n:
	default:
		panic("channel overflow: qForPhase1")
	}
}

func (p *populationEventHandler) queueToPhase2(n *population.NodeAppearance) {
	select {
	case p.qForPhase2 <- n:
	default:
		panic("channel overflow: qForPhase2")
	}
}

func (p *populationEventHandler) queueToPhase3(v ph2ctl.UpdateSignal) {
	select {
	case p.qForPhase3 <- v:
	default:
		panic("channel overflow: qForPhase3")
	}
}

func (p *populationEventHandler) OnPurgatoryNodeUpdate(populationVersion uint32, n population.MemberPacketSender, flags population.UpdateFlags) {

	// if flags&population.FlagCreated != 0 {
	//	p.qForPhase1 <- n
	// }

	if flags&population.FlagUpdatedProfile != 0 {
		p.queueToPhase1(n)
	}
}

func (p *populationEventHandler) OnDynamicNodeUpdate(populationVersion uint32, n *population.NodeAppearance, flags population.UpdateFlags) {

	if flags&population.FlagFixedInit != 0 {
		return // not a dynamic node
	}

	if flags&population.FlagUpdatedProfile != 0 {
		p.queueToPhase1(n)
	}
	// if flags&population.FlagUpdatedProfile != 0 {
	//	p.qForPhase3 <- ph2ctl.NewDynamicNodeReady(n)
	// }
}

func (p *populationEventHandler) OnDynamicPopulationCompleted(populationVersion uint32, indexedCount int) {
}

func (p *populationEventHandler) OnCustomEvent(populationVersion uint32, n *population.NodeAppearance, event interface{}) {
	if te, ok := event.(ph2ctl.UpdateSignal); ok && te.IsPingSignal() {
		p.queueToPhase3(te)
		return
	}
	panic(fmt.Sprintf("unknown custom event: %v", event))
}

func (p *populationEventHandler) OnTrustUpdated(populationVersion uint32, n *population.NodeAppearance,
	trustBefore, trustAfter member.TrustLevel, hasFullProfile bool) {

	// TODO ignore positive trust above member.TrustBySelf while hasFullProfile == false

	switch {
	case trustBefore == trustAfter:
		return
	case trustAfter.IsNegative():
		if !trustBefore.IsNegative() {
			p.queueToPhase3(ph2ctl.UpdateSignal{NewTrustLevel: trustAfter, UpdatedNode: n})
		}
		return
	default:
		if trustBefore == member.UnknownTrust && trustAfter >= member.TrustBySelf {
			n.UnsafeEnsureStateAvailable()
			p.queueToPhase2(n)
			p.queueToPhase3(ph2ctl.UpdateSignal{NewTrustLevel: member.TrustBySelf, UpdatedNode: n})
		}
		if trustBefore < member.TrustBySome && trustAfter >= member.TrustBySome {
			p.queueToPhase3(ph2ctl.UpdateSignal{NewTrustLevel: member.TrustBySome, UpdatedNode: n})
		}
		if trustBefore < member.TrustByNeighbors && trustAfter >= member.TrustByNeighbors {
			p.queueToPhase3(ph2ctl.UpdateSignal{NewTrustLevel: member.TrustByNeighbors, UpdatedNode: n})
		}
	}
}

func (p *populationEventHandler) OnNodeStateAssigned(populationVersion uint32, n *population.NodeAppearance) {
}
