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
