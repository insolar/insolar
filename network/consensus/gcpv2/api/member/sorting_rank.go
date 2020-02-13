// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package member

import (
	"github.com/insolar/insolar/insolar"
)

type SortingRank struct {
	nodeID    insolar.ShortNodeID
	powerRole uint16
}

func NewSortingRank(nodeID insolar.ShortNodeID, role PrimaryRole, pw Power, mode OpMode) SortingRank {
	return SortingRank{nodeID, SortingPowerRole(role, pw, mode)}
}

func (v SortingRank) GetNodeID() insolar.ShortNodeID {
	return v.nodeID
}

func (v SortingRank) IsWorking() bool {
	return v.powerRole != 0
}

func (v SortingRank) GetWorkingRole() PrimaryRole {
	return PrimaryRole(v.powerRole >> 8)
}

func (v SortingRank) GetPower() Power {
	return Power(v.powerRole)
}

// NB! Sorting is REVERSED
func (v SortingRank) Less(o SortingRank) bool {
	if o.powerRole < v.powerRole {
		return true
	}
	if o.powerRole > v.powerRole {
		return false
	}
	return o.nodeID < v.nodeID
}

// NB! Sorting is REVERSED
func LessByID(vNodeID, oNodeID insolar.ShortNodeID) bool {
	return oNodeID < vNodeID
}

func SortingPowerRole(role PrimaryRole, pw Power, mode OpMode) uint16 {
	if role == 0 {
		panic("illegal value")
	}
	if pw == 0 || mode.IsPowerless() {
		return 0
	}
	return uint16(role)<<8 | uint16(pw)
}
