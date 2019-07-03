package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMaterial(t *testing.T) {
	v := PrimaryRoleInactive
	require.False(t, v.IsMaterial())
	v = PrimaryRoleNeutral
	require.False(t, v.IsMaterial())
	v = PrimaryRoleHeavyMaterial
	require.True(t, v.IsMaterial())
	v = PrimaryRoleLightMaterial
	require.True(t, v.IsMaterial())
	v = PrimaryRoleVirtual
	require.False(t, v.IsMaterial())
}

func TestIsHeavyMaterial(t *testing.T) {
	v := PrimaryRoleInactive
	require.False(t, v.IsHeavyMaterial())
	v = PrimaryRoleNeutral
	require.False(t, v.IsHeavyMaterial())
	v = PrimaryRoleHeavyMaterial
	require.True(t, v.IsHeavyMaterial())
	v = PrimaryRoleLightMaterial
	require.False(t, v.IsHeavyMaterial())
	v = PrimaryRoleVirtual
	require.False(t, v.IsHeavyMaterial())
}

func TestIsLightMaterial(t *testing.T) {
	v := PrimaryRoleInactive
	require.False(t, v.IsLightMaterial())
	v = PrimaryRoleNeutral
	require.False(t, v.IsLightMaterial())
	v = PrimaryRoleHeavyMaterial
	require.False(t, v.IsLightMaterial())
	v = PrimaryRoleLightMaterial
	require.True(t, v.IsLightMaterial())
	v = PrimaryRoleVirtual
	require.False(t, v.IsLightMaterial())
}

func TestIsVirtual(t *testing.T) {
	v := PrimaryRoleInactive
	require.False(t, v.IsVirtual())
	v = PrimaryRoleNeutral
	require.False(t, v.IsVirtual())
	v = PrimaryRoleHeavyMaterial
	require.False(t, v.IsVirtual())
	v = PrimaryRoleLightMaterial
	require.False(t, v.IsVirtual())
	v = PrimaryRoleVirtual
	require.True(t, v.IsVirtual())
}

func TestIsNeutral(t *testing.T) {
	v := PrimaryRoleInactive
	require.False(t, v.IsNeutral())
	v = PrimaryRoleNeutral
	require.True(t, v.IsNeutral())
	v = PrimaryRoleHeavyMaterial
	require.False(t, v.IsNeutral())
	v = PrimaryRoleLightMaterial
	require.False(t, v.IsNeutral())
	v = PrimaryRoleVirtual
	require.False(t, v.IsNeutral())
}

func TestIsInactive(t *testing.T) {
	v := PrimaryRoleInactive
	require.True(t, v.IsInactive())
	v = PrimaryRoleNeutral
	require.False(t, v.IsInactive())
	v = PrimaryRoleHeavyMaterial
	require.False(t, v.IsInactive())
	v = PrimaryRoleLightMaterial
	require.False(t, v.IsInactive())
	v = PrimaryRoleVirtual
	require.False(t, v.IsInactive())
}

func TestIsDiscovery(t *testing.T) {
	v := SpecialRoleNone
	require.False(t, v.IsDiscovery())
	v = SpecialRoleDiscovery
	require.True(t, v.IsDiscovery())
}
