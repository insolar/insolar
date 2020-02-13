// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package nodeset

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
)

/*
Contains copy of NodeAppearance fields that can be changed, to avoid possible racing
*/

type VectorEntryData struct {
	RequestedPower member.Power
	RequestedMode  member.OpMode
	TrustLevel     member.TrustLevel
	Profile        profiles.ActiveNode

	// Node *NodeAppearance
	// common.MembershipAnnouncement
	proofs.NodeAnnouncedState
}

type EntryApplyFunc func(nodeData VectorEntryData, postponed bool, filter uint32)
type EntryFilterFunc func(index int, nodeData VectorEntryData, parentFilter uint32) (bool, uint32)

type VectorEntryScanner interface {
	GetIndexedCount() int
	GetSortedCount() int
	ScanIndexed(apply func(index int, nodeData VectorEntryData))
	ScanSorted(apply EntryApplyFunc, filterValue uint32)
	ScanSortedWithFilter(parentFilter uint32, apply EntryApplyFunc, filter EntryFilterFunc)
}

type VectorEntryDigester interface {
	AddNext(nodeData VectorEntryData, zeroPower bool)
	ForkSequence() VectorEntryDigester
}
