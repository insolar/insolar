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

package consensus

import (
	"github.com/insolar/insolar/network/consensus/common/consensuskit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

func NewSimpleSelectionStrategyFactory() SelectionStrategyFactory {
	return &simpleSelectionStrategyFactory{}
}

type simpleSelectionStrategyFactory struct{}

func (*simpleSelectionStrategyFactory) CreateSelectionStrategy(aggressivePhasing bool) SelectionStrategy {
	return &simpleSelectionStrategy{aggressivePhasing}
}

type simpleSelectionStrategy struct {
	aggressivePhasing bool
}

func (p *simpleSelectionStrategy) CanStartVectorsEarly(consensusMembers int, countFraud int, countTrustBySome int, countTrustByNeighbors int) bool {
	if countFraud != 0 {
		return false
	}
	if p.aggressivePhasing {
		return true
	}
	return countTrustBySome >= consensuskit.BftMajority(consensusMembers) || countTrustByNeighbors >= 1+consensusMembers>>1
}

func (*simpleSelectionStrategy) TrySelectOnAdded(globulaStats *nodeset.ConsensusStatTable, addedNode profiles.StaticProfile,
	nodeStats *nodeset.ConsensusStatRow) Selection {
	return nil
}

func (*simpleSelectionStrategy) SelectOnStopped(globulaStats *nodeset.ConsensusStatTable, timeIsOut bool, bftMajorityArg int) Selection {

	absMajority := true
	if globulaStats.RowCount() < bftMajorityArg {
		//bftMajority = uint16(consensuskit.BftMajority(globulaStats.RowCount()))
		absMajority = false
	}
	bftMajority := uint16(bftMajorityArg)

	resultSet := nodeset.NewConsensusBitsetRow(globulaStats.ColumnCount())
	for i := 0; i < resultSet.ColumnCount(); i++ {
		tc := globulaStats.GetColumn(i)
		decision := consensusDecisionOfNode(tc, absMajority, bftMajority)
		resultSet.Set(i, decision)
	}

	return NewSelection(!absMajority, resultSet)
}

func consensusDecisionOfNode(tc *nodeset.ConsensusStatColumn, absMajority bool, bftMajority uint16) nodeset.ConsensusBitsetEntry {

	switch {
	case tc.GetSummaryByValue(nodeset.ConsensusStatTrusted)+tc.GetSummaryByValue(nodeset.ConsensusStatDoubted) >= bftMajority:
		return nodeset.CbsIncluded
	case tc.GetSummaryByValue(nodeset.ConsensusStatFraud)+ //tc.GetSummaryByValue(nodeset.ConsensusStatMissingThere)+
		tc.GetSummaryByValue(nodeset.ConsensusStatFraudSuspect) >= bftMajority:
		//if absMajority {
		//	return nodeset.CbsExcluded
		//}
		return nodeset.CbsFraud
	default:
		return nodeset.CbsSuspected
	}
}
