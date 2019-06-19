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
	"fmt"
	"sort"

	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type copyOnlinePopulationTo interface {
	OnlinePopulation
	copyTo(p copyOnlinePopulation)
}

type copyOnlinePopulation interface {
	OnlinePopulation
	makeCopyOf(slots []updatableSlot, local *updatableSlot)
}

var _ OnlinePopulation = &OneNodePopulation{}

func NewOneNodePopulation(localNode common.NodeIntroProfile, verifier common2.SignatureVerifier) OneNodePopulation {
	localNode.GetShortNodeId()
	return OneNodePopulation{
		localNode: updatableSlot{
			nodeSlot: newNodeSlot(0, localNode, verifier),
		},
	}
}

func NewManyNodePopulation(localNode common.NodeIntroProfile, nodes []common.NodeIntroProfile, joiners bool) ManyNodePopulation {
	localNode.GetShortNodeId()
	r := ManyNodePopulation{}
	r.makeOfProfiles(nodes, localNode, joiners)
	return r
}

type OneNodePopulation struct {
	localNode updatableSlot
}

func (c *OneNodePopulation) Copy() ManyNodePopulation {
	r := ManyNodePopulation{}
	v := []updatableSlot{c.localNode}
	r.makeCopyOf(v, &v[0])
	return r
}

func (c *OneNodePopulation) copyTo(p copyOnlinePopulation) {
	p.makeCopyOf([]updatableSlot{c.localNode}, &c.localNode)
}

func (c *OneNodePopulation) FindProfile(nodeId common2.ShortNodeID) common.NodeProfile {
	if c.localNode.GetShortNodeId() != nodeId {
		return nil
	}
	return &c.localNode
}

func (c *OneNodePopulation) GetCount() int {
	return 1
}

func (c *OneNodePopulation) GetProfiles() []common.NodeProfile {
	return []common.NodeProfile{&c.localNode.nodeSlot}
}

func (c *OneNodePopulation) GetLocalProfile() common.LocalNodeProfile {
	return &c.localNode.nodeSlot
}

var _ copyOnlinePopulationTo = &ManyNodePopulation{}

type ManyNodePopulation struct {
	slots    []updatableSlot
	slotById map[common2.ShortNodeID]*updatableSlot
	local    *updatableSlot
}

func (c *ManyNodePopulation) copyTo(p copyOnlinePopulation) {
	p.makeCopyOf(c.slots, c.local)
}

func (c *ManyNodePopulation) makeCopyOf(slots []updatableSlot, local *updatableSlot) {
	c.slots = append(make([]updatableSlot, 0, len(slots)), slots...)
	c.slotById = make(map[common2.ShortNodeID]*updatableSlot, len(slots))

	for i := range c.slots {
		v := &c.slots[i]
		id := v.GetShortNodeId()
		if _, ok := c.slotById[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		c.slotById[id] = v
		if local.GetShortNodeId() == id {
			c.local = v
		}
	}
}

func (c *ManyNodePopulation) makeCopyOfMap(slots map[common2.ShortNodeID]*updatableSlot, local *updatableSlot, less LessFunc) {
	c.slots = append(make([]updatableSlot, len(slots)))
	c.slotById = make(map[common2.ShortNodeID]*updatableSlot, len(slots))

	if less == nil {
		for id, vv := range slots {
			idx := vv.GetIndex()
			if c.slots[idx].NodeIntroProfile != nil {
				panic(fmt.Sprintf("duplicate index: %v", idx))
			}
			c.slots[idx] = *vv
			v := &c.slots[idx]
			c.slotById[id] = v
		}
	} else {
		idx := 0
		for _, vv := range slots {
			c.slots[idx] = *vv
			idx++
		}
		sort.Sort(&slotArraySorter{values: c.slots, lessFn: less})
		for i := range c.slots {
			v := &c.slots[i]
			v.SetIndex(i)
			c.slotById[v.GetShortNodeId()] = v
		}
	}
	c.local = c.slotById[local.GetShortNodeId()]
	if c.local == nil {
		panic("illegal state")
	}
}

func (c *ManyNodePopulation) makeOfProfiles(nodes []common.NodeIntroProfile, localNode common.NodeIntroProfile, joiners bool) {
	buf := make([]updatableSlot, len(nodes)+1) // +1 local node may not be on the list
	c.slotById = make(map[common2.ShortNodeID]*updatableSlot, len(nodes)+1)

	c.local = &buf[0]
	c.local.index = 0
	c.local.NodeIntroProfile = localNode
	c.local.setJoiner(joiners)
	c.slotById[localNode.GetShortNodeId()] = c.local

	slotIndex := 1

	for _, n := range nodes {
		if n == localNode {
			continue
		}
		id := n.GetShortNodeId()
		if _, ok := c.slotById[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		buf[slotIndex].nodeSlot = newNodeSlot(slotIndex, n, nil)
		buf[slotIndex].setJoiner(joiners)
		c.slotById[id] = &buf[slotIndex]

		slotIndex++
	}
	c.slots = buf[:slotIndex]
}

func (c *ManyNodePopulation) FindProfile(nodeId common2.ShortNodeID) common.NodeProfile {
	return &c.slotById[nodeId].nodeSlot
}

func (c *ManyNodePopulation) GetCount() int {
	return len(c.slots)
}

func (c *ManyNodePopulation) GetProfiles() []common.NodeProfile {
	r := make([]common.NodeProfile, len(c.slots))
	for i := range c.slots {
		r[i] = &c.slots[i].nodeSlot
	}
	return r
}

func (c *ManyNodePopulation) GetLocalProfile() common.LocalNodeProfile {
	return c.local
}

func (c *ManyNodePopulation) Copy() ManyNodePopulation {
	r := ManyNodePopulation{}
	r.makeCopyOf(c.slots, c.local)
	return r
}

var _ OnlinePopulation = &DynamicPopulation{}

type DynamicPopulation struct {
	slotById map[common2.ShortNodeID]*updatableSlot
	local    *updatableSlot
}

func NewDynamicPopulation(src copyOnlinePopulationTo) DynamicPopulation {
	r := DynamicPopulation{}
	src.copyTo(&r)
	return r
}

func (c *DynamicPopulation) makeCopyOf(slots []updatableSlot, local *updatableSlot) {
	c.slotById = make(map[common2.ShortNodeID]*updatableSlot, len(slots))

	for i := range slots {
		v := slots[i]
		id := v.GetShortNodeId()
		if _, ok := c.slotById[id]; ok {
			panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
		}
		c.slotById[id] = &v
	}
	c.local = c.slotById[local.GetShortNodeId()]
	if c.local == nil {
		panic("illegal state")
	}
}

func (c *DynamicPopulation) FindProfile(nodeId common2.ShortNodeID) common.NodeProfile {
	return &c.slotById[nodeId].nodeSlot
}

func (c *DynamicPopulation) FindUpdatableProfile(nodeId common2.ShortNodeID) common.UpdatableNodeProfile {
	return c.slotById[nodeId]
}

func (c *DynamicPopulation) GetCount() int {
	return len(c.slotById)
}

func (c *DynamicPopulation) SortDefault() {
	c.Sort(common.LessForNodeProfile)
}

type LessFunc func(c common.NodeProfile, o common.NodeProfile) bool

func (c *DynamicPopulation) Sort(lessFn LessFunc) {
	sorter := slotSorter{values: c.getUnorderedSlots(), lessFn: lessFn}
	sort.Sort(&sorter)
	for i, v := range sorter.values {
		v.SetIndex(i)
	}
}

func (c *DynamicPopulation) GetProfiles() []common.NodeProfile {
	r := make([]common.NodeProfile, len(c.slotById))
	for _, v := range c.slotById {
		idx := v.GetIndex()
		if r[idx] != nil {
			panic(fmt.Sprintf("duplicate index: %v", idx))
		}
		r[idx] = &v.nodeSlot
	}
	return r
}

func (c *DynamicPopulation) GetUnorderedProfiles() []common.UpdatableNodeProfile {
	r := make([]common.UpdatableNodeProfile, len(c.slotById))
	idx := 0
	for _, v := range c.slotById {
		r[idx] = v
		idx++
	}
	return r
}

func (c *DynamicPopulation) getUnorderedSlots() []*updatableSlot {
	r := make([]*updatableSlot, len(c.slotById))
	idx := 0
	for _, v := range c.slotById {
		r[idx] = v
		idx++
	}
	return r
}

func (c *DynamicPopulation) GetLocalProfile() common.LocalNodeProfile {
	return c.local
}

func (c *DynamicPopulation) CopyAndSort(less LessFunc) ManyNodePopulation {
	if less == nil {
		panic("lessFunc is nil")
	}
	r := ManyNodePopulation{}
	r.makeCopyOfMap(c.slotById, c.local, less)
	return r
}

func (c *DynamicPopulation) CopyAndSortDefault() ManyNodePopulation {
	r := ManyNodePopulation{}
	r.makeCopyOfMap(c.slotById, c.local, common.LessForNodeProfile)
	return r
}

func (c *DynamicPopulation) CopyUnsorted() ManyNodePopulation {
	r := ManyNodePopulation{}
	r.makeCopyOfMap(c.slotById, c.local, nil)
	return r
}

func (c *DynamicPopulation) AddJoinerProfile(n common.NodeIntroProfile) common.UpdatableNodeProfile {
	id := n.GetShortNodeId()
	if _, ok := c.slotById[id]; ok {
		panic(fmt.Sprintf("duplicate ShortNodeID: %v", id))
	}
	v := updatableSlot{newNodeSlot(0 /* force later collision without sorting */, n, nil)}
	c.slotById[id] = &v
	return &v
}

func (c *DynamicPopulation) RemoveProfile(id common2.ShortNodeID) {
	delete(c.slotById, id)
}

var _ sort.Interface = &slotSorter{}

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

var _ sort.Interface = &slotArraySorter{}

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
