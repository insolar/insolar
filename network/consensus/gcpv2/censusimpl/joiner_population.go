// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func NewJoinerPopulation(localNode profiles.StaticProfile, vf cryptkit.SignatureVerifierFactory) OneJoinerPopulation {
	localNode.GetStaticNodeID()

	verifier := vf.CreateSignatureVerifierWithPKS(localNode.GetPublicKeyStore())
	return OneJoinerPopulation{
		localNode: updatableSlot{
			NodeProfileSlot: NewJoinerProfile(localNode, verifier),
		},
	}
}

var _ census.OnlinePopulation = &OneJoinerPopulation{}

type OneJoinerPopulation struct {
	localNode updatableSlot
}

func (c *OneJoinerPopulation) GetSuspendedCount() int {
	return 0
}

func (c *OneJoinerPopulation) GetMistrustedCount() int {
	return 0
}

func (c *OneJoinerPopulation) GetIdleProfiles() []profiles.ActiveNode {
	return nil
}

func (c *OneJoinerPopulation) GetIdleCount() int {
	return 0
}

func (c *OneJoinerPopulation) GetIndexedCount() int {
	return 0 // joiner is not counted
}

func (c *OneJoinerPopulation) GetIndexedCapacity() int {
	return 0 // joiner is not counted
}

func (c *OneJoinerPopulation) IsValid() bool {
	return true
}

func (c *OneJoinerPopulation) IsClean() bool {
	return c.localNode.GetOpMode().IsClean()
}

func (c *OneJoinerPopulation) GetRolePopulation(role member.PrimaryRole) census.RolePopulation {
	return nil
}

func (c *OneJoinerPopulation) GetWorkingRoles() []member.PrimaryRole {
	return nil
}

func (c *OneJoinerPopulation) copyTo(p copyFromPopulation) {
	v := []updatableSlot{c.localNode}
	v[0].index = 0 // removes Joiner status

	p.makeCopyOf(v, &v[0])
}

func (c *OneJoinerPopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.ActiveNode {
	if c.localNode.GetNodeID() != nodeID {
		return nil
	}
	return &c.localNode
}

func (c *OneJoinerPopulation) GetProfiles() []profiles.ActiveNode {
	return []profiles.ActiveNode{}
}

func (c *OneJoinerPopulation) GetLocalProfile() profiles.LocalNode {
	return &c.localNode.NodeProfileSlot
}
