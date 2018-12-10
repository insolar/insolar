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
	val := []byte{0xD5} // 11010101

	jet, depth := tree.Find(val)
	assert.Equal(t, tree.Head.Right.Right.Left.Right, jet)
	assert.Equal(t, depth, 4)
}
