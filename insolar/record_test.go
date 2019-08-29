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

package insolar_test

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	base58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ID and Reference serialization tests

func TestNewIDFromBytes(t *testing.T) {
	id := gen.ID()
	actualID := *insolar.NewIDFromBytes(id.Bytes())
	require.Equal(t, id, actualID)

	insolar.NewIDFromBytes(nil)
}

func TestNewIDFromBase58(t *testing.T) {
	id := gen.ID()
	idStr := "1" + base58.Encode(id.Bytes())
	id2, err := insolar.NewIDFromBase58(idStr)
	require.NoError(t, err)

	assert.Equal(t, id, *id2)
}

func TestRecordID_String(t *testing.T) {
	id := gen.ID()
	idStr := "1" + base58.Encode(id.Bytes()) + ".record"

	assert.Equal(t, idStr, id.String())
}

func TestNewRefFromBase58(t *testing.T) {
	recordID := gen.ID()
	domainID := gen.ID()
	refStr := "1" + base58.Encode(recordID.Bytes()) + insolar.RecordRefIDSeparator + "1" + base58.Encode(domainID.Bytes())

	expectedRef := insolar.NewGlobalReference(recordID, domainID)
	actualRef, err := insolar.NewReferenceFromBase58(refStr)
	require.NoError(t, err)

	assert.Equal(t, expectedRef, actualRef)
}

func TestRecordRef_String(t *testing.T) {
	ref := gen.Reference()
	expectedRefStr := "1" + base58.Encode(ref.GetLocal().Bytes())

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

func TestNewReferenceFromBase58(t *testing.T) {
	origin := gen.Reference()
	decoded, err := insolar.NewReferenceFromBase58(origin.String())
	require.NoError(t, err)
	assert.Equal(t, origin, *decoded)
}
