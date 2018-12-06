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

// StaticRole holds role of node.
type StaticRole int

const (
	StaticRoleUnknown = StaticRole(iota)
	StaticRoleVirtual
	StaticRoleHeavyMaterial
	StaticRoleLightMaterial
)

// AllStaticRoles is an array of all possible StaticRoles.
var AllStaticRoles = []StaticRole{
	StaticRoleVirtual,
	StaticRoleLightMaterial,
	StaticRoleHeavyMaterial,
}

// GetStaticRoleFromString converts role from string to StaticRole.
func GetStaticRoleFromString(role string) StaticRole {
	switch role {
	case "virtual":
		return StaticRoleVirtual
	case "heavy_material":
		return StaticRoleHeavyMaterial
	case "light_material":
		return StaticRoleLightMaterial
	}

	return StaticRoleUnknown
}

func (nr StaticRole) String() string {
	switch nr {
	case StaticRoleVirtual:
		return "virtual"
	case StaticRoleHeavyMaterial:
		return "heavy_material"
	case StaticRoleLightMaterial:
		return "light_material"
	}

	return "unknown"
}
