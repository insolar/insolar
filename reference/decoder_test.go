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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/pulse"
)

func TestDecoder_Decode_legacy(t *testing.T) {
	t.Parallel()

	legacyReference_ok := "1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11111111111111111111111111111111"
	{ // good old reference, ok to parse
		dec := NewDefaultDecoder(AllowLegacy)
		global, err := dec.Decode(legacyReference_ok)
		if assert.NoError(t, err) {
			assert.Equal(t, global.addressLocal, global.addressBase)
			assert.Equal(t, pulse.Number(65537), global.addressLocal.GetPulseNumber())
			assert.Equal(t, uint8(0x0), global.addressBase.getScope())
		}
	}
	{ // good old reference, disallow parsing
		dec := NewDefaultDecoder(0)
		_, err := dec.Decode(legacyReference_ok)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, legacy domain name")
		}
	}

	legacyReference_bad := "1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs"
	{ // bad legacy reference (domain isn't empty)
		dec := NewDefaultDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_bad)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, insufficient address length")
		}
	}

	legacyReference_empty := "11111111111111111111111111111111.11111111111111111111111111111111"
	{ // empty legacy reference
		dec := NewDefaultDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_empty)
		assert.NoError(t, err)
	}

	legacyReference_notFull := "115Ltamw9sE7JyRPGtz53j8FUbhdipmJ.11111111111111111111111111111111"
	{
		dec := NewDefaultDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_notFull)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "insufficient length")
		}
	}
	legacyReference_badSymbols := "1tJEBzbVurpgUrtyoloyAM3hCsSAxKLJ5U8LTb1EaerkZs.11111111111111111111111111111111"
	{ // good old reference, ok to parse
		dec := NewDefaultDecoder(AllowLegacy)
		_, err := dec.Decode(legacyReference_badSymbols)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "input string contains bad charachters")
		}
	}
}

func TestDecoder_Decode_new(t *testing.T) {
	t.Parallel()

	newReference_fixed := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.record"
	{ //
		dec := NewDefaultDecoder(AllowRecords)
		_, err := dec.Decode(newReference_fixed)
		assert.NoError(t, err)
	}
	{ //
		dec := NewDefaultDecoder(0)
		_, err := dec.Decode(newReference_fixed)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, record reference is not allowed")
		}
	}

	newReference_var_notReally := "base58+insolarv1:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_var_really := "base58+insolarv1:0115Ltamw9sE7JyRPGtz53j8FUbhdipmJ.record"
	{
		dec := NewDefaultDecoder(AllowRecords)

		g1, err1 := dec.Decode(newReference_var_notReally)
		g2, err2 := dec.Decode(newReference_var_really)

		if assert.NoError(t, err1) &&
			assert.NoError(t, err2) {

			assert.Equal(t, uint8(0), g1.addressLocal.getScope())
			assert.Equal(t, pulse.Number(0x1000), g1.addressLocal.GetPulseNumber())
			assert.Equal(t, g1.addressLocal.GetPulseNumber(), g2.addressLocal.GetPulseNumber())
			assert.Equal(t, g1.addressLocal.getScope(), g2.addressLocal.getScope())
			assert.Equal(t, g1.addressLocal.GetHash(), g2.addressLocal.GetHash())
		}
	}

	newReference_wo_part1 := "insolarv1:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_wo_part2 := "base58:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_wo_part1)
		assert.NoError(t, err)

		_, err = dec.Decode(newReference_wo_part2)
		assert.NoError(t, err)
	}

	newReference_part_switched := "insolarv1+base58:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_part_switched)
		assert.NoError(t, err)
	}

	newReference_bad_parts1 := "insolarv0+base58:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_bad_parts2 := "insolarv1+base59:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_bad_parts3 := "insolarv1+base58+bad:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_bad_parts1)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unsupported schema")
		}

		_, err = dec.Decode(newReference_bad_parts2)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unknown decoder")
		}

		_, err = dec.Decode(newReference_bad_parts3)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid schema")
		}
	}

	newReference_with_authority := "insolarv1://authority/0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_with_empty_authority := "insolarv1:///0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	newReference_with_bad_authority := "insolarv1://0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.record"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_with_authority)
		assert.NoError(t, err)

		_, err = dec.Decode(newReference_with_empty_authority)
		assert.NoError(t, err)

		_, err = dec.Decode(newReference_with_bad_authority)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "empty reference body")
		}
	}

	newReference_empty_body := "insolarv1:"
	newReference_empty_first := "insolarv1:.record"
	newReference_empty_second := "insolarv1:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3."
	newReference_legacy_domain := "insolarv1:0114CxjQofp9Rrh2jwVAdqaqVPfZEsrP27WaP8dgnHY3.11111111111111111111111111111111"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_empty_body)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "empty reference body")
		}

		_, err = dec.Decode(newReference_empty_first)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, empty reference body")
		}

		_, err = dec.Decode(newReference_empty_second)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, empty domain name")
		}

		_, err = dec.Decode(newReference_legacy_domain)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid reference, legacy domain name")
		}
	}

	newReference_badPrefix1 := "base58+insolarv1:21tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.record"
	newReference_badPrefix2 := "base58+insolarv1:91tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.record"
	newReference_badPrefix3 := "base58+insolarv1:a1tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.record"
	{
		var err error
		dec := NewDefaultDecoder(AllowRecords)

		_, err = dec.Decode(newReference_badPrefix1)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unsupported address prefix")
		}

		_, err = dec.Decode(newReference_badPrefix2)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unsupported address prefix")
		}

		_, err = dec.Decode(newReference_badPrefix3)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "record reference alias")
		}
	}

	newReference_full := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs"
	newReference_empty := "base58+insolarv1:0.0"
	{
		var err error
		dec := NewDefaultDecoder(0)

		_, err = dec.Decode(newReference_full)
		assert.NoError(t, err)

		_, err = dec.Decode(newReference_empty)
		assert.NoError(t, err)
	}

	newReference_brokenBody1 := "base58+insolarv1:11tJEBzbVurpgUolortyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs"
	newReference_brokenBody2 := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hyoloCsSAxKLJ5U8LTb1EaerkZs"
	newReference_brokenBody3 := "base58+insolarv1:01tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.01tJEBzbVurpgUrtyAM3hyoloCsSAxKLJ5U8LTb1EaerkZs"
	{
		var err error
		dec := NewDefaultDecoder(IgnoreParity)

		_, err = dec.Decode(newReference_brokenBody1)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "input string contains bad charachters")
		}

		_, err = dec.Decode(newReference_brokenBody2)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "input string contains bad charachters")
		}

		_, err = dec.Decode(newReference_brokenBody3)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "input string contains bad charachters")
		}
	}
}

func TestDecoder_Decode_parity(t *testing.T) {
	t.Parallel()

	newReference_with_parity := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs/22FwFmj"
	newReference_with_badParity := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs/2ololo"
	newReference_with_badParityPrefix := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs/FwFmj"
	newReference_with_emptyBodyAndParity := "base58+insolarv1:/ololo"
	{
		var err error
		dec := NewDefaultDecoder(IgnoreParity)

		_, err = dec.Decode(newReference_with_parity)
		assert.NoError(t, err)

		_, err = dec.Decode(newReference_with_badParity)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "unable to decode parity")
		}

		_, err = dec.Decode(newReference_with_badParityPrefix)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid parity prefix")
		}

		_, err = dec.Decode(newReference_with_emptyBodyAndParity)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "empty reference body")
		}
	}
}

func TestDecoder_Decode_aliases(t *testing.T) {
	t.Parallel()

	newReference_withDomainNameDecoder := "base58+insolarv1:11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs.help.somebody"
	newReference_withRecordNameDecoder := "base58+insolarv1:help.11tJEBzbVurpgUrtyAM3hCsSAxKLJ5U8LTb1EaerkZs"
	{
		var err error
		dec := NewDefaultDecoder(0)

		_, err = dec.Decode(newReference_withDomainNameDecoder)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "aliases are not allowed")
		}

		_, err = dec.Decode(newReference_withRecordNameDecoder)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "aliases are not allowed")
		}
	}
}

func TestCycle(t *testing.T) {
	inp := "15qwjxArNbE36WiUvT3NRQ9JqyUdW3mCRMBHKY31dJh6"

	var err error
	dec := NewDefaultDecoder(0)
	enc := NewBase58Encoder(0)

	gl, err := dec.Decode(inp)
	assert.NoError(t, err)

	out, err := enc.Encode(&gl)
	assert.NoError(t, err)

	assert.Equal(t, inp, out)

}
