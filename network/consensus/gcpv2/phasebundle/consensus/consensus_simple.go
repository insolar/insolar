// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
		// bftMajority = uint16(consensuskit.BftMajority(globulaStats.RowCount()))
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
	case tc.GetSummaryByValue(nodeset.ConsensusStatFraud)+ // tc.GetSummaryByValue(nodeset.ConsensusStatMissingThere)+
		tc.GetSummaryByValue(nodeset.ConsensusStatFraudSuspect) >= bftMajority:
		// if absMajority {
		//	return nodeset.CbsExcluded
		// }
		return nodeset.CbsFraud
	default:
		return nodeset.CbsSuspected
	}
}
