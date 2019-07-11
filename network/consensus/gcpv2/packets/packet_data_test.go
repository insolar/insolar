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

package packets

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/stretchr/testify/require"
)

func TestGetPower(t *testing.T) {
	require.Equal(t, MembershipRank(1).GetPower(), gcp_types.MemberPower(1))
}

func TestGetIndex(t *testing.T) {
	require.Equal(t, uint16(0), MembershipRank((1<<8)-1).GetIndex())

	require.Equal(t, uint16(1), MembershipRank(1<<8).GetIndex())
}

func TestGetTotalCount(t *testing.T) {
	require.Equal(t, uint16(0), MembershipRank((1<<18)-1).GetTotalCount())

	require.Equal(t, uint16(1), MembershipRank(1<<18).GetTotalCount())
}

func TestIsJoiner(t *testing.T) {
	require.False(t, MembershipRank(1).IsJoiner())

	require.True(t, JoinerMembershipRank.IsJoiner())
}

func TestString(t *testing.T) {
	joiner := "{joiner}"

	require.Equal(t, JoinerMembershipRank.String(), joiner)
	require.NotEqual(t, MembershipRank(1).String(), joiner)
	require.Equal(t, joiner, JoinerMembershipRank.String())
	require.NotEqual(t, joiner, MembershipRank(1).String())
}

func TestNewMembershipRank(t *testing.T) {
	require.Panics(t, func() { NewMembershipRank(gcp_types.MemberModeNormal, gcp_types.MemberPower(1), 1, 1) })

	require.Panics(t, func() { NewMembershipRank(gcp_types.MemberModeNormal, gcp_types.MemberPower(1), 0x03FF+1, 1) })

	require.Panics(t, func() { NewMembershipRank(gcp_types.MemberModeNormal, gcp_types.MemberPower(1), 1, 0x03FF+1) })

	require.Panics(t, func() { NewMembershipRank(gcp_types.MemberModeNormal, gcp_types.MemberPower(1), 1, 0x03FF+1) })

	require.Equal(t, MembershipRank(0x80101), NewMembershipRank(gcp_types.MemberModeNormal, gcp_types.MemberPower(1), 1, 2))
}

func TestEnsureNodeIndex(t *testing.T) {
	require.Panics(t, func() { ensureNodeIndex(0x03FF + 1) })

	require.Equal(t, uint32(2), ensureNodeIndex(2))
}
