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

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemberPowerOf(t *testing.T) {
	require.Equal(t, MemberPowerOf(1), MemberPower(1))

	require.Equal(t, MemberPowerOf(0x1F), MemberPower(0x1F))

	require.Equal(t, MemberPowerOf(MaxLinearMemberPower), MemberPower(0xFF))

	require.Equal(t, MemberPowerOf(MaxLinearMemberPower+1), MemberPower(0xFF))

	require.Equal(t, MemberPowerOf(0x1F+1), MemberPower(0x1F+1))

	require.Equal(t, MemberPowerOf(0x1F<<1), MemberPower(0x2F))
}

func TestToLinearValue(t *testing.T) {
	require.Equal(t, MemberPowerOf(0).ToLinearValue(), uint16(0))

	require.Equal(t, MemberPowerOf(0x1F).ToLinearValue(), uint16(0x1F))

	require.Equal(t, MemberPowerOf(0x1F+1).ToLinearValue(), uint16(0x1F+1))

	require.Equal(t, MemberPowerOf(0x1F<<1).ToLinearValue(), uint16(0x3e))
}

func TestPercentAndMin(t *testing.T) {
	require.Equal(t, MemberPowerOf(MaxLinearMemberPower).PercentAndMin(100, MemberPowerOf(0)), ^MemberPower(0))

	require.Equal(t, MemberPowerOf(3).PercentAndMin(1, MemberPowerOf(2)), MemberPower(2))

	require.Equal(t, MemberPowerOf(3).PercentAndMin(80, MemberPowerOf(1)), MemberPower(2))
}

func TestNormalize(t *testing.T) {
	zero := MemberPowerSet([...]MemberPower{0, 0, 0, 0})
	require.Equal(t, zero.Normalize(), zero)

	require.Equal(t, MemberPowerSet([...]MemberPower{1, 0, 0, 0}).Normalize(), zero)

	m := MemberPowerSet([...]MemberPower{1, 1, 1, 1})
	require.Equal(t, m.Normalize(), m)
}

// Illegal cases:
// [ x,  y,  z,  0] - when any !=0 value of x, y, z
// [ 0,  x,  0,  y] - when x != 0 and y != 0
// any combination of non-zero x, y such that x > y and y > 0 and position(x) < position(y)
// And cases from the function logic.
func TestIsValid(t *testing.T) {
	require.True(t, MemberPowerSet([...]MemberPower{0, 0, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{1, 0, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 1, 0, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 0, 1, 0}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{0, 1, 0, 1}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{2, 1, 2, 2}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 0, 0, 1}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 1, 0, 1}).IsValid())

	require.False(t, MemberPowerSet([...]MemberPower{1, 1, 2, 1}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 0, 2, 2}).IsValid())

	require.True(t, MemberPowerSet([...]MemberPower{1, 1, 2, 2}).IsValid())
}
