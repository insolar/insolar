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

package censusimpl

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

type copyToPopulation interface {
	copyTo(p copyFromPopulation)
}

type copyFromPopulation interface {
	makeCopyOf(slots []updatableSlot, local *updatableSlot)
}

func NewManyNodePopulation(nodes []profiles.StaticProfile, localID insolar.ShortNodeID,
	vf cryptkit.SignatureVerifierFactory) ManyNodePopulation {

	r := ManyNodePopulation{}
	r.makeOfProfiles(nodes, localID, vf)
	return r
}

var _ copyToPopulation = &ManyNodePopulation{}
var _ census.OnlinePopulation = &ManyNodePopulation{}

type ManyNodePopulation struct {
	slots             []updatableSlot
	slotByID          map[insolar.ShortNodeID]*updatableSlot
	local             *updatableSlot
	roles             []roleRecord
	workingRoles      []member.PrimaryRole
	assignedSlotCount uint16
	isInvalid         bool
	suspendedCount    uint16
	mistrustedCount   uint16
}

func (c ManyNodePopulation) String() string {
	b := strings.Builder{}
	if c.isInvalid {
		b.WriteString("invalid ")
	}
	if c.local == nil {
		b.WriteString("local:<nil> ")
	} else {
		b.WriteString(fmt.Sprintf("local:%d ", c.local.GetNodeID()))
	}
	if c.suspendedCount > 0 {
		b.WriteString(fmt.Sprintf("susp:%d ", c.suspendedCount))
	}
	if c.mistrustedCount > 0 {
		b.WriteString(fmt.Sprintf("mistr:%d ", c.mistrustedCount))
	}
	if len(c.slots) == len(c.slotByID) && len(c.slots) == int(c.assignedSlotCount) {
		b.WriteString(fmt.Sprintf("profiles:%d[", c.assignedSlotCount))
	} else {
		b.WriteString(fmt.Sprintf("profiles:%d/%d/%d[", c.assignedSlotCount, len(c.slots), len(c.slotByID)))
	}
	if len(c.slots) < 50 {
		for _, slot := range c.slots {
			if slot.IsEmpty() {
				b.WriteString(" ____ ")
				continue
			}

			id := slot.GetNodeID()
			switch {
			case slot.IsJoiner():
				b.WriteString(fmt.Sprintf("+%04d ", id))
			case slot.mode.IsEvictedGracefully():
				b.WriteString(fmt.Sprintf("-%04d ", id))
			case slot.mode.IsEvicted():
				b.WriteString(fmt.Sprintf("!%04d ", id))
			case slot.mode.IsMistrustful():
				b.WriteString(fmt.Sprintf("?%04d ", id))
			case slot.mode.IsSuspended():
				b.WriteString(fmt.Sprintf("s%04d ", id))
			default:
				b.WriteString(fmt.Sprintf(" %04d ", id))
			}
		}
	} else {
		b.WriteString("too many")
	}
	b.WriteRune(']')
	return b.String()
}

func (c *ManyNodePopulation) GetSuspendedCount() int {
	return int(c.suspendedCount)
}

func (c *ManyNodePopulation) GetMistrustedCount() int {
	return int(c.mistrustedCount)
}

func (c *ManyNodePopulation) GetIdleProfiles() []profiles.ActiveNode {
	if len(c.roles) == 0 {
		return nil
	}
	return c.roles[member.PrimaryRoleInactive].GetProfiles()
}

func (c *ManyNodePopulation) GetIdleCount() int {
	if len(c.roles) == 0 {
		return 0
	}
	return int(c.roles[member.PrimaryRoleInactive].roleCount)
}

func (c *ManyNodePopulation) GetIndexedCount() int {
	return int(c.assignedSlotCount)
}

func (c *ManyNodePopulation) GetIndexedCapacity() int {
	return len(c.slots)
}

func (c *ManyNodePopulation) IsValid() bool {
	return !c.isInvalid
}

func (c *ManyNodePopulation) IsClean() bool {
	return !c.isInvalid && c.suspendedCount == 0 && c.mistrustedCount == 0 && c.local.GetOpMode().IsClean()
}

func (c *ManyNodePopulation) GetRolePopulation(role member.PrimaryRole) census.RolePopulation {
	if role == member.PrimaryRoleInactive || int(role) >= len(c.workingRoles) {
		return nil
	}
	if c.roles[role].container == nil && c.roles[role].idleCount == 0 {
		return nil
	}
	return &c.roles[role]
}

func (c *ManyNodePopulation) GetWorkingRoles() []member.PrimaryRole {
	return append(make([]member.PrimaryRole, 0, len(c.workingRoles)), c.workingRoles...)
}

func (c *ManyNodePopulation) copyTo(p copyFromPopulation) {
	p.makeCopyOf(c.slots, c.local)
}

type RecoverableReport func(e census.RecoverableErrorTypes, msg string, args ...interface{})

func panicOnRecoverable(e census.RecoverableErrorTypes, msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

// TODO it needs extensive testing on detection/tolerance when an invalid population is provided
func (c *ManyNodePopulation) makeCopyOfMapAndSeparateEvicts(slots map[insolar.ShortNodeID]*updatableSlot,
	local *updatableSlot, fail RecoverableReport) []*updatableSlot {

	if fail == nil {
		fail = panicOnRecoverable
	}

	if len(slots) == 0 {
		c.isInvalid = true
		fail(census.EmptyPopulation, "empty node population")
		return nil
	}

	localID := local.GetNodeID()

	evicts, slotCount := c._filterAndFillInSlots(slots, fail)
	c._fillInRoleStatsAndMap(localID, slotCount, true, false, fail)
	evicts = c._adjustSlotsAndCopyEvicts(localID, evicts)

	return evicts
}

func (c *ManyNodePopulation) _filterAndFillInSlots(slots map[insolar.ShortNodeID]*updatableSlot,
	fail RecoverableReport) ([]*updatableSlot, int) {

	if len(slots) > member.MaxNodeIndex {
		panic("too many nodes")
	}

	c.slots = make([]updatableSlot, len(slots))
	evicts := make([]*updatableSlot, 0, len(slots))

	slotCount := 0
	for id, vv := range slots {
		if vv == nil || vv.IsEmpty() || id == insolar.AbsentShortNodeID {
			c.isInvalid = true
			fail(census.EmptySlot, "invalid slot: id:%d", id)
			continue
		}
		switch {
		case vv.GetPrimaryRole() == member.PrimaryRoleInactive:
			c.isInvalid = true
			fail(census.IllegalRole, "invalid role: id:%d", id)
		case vv.IsJoiner():
			c.isInvalid = true
			fail(census.IllegalMode, "invalid mode: id:%d joiner", id)
		case vv.mode.IsEvicted():
			//
		case int(vv.index) /* avoid panic */ >= len(c.slots):
			c.isInvalid = true
			fail(census.IllegalIndex, "index out of bound: id:%d %d", id, vv.index)
		case c.slots[vv.index].StaticProfile != nil:
			c.isInvalid = true
			fail(census.DuplicateIndex, "duplicate index: id:%d %d", id, vv.index)
		default:
			if vv.GetExtension() == nil {
				c.isInvalid = true
				fail(census.BriefProfile, "incomplete index: id:%d %d", id, vv.StaticProfile)
			}
			c.slots[vv.index] = *vv
			slotCount++
			continue
		}
		evicts = append(evicts, vv)
	}

	if slotCount != 0 {
		c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, slotCount)
	}

	return evicts, slotCount
}

func (c *ManyNodePopulation) _fillInRoleStatsAndMap(localID insolar.ShortNodeID, slotCount int,
	compactIndex bool, checkUniqueID bool, fail RecoverableReport) {

	if slotCount > 0 {
		if slotCount > member.MaxNodeIndex {
			panic("too many nodes")
		}
		c.roles = make([]roleRecord, member.PrimaryRoleCount)
	}

	lastRole := member.PrimaryRole(0xFF)
	j := member.Index(0)
	for i := range c.slots {
		if slotCount == 0 {
			break
		}

		if c.slots[i].IsEmpty() {
			if !compactIndex {
				j++
			}
			continue
		}

		vv := &c.slots[j]
		if i != int(j) {
			*vv = c.slots[i]
			vv.index = j
		}

		nodeID := vv.GetNodeID()
		if checkUniqueID && c.slotByID[nodeID] != nil {
			// NB! this flag is only used when strictChecks == true, it should panic always
			c.isInvalid = true
			fail(census.DuplicateID, "duplicate ShortNodeID: id:%d idx:%d", nodeID, i)
		}
		c.slotByID[nodeID] = vv

		role := vv.GetPrimaryRole()
		if role == member.PrimaryRoleInactive {
			c.isInvalid = true
			fail(census.IllegalRole, "invalid role: id:%d idx:%d", nodeID, i)
		}

		if vv.power == 0 || vv.mode.IsPowerless() || role == member.PrimaryRoleInactive {
			c.roles[role].idleCount++
			if c.roles[member.PrimaryRoleInactive].container == nil {
				c.roles[member.PrimaryRoleInactive].container = c
				c.roles[member.PrimaryRoleInactive].firstNode = j.AsUint16()
			}
			lastRole = member.PrimaryRoleInactive
		} else {
			if lastRole < role {
				c.isInvalid = true
				fail(census.IllegalSorting, "invalid population order: id:%d idx:%d prev:%v this:%v", nodeID, i, lastRole, role)
			}

			if c.roles[role].role == member.PrimaryRoleInactive {
				c.roles[role].container = c
				c.roles[role].role = role
				c.roles[role].firstNode = j.AsUint16()
				c.workingRoles = append(c.workingRoles, role)
			}
			c.roles[role].roleCount++
			c.roles[role].rolePower += uint32(vv.power.ToLinearValue())

			lastRole = role
		}

		if vv.mode.IsSuspended() {
			c.suspendedCount++
		}
		if vv.mode.IsMistrustful() {
			c.mistrustedCount++
		}

		slotCount--
		j++
	}
	c.assignedSlotCount = uint16(j)

	for i := range c.roles {
		if c.roles[i].role != member.PrimaryRoleInactive {
			c.roles[i].prepare()
		}
	}

	c.local = c.slotByID[localID]
	if c.local == nil {
		c.isInvalid = true
		fail(census.MissingSelf, "missing self: id:%d", localID)
	}
}

func (c *ManyNodePopulation) _adjustSlotsAndCopyEvicts(localID insolar.ShortNodeID, evicts []*updatableSlot) []*updatableSlot {

	evictCopies := c.slots[c.assignedSlotCount:] // reuse remaining capacity for copies of evicts
	if c.assignedSlotCount == 0 {
		c.slots = nil
	} else {
		c.slots = c.slots[:c.assignedSlotCount]
	}
	if len(evictCopies) < len(evicts) {
		evictCopies = make([]updatableSlot, len(evicts))
	}

	for i := range evicts {
		evictCopies[i] = *evicts[i]
		evicts[i] = &evictCopies[i]
		if c.local == nil && evictCopies[i].GetNodeID() == localID {
			c.local = &evictCopies[i]
		}
	}

	return evicts
}

func (c *ManyNodePopulation) makeOfProfiles(nodes []profiles.StaticProfile, localNodeID insolar.ShortNodeID,
	vf cryptkit.SignatureVerifierFactory) {

	/*
			Sorting of nodes aren't necessary here as they all will be zero power, and in this case ordering is ignored
		    by internal procedures of ManyNodePopulation
	*/

	if len(nodes) == 0 {
		panic("empty node population")
	}

	buf := make([]updatableSlot, len(nodes)) // local node MUST be on the list
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(nodes))

	for i, n := range nodes {
		if n.GetStaticNodeID().IsAbsent() {
			panic("illegal value")
		}
		verifier := vf.CreateSignatureVerifierWithPKS(n.GetPublicKeyStore())
		buf[i].NodeProfileSlot = NewNodeProfile(member.Index(i), n, verifier, 0) // Power MUST BE zero, index will be assigned later
	}
	c.slots = buf
	c._fillInRoleStatsAndMap(localNodeID, len(c.slots), false, true, panicOnRecoverable)
}

func (c *ManyNodePopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode {
	slot := c.slotByID[nodeID]
	if slot == nil {
		return nil
	}
	return &slot.NodeProfileSlot
}

func (c *ManyNodePopulation) GetProfiles() []profiles.ActiveNode {
	r := make([]profiles.ActiveNode, len(c.slots))
	for i := range c.slots {
		r[i] = &c.slots[i].NodeProfileSlot
	}
	return r
}

func (c *ManyNodePopulation) GetLocalProfile() profiles.LocalNode {
	return c.local
}

func (c *ManyNodePopulation) setInvalid() {
	c.isInvalid = true
}

type DynamicPopulation struct {
	slotByID map[insolar.ShortNodeID]*updatableSlot
	local    *updatableSlot
}

func NewDynamicPopulationCopySelf(src copyToPopulation) DynamicPopulation {
	r := DynamicPopulation{}
	src.copyTo(&r)
	return r
}

func (c *DynamicPopulation) makeCopyOf(slots []updatableSlot, local *updatableSlot) {
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(slots))
	v := *local
	v.index = 0
	c.local = &v
	c.slotByID[v.GetNodeID()] = c.local
}

func (c *DynamicPopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode {
	return &c.slotByID[nodeID].NodeProfileSlot
}

func (c *DynamicPopulation) FindUpdatableProfile(nodeID insolar.ShortNodeID) profiles.Updatable {
	return c.slotByID[nodeID]
}

func (c *DynamicPopulation) GetCount() int {
	return len(c.slotByID)
}

type LessFunc func(c profiles.ActiveNode, o profiles.ActiveNode) bool

func (c *DynamicPopulation) GetProfiles() []profiles.ActiveNode {
	r := make([]profiles.ActiveNode, len(c.slotByID))
	for _, v := range c.slotByID {
		idx := v.GetIndex()
		if r[idx] != nil {
			panic(fmt.Sprintf("duplicate index: %v", idx))
		}
		r[idx] = &v.NodeProfileSlot
	}
	return r
}

func (c *DynamicPopulation) GetUnorderedProfiles() []profiles.Updatable {
	r := make([]profiles.Updatable, len(c.slotByID))
	idx := 0
	for _, v := range c.slotByID {
		r[idx] = v
		idx++
	}
	return r
}

func (c *DynamicPopulation) getUnorderedSlots() []*updatableSlot {
	r := make([]*updatableSlot, len(c.slotByID))
	idx := 0
	for _, v := range c.slotByID {
		r[idx] = v
		idx++
	}
	return r
}

func (c *DynamicPopulation) GetLocalProfile() profiles.LocalNode {
	return c.local
}

func (c *DynamicPopulation) CopyAndSeparate(forceInvalid bool, report RecoverableReport) (*ManyNodePopulation, census.EvictedPopulation) {

	r := ManyNodePopulation{}
	var handler RecoverableReport
	var issues census.RecoverableErrorTypes

	if report != nil {
		handler = func(e census.RecoverableErrorTypes, msg string, args ...interface{}) {
			issues |= e
			report(e, msg, args)
		}
	}
	evicts := r.makeCopyOfMapAndSeparateEvicts(c.slotByID, c.local, handler)
	if forceInvalid {
		r.setInvalid()
		issues |= census.External
	}
	evPop := newEvictedPopulation(evicts, issues)
	return &r, &evPop
}

func (c *DynamicPopulation) AddProfile(n profiles.StaticProfile) profiles.Updatable {
	id := n.GetStaticNodeID()
	if _, ok := c.slotByID[id]; ok {
		panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
	}
	v := updatableSlot{NewNodeProfile(0, n, nil, 0), 0}
	c.slotByID[id] = &v
	return &v
}

func (c *DynamicPopulation) RemoveProfile(id insolar.ShortNodeID) {
	delete(c.slotByID, id)
}

func (c *DynamicPopulation) RemoveOthers() {
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot)
	c.slotByID[c.local.GetNodeID()] = c.local
}

type slotSorter struct {
	values []*updatableSlot
	lessFn LessFunc
}

func (c *slotSorter) Len() int {
	return len(c.values)
}

func (c *slotSorter) Less(i, j int) bool {
	return c.lessFn(c.values[i], c.values[j])
}

func (c *slotSorter) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}

var _ census.RolePopulation = &roleRecord{}

type roleRecord struct {
	container      *ManyNodePopulation
	role           member.PrimaryRole
	rolePower      uint32
	firstNode      uint16
	roleCount      uint16
	idleCount      uint16
	powerPositions []unitizedPowerPosition
}

type unitizedPowerPosition struct {
	powerStartsAt uint32
	indexStartsAt uint16
	powerUnit     uint16 // linear power
	unitCount     uint16
}

func (p *roleRecord) prepare() {
	if p.container == nil || len(p.powerPositions) > 0 {
		panic("illegal state")
	}
	if p.rolePower == 0 {
		return
	}

	roleSlots := p.container.slots[p.firstNode : p.firstNode+p.roleCount]
	p.powerPositions = make([]unitizedPowerPosition, len(roleSlots))

	powerPosition := uint32(0)
	lastPowerUnit := uint16(math.MaxUint16)
	lastPosition := -1

	for i := range roleSlots {
		slotPower := roleSlots[i].power.ToLinearValue()
		if lastPowerUnit != slotPower {
			if lastPowerUnit < slotPower {
				panic("illegal state")
			}
			lastPowerUnit = slotPower
			lastPosition++
			p.powerPositions[lastPosition].powerUnit = slotPower
			p.powerPositions[lastPosition].powerStartsAt = powerPosition
			p.powerPositions[lastPosition].indexStartsAt = uint16(i)
		}
		p.powerPositions[lastPosition].unitCount++
		powerPosition += uint32(slotPower)
	}
	lastPosition++

	p.powerPositions = p.powerPositions[:lastPosition]
	if lastPosition > 10 && lastPosition < cap(p.powerPositions)>>1 {
		p.powerPositions = append(make([]unitizedPowerPosition, 0, lastPosition), p.powerPositions...)
	}

	if p.rolePower != powerPosition {
		panic("illegal state")
	}
}

func (p *roleRecord) IsValid() bool {
	return p.container != nil && p.container.IsValid()
}

func (p *roleRecord) GetPrimaryRole() member.PrimaryRole {
	return p.role
}

func (p *roleRecord) GetWorkingPower() uint32 {
	return p.rolePower
}

func (p *roleRecord) GetWorkingCount() int {
	return int(p.roleCount)
}

func (p *roleRecord) GetIdleCount() int {
	return int(p.idleCount)
}

func (p *roleRecord) GetProfiles() []profiles.ActiveNode {
	if !p.IsValid() {
		panic("illegal state")
	}
	if p.roleCount == 0 {
		return nil
	}
	nodes := make([]profiles.ActiveNode, p.roleCount)
	for i := range nodes {
		nodes[i] = p.getByIndex(uint16(i))
	}
	return nodes
}

func (p *roleRecord) GetAssignmentByPower(metric uint64,
	excludeID insolar.ShortNodeID) (assigned, excluded profiles.ActiveNode) {

	if p.roleCount == 0 || p.rolePower == 0 || !p.IsValid() {
		return nil, nil
	}
	if p.roleCount == 1 {
		assigned = p.getByIndex(0)
		if assigned.GetNodeID() == excludeID {
			return assigned, assigned
		}
		return assigned, nil
	}

	selector0 := uint32(metric % uint64(p.rolePower))
	assigned = p.getByIndex(p.getIndexByPower(selector0))
	if assigned.GetNodeID() != excludeID {
		return assigned, nil
	}
	excluded = assigned

	excludedPower := uint32(excluded.GetDeclaredPower().ToLinearValue())
	selector1 := uint32(metric%uint64(p.rolePower-excludedPower)) + selector0 + 1
	if selector1 >= p.rolePower {
		selector1 -= p.rolePower
	}
	assigned = p.getByIndex(p.getIndexByPower(selector1))
	if assigned.GetNodeID() != excludeID {
		return assigned, excluded
	}
	panic("not possible")
}

func (p *roleRecord) GetAssignmentByCount(metric uint64,
	excludeID insolar.ShortNodeID) (assigned, excluded profiles.ActiveNode) {

	if p.roleCount == 0 || !p.IsValid() {
		return nil, nil
	}
	if p.roleCount == 1 {
		assigned = p.getByIndex(0)
		if assigned.GetNodeID() == excludeID {
			return assigned, assigned
		}
		return assigned, nil
	}

	selector0 := uint16(metric % uint64(p.roleCount))
	assigned = p.getByIndex(selector0)
	if assigned.GetNodeID() != excludeID {
		return assigned, nil
	}
	excluded = assigned

	selector1 := uint16(metric%uint64(p.roleCount-1)) + selector0 + 1
	if selector1 >= p.roleCount {
		selector1 -= p.roleCount
	}
	assigned = p.getByIndex(selector1)
	if assigned.GetNodeID() != excludeID {
		return assigned, excluded
	}
	panic("not possible")
}

// may return garbage if used without proper checks
func (p *roleRecord) getByIndex(index uint16) profiles.ActiveNode {

	return p.container.slots[p.firstNode+index].AsActiveNode()
}

// may return garbage if used without proper checks
func (p *roleRecord) getIndexByPower(powerPosition uint32) uint16 {

	foldedPos := sort.Search(len(p.powerPositions),
		func(i int) bool { return p.powerPositions[i].powerStartsAt >= powerPosition })

	pp := p.powerPositions[foldedPos]

	if pp.unitCount == 1 {
		return pp.indexStartsAt
	}
	return pp.indexStartsAt + uint16((powerPosition-pp.powerStartsAt)/uint32(pp.powerUnit))
}
