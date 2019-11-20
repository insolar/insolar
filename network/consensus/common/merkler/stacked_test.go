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

package merkler

import (
	"fmt"
	"math"
	"math/bits"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

func TestStackedCalculator_Unbalanced_AllEntries(t *testing.T) {
	md := NewStackedCalculator(xorPairDigester{}, cryptkit.Digest{})

	for bit := uint64(1); bit != 0; bit <<= 1 {
		md.AddNext(newBits64(bit))
		require.Equal(t, bit<<1-1, md.ForkSequence().FinishSequence().FoldToUint64())
	}
	require.Equal(t, 64, md.Count())
	require.Equal(t, uint64(math.MaxUint64), md.ForkSequence().FinishSequence().FoldToUint64())

	md2 := md.ForkSequence()
	for bit := uint64(1) << 63; bit != 0; bit >>= 1 {
		md2.AddNext(newBits64(bit))
		require.Equal(t, bit-1, md2.ForkSequence().FinishSequence().FoldToUint64())
	}
	require.Equal(t, uint64(0), md2.FinishSequence().FoldToUint64())

	for bit := uint64(1); bit != 0; bit <<= 1 {
		md.AddNext(newBits64(bit))
		require.Equal(t, ^(bit<<1 - 1), md.ForkSequence().FinishSequence().FoldToUint64())
	}
	require.Equal(t, 128, md.Count())
	require.Equal(t, uint64(0), md.FinishSequence().FoldToUint64())
}

func TestStackedCalculator_EntryPos(t *testing.T) {
	// this table is a number of right-to-left transitions
	expectedPos := [0x3b]byte{
		//  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
		0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
		2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, // overflow
	}

	for balanced := 0; balanced <= 1; balanced++ {
		t.Run(fmt.Sprintf("balanced=%v", balanced != 0), func(t *testing.T) {
			for markBit := uint64(1); markBit < uint64(1)<<uint8(len(expectedPos)); /* to avoid overflow in this test */ markBit <<= 1 {
				var md StackedCalculator
				if balanced != 0 {
					md = NewStackedCalculator(xorShiftPairDigester{}, cryptkit.NewDigest(newBits64(0), "uint64"))
				} else {
					md = NewStackedCalculator(xorShiftPairDigester{}, cryptkit.Digest{})
				}

				for bit := uint64(1); bit != markBit; bit <<= 1 {
					md.AddNext(newBits64(0))
				}
				md.AddNext(newBits64(markBit))

				markBitPos := bits.Len64(markBit) - 1
				expected := byte(markBitPos) + expectedPos[markBitPos]

				v := md.ForkSequence().FinishSequence().FoldToUint64()
				require.Equal(t, 1, bits.OnesCount64(v), "0x%x", markBitPos)
				require.Equal(t, expected, byte(bits.Len64(v)-1), "0x%x", markBitPos)

				for bit := uint64(markBit) << 1; bit != 0; bit <<= 1 {
					md.AddNext(newBits64(0))
					v := md.ForkSequence().FinishSequence().FoldToUint64()
					require.Equal(t, expected, byte(bits.Len64(v)-1), "0x%x", markBitPos)
				}
				require.Equal(t, 64, md.Count())
			}
		})
	}
}

func TestStackedCalculator_Balanced_AllEntries(t *testing.T) {
	const unbalanced = uint64(1 << 56)
	expectedUnbalanced := [32]byte{
		1, 1, 1, 0, 2, 1, 1, 0, 3, 2, 2, 1, 2, 1, 1, 0,
		4, 3, 3, 2, 3, 2, 2, 1, 3, 2, 2, 1, 2, 1, 1, 0}

	md := NewStackedCalculator(xorCountPairDigester{}, cryptkit.NewDigest(newBits64(unbalanced), "uint64"))

	require.Equal(t, unbalanced, md.ForkSequence().FinishSequence().FoldToUint64())

	for bit := uint64(1); bit <= 1<<31; bit <<= 1 {
		md.AddNext(newBits64(bit))
		v := md.ForkSequence().FinishSequence().FoldToUint64()

		require.Equal(t, uint64(expectedUnbalanced[bits.Len64(bit)-1]), v>>56)
		require.Equal(t, bit<<1-1, v&math.MaxUint32)
	}
	require.Equal(t, 32, md.Count())
	require.Equal(t, uint64(math.MaxUint32), md.ForkSequence().FinishSequence().FoldToUint64())

	md2 := md.ForkSequence()
	for bit := uint64(1) << 31; bit != 0; bit >>= 1 {
		md2.AddNext(newBits64(bit))
		v := md2.ForkSequence().FinishSequence().FoldToUint64()
		//fmt.Println(v>>56)
		require.Equal(t, bit-1, v&math.MaxUint32)
	}
	require.Equal(t, uint64(0), md2.FinishSequence().FoldToUint64())

	for bit := uint64(1); bit <= 1<<31; bit <<= 1 {
		md.AddNext(newBits64(bit))
		v := md.ForkSequence().FinishSequence().FoldToUint64()
		//fmt.Println(v>>56)
		require.Equal(t, ^uint32(bit<<1-1), uint32(v))
	}
	require.Equal(t, 64, md.Count())
	require.Equal(t, uint64(0), md.FinishSequence().FoldToUint64())
}

func newBits64(v uint64) *longbits.Bits64 {
	v64 := longbits.NewBits64(v)
	return &v64
}

type xorPairDigester struct{}

func (p xorPairDigester) GetDigestSize() int {
	return 8
}

func (p xorPairDigester) DigestPair(digest0 longbits.FoldableReader, digest1 longbits.FoldableReader) cryptkit.Digest {
	return cryptkit.NewDigest(newBits64(digest0.FoldToUint64()^digest1.FoldToUint64()), "uint64")
}

func (p xorPairDigester) GetDigestMethod() cryptkit.DigestMethod {
	return "xor64"
}

type xorCountPairDigester struct {
	xorPairDigester
}

func (p xorCountPairDigester) DigestPair(digest0 longbits.FoldableReader, digest1 longbits.FoldableReader) cryptkit.Digest {
	const topByteMask = ^uint64(math.MaxUint64 >> 8)

	v0 := digest0.FoldToUint64()
	v1 := digest1.FoldToUint64()
	xored := (v0 ^ v1) &^ topByteMask
	//	counted := uint64(0)
	counted := v0&topByteMask + v1&topByteMask
	return cryptkit.NewDigest(newBits64(counted|xored), "uint64")
}

type xorShiftPairDigester struct {
	xorPairDigester
}

func (p xorShiftPairDigester) DigestPair(digest0 longbits.FoldableReader, digest1 longbits.FoldableReader) cryptkit.Digest {
	v0 := digest0.FoldToUint64()
	v1 := digest1.FoldToUint64()
	return cryptkit.NewDigest(newBits64(v0^(v1<<1)), "uint64")
}
