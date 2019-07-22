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

package profiles

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"sort"
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

//func AsSortingPowerRole(np ActiveNode) uint16 {
//	st := np.GetStatic()
//	return member.SortingPowerRole(st.GetPrimaryRole(), np.GetDeclaredPower(), np.GetOpMode())
//}
//
//func AsSortingPowerRoleOfStatic(st StaticProfile, enableStartPower bool) uint16 {
//	if enableStartPower {
//		return member.SortingPowerRole(st.GetPrimaryRole(), st.GetStartPower(), member.ModeNormal)
//	}
//	return 0
//}

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
