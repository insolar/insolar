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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/utils"
)

// NewRules creates network Rules component
func NewRules() network.Rules {
	return &rules{}
}

type rules struct {
	CertificateManager core.CertificateManager `inject:""`
	NodeKeeper         network.NodeKeeper      `inject:""`
}

// CheckMajorityRule returns true id MajorityRule check passed, also returns active discovery nodes count
func (r *rules) CheckMajorityRule() (bool, int) {
	// activeNodes []core.Node
	cert := r.CertificateManager.GetCertificate()
	majorityRule := cert.GetMajorityRule()
	activeDiscoveryNodesLen := len(utils.FindDiscoveriesInNodeList(r.NodeKeeper.GetActiveNodes(), cert))
	return activeDiscoveryNodesLen >= majorityRule, activeDiscoveryNodesLen
}

// CheckMinRole returns true if MinRole check passed
func (r *rules) CheckMinRole() bool {
	cert := r.CertificateManager.GetCertificate()

	nodes := r.NodeKeeper.GetActiveNodes()

	var virtualCount, heavyCount, lightCount uint
	for _, n := range nodes {
		switch n.Role() {
		case core.StaticRoleVirtual:
			virtualCount++
		case core.StaticRoleHeavyMaterial:
			heavyCount++
		case core.StaticRoleLightMaterial:
			lightCount++
		default:
			log.Warn("unknown node role")
		}
	}

	v, h, l := cert.GetMinRoles()
	return virtualCount >= v &&
		heavyCount >= h &&
		lightCount >= l
}
