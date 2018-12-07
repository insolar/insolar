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
)

func TestTree_Find(t *testing.T) {
	// Pulse in ID is equal to depth.
	tree := Tree{
		Head: &Jet{
			ID: *core.NewRecordID(0, nil),
			Right: &Jet{
				ID: *core.NewRecordID(1, nil),
				Right: &Jet{
					ID: *core.NewRecordID(2, nil),
					Left: &Jet{
						ID: *core.NewRecordID(3, nil),
						Right: &Jet{
							ID: *core.NewRecordID(4, nil),
						},
						Left: &Jet{
							ID: *core.NewRecordID(4, nil),
						},
					},
					Right: &Jet{
						ID: *core.NewRecordID(3, nil),
					},
				},
			},
			Left: &Jet{
				ID: *core.NewRecordID(1, nil),
			},
		},
	}
	val := []byte{0xD5} // 11010101

	jet, depth := tree.Find(val, 0)
	assert.Equal(t, tree.Head, jet)
	assert.Equal(t, depth, 0)
	jet, depth = tree.Find(val, 1)
	assert.Equal(t, tree.Head.Right, jet)
	assert.Equal(t, depth, 1)
	jet, depth = tree.Find(val, 2)
	assert.Equal(t, tree.Head.Right.Right, jet)
	assert.Equal(t, depth, 2)
	jet, depth = tree.Find(val, 3)
	assert.Equal(t, tree.Head.Right.Right.Left, jet)
	assert.Equal(t, depth, 3)
	jet, depth = tree.Find(val, 4)
	assert.Equal(t, tree.Head.Right.Right.Left.Right, jet)
	assert.Equal(t, depth, 4)
}
