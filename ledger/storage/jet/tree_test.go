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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree_Find(t *testing.T) {
	// Pulse in ID is equal to depth.
	tree := Tree{
		Head: &Jet{
			Right: &Jet{
				Right: &Jet{
					Left: &Jet{
						Right: &Jet{},
						Left:  &Jet{},
					},
					Right: &Jet{},
				},
			},
			Left: &Jet{},
		},
	}
	val := make([]byte, 32)
	val[0] = 0xD5 // 11010101

	jet, depth := tree.Find(val)
	assert.Equal(t, tree.Head.Right.Right.Left.Right, jet)
	assert.Equal(t, depth, 4)
}

func TestTree_Update(t *testing.T) {
	tree := Tree{Head: &Jet{}}

	val := make([]byte, 32)
	val[0] = 0xD5 // 11010101

	jet, depth := tree.Find(val)
	require.Equal(t, tree.Head, jet)
	assert.Equal(t, 0, depth)

	tree.Update([]byte{1 << 7})
	jet, depth = tree.Find(val)
	require.Equal(t, tree.Head.Right, jet)
	assert.Equal(t, 1, depth)

	tree.Update([]byte{val[0]})
	jet, depth = tree.Find(val)
	require.Equal(t, tree.Head.Right.Right.Left.Right.Left.Right.Left.Right, jet)
	assert.Equal(t, 8, depth)
}
