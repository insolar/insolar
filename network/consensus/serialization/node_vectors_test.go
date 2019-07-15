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

package serialization

import (
	"bytes"
	"context"
	"crypto/rand"
	"math"
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/stretchr/testify/require"
)

func TestNodeVectors_SerializeTo(t *testing.T) {
	nv := NodeVectors{}

	header := Header{}
	pctx := newPacketContext(context.Background(), &header)
	sctx := newSerializeContext(pctx, nil, nil, nil, nil)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := nv.SerializeTo(sctx, buf)
	require.NoError(t, err)
}

func TestNodeVectors_DeserializeFrom(t *testing.T) {
	nv := NodeVectors{}

	header := Header{}
	pctx := newPacketContext(context.Background(), &header)
	sctx := newSerializeContext(pctx, nil, nil, nil, nil)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := nv.SerializeTo(sctx, buf)
	require.NoError(t, err)

	dcxt := newDeserializeContext(pctx, nil, nil)
	nv2 := NodeVectors{}
	err = nv2.DeserializeFrom(dcxt, buf)
	require.NoError(t, err)

	require.Equal(t, nv, nv2)
}

func TestNodeVectors_AdditionalVectors(t *testing.T) {
	nv := NodeVectors{
		AdditionalStateVectors: make([]GlobulaStateVector, 2),
	}

	header := Header{}

	// Make bit range = 2 :(
	header.ClearFlag(1)
	header.SetFlag(2)

	pctx := newPacketContext(context.Background(), &header)
	sctx := newSerializeContext(pctx, nil, nil, nil, nil)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := nv.SerializeTo(sctx, buf)
	require.NoError(t, err)

	dcxt := newDeserializeContext(pctx, nil, nil)
	nv2 := NodeVectors{}
	err = nv2.DeserializeFrom(dcxt, buf)
	require.NoError(t, err)

	require.Equal(t, nv, nv2)
	require.Equal(t, 2, len(nv2.AdditionalStateVectors))
}

func TestNodeAppearanceBitset_isCompressed(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.False(t, b.isCompressed())

	b.FlagsAndLoLength = 64 // 0b01000000
	require.True(t, b.isCompressed())
}

func TestNodeAppearanceBitset_setIsCompressed(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.False(t, b.isCompressed())

	b.setIsCompressed(true)
	require.True(t, b.isCompressed())

	b.setIsCompressed(false)
	require.False(t, b.isCompressed())
}

func TestNodeAppearanceBitset_hasHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.False(t, b.hasHiLength())

	b.FlagsAndLoLength = 128 // 0b10000000
	require.True(t, b.hasHiLength())
}

func TestNodeAppearanceBitset_setHasHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.False(t, b.hasHiLength())

	b.setHasHiLength(true)
	require.True(t, b.hasHiLength())

	b.setHasHiLength(false)
	require.False(t, b.hasHiLength())
}

func TestNodeAppearanceBitset_getLoLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.EqualValues(t, 0, b.getLoLength())

	b.FlagsAndLoLength = 4 // 0b00000100
	require.EqualValues(t, 4, b.getLoLength())
}

func TestNodeAppearanceBitset_setLoLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.EqualValues(t, 0, b.getLoLength())

	b.setLoLength(50)
	require.EqualValues(t, 50, b.getLoLength())
}

func TestNodeAppearanceBitset_setLoLength_Panic(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.Panics(t, func() { b.setLoLength(loByteLengthMax + 1) })
}

func TestNodeAppearanceBitset_clearLoLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	b.setLoLength(50)
	require.EqualValues(t, 50, b.getLoLength())

	b.clearLoLength()
	require.EqualValues(t, 0, b.getLoLength())
}

func TestNodeAppearanceBitset_getHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.EqualValues(t, 0, b.getHiLength())

	b.HiLength = 4 // 0b00000100
	require.EqualValues(t, 4, b.getHiLength())
}

func TestNodeAppearanceBitset_setHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.EqualValues(t, 0, b.getHiLength())

	b.setHiLength(100)
	require.EqualValues(t, 100, b.getHiLength())
}

func TestNodeAppearanceBitset_setHiLength_Panic(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.Panics(t, func() { b.setHiLength(hiByteLengthMax + 1) })
}

func TestNodeAppearanceBitset_clearHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	b.setHiLength(50)
	require.EqualValues(t, 50, b.getHiLength())

	b.clearHiLength()
	require.EqualValues(t, 0, b.getHiLength())
}

func TestNodeAppearanceBitset_getLength(t *testing.T) {
	b := NodeAppearanceBitset{}
	b.FlagsAndLoLength = 4 // 0b00000100

	require.EqualValues(t, 4, b.getLength())

	b.HiLength = 5                           // 0b00000101
	require.EqualValues(t, 4, b.getLength()) // 0b00000100

	b.setHasHiLength(true)
	require.EqualValues(t, 164, b.getLength()) // 0b00010100100
}

func TestNodeAppearanceBitset_setLength(t *testing.T) {
	b := NodeAppearanceBitset{}
	require.EqualValues(t, 0, b.getLength())

	b.setLength(loByteLengthMax)
	require.EqualValues(t, loByteLengthMax, b.getLength())
	require.False(t, b.hasHiLength())

	b.setLength(1000)
	require.EqualValues(t, 1000, b.getLength())
	require.True(t, b.hasHiLength())
}

func TestNodeAppearanceBitset_setLength_Panic(t *testing.T) {
	b := NodeAppearanceBitset{}

	require.Panics(t, func() { b.setLength(byteLengthMax + 1) })
}

func TestNodeAppearanceBitset_SetBitset(t *testing.T) {
	b := NodeAppearanceBitset{}

	bitset := member.StateBitset{
		member.BeHighTrust,
		member.BeHighTrust,
		member.BeTimeout,
		member.BeFraud,
		member.BeBaselineTrust,
	}

	b.SetBitset(bitset)

	require.EqualValues(t, 5, b.getLength())
	require.NotEmpty(t, 5, b.Bytes)
}

func TestNodeAppearanceBitset_SetBitset_Panic(t *testing.T) {
	b := NodeAppearanceBitset{}

	bitset := member.StateBitset(make([]member.BitsetEntry, math.MaxUint16+1))

	require.Panics(t, func() { b.SetBitset(bitset) })
}

func TestNodeAppearanceBitset_GetBitset(t *testing.T) {
	b := NodeAppearanceBitset{}

	bitset := member.StateBitset{
		member.BeHighTrust,
		member.BeHighTrust,
		member.BeTimeout,
		member.BeFraud,
		member.BeBaselineTrust,
	}

	b.SetBitset(bitset)

	require.Equal(t, bitset, b.GetBitset())
}

func TestNodeAppearanceBitset_Empty(t *testing.T) {
	b := NodeAppearanceBitset{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := b.SerializeTo(nil, buf)
	require.NoError(t, err)

	require.EqualValues(t, buf.Len(), 1)

	b2 := NodeAppearanceBitset{}
	err = b2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, b, b2)
}

func TestNodeAppearanceBitset_NoHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	bitset := member.StateBitset{
		member.BeHighTrust,
		member.BeHighTrust,
		member.BeTimeout,
		member.BeFraud,
		member.BeBaselineTrust,
	}

	b.SetBitset(bitset)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := b.SerializeTo(nil, buf)
	require.NoError(t, err)

	require.EqualValues(t, buf.Len(), 1+len(bitset))

	b2 := NodeAppearanceBitset{}
	err = b2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, b, b2)
	require.Equal(t, bitset, b2.GetBitset())
}

func TestNodeAppearanceBitset_HasHiLength(t *testing.T) {
	b := NodeAppearanceBitset{}

	bitset := member.StateBitset(make([]member.BitsetEntry, loByteLengthMax+1))

	b.SetBitset(bitset)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := b.SerializeTo(nil, buf)
	require.NoError(t, err)

	require.EqualValues(t, buf.Len(), 1+1+len(bitset))

	b2 := NodeAppearanceBitset{}
	err = b2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, b, b2)
	require.Equal(t, bitset, b2.GetBitset())
}

func TestGlobulaStateVector_SerializeTo(t *testing.T) {
	v := GlobulaStateVector{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := v.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 132, buf.Len())
}

func TestGlobulaStateVector_DeserializeFrom(t *testing.T) {
	v1 := GlobulaStateVector{
		ExpectedRank: 2,
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(v1.VectorHash[:], b)
	copy(v1.SignedGlobulaStateHash[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := v1.SerializeTo(nil, buf)
	require.NoError(t, err)

	v2 := GlobulaStateVector{}
	err = v2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, v1, v2)
}

// func TestBitsetByteSize(t *testing.T) {
// 	require.EqualValues(t, 6, bitsetByteSize(16))
// 	require.EqualValues(t, 5, bitsetByteSize(12))
// 	require.EqualValues(t, 1, bitsetByteSize(2))
// 	require.EqualValues(t, 2, bitsetByteSize(3))
// }
