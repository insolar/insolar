package core

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

func (p *nodeCallback) init() {
	p.fraudFactory = errors.NewFraudFactory(p.captureMisbehavior)
	p.blameFactory = errors.NewBlameFactory(p.captureMisbehavior)
}

func (p *nodeCallback) setNodeToPhaseCallback(phaseControllerCallback NodeUpdateCallback) {
	p.phaseControllerCallback = phaseControllerCallback
}

type nodeCallback struct {
	fraudFactory            errors.FraudFactory
	blameFactory            errors.BlameFactory
	phaseControllerCallback NodeUpdateCallback
}

func (p *nodeCallback) GetFraudFactory() errors.FraudFactory {
	return p.fraudFactory
}

func (p *nodeCallback) GetBlameFactory() errors.BlameFactory {
	return p.blameFactory
}

func (p *nodeCallback) captureMisbehavior(r errors.MisbehaviorReport) interface{} {
	return nil
}

func (p *nodeCallback) onTrustUpdated(n *NodeAppearance, before packets.NodeTrustLevel, after packets.NodeTrustLevel) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnTrustUpdated(n, before, after)
}

func (p *nodeCallback) onNodeStateAssigned(n *NodeAppearance) {
	if p.phaseControllerCallback == nil {
		return
	}
	p.phaseControllerCallback.OnNodeStateAssigned(n)
}
