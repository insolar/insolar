/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */
 
package rules

import (
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestRules_CheckMinRole(t *testing.T) {
	cm := component.Manager{}
	r := NewRules()
	certManager := testutils.NewCertificateManagerMock(t)
	cert := testutils.NewCertificateMock(t)
	nodeKeeper := network.NewNodeKeeperMock(t)
	cm.Inject(r, certManager, nodeKeeper)

	certManager.GetCertificateMock.Set(func() (r core.Certificate) {
		return cert
	})

	cert.GetMinRolesMock.Set(func() (r uint, r1 uint, r2 uint) {
		return 1,0,0
	})

	nodeKeeper.GetActiveNodesMock.Set(func() (r []core.Node) {
		nodes, _ := getDiscoveryNodes(5)
		nodes = append(nodes, newNode(250))
		return nodes
	})

	result := r.CheckMinRole()
	assert.True(t, result)
}

func TestRules_CheckMajorityRule(t *testing.T) {
	cm := component.Manager{}
	r := NewRules()
	certManager := testutils.NewCertificateManagerMock(t)
	cert := testutils.NewCertificateMock(t)

	certManager.GetCertificateMock.Set(func() (r core.Certificate) {
		return cert
	})

	cert.GetDiscoveryNodesMock.Set(func() (r []core.DiscoveryNode) {
		_, nodes := getDiscoveryNodes(5)
		return nodes
	})

	cert.GetMajorityRuleMock.Set(func() (r int) {
		return 4
	})

	nodeKeeper := network.NewNodeKeeperMock(t)
	nodeKeeper.GetActiveNodesMock.Set(func() (r []core.Node) {
		nodes, _ := getDiscoveryNodes(5)
		nodes = append(nodes, newNode(250))
		return nodes
	})

	cm.Inject(r, certManager, nodeKeeper)

	result, count := r.CheckMajorityRule()
	assert.True(t, result)
	assert.Equal(t, 5, count)
}

func getDiscoveryNodes(count int) ([]core.Node, []core.DiscoveryNode) {
	result1 := make([]core.Node, count)
	result2 := make([]core.DiscoveryNode, count)

	for i := 0; i < count; i++ {
		n := newNode(i)
		d := certificate.NewBootstrapNode(nil, "", "127.0.0.1:3000", n.ID().String())
		result1[i] = n
		result2[i] = d
	}

	return result1, result2
}

func newNode(id int) core.Node {
	recordRef := core.RecordRef{byte(id)}
	node := nodenetwork.NewNode(recordRef, core.StaticRoleVirtual, nil, "127.0.0.1:3000", "")
	return node
}
