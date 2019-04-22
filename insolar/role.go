//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

// DynamicRole is number representing a node role.
type DynamicRole int

const (
	// DynamicRoleUndefined is used for special cases.
	DynamicRoleUndefined = DynamicRole(iota)
	// DynamicRoleVirtualExecutor is responsible for current pulse CPU operations.
	DynamicRoleVirtualExecutor
	// DynamicRoleVirtualValidator is responsible for previous pulse CPU operations.
	DynamicRoleVirtualValidator
	// DynamicRoleLightExecutor is responsible for current pulse Disk operations.
	DynamicRoleLightExecutor
	// DynamicRoleLightValidator is responsible for previous pulse Disk operations.
	DynamicRoleLightValidator
	// DynamicRoleHeavyExecutor is responsible for permanent Disk operations.
	DynamicRoleHeavyExecutor
)

// IsVirtualRole checks if node role is virtual (validator or executor).
func (r DynamicRole) IsVirtualRole() bool {
	switch r {
	case DynamicRoleVirtualExecutor:
		return true
	case DynamicRoleVirtualValidator:
		return true
	}
	return false
}
