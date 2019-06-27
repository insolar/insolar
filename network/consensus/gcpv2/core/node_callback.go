package core

import "github.com/insolar/insolar/network/consensus/gcpv2/errors"

func (p *nodeCallback) init() {
	p.fraudFactory = errors.NewFraudFactory(p.captureMisbehavior)
	p.blameFactory = errors.NewBlameFactory(p.captureMisbehavior)
}

type nodeCallback struct {
	fraudFactory errors.FraudFactory
	blameFactory errors.BlameFactory
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
