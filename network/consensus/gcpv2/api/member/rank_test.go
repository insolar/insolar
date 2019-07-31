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

func TestGetPower(t *testing.T) {
	require.Equal(t, Rank(1).GetPower(), Power(1))
}

func TestGetIndex(t *testing.T) {
	require.Equal(t, JoinerIndex, JoinerRank.GetIndex())

	require.Zero(t, Rank((1<<8)-1).GetIndex())

	require.Equal(t, Index(1), Rank(1<<8).GetIndex())
}

func TestGetTotalCount(t *testing.T) {
	require.Zero(t, Rank((1<<18)-1).GetTotalCount())

	require.Equal(t, uint16(1), Rank(1<<18).GetTotalCount())
}

func TestIsJoiner(t *testing.T) {
	require.False(t, Rank(1).IsJoiner())

	require.True(t, JoinerRank.IsJoiner())
}

func TestString(t *testing.T) {
	joiner := "{joiner}"

	require.Equal(t, JoinerRank.String(), joiner)
	require.NotEqual(t, Rank(1).String(), joiner)
	require.Equal(t, joiner, JoinerRank.String())
	require.NotEqual(t, joiner, Rank(1).String())
}

func TestNewMembershipRank(t *testing.T) {
	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 0x03FF+1, 1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 0x03FF+1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 0x03FF+1) })

	require.Equal(t, Rank(0x80101), NewMembershipRank(ModeNormal, Power(1), 1, 2))
}

func TestAsMembershipRank(t *testing.T) {
	fr := FullRank{}
	fr.OpMode = ModeNormal
	fr.Power = 1
	fr.TotalIndex = 1
	require.Equal(t, Rank(0x80101), fr.AsMembershipRank(2))

	require.Panics(t, func() { fr.AsMembershipRank(0) })
}
