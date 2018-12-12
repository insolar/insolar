/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package jet

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree_Find(t *testing.T) {
	// Pulse in ID is equal to depth.
	tree := Tree{
		Head: &jet{
			Right: &jet{
				Right: &jet{
					Left: &jet{
						Right: &jet{},
						Left:  &jet{},
					},
					Right: &jet{},
				},
			},
			Left: &jet{},
		},
	}
	val := make([]byte, core.RecordIDSize-core.PulseNumberSize)
	val[0] = 0xD5 // 11010101
	expectedPrefix := make([]byte, core.RecordIDSize-core.PulseNumberSize-1)
	expectedPrefix[0] = 0xD0 // 11010000

	id := tree.Find(val)
	depth, prefix := Jet(*id)
	assert.Equal(t, depth, uint8(4))
	assert.Equal(t, expectedPrefix, prefix)
}

func TestTree_Update(t *testing.T) {
	tree := Tree{Head: &jet{}}

	val := make([]byte, core.RecordHashSize)
	val[0] = 0xD5 // 11010101

	id := tree.Find(val)
	depth, prefix := Jet(*id)
	assert.Equal(t, depth, uint8(0))
	assert.Equal(t, prefix, make([]byte, core.RecordHashSize-1))

	tree.Update(*NewID(1, []byte{1 << 7}))
	id = tree.Find(val)
	depth, prefix = Jet(*id)
	expectedPrefix := make([]byte, core.RecordHashSize-1)
	expectedPrefix[0] = 0x80
	require.Equal(t, uint8(1), depth)
	assert.Equal(t, expectedPrefix, prefix)

	tree.Update(*NewID(8, val))
	id = tree.Find(val)
	depth, prefix = Jet(*id)
	assert.Equal(t, uint8(8), depth)
	assert.Equal(t, val[:core.RecordHashSize-1], prefix)
}
