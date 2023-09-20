package insolar

// StaticRole holds role of node.
type StaticRole uint32

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

func (nr StaticRole) Equal(anr StaticRole) bool {
	return nr == anr
}
