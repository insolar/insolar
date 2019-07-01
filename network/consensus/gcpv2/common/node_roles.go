package common

type NodePrimaryRole uint8

const (
	PrimaryRoleNeutral NodePrimaryRole = iota
	PrimaryRoleHeavyMaterial
	PrimaryRoleLightMaterial
	PrimaryRoleVirtual
)

func (v NodePrimaryRole) IsMaterial() bool {
	return v == PrimaryRoleHeavyMaterial || v == PrimaryRoleLightMaterial
}

func (v NodePrimaryRole) IsHeavyMaterial() bool {
	return v == PrimaryRoleHeavyMaterial
}

func (v NodePrimaryRole) IsLightMaterial() bool {
	return v == PrimaryRoleLightMaterial
}

func (v NodePrimaryRole) IsVirtual() bool {
	return v == PrimaryRoleVirtual
}

func (v NodePrimaryRole) IsNeutral() bool {
	return v == PrimaryRoleNeutral
}

type NodeSpecialRole uint8

const (
	SpecialRoleNoRole    NodeSpecialRole = 0
	SpecialRoleDiscovery NodeSpecialRole = 1 << iota
)

func (v NodeSpecialRole) IsDiscovery() bool {
	return v == SpecialRoleDiscovery
}
