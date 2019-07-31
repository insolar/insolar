//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
