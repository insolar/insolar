// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package consensus

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

type Selection interface {
	/* When false - disables chasing timeout */
	CanBeImproved() bool
	/* This bitset only allows values of NbsConsensus[*] */
	GetConsensusVector() nodeset.ConsensusBitsetRow
}

type SelectionStrategy interface {
	CanStartVectorsEarly(consensusMembers int, countFraud int, countTrustBySome int, countTrustByNeighbors int) bool
	/* Result can be nil - it means no-decision */
	TrySelectOnAdded(globulaStats *nodeset.ConsensusStatTable, addedNode profiles.StaticProfile,
		nodeStats *nodeset.ConsensusStatRow) Selection
	SelectOnStopped(globulaStats *nodeset.ConsensusStatTable, timeIsOut bool, bftMajority int) Selection
}

type SelectionStrategyFactory interface {
	CreateSelectionStrategy(aggressivePhasing bool) SelectionStrategy
}

func NewSelection(canBeImproved bool, bitset nodeset.ConsensusBitsetRow) Selection {
	return &selectionTemplate{canBeImproved: canBeImproved, bitset: bitset}
}

type selectionTemplate struct {
	canBeImproved bool
	bitset        nodeset.ConsensusBitsetRow
}

func (c *selectionTemplate) CanBeImproved() bool {
	return c.canBeImproved
}

func (c *selectionTemplate) GetConsensusVector() nodeset.ConsensusBitsetRow {
	return c.bitset
}
