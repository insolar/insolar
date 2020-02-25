// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
)

// ID and Reference serialization tests

func TestNewIDFromBytes(t *testing.T) {
	id := gen.ID()
	actualID := *insolar.NewIDFromBytes(id.Bytes())
	require.Equal(t, id, actualID)

	insolar.NewIDFromBytes(nil)
}

func TestNewIDFromString(t *testing.T) {
	id := gen.ID()
	idStr := "insolar:1" + base64.RawURLEncoding.EncodeToString(id.Bytes())
	id2, err := insolar.NewIDFromString(idStr)
	require.NoError(t, err)

	assert.Equal(t, id, *id2)
}

func TestRecordID_String(t *testing.T) {
	id := gen.ID()
	idStr := "insolar:1" + base64.RawURLEncoding.EncodeToString(id.Bytes()) + ".record"

	assert.Equal(t, idStr, id.String())
}

func TestNewRefFromString(t *testing.T) {
	recordID := gen.ID()
	domainID := gen.ID()
	refStr := "insolar:1" +
		base64.RawURLEncoding.EncodeToString(recordID.Bytes()) +
		insolar.RecordRefIDSeparator + "1" +
		base64.RawURLEncoding.EncodeToString(domainID.Bytes())

	expectedRef := insolar.NewGlobalReference(recordID, domainID)
	actualRef, err := insolar.NewReferenceFromString(refStr)
	require.NoError(t, err)

	assert.Equal(t, expectedRef, actualRef)
}

func TestRecordRef_String(t *testing.T) {
	ref := gen.Reference()
	expectedRefStr := "insolar:1" + base64.RawURLEncoding.EncodeToString(ref.GetLocal().Bytes())

	assert.Equal(t, expectedRefStr, ref.String())
}

func TestRecordID_DebugString_Jet(t *testing.T) {
	j := insolar.ID(*insolar.NewJetID(0, []byte{}))
	assert.Equal(t, "[JET 0 -]", j.DebugString())

	j = insolar.ID(*insolar.NewJetID(1, []byte{}))
	assert.Equal(t, "[JET 1 0]", j.DebugString())
	j = insolar.ID(*insolar.NewJetID(2, []byte{}))
	assert.Equal(t, "[JET 2 00]", j.DebugString())

	j = insolar.ID(*insolar.NewJetID(1, []byte{128}))
	assert.Equal(t, "[JET 1 1]", j.DebugString())
	j = insolar.ID(*insolar.NewJetID(2, []byte{192}))
	assert.Equal(t, "[JET 2 11]", j.DebugString())
}

func BenchmarkRecordID_DebugString_ZeroDepth(b *testing.B) {
	jet := insolar.ID(*insolar.NewJetID(0, []byte{}))
	for n := 0; n < b.N; n++ {
		jet.DebugString()
	}
}

func BenchmarkRecordID_DebugString_Depth1(b *testing.B) {
	jet := insolar.ID(*insolar.NewJetID(1, []byte{128}))
	for n := 0; n < b.N; n++ {
		jet.DebugString()
	}
}

func BenchmarkRecordID_DebugString_Depth5(b *testing.B) {
	jet := insolar.ID(*insolar.NewJetID(5, []byte{128}))
	for n := 0; n < b.N; n++ {
		jet.DebugString()
	}
}

func TestNewReferenceFromString(t *testing.T) {
	origin := gen.Reference()
	decoded, err := insolar.NewReferenceFromString(origin.String())
	require.NoError(t, err)
	assert.Equal(t, origin, *decoded)
}
