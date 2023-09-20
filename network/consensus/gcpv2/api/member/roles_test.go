package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMaterial(t *testing.T) {
	require.False(t, PrimaryRoleInactive.IsMaterial())

	require.False(t, PrimaryRoleNeutral.IsMaterial())

	require.True(t, PrimaryRoleHeavyMaterial.IsMaterial())

	require.True(t, PrimaryRoleLightMaterial.IsMaterial())

	require.False(t, PrimaryRoleVirtual.IsMaterial())
}

func TestIsHeavyMaterial(t *testing.T) {
	require.False(t, PrimaryRoleInactive.IsHeavyMaterial())

	require.False(t, PrimaryRoleNeutral.IsHeavyMaterial())

	require.True(t, PrimaryRoleHeavyMaterial.IsHeavyMaterial())

	require.False(t, PrimaryRoleLightMaterial.IsHeavyMaterial())

	require.False(t, PrimaryRoleVirtual.IsHeavyMaterial())
}

func TestIsLightMaterial(t *testing.T) {
	require.False(t, PrimaryRoleInactive.IsLightMaterial())

	require.False(t, PrimaryRoleNeutral.IsLightMaterial())

	require.False(t, PrimaryRoleHeavyMaterial.IsLightMaterial())

	require.True(t, PrimaryRoleLightMaterial.IsLightMaterial())

	require.False(t, PrimaryRoleVirtual.IsLightMaterial())
}

func TestIsVirtual(t *testing.T) {
	require.False(t, PrimaryRoleInactive.IsVirtual())

	require.False(t, PrimaryRoleNeutral.IsVirtual())

	require.False(t, PrimaryRoleHeavyMaterial.IsVirtual())

	require.False(t, PrimaryRoleLightMaterial.IsVirtual())

	require.True(t, PrimaryRoleVirtual.IsVirtual())
}

func TestIsNeutral(t *testing.T) {
	require.False(t, PrimaryRoleInactive.IsNeutral())

	require.True(t, PrimaryRoleNeutral.IsNeutral())

	require.False(t, PrimaryRoleHeavyMaterial.IsNeutral())

	require.False(t, PrimaryRoleLightMaterial.IsNeutral())

	require.False(t, PrimaryRoleVirtual.IsNeutral())
}

func TestIsDiscovery(t *testing.T) {
	require.False(t, SpecialRoleNone.IsDiscovery())

	require.True(t, SpecialRoleDiscovery.IsDiscovery())
}
