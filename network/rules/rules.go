// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package rules

import (
	"fmt"

	"github.com/pkg/errors"

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
	strErr := fmt.Sprintf("Active discovery nodes len actual %d, expected %d. Not active ",
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
			strErr += d.GetHost() + " " + d.GetRole().String() + " "
		}
	}
	return activeDiscoveryNodesLen, errors.Wrap(errors.New(strErr), "MajorityRule failed")
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

	err := errors.New(fmt.Sprintf("%s actual %d expected %d, %s actual %d expected %d, %s actual %d expected %d",
		insolar.StaticRoleVirtual.String(), virtualCount, v,
		insolar.StaticRoleHeavyMaterial.String(), heavyCount, h,
		insolar.StaticRoleLightMaterial.String(), lightCount, l))
	return errors.Wrap(err, "MinRoles failed")
}
