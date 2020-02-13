// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
