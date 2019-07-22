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

package inspectors

import (
	"context"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
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

func (*bypassVectorInspector) PrepareForInspection(ctx context.Context) {
	panic("illegal state")
}

func (*bypassVectorInspector) CreateVector(cryptkit.DigestSigner) statevector.Vector {
	panic("illegal state")
}

func (*bypassVectorInspector) InspectVector(ctx context.Context, sender *core.NodeAppearance, customOptions uint32,
	otherData statevector.Vector) InspectedVector {

	return &bypassVector{sender, customOptions, otherData}
}

func (*bypassVectorInspector) GetBitset() member.StateBitset {
	panic("illegal state")
}

type bypassVector struct {
	n             *core.NodeAppearance
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

func (p *bypassVector) GetNode() *core.NodeAppearance {
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
