// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package population

import (
	"context"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"

	"github.com/insolar/insolar/insolar"
)

type RealmPopulation interface {
	GetIndexedCount() int
	GetJoinersCount() int
	// GetVotersCount() int

	GetSealedCapacity() (int, bool)
	SealIndexed(indexedCapacity int) bool

	GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetJoinerNodeAppearance(id insolar.ShortNodeID) *NodeAppearance
	GetNodeAppearanceByIndex(idx int) *NodeAppearance

	GetShuffledOtherNodes() []*NodeAppearance /* excludes joiners and self */
	GetIndexedNodes() []*NodeAppearance       /* no joiners included */
	GetIndexedNodesAndHasNil() ([]*NodeAppearance, bool)
	GetIndexedCountAndCompleteness() (int, bool)

	GetSelf() *NodeAppearance

	// CreateNodeAppearance(ctx context.Context, inp profiles.ActiveNode) *NodeAppearance
	AddReservation(id insolar.ShortNodeID) (bool, *NodeAppearance)
	FindReservation(id insolar.ShortNodeID) (bool, *NodeAppearance)

	AddToDynamics(ctx context.Context, n *NodeAppearance) (*NodeAppearance, error)
	GetAnyNodes(includeIndexed bool, shuffle bool) []*NodeAppearance

	CreateVectorHelper() *RealmVectorHelper
	CreatePacketLimiter(isJoiner bool) phases.PacketLimiter

	GetTrustCounts() (fraudCount, bySelfCount, bySomeCount, byNeighborsCount uint16)
	GetDynamicCounts() (briefCount, fullCount uint16)
	GetPurgatoryCounts() (addedCount, ascentCount uint16)
}
