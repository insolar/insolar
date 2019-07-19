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
	"sort"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

type copyToPopulation interface {
	copyTo(p copyFromPopulation, fullCopy bool)
}

type copyFromPopulation interface {
	makeFullCopyOf(slots []updatableSlot, local *updatableSlot)
	makeSelfCopyOf(slots []updatableSlot, local *updatableSlot)
}

var _ census.OnlinePopulation = &OneNodePopulation{}

func NewOneNodePopulation(localNode profiles.StaticProfile, verifier cryptkit.SignatureVerifier) OneNodePopulation {
	localNode.GetStaticNodeID()
	return OneNodePopulation{
		localNode: updatableSlot{
			NodeProfileSlot: NewNodeProfile(0, localNode, verifier, localNode.GetStartPower()),
		},
	}
}

func NewManyNodePopulation(localNode profiles.StaticProfile, nodes []profiles.StaticProfile) ManyNodePopulation {
	localNode.GetStaticNodeID()
	r := ManyNodePopulation{}
	r.makeOfProfiles(nodes, localNode)
	return r
}

type OneNodePopulation struct {
	localNode updatableSlot
}

func (c *OneNodePopulation) Copy() ManyNodePopulation {
	r := ManyNodePopulation{}
	v := []updatableSlot{c.localNode}
	r.makeFullCopyOf(v, &v[0])
	return r
}

func (c *OneNodePopulation) copyTo(p copyFromPopulation, fullCopy bool) {
	if fullCopy {
		p.makeFullCopyOf([]updatableSlot{c.localNode}, &c.localNode)
	} else {
		p.makeSelfCopyOf([]updatableSlot{c.localNode}, &c.localNode)
	}
}

func (c *OneNodePopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode {
	if c.localNode.GetNodeID() != nodeID {
		return nil
	}
	return &c.localNode
}

func (c *OneNodePopulation) GetCount() int {
	return 1
}

func (c *OneNodePopulation) GetProfiles() []profiles.ActiveNode {
	return []profiles.ActiveNode{&c.localNode.NodeProfileSlot}
}

func (c *OneNodePopulation) GetLocalProfile() profiles.LocalNode {
	return &c.localNode.NodeProfileSlot
}

var _ copyToPopulation = &ManyNodePopulation{}

type ManyNodePopulation struct {
	slots    []updatableSlot
	slotByID map[insolar.ShortNodeID]*updatableSlot
	local    *updatableSlot
}

func (c *ManyNodePopulation) copyTo(p copyFromPopulation, fullCopy bool) {
	if fullCopy {
		p.makeFullCopyOf(c.slots, c.local)
	} else {
		p.makeSelfCopyOf(c.slots, c.local)
	}
}

func (c *ManyNodePopulation) makeFullCopyOf(slots []updatableSlot, local *updatableSlot) {
	c.slots = append(make([]updatableSlot, 0, len(slots)), slots...)
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(slots))

	for i := range c.slots {
		v := &c.slots[i]
		id := v.GetNodeID()
		if _, ok := c.slotByID[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		c.slotByID[id] = v
		if local.GetNodeID() == id {
			c.local = v
		}
	}
}

func (c *ManyNodePopulation) makeCopyOfMapAndSeparateEvicts(slots map[insolar.ShortNodeID]*updatableSlot, local *updatableSlot) []*updatableSlot {

	var evicts []*updatableSlot
	// TODO HACK - must use vector-based ordering
	slotCount := len(slots)
	indexed := make([]*updatableSlot, slotCount)
	c.slots = make([]updatableSlot, slotCount)

	maxSlotCount := slotCount
	for _, vv := range slots {
		switch {
		case vv.IsJoiner():
			panic(fmt.Sprintf("unsorted index: %v", vv))
		case vv.GetOpMode().IsEvicted():
			maxSlotCount--
			if evicts == nil {
				evicts = make([]*updatableSlot, 0, slotCount)
			}
			c.slots[maxSlotCount] = *vv
			evicts = append(evicts, &c.slots[maxSlotCount])
		default:
			slotIndex := vv.index
			if indexed[slotIndex] != nil {
				panic("illegal state - duplicate index")
			}
			indexed[slotIndex] = vv
		}
	}

	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, maxSlotCount)
	c.slots = c.slots[:maxSlotCount]

	i := member.Index(0)
	for _, vv := range indexed {
		if vv == nil {
			continue
		}
		c.slots[i] = *vv
		c.slots[i].index = i
		c.slotByID[vv.GetNodeID()] = &c.slots[i]
		i++
	}

	c.local = c.slotByID[local.GetNodeID()]
	if c.local == nil {
		panic("illegal state")
	}

	return evicts
}

func (c *ManyNodePopulation) makeCopyOfMapAndSort(slots map[insolar.ShortNodeID]*updatableSlot, local *updatableSlot, less LessFunc) {
	c.slots = append(make([]updatableSlot, len(slots)))
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(slots))

	idx := 0
	for _, vv := range slots {
		c.slots[idx] = *vv
		idx++
	}
	sort.Sort(&slotArraySorter{values: c.slots, lessFn: less})
	for i := range c.slots {
		v := &c.slots[i]
		v.SetIndex(member.AsIndex(i))
		c.slotByID[v.GetNodeID()] = v
	}

	c.local = c.slotByID[local.GetNodeID()]
	if c.local == nil {
		panic("illegal state")
	}
}

func (c *ManyNodePopulation) makeOfProfiles(nodes []profiles.StaticProfile, localNode profiles.StaticProfile) {
	buf := make([]updatableSlot, len(nodes)+1) // +1 local node may not be on the list
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(nodes)+1)

	c.local = &buf[0]
	c.local.index = 0
	c.local.StaticProfile = localNode
	c.slotByID[localNode.GetStaticNodeID()] = c.local

	slotIndex := member.AsIndex(1)

	for _, n := range nodes {
		id := n.GetStaticNodeID()
		if id == localNode.GetStaticNodeID() {
			continue
		}
		if _, ok := c.slotByID[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		buf[slotIndex].NodeProfileSlot = NewNodeProfile(slotIndex, n, nil, 0)
		c.slotByID[id] = &buf[slotIndex]

		slotIndex++
	}
	c.slots = buf[:slotIndex]
}

func (c *ManyNodePopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode {
	slot := c.slotByID[nodeID]
	if slot == nil {
		return nil
	}
	return &slot.NodeProfileSlot
}

func (c *ManyNodePopulation) GetCount() int {
	return len(c.slots)
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

func (c *ManyNodePopulation) Copy() ManyNodePopulation {
	r := ManyNodePopulation{}
	r.makeFullCopyOf(c.slots, c.local)
	return r
}

type DynamicPopulation struct {
	slotByID map[insolar.ShortNodeID]*updatableSlot
	local    *updatableSlot
}

func NewDynamicPopulation(src copyToPopulation) DynamicPopulation {
	r := DynamicPopulation{}
	src.copyTo(&r, true)
	return r
}

func NewDynamicPopulationCopySelf(src copyToPopulation) DynamicPopulation {
	r := DynamicPopulation{}
	src.copyTo(&r, false)
	return r
}

func (c *DynamicPopulation) makeFullCopyOf(slots []updatableSlot, local *updatableSlot) {
	c.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, len(slots))

	localID := local.GetNodeID()

	for i := range slots {
		v := slots[i]
		id := v.GetNodeID()
		if _, ok := c.slotByID[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		c.slotByID[id] = &v
	}
	c.local = c.slotByID[localID]
	if c.local == nil {
		panic("illegal state")
	}
}

func (c *DynamicPopulation) makeSelfCopyOf(slots []updatableSlot, local *updatableSlot) {
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

func (c *DynamicPopulation) Sort(lessFn LessFunc) {
	sorter := slotSorter{values: c.getUnorderedSlots(), lessFn: lessFn}
	sort.Sort(&sorter)
	for i, v := range sorter.values {
		v.SetIndex(member.AsIndex(i))
	}
}

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

func (c *DynamicPopulation) CopyAndSeparate() (*ManyNodePopulation, census.EvictedPopulation) {
	r := ManyNodePopulation{}
	evicts := r.makeCopyOfMapAndSeparateEvicts(c.slotByID, c.local)
	evPop := newEvictedPopulation(evicts)
	return &r, &evPop
}

func (c *DynamicPopulation) AddProfile(n profiles.StaticProfile) profiles.Updatable {
	id := n.GetStaticNodeID()
	if _, ok := c.slotByID[id]; ok {
		panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
	}
	v := updatableSlot{NewJoinerProfile(n, nil, n.GetStartPower()), 0}
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

type slotArraySorter struct {
	values []updatableSlot
	lessFn LessFunc
}

func (c *slotArraySorter) Len() int {
	return len(c.values)
}

func (c *slotArraySorter) Less(i, j int) bool {
	return c.lessFn(&c.values[i], &c.values[j])
}

func (c *slotArraySorter) Swap(i, j int) {
	c.values[i], c.values[j] = c.values[j], c.values[i]
}
