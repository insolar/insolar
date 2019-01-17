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
	"strings"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree_Find(t *testing.T) {
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
	lookup := core.NewRecordID(0, []byte{0xD5}) // 11010101
	jetLookup := NewID(15, []byte{1, 2, 3})
	expectedPrefix := make([]byte, core.RecordIDSize-core.PulseNumberSize-1)
	expectedPrefix[0] = 0xD0 // 11010000

	id, actual := tree.Find(*lookup)
	depth, prefix := Jet(*id)
	assert.Equal(t, depth, uint8(4))
	assert.Equal(t, expectedPrefix, prefix)
	assert.False(t, actual)

	jetID, actual := tree.Find(*jetLookup)
	assert.Equal(t, jetLookup, jetID)
	assert.True(t, actual)
}

func TestTree_Update(t *testing.T) {
	tree := Tree{Head: &jet{}}

	lookup := core.NewRecordID(0, []byte{0xD5}) // 11010101

	id, actual := tree.Find(*lookup)
	depth, prefix := Jet(*id)
	assert.Equal(t, depth, uint8(0))
	assert.Equal(t, prefix, make([]byte, core.RecordHashSize-1))
	assert.Equal(t, false, actual)

	tree.Update(*NewID(1, []byte{1 << 7}), false)
	id, actual = tree.Find(*lookup)
	depth, prefix = Jet(*id)
	expectedPrefix := make([]byte, core.RecordHashSize-1)
	expectedPrefix[0] = 0x80
	require.Equal(t, uint8(1), depth)
	assert.Equal(t, expectedPrefix, prefix)
	assert.Equal(t, false, actual)

	tree.Update(*NewID(8, lookup.Hash()), false)
	id, actual = tree.Find(*lookup)
	depth, prefix = Jet(*id)
	assert.Equal(t, uint8(8), depth)
	assert.Equal(t, lookup.Hash()[:core.RecordHashSize-1], prefix)
	assert.Equal(t, false, actual)

	tree.Update(*NewID(8, lookup.Hash()), true)
	id, actual = tree.Find(*lookup)
	depth, prefix = Jet(*id)
	assert.Equal(t, uint8(8), depth)
	assert.Equal(t, lookup.Hash()[:core.RecordHashSize-1], prefix)
	assert.Equal(t, true, actual)
}

func TestTree_Split(t *testing.T) {
	tree := Tree{
		Head: &jet{
			Right: &jet{
				Right: &jet{},
			},
			Left: &jet{},
		},
	}
	tooDeep := NewID(6, []byte{0xD5}) // 11010101
	ok := NewID(2, []byte{0xD5})      // 11010101

	t.Run("not existing jet returns error", func(t *testing.T) {
		_, _, err := tree.Split(*tooDeep)
		assert.Error(t, err)
	})

	t.Run("splits jet", func(t *testing.T) {
		okDepth, okPrefix := Jet(*ok)
		lExpectedPrefix := make([]byte, len(okPrefix))
		copy(lExpectedPrefix, okPrefix)
		lExpectedPrefix[0] = 0xC0 // 11000000
		rExpectedPrefix := make([]byte, len(okPrefix))
		copy(rExpectedPrefix, okPrefix)
		rExpectedPrefix[0] = 0xE0 // 11100000

		left, right, err := tree.Split(*ok)
		require.NoError(t, err)
		lDepth, lPrefix := Jet(*left)
		rDepth, rPrefix := Jet(*right)
		assert.Equal(t, uint8(okDepth+1), lDepth)
		assert.Equal(t, uint8(okDepth+1), rDepth)
		assert.Equal(t, lExpectedPrefix, lPrefix)
		assert.Equal(t, rExpectedPrefix, rPrefix)
	})
}

func TestTree_String(t *testing.T) {
	tree := Tree{
		Head: &jet{
			Left: &jet{
				Actual: true,
				Right: &jet{
					Actual: true,
					Left:   &jet{Actual: true},
					Right:  &jet{},
				},
			},
			Right: &jet{
				Left:  &jet{},
				Right: &jet{},
			},
		},
	}
	treeOut := strings.Join([]string{
		"root (level=0 actual=false)",
		" 0 (level=1 actual=true)",
		"  01 (level=2 actual=true)",
		"   010 (level=3 actual=true)",
		"   011 (level=3 actual=false)",
		" 1 (level=1 actual=false)",
		"  10 (level=2 actual=false)",
		"  11 (level=2 actual=false)",
	}, "\n") + "\n"
	assert.Equal(t, treeOut, tree.String())

	emptyTree := Tree{
		Head: &jet{},
	}
	assert.Equal(t, "root (level=0 actual=false)\n", emptyTree.String())
}

func TestTree_Merge(t *testing.T) {
	t.Run("One level", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: true,
				},
				Right: &jet{
					Actual: false,
				},
				Actual: true,
			},
		}

		newTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: false,
				},
				Right: &jet{
					Actual: true,
				},
				Actual: false,
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t,
			Tree{
				Head: &jet{
					Left: &jet{
						Actual: true,
					},
					Right: &jet{
						Actual: true,
					},
					Actual: true,
				},
			}.String(),
			result.String())
	})

	t.Run("Five levels saved is false active is true", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: false,
					Left: &jet{
						Actual: false,
					},
					Right: &jet{
						Actual: false,
					},
				},
				Right: &jet{
					Actual: false,
					Left: &jet{
						Actual: false,
						Left: &jet{
							Actual: false,
						},
						Right: &jet{
							Actual: false,
							Left: &jet{
								Actual: false,
							},
							Right: &jet{
								Actual: false,
							},
						},
					},
					Right: &jet{
						Actual: false,
					},
				},
				Actual: false,
			},
		}

		newTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
					},
					Right: &jet{
						Actual: true,
					},
				},
				Right: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
						Left: &jet{
							Actual: true,
						},
						Right: &jet{
							Actual: true,
							Left: &jet{
								Actual: true,
							},
							Right: &jet{
								Actual: true,
							},
						},
					},
					Right: &jet{
						Actual: true,
					},
				},
				Actual: true,
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t, newTree.String(), result.String())
	})

	t.Run("Two levels saved is empty", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Actual: false,
			},
		}

		newTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
					},
					Right: &jet{
						Actual: true,
					},
				},
				Right: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
						Left: &jet{
							Actual: true,
						},
						Right: &jet{
							Actual: true,
							Left: &jet{
								Actual: true,
							},
							Right: &jet{
								Actual: true,
							},
						},
					},
					Right: &jet{
						Actual: true,
					},
				},
				Actual: true,
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t, newTree.String(), result.String())
	})

	t.Run("Three levels saved with left and new with right", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Actual: true,
				Left: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
						Left: &jet{
							Actual: true,
						},
					},
					Right: &jet{
						Actual: false,
					},
				},
			},
		}

		newTree := Tree{
			Head: &jet{
				Right: &jet{
					Right: &jet{
						Actual: true,
						Right: &jet{
							Actual: true,
							Right: &jet{
								Actual: false,
							},
						},
					},
				},
				Actual: true,
			},
		}

		expectedTree := Tree{
			Head: &jet{
				Left: &jet{
					Actual: true,
					Left: &jet{
						Actual: true,
						Left: &jet{
							Actual: true,
						},
					},
					Right: &jet{
						Actual: false,
					},
				},

				Actual: true,

				Right: &jet{
					Right: &jet{
						Actual: true,
						Right: &jet{
							Actual: true,
							Right: &jet{
								Actual: false,
							},
						},
					},
				},
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t, expectedTree.String(), result.String())
	})

	t.Run("heads saved false new true", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Actual: false,
			},
		}

		newTree := Tree{
			Head: &jet{
				Actual: true,
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t, newTree.String(), result.String())
	})

	t.Run("heads saved false new true", func(t *testing.T) {
		savedTree := Tree{
			Head: &jet{
				Actual: true,
			},
		}

		newTree := Tree{
			Head: &jet{
				Actual: false,
			},
		}

		result := savedTree.Merge(&newTree)

		require.Equal(t, savedTree.String(), result.String())
	})
}

func TestTree_LeafIDs(t *testing.T) {
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

	leafIDs := tree.LeafIDs()

	require.Equal(t, len(leafIDs), 4)
	assert.Equal(t, leafIDs[0], *NewID(1, nil))          // 0000
	assert.Equal(t, leafIDs[1], *NewID(4, []byte{0xC0})) // 1100
	assert.Equal(t, leafIDs[2], *NewID(4, []byte{0xD0})) // 1101
	assert.Equal(t, leafIDs[3], *NewID(3, []byte{0xE0})) // 1110
}
