// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package member

type PrimaryRole uint8 // MUST BE 6-bit

const (
	PrimaryRoleInactive PrimaryRole = iota
	PrimaryRoleNeutral
	PrimaryRoleHeavyMaterial
	PrimaryRoleLightMaterial
	PrimaryRoleVirtual
	// PrimaryRoleCascade
	// PrimaryRoleRecrypt
	maxPrimaryRole
)

const PrimaryRoleCount = int(maxPrimaryRole)

func (v PrimaryRole) Equal(other PrimaryRole) bool {
	return v == other
}

func (v PrimaryRole) IsMaterial() bool {
	return v == PrimaryRoleHeavyMaterial || v == PrimaryRoleLightMaterial
}

func (v PrimaryRole) IsHeavyMaterial() bool {
	return v == PrimaryRoleHeavyMaterial
}

func (v PrimaryRole) IsLightMaterial() bool {
	return v == PrimaryRoleLightMaterial
}

func (v PrimaryRole) IsVirtual() bool {
	return v == PrimaryRoleVirtual
}

func (v PrimaryRole) IsNeutral() bool {
	return v == PrimaryRoleNeutral
}

type SpecialRole uint8

const (
	SpecialRoleNone      SpecialRole = 0
	SpecialRoleDiscovery SpecialRole = 1 << iota
)

func (v SpecialRole) IsDiscovery() bool {
	return v == SpecialRoleDiscovery
}

func (v SpecialRole) Equal(other SpecialRole) bool {
	return v == other
}
