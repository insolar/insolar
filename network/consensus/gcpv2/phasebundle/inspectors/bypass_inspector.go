// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package inspectors

import (
	"context"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

func NewBypassInspector() VectorInspector {
	return &bypassVectorInspector{}
}

type bypassVectorInspector struct {
}

func (*bypassVectorInspector) CreateNextPopulation(nodeset.ConsensusBitsetRow) ([]profiles.PopulationRank, proofs.CloudStateHash, proofs.GlobulaStateHash) {
	panic("illegal state")
}

func (*bypassVectorInspector) PrepareForInspection(ctx context.Context) bool {
	panic("illegal state")
}

func (*bypassVectorInspector) CreateVector(cryptkit.DigestSigner) statevector.Vector {
	panic("illegal state")
}

func (*bypassVectorInspector) InspectVector(ctx context.Context, sender *population.NodeAppearance, customOptions uint32,
	otherData statevector.Vector) InspectedVector {

	return &bypassVector{sender, customOptions, otherData}
}

func (*bypassVectorInspector) GetBitset() member.StateBitset {
	panic("illegal state")
}

type bypassVector struct {
	n             *population.NodeAppearance
	customOptions uint32
	otherData     statevector.Vector
}

func (p *bypassVector) GetCustomOptions() uint32 {
	return p.customOptions
}

func (p *bypassVector) HasSenderFault() bool {
	return false
}

func (p *bypassVector) GetInspectionResults() (*nodeset.ConsensusStatRow, nodeset.NodeVerificationResult) {
	return nil, nodeset.NvrNotVerified
}

func (p *bypassVector) GetBitset() member.StateBitset {
	return p.otherData.Bitset
}

func (p *bypassVector) GetNode() *population.NodeAppearance {
	return p.n
}

func (p *bypassVector) Reinspect(ctx context.Context, inspector VectorInspector) InspectedVector {
	iv := inspector.InspectVector(ctx, p.n, p.customOptions, p.otherData)
	if _, ok := iv.(*bypassVector); ok {
		panic("illegal state")
	}
	return iv
}

func (*bypassVector) Inspect(ctx context.Context) {
	panic("illegal state")
}

func (*bypassVector) IsInspected() bool {
	return false
}

func (*bypassVector) HasMissingMembers() bool {
	return false
}
