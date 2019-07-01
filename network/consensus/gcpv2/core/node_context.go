package core

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

func (p *nodeContext) initPrep(capture errors.MisbehaviorReportFunc) {
	p.fraudFactory = errors.NewFraudFactory(capture)
	p.blameFactory = errors.NewBlameFactory(capture)
}

func (p *nodeContext) initFull(neighborTrustThreshold uint8, capture errors.MisbehaviorReportFunc) {

	p.fraudFactory = errors.NewFraudFactory(capture)
	p.blameFactory = errors.NewBlameFactory(capture)
	p.nbTrustThreshold = neighborTrustThreshold
}

func (p *nodeContext) setNodeToPhaseCallback(phaseControllerCallback NodeUpdateCallback) {
	p.phaseControllerCallback = phaseControllerCallback
}

type NodeContextHolder *nodeContext

type nodeContext struct {
	fraudFactory            errors.FraudFactory
	blameFactory            errors.BlameFactory
	nbTrustThreshold        uint8
	phaseControllerCallback NodeUpdateCallback
}

func (p *nodeContext) GetNeighbourhoodTrustThreshold() uint8 {
	if p.nbTrustThreshold == 0 {
		panic("illegal state: not allowed for PrepRealm")
	}
	return p.nbTrustThreshold
}

func (p *nodeContext) GetFraudFactory() errors.FraudFactory {
	return p.fraudFactory
}

func (p *nodeContext) GetBlameFactory() errors.BlameFactory {
	return p.blameFactory
}

func (p *nodeContext) captureMisbehavior(r errors.MisbehaviorReport) interface{} {
	return nil
}

func (p *nodeContext) onTrustUpdated(n *NodeAppearance, before packets.NodeTrustLevel, after packets.NodeTrustLevel) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnTrustUpdated(n, before, after)
}

func (p *nodeContext) onNodeStateAssigned(n *NodeAppearance) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnNodeStateAssigned(n)
}

func (p *nodeContext) onCustomEvent(n *NodeAppearance, event interface{}) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnCustomEvent(n, event)
}
