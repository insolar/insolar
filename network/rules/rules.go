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
	"errors"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

// CheckMajorityRule returns error if MajorityRule check not passed, also returns active discovery nodes count
func CheckMajorityRule(cert insolar.Certificate, nodes []insolar.NetworkNode) (int, error) {
	majorityRule := cert.GetMajorityRule()
	discoveriesInList := network.FindDiscoveriesInNodeList(nodes, cert)
	activeDiscoveryNodesLen := len(discoveriesInList)
	if activeDiscoveryNodesLen >= majorityRule {
		return activeDiscoveryNodesLen, nil
	}
	strErr := fmt.Sprintf("MajorityRule failed. Active discovery nodes len %d of %d (majorityRule). Diff:\n",
		activeDiscoveryNodesLen, majorityRule)
	discoveries := cert.GetDiscoveryNodes()
	for _, d := range discoveries {
		var found bool
		for _, n := range nodes {
			if d.GetNodeRef().Equal(n.ID()) {
				found = true
				break
			}
		}
		if !found {
			strErr += "host: " + d.GetHost() + " role: " + d.GetRole().String() + "\n"
		}
	}
	return activeDiscoveryNodesLen, errors.New(strErr)
}

// CheckMinRole returns true if MinRole check passed
func CheckMinRole(cert insolar.Certificate, nodes []insolar.NetworkNode) error {
	var virtualCount, heavyCount, lightCount uint
	for _, n := range nodes {
		switch n.Role() {
		case insolar.StaticRoleVirtual:
			virtualCount++
		case insolar.StaticRoleHeavyMaterial:
			heavyCount++
		case insolar.StaticRoleLightMaterial:
			lightCount++
		default:
			log.Warn("unknown node role")
		}
	}

	v, h, l := cert.GetMinRoles()
	if virtualCount >= v &&
		heavyCount >= h &&
		lightCount >= l {
		return nil
	}

	strErr := "MinRoles failed. " + checkMinRoleError(nodes, insolar.StaticRoleVirtual, virtualCount, v) +
		checkMinRoleError(nodes, insolar.StaticRoleHeavyMaterial, heavyCount, h) +
		checkMinRoleError(nodes, insolar.StaticRoleLightMaterial, lightCount, l)
	return errors.New(strErr)
}

func checkMinRoleError(nodes []insolar.NetworkNode, role insolar.StaticRole, count uint, minRole uint) string {
	var strErr string
	if count < minRole {
		strErr += fmt.Sprintf(role.String()+" %d of %d\n", count, minRole)
		for _, node := range nodes {
			if role == node.Role() {
				strErr += "host: " + node.Address() + "\n"
			}
		}
	}
	return strErr
}
