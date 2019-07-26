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

func TestLen(t *testing.T) {
	sb := StateBitset{1, 2}
	require.Equal(t, len(sb), sb.Len())
}

func TestIsTrusted(t *testing.T) {
	require.True(t, BeHighTrust.IsTrusted())

	require.True(t, BeLimitedTrust.IsTrusted())

	require.False(t, BeBaselineTrust.IsTrusted())

	require.False(t, BeTimeout.IsTrusted())

	require.False(t, BeFraud.IsTrusted())

	require.False(t, maxBitsetEntry.IsTrusted())
}

func TestIsTimeout(t *testing.T) {
	require.False(t, BeHighTrust.IsTimeout())

	require.False(t, BeLimitedTrust.IsTimeout())

	require.False(t, BeBaselineTrust.IsTimeout())

	require.True(t, BeTimeout.IsTimeout())

	require.False(t, BeFraud.IsTimeout())

	require.False(t, maxBitsetEntry.IsTimeout())
}

func TestIsFraud(t *testing.T) {
	require.False(t, BeHighTrust.IsFraud())

	require.False(t, BeLimitedTrust.IsFraud())

	require.False(t, BeBaselineTrust.IsFraud())

	require.False(t, BeTimeout.IsFraud())

	require.True(t, BeFraud.IsFraud())

	require.False(t, maxBitsetEntry.IsFraud())
}

func TestFmtBitsetEntry(t *testing.T) {
	require.NotEmpty(t, FmtBitsetEntry(0))

	require.NotEmpty(t, FmtBitsetEntry(1))

	require.NotEmpty(t, FmtBitsetEntry(2))

	require.NotEmpty(t, FmtBitsetEntry(3))

	require.NotEmpty(t, FmtBitsetEntry(4))

	require.NotEmpty(t, FmtBitsetEntry(5))
}

func TestBitsetEntryString(t *testing.T) {
	require.NotEmpty(t, BeHighTrust.String())

	require.NotEmpty(t, BeLimitedTrust.String())

	require.NotEmpty(t, BeBaselineTrust.String())

	require.NotEmpty(t, BeTimeout.String())

	require.NotEmpty(t, BeFraud.String())

	require.NotEmpty(t, maxBitsetEntry.String())
}
