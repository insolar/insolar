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

package ph2ctl

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNeighbourWeightScalerInt64(t *testing.T) {
	require.Panics(t, func() { NewScalerInt64(-1) })

	fullRange := int64(0)
	n1 := NewScalerInt64(fullRange)
	require.Equal(t, uint32(fullRange), n1.max)

	require.Equal(t, uint8(0), n1.shift)

	fullRange = int64(1 << 32)
	n2 := NewScalerInt64(fullRange)
	require.Equal(t, uint8(1), n2.shift)

	require.Equal(t, uint32(fullRange>>1), n2.max)
}

func TestNewNeighbourWeightScalerUint64(t *testing.T) {
	fullRange := uint64(0)
	n1 := NewScalerUint64(0, fullRange)
	require.Equal(t, uint32(fullRange), n1.max)

	require.Equal(t, uint8(0), n1.shift)

	fullRange = uint64(1 << 32)
	n2 := NewScalerUint64(0, fullRange)
	require.Equal(t, uint8(1), n2.shift)

	require.Equal(t, uint32(fullRange>>1), n2.max)
}

func TestScaleInt64(t *testing.T) {
	n1 := NewScalerInt64(0)
	require.Equal(t, uint32(0), n1.ScaleInt64(-1))

	require.Equal(t, uint32(0), n1.ScaleInt64(0))
	require.Equal(t, uint32(math.MaxUint32), n1.ScaleInt64(1))

	n2 := NewScalerInt64(1 << 32)
	require.Equal(t, uint32(math.MaxUint32), n2.ScaleInt64(1<<32))

	require.Equal(t, uint32(0x3fffffff), n2.ScaleInt64(1<<30))
}

func TestScaleUint64(t *testing.T) {
	n1 := NewScalerUint64(0, 0)
	require.Equal(t, uint32(0), n1.ScaleUint64(0))
	require.Equal(t, uint32(math.MaxUint32), n1.ScaleUint64(1))

	n2 := NewScalerUint64(0, 1<<32)
	require.Equal(t, uint32(math.MaxUint32), n2.ScaleUint64(1<<32))

	require.Equal(t, uint32(0x3fffffff), n2.ScaleUint64(1<<30))
}
