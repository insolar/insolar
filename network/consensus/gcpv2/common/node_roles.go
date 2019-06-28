package common

type NodePrimaryRole uint8 //MUST BE 6-bit

const (
	PrimaryRoleUnknown NodePrimaryRole = iota
	PrimaryRoleNeutral
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

func (v NodePrimaryRole) IsUnknown() bool {
	return v == PrimaryRoleUnknown
}

type NodeSpecialRole uint8

const (
	SpecialRoleNoRole    NodeSpecialRole = 0
	SpecialRoleDiscovery NodeSpecialRole = 1 << iota
)

func (v NodeSpecialRole) IsDiscovery() bool {
	return v == SpecialRoleDiscovery
}
