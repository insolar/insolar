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

package census

import (
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/census.OfflinePopulation -o . -s _mock.go -g

type OfflinePopulation interface {
	FindRegisteredProfile(identity endpoints.Inbound) profiles.Host
	// FindPulsarProfile(pulsarId PulsarId) PulsarProfile
}

type OnlinePopulation interface {
	FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode
	/* indicates that this population was built without issues */
	IsValid() bool
	IsClean() bool

	GetIndexedCount() int
	/*
		Indicates the expected size of population.
		When IsValid()==true then GetIndexedCapacity() == GetIndexedCount(), otherwise GetIndexedCapacity() >= GetIndexedCount()
	*/
	GetIndexedCapacity() int

	GetSuspendedCount() int
	GetMistrustedCount() int

	/*
		It returns nil for (1) PrimaryRoleInactive, (2) for roles without any members, working or idle
	*/
	GetRolePopulation(role member.PrimaryRole) RolePopulation

	/*
		It returns a list of idle members, irrelevant of their roles. Will return nil when !IsValid().
		The returned slice never contains nil.
	*/
	GetIdleProfiles() []profiles.ActiveNode

	/* returns a total number of idle members in the population, irrelevant of their roles */
	GetIdleCount() int

	/* returns a sorted (starting from the highest role) list of roles with non-idle members  */
	GetWorkingRoles() []member.PrimaryRole

	/* when !IsValid() resulting slice will contain nil for missing members, when GetIndexedCapacity() > GetIndexedCount() */
	GetProfiles() []profiles.ActiveNode

	/* can be nil when !IsValid() */
	GetLocalProfile() profiles.LocalNode
}

type RecoverableErrorTypes uint32

const EmptyPopulation RecoverableErrorTypes = 0

const (
	External RecoverableErrorTypes = 1 << iota
	EmptySlot
	IllegalRole
	IllegalMode
	IllegalIndex
	DuplicateIndex
	BriefProfile
	DuplicateID
	IllegalSorting
	MissingSelf
)

func (v RecoverableErrorTypes) String() string {
	b := strings.Builder{}
	b.WriteRune('[')
	appendByBit(&b, &v, "External")
	appendByBit(&b, &v, "EmptySlot")
	appendByBit(&b, &v, "IllegalRole")
	appendByBit(&b, &v, "IllegalMode")
	appendByBit(&b, &v, "IllegalIndex")
	appendByBit(&b, &v, "DuplicateIndex")
	appendByBit(&b, &v, "BriefProfile")
	appendByBit(&b, &v, "DuplicateID")
	appendByBit(&b, &v, "IllegalSorting")
	appendByBit(&b, &v, "MissingSelf")
	b.WriteRune(']')

	return b.String()
}

func appendByBit(b *strings.Builder, v *RecoverableErrorTypes, s string) {
	if *v&1 != 0 {
		b.WriteString(s)
		b.WriteByte(' ')
	}
	*v >>= 1
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation -o . -s _mock.go -g

type EvictedPopulation interface {
	/* when the relevant online population is !IsValid() then not all nodes can be accessed by nodeID */
	FindProfile(nodeID insolar.ShortNodeID) profiles.EvictedNode
	GetCount() int
	/* slice will never contain nil. when the relevant online population is !IsValid() then it will also include erroneous nodes */
	GetProfiles() []profiles.EvictedNode

	IsValid() bool
	GetDetectedErrors() RecoverableErrorTypes
}

type PopulationBuilder interface {
	GetCount() int
	// SetCapacity
	AddProfile(intro profiles.StaticProfile) profiles.Updatable
	RemoveProfile(nodeID insolar.ShortNodeID)
	GetUnorderedProfiles() []profiles.Updatable
	FindProfile(nodeID insolar.ShortNodeID) profiles.Updatable
	GetLocalProfile() profiles.Updatable
	RemoveOthers()
}

type RolePopulation interface {
	/* != PrimaryRoleInactive */
	GetPrimaryRole() member.PrimaryRole
	/* true when the relevant population is valid and there are some members in this role */
	IsValid() bool
	/* total power of all members in this role */
	GetWorkingPower() uint32
	/* total number of working members in this role */
	GetWorkingCount() int
	/* number of idle members in this role */
	GetIdleCount() int

	/* list of working members in this role, will return nil when !IsValid() or GetWorkingPower()==0 */
	GetProfiles() []profiles.ActiveNode

	/*
		Returns a member (assigned) that can be assigned to to a task with the given (metric).
		It does flat distribution (all members with non-zero power are considered of same weight).

		If a default distribution falls the a member given as excludeID, then such member will be returned as (excluded)
		and the function will provide an alternative member as (assigned).

		When it was not possible to provide an alternative member then the same member will be returned as (assigned) and (excluded).

		When population is empty or invalid, then (nil, nil) is returned.
	*/
	GetAssignmentByCount(metric uint64, excludeID insolar.ShortNodeID) (assigned, excluded profiles.ActiveNode)
	/*
		Similar to GetAssignmentByCount, but it does weighed distribution across non-zero power members based on member's power.
	*/
	GetAssignmentByPower(metric uint64, excludeID insolar.ShortNodeID) (assigned, excluded profiles.ActiveNode)
}
