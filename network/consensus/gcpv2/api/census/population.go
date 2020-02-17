// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
