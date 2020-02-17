// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package profiles

import (
	"sort"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

func AsRank(np ActiveNode, count member.Index) member.Rank {
	if np.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(np.GetOpMode(), np.GetDeclaredPower(), np.GetIndex(), count)
}

func AsSortingRank(np ActiveNode) member.SortingRank {
	return member.NewSortingRank(np.GetNodeID(), np.GetStatic().GetPrimaryRole(), np.GetDeclaredPower(), np.GetOpMode())
}

func AsSortingRankOfStatic(st StaticProfile, enableStartPower bool) member.SortingRank {
	if enableStartPower {
		return member.NewSortingRank(st.GetStaticNodeID(), st.GetPrimaryRole(), st.GetStartPower(), member.ModeNormal)
	}
	return member.NewSortingRank(st.GetStaticNodeID(), st.GetPrimaryRole(), 0, member.ModeNormal)
}

// func AsSortingPowerRole(np ActiveNode) uint16 {
//	st := np.GetStatic()
//	return member.SortingPowerRole(st.GetPrimaryRole(), np.GetDeclaredPower(), np.GetOpMode())
// }
//
// func AsSortingPowerRoleOfStatic(st StaticProfile, enableStartPower bool) uint16 {
//	if enableStartPower {
//		return member.SortingPowerRole(st.GetPrimaryRole(), st.GetStartPower(), member.ModeNormal)
//	}
//	return 0
// }

func LessForActiveNodes(vN, oN ActiveNode) bool {
	return AsSortingRank(vN).Less(AsSortingRank(oN))
}

func LessForStaticProfiles(vN, oN StaticProfile, enableStartPower bool) bool {
	return AsSortingRankOfStatic(vN, enableStartPower).Less(AsSortingRankOfStatic(oN, enableStartPower))
}

func SortActiveNodes(nodes []ActiveNode) {
	sort.Sort(&sorterActiveNode{nodes})
}

func SortStaticProfiles(nodes []StaticProfile, enableStartPower bool) {
	sort.Sort(&sorterStaticProfile{nodes, enableStartPower})
}

type sorterActiveNode struct {
	values []ActiveNode
}

func (c *sorterActiveNode) Len() int {
	return len(c.values)
}

func (c *sorterActiveNode) Less(i, j int) bool {
	return LessForActiveNodes(c.values[i], c.values[j])
}

func (c *sorterActiveNode) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}

type sorterStaticProfile struct {
	values           []StaticProfile
	enableStartPower bool
}

func (c *sorterStaticProfile) Len() int {
	return len(c.values)
}

func (c *sorterStaticProfile) Less(i, j int) bool {
	return LessForStaticProfiles(c.values[i], c.values[j], c.enableStartPower)
}

func (c *sorterStaticProfile) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}
