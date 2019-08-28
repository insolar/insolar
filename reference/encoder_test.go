//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reference

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

func generateBits224(zeroBytes int) longbits.Bits224 {
	var hash [28]byte
	if _, err := rand.Read(hash[:]); err != nil {
		panic(err)
	}
	if zeroBytes < 0 || zeroBytes > 28 {
		panic("bad number of zeroBytes")
	} else {
		for i := 0; i < zeroBytes; i++ {
			hash[28-(i+1)] = 0
		}
	}
	return hash
}

func populateUsualPulse(number *uint32) {
	*number = uint32(rand.Int31n(pulse.MaxTimePulse-pulse.MinTimePulse)) + pulse.MinTimePulse
}

func populateSpecialPulse(number *uint32) {
	*number = uint32(rand.Int31n(pulse.MinTimePulse - 1))
}

func createRandomSelfReference() *Global {
	g := &Global{}
	g.addressLocal.hash = generateBits224(0)
	populateUsualPulse(&g.addressLocal.pulseAndScope)
	g.tryConvertToSelf()

	return g
}

func createRandomSelfReferenceWithSpecialPulse() *Global {
	g := &Global{}
	g.addressLocal.hash = generateBits224(0)
	populateSpecialPulse(&g.addressLocal.pulseAndScope)
	g.tryConvertToSelf()

	return g
}

func fixedSelfReference() *Global {
	g := &Global{
		addressLocal: Local{
			pulseAndScope: 0x16ee3f28,
			hash: longbits.Bits224{
				0x0b, 0xad, 0xb3, 0x7c,
				0x58, 0x21, 0xb6, 0xd9,
				0x55, 0x26, 0xa4, 0x1a,
				0x95, 0x04, 0x68, 0x0b,
				0x4e, 0x7c, 0x8b, 0x76,
				0x3a, 0x1b, 0x1d, 0x49,
				0xd4, 0x95, 0x5c, 0x84,
			},
		},
	}
	g.tryConvertToSelf()
	return g
}

func fixedReference() *Global {
	g := &Global{
		addressLocal: Local{
			pulseAndScope: 0x16ee3f28,
			hash: longbits.Bits224{
				0x0b, 0xad, 0xb3, 0x7c,
				0x58, 0x21, 0xb6, 0xd9,
				0x55, 0x26, 0xa4, 0x1a,
				0x95, 0x04, 0x68, 0x0b,
				0x4e, 0x7c, 0x8b, 0x76,
				0x3a, 0x1b, 0x1d, 0x49,
				0xd4, 0x95, 0x5c, 0x84,
			},
		},
		addressBase: Local{
			pulseAndScope: 0x16ee3f28,
			hash: longbits.Bits224{
				0xcb, 0xe0, 0x25, 0x5a,
				0xa5, 0xb7, 0xd4, 0x4b,
				0xec, 0x40, 0xf8, 0x4c,
				0x89, 0x2b, 0x9b, 0xff,
				0xd4, 0x36, 0x29, 0xb0,
				0x22, 0x3b, 0xee, 0xa5,
				0xf4, 0xf7, 0x43, 0x91,
			},
		},
	}
	g.tryConvertToSelf()
	return g
}

func fixedSelfReferenceWithSpecialPulse() *Global {
	g := &Global{
		addressLocal: Local{
			pulseAndScope: 0x42fd,
			hash: longbits.Bits224{
				0xcb, 0xe0, 0x25, 0x5a,
				0xa5, 0xb7, 0xd4, 0x4b,
				0xec, 0x40, 0xf8, 0x4c,
				0x89, 0x2b, 0x9b, 0xff,
				0xd4, 0x36, 0x29, 0xb0,
				0x22, 0x3b, 0xee, 0xa5,
				0xf4, 0xf7, 0x43, 0x91,
			},
		},
	}
	g.tryConvertToSelf()
	return g

}

func fixedSelfReferenceWithSpecialPulseZeroed() *Global {
	g := &Global{
		addressLocal: Local{
			pulseAndScope: 0x42fd,
			hash: longbits.Bits224{
				0x04, 0x03, 0x74, 0xf6,
				0x92, 0x4b, 0x98, 0xcb,
				0xf8, 0x71, 0x3f, 0x8d,
				0x96, 0x2d, 0x7c, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
	}
	g.tryConvertToSelf()
	return g
}

func TestEncoder_Encode(t *testing.T) {
	t.Parallel()

	{
		enc := NewBase58Encoder(0)
		result, _ := enc.Encode(nil)
		assert.Equal(t, NilRef, result)
	}

	{
		enc := NewBase58Encoder(0)
		result, _ := enc.Encode(&Global{})
		assert.Equal(t, "00", result)
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(0)
		result, _ := enc.Encode(g)
		assert.NotContains(t, result, '.')
		assert.NotContains(t, result, ':')
		assert.NotContains(t, result, '/')
		assert.Regexp(t, "^1.*", result)
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(FormatSchema)
		result, _ := enc.Encode(g)
		assert.NotContains(t, result, '.')
		assert.Contains(t, result, "base58+insolarv1:")
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(EncodingSchema)
		result, _ := enc.Encode(g)
		assert.NotContains(t, result, '.')
		assert.Contains(t, result, "base58:")
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(0)
		enc.authorityName = "someAuthority"
		result, _ := enc.Encode(g)
		assert.NotContains(t, result, '.')
		assert.Contains(t, result, "//someAuthority/")
	}

	{
		g := createRandomSelfReferenceWithSpecialPulse()
		enc := NewBase58Encoder(0)
		result, _ := enc.Encode(g)
		assert.Regexp(t, "^0.*", result)
	}
}

func TestEncoder_EncodeRecord(t *testing.T) {
	t.Parallel()

	{
		enc := NewBase58Encoder(0)
		result, _ := enc.EncodeRecord(nil)
		assert.Equal(t, NilRef, result)
	}

	{
		enc := NewBase58Encoder(0)

		g := Global{}
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Equal(t, "0.record", result)

		g.addressLocal.hash = generateBits224(0)
		result, _ = enc.EncodeRecord(&g.addressLocal)
		assert.Equal(t, "0.record", result)

	}

	{
		g := createRandomSelfReferenceWithSpecialPulse()
		enc := NewBase58Encoder(0)
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Regexp(t, "^0.*.record", result)
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(0)
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Regexp(t, "^1.*.record", result)
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(EncodingSchema)
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Contains(t, result, "base58:")
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(FormatSchema)
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Contains(t, result, "base58+insolarv1:")
	}

	{
		g := createRandomSelfReference()
		enc := NewBase58Encoder(0)
		enc.authorityName = "someAuthority"
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.NotContains(t, result, '.')
		assert.Contains(t, result, "//someAuthority/")
	}

	{
		g := createRandomSelfReferenceWithSpecialPulse()
		enc := NewBase58Encoder(0)

		res1, _ := enc.EncodeRecord(&g.addressLocal)
		maxLen := len(res1)

		g.addressLocal.hash = generateBits224(13)
		res2, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Greater(t, maxLen, len(res2))
	}
}

func TestEncoder_FixedEncode(t *testing.T) {
	t.Parallel()

	{
		g := fixedSelfReference()
		enc := NewBase58Encoder(0)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "12YWgCsqDVGBrcGgNPw2F7KevKXRABqx1uupHu35a325R", val)
	}

	{
		g := fixedSelfReferenceWithSpecialPulse()
		enc := NewBase58Encoder(0)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "011ERYsdJcrfuNhhtZAscud993Tm8K4wmtZ6DPkTWHyS", val)
	}

	{
		g := fixedSelfReferenceWithSpecialPulseZeroed()
		enc := NewBase58Encoder(0)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "011dWdkKLrn91P6sfRJgRSmk8j", val)
	}

	{
		g := fixedReference()
		enc := NewBase58Encoder(0)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "12YWgCsqDVGBrcGgNPw2F7KevKXRABqx1uupHu35a325R.12YWgCuosS34u2DHtQfk1tLAEvHJyUSRmrkaNarAEJgnC", val)

		g.addressLocal.pulseAndScope = 0
		val, err = enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "0.12YWgCuosS34u2DHtQfk1tLAEvHJyUSRmrkaNarAEJgnC", val)

		g.addressLocal.hash = [28]byte{}
		val, err = enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "0.12YWgCuosS34u2DHtQfk1tLAEvHJyUSRmrkaNarAEJgnC", val)
	}

	{
		g := fixedReference()
		enc := NewBase58Encoder(0)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "12YWgCsqDVGBrcGgNPw2F7KevKXRABqx1uupHu35a325R.12YWgCuosS34u2DHtQfk1tLAEvHJyUSRmrkaNarAEJgnC", val)

		g.addressBase.pulseAndScope = 0
		val, err = enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "12YWgCsqDVGBrcGgNPw2F7KevKXRABqx1uupHu35a325R.record", val)

		g.addressBase.hash = [28]byte{}
		val, err = enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "12YWgCsqDVGBrcGgNPw2F7KevKXRABqx1uupHu35a325R.record", val)
	}

	{
		g := fixedSelfReferenceWithSpecialPulseZeroed()
		enc := NewBase58Encoder(FormatSchema)

		val, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, "base58+insolarv1:011dWdkKLrn91P6sfRJgRSmk8j", val)
	}
}

func TestNewBase64Encoder(t *testing.T) {
	t.Parallel()

	{
		g := createRandomSelfReference()
		enc := NewBase64Encoder(FormatSchema)
		result, err := enc.EncodeRecord(&g.addressLocal)
		assert.NoError(t, err)
		assert.Contains(t, result, "base64+insolarv1:")
	}

	{
		g := fixedReference()
		enc := NewBase64Encoder(FormatSchema)
		result, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, result, "base64+insolarv1:1Fu4_KAuts3xYIbbZVSakGpUEaAtOfIt2OhsdSdSVXIQ.1Fu4_KMvgJVqlt9RL7ED4TIkrm__UNimwIjvupfT3Q5E")
	}

	{
		g := fixedSelfReferenceWithSpecialPulseZeroed()
		enc := NewBase64Encoder(FormatSchema)
		result, err := enc.Encode(g)
		assert.NoError(t, err)
		assert.Equal(t, result, "base64+insolarv1:0AABC_QQDdPaSS5jL-HE_jZYtfA")
	}
}
