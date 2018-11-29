/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

// NodeRole holds role of node
type NodeRole int

const (
	RoleUnknown = NodeRole(iota)
	RoleVirtual
	RoleHeavyMaterial
	RoleLightMaterial
)

var AllNodeRoles = []NodeRole{
	RoleVirtual,
	RoleLightMaterial,
	RoleHeavyMaterial,
}

// GetRoleFromString converts role from string to NodeRole
func GetRoleFromString(role string) NodeRole {
	switch role {
	case "virtual":
		return RoleVirtual
	case "heavy_material":
		return RoleHeavyMaterial
	case "light_material":
		return RoleLightMaterial
	}

	return RoleUnknown
}

func (nr NodeRole) String() string {
	switch nr {
	case RoleVirtual:
		return "virtual"
	case RoleHeavyMaterial:
		return "heavy_material"
	case RoleLightMaterial:
		return "light_material"
	}

	return "unknown"
}
