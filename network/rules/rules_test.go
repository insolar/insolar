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

package rules

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
)

func TestRules_CheckMinRole(t *testing.T) {
	cert := testutils.NewCertificateMock(t)
	nodes := []insolar.NetworkNode{
		node.NewNode(gen.Reference(), insolar.StaticRoleHeavyMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleLightMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleLightMaterial, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "", ""),
		node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "", ""),
	}
	cert.GetMinRolesMock.Set(func() (r uint, r1 uint, r2 uint) {
		return 1, 0, 0
	})
	err := CheckMinRole(cert, nodes)
	require.NoError(t, err)

	cert.GetMinRolesMock.Set(func() (r uint, r1 uint, r2 uint) {
		return 3, 2, 4
	})
	err = CheckMinRole(cert, nodes)
	require.Error(t, err)
}

func TestRules_CheckMajorityRule(t *testing.T) {
	discNodesCount := 5
	netNodes, discoveryNodes := getDiscoveryNodes(discNodesCount)
	cert := testutils.NewCertificateMock(t)
	cert.GetDiscoveryNodesMock.Set(func() (r []insolar.DiscoveryNode) {
		return discoveryNodes
	})
	cert.GetMajorityRuleMock.Set(func() (r int) {
		return discNodesCount
	})
	count, err := CheckMajorityRule(cert, netNodes)
	require.NoError(t, err)

	require.Equal(t, discNodesCount, count)

	netNodes = netNodes[:len(netNodes)-len(netNodes)/2]
	count, err = CheckMajorityRule(cert, netNodes)
	require.Error(t, err)

	require.Equal(t, len(netNodes), count)
}

func TestRules_checkMinRoleError(t *testing.T) {
	nodes := []insolar.NetworkNode{
		node.NewNode(gen.Reference(), insolar.StaticRoleHeavyMaterial, nil, "", ""),
	}
	require.NotEmpty(t, checkMinRoleError(nodes, insolar.StaticRoleHeavyMaterial, 1, 2))

	require.Empty(t, checkMinRoleError(nodes, insolar.StaticRoleHeavyMaterial, 2, 2))

	require.NotEmpty(t, checkMinRoleError(nodes, insolar.StaticRoleLightMaterial, 1, 2))
}

func getDiscoveryNodes(count int) ([]insolar.NetworkNode, []insolar.DiscoveryNode) {
	netNodes := make([]insolar.NetworkNode, count)
	discoveryNodes := make([]insolar.DiscoveryNode, count)
	for i := 0; i < count; i++ {
		n := newNode(gen.Reference(), i)
		d := certificate.NewBootstrapNode(nil, "", n.Address(), n.ID().String(), n.Role().String())
		netNodes[i] = n
		discoveryNodes[i] = d
	}
	return netNodes, discoveryNodes
}

func newNode(ref insolar.Reference, i int) insolar.NetworkNode {
	return node.NewNode(ref, insolar.AllStaticRoles[i%len(insolar.AllStaticRoles)], nil,
		"127.0.0.1:"+strconv.Itoa(30000+i), "")
}
