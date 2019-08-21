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

func createSelfReferenceWithUnknownPulseNumber() *Global {
	g := &Global{}
	g.addressLocal.hash = generateBits224(0)
	g.addressLocal.pulseAndScope = uint32(pulse.Unknown)

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

func TestNewBase64Encoder(t *testing.T) {
	t.Parallel()

	{
		g := createRandomSelfReference()
		enc := NewBase64Encoder(FormatSchema)
		result, _ := enc.EncodeRecord(&g.addressLocal)
		assert.Contains(t, result, "base64+insolarv1:")
	}
}
