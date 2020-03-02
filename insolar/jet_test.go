// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJet_DebugString(t *testing.T) {
	var id *JetID
	id = NewJetID(0, []byte{})
	assert.Equal(t, "[JET 0 -]", id.DebugString())

	id = NewJetID(1, []byte{})
	assert.Equal(t, "[JET 1 0]", id.DebugString())
	id = NewJetID(2, []byte{})
	assert.Equal(t, "[JET 2 00]", id.DebugString())

	id = NewJetID(1, []byte{128})
	assert.Equal(t, "[JET 1 1]", id.DebugString())
	id = NewJetID(2, []byte{192})
	assert.Equal(t, "[JET 2 11]", id.DebugString())
}

func BenchmarkJet_DebugString_ZeroDepth(b *testing.B) {
	id := NewJetID(0, []byte{})
	for n := 0; n < b.N; n++ {
		id.DebugString()
	}
}

func BenchmarkJet_DebugString_Depth1(b *testing.B) {
	id := NewJetID(1, []byte{128})
	for n := 0; n < b.N; n++ {
		id.DebugString()
	}
}

func BenchmarkJet_DebugString_Depth5(b *testing.B) {
	id := NewJetID(5, []byte{128})
	for n := 0; n < b.N; n++ {
		id.DebugString()
	}
}
