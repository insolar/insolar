package jet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
)

func TestTree_Update(t *testing.T) {
	tree := Tree{Head: &Jet{}}
	var (
		depth  uint8
		prefix []byte
	)

	lookup := insolar.NewID(gen.PulseNumber(), []byte{0xD5}) // 11010101

	id, actual := tree.Find(*lookup)
	depth, prefix = id.Depth(), id.Prefix()
	assert.Equal(t, depth, uint8(0))
	assert.Equal(t, prefix, make([]byte, insolar.RecordHashSize-1))
	assert.Equal(t, false, actual)

	tree.Update(*insolar.NewJetID(1, []byte{1 << 7}), false)
	id, actual = tree.Find(*lookup)
	depth, prefix = id.Depth(), id.Prefix()
	expectedPrefix := make([]byte, insolar.RecordHashSize-1)
	expectedPrefix[0] = 0x80
	require.Equal(t, uint8(1), depth)
	assert.Equal(t, expectedPrefix, prefix)
	assert.Equal(t, false, actual)

	tree.Update(*insolar.NewJetID(8, lookup.Hash()), false)
	id, actual = tree.Find(*lookup)
	depth, prefix = id.Depth(), id.Prefix()
	assert.Equal(t, uint8(8), depth)
	assert.Equal(t, lookup.Hash()[:insolar.RecordHashSize-1], prefix)
	assert.Equal(t, false, actual)

	tree.Update(*insolar.NewJetID(8, lookup.Hash()), true)
	id, actual = tree.Find(*lookup)
	depth, prefix = id.Depth(), id.Prefix()
	assert.Equal(t, uint8(8), depth)
	assert.Equal(t, lookup.Hash()[:insolar.RecordHashSize-1], prefix)
	assert.Equal(t, true, actual)
}

func TestTree_Find(t *testing.T) {
	tree := Tree{
		Head: &Jet{
			Left: &Jet{},
			Right: &Jet{
				Right: &Jet{
					Left: &Jet{
						Left:  &Jet{},
						Right: &Jet{},
					},
					Right: &Jet{},
				},
			},
		},
	}
	lookup := insolar.NewID(gen.PulseNumber(), []byte{0xD5}) // 11010101
	jetLookup := insolar.NewJetID(15, []byte{1, 2, 3})
	expectedPrefix := make([]byte, insolar.RecordIDSize-insolar.PulseNumberSize-1)
	expectedPrefix[0] = 0xD0 // 11010000

	id, actual := tree.Find(*lookup)
	depth, prefix := id.Depth(), id.Prefix()
	assert.Equal(t, depth, uint8(4))
	assert.Equal(t, expectedPrefix, prefix)
	assert.False(t, actual)

	jetID, actual := tree.Find(insolar.ID(*jetLookup))
	assert.Equal(t, *jetLookup, jetID)
	assert.True(t, actual)
}

func TestTree_Split(t *testing.T) {
	tree := Tree{
		Head: &Jet{
			Left: &Jet{},
			Right: &Jet{
				Right: &Jet{},
			},
		},
	}
	tooDeep := insolar.NewJetID(6, []byte{0xD5}) // 11010101
	ok := insolar.NewJetID(2, []byte{0xD5})      // 11010101

	t.Run("not existing jet returns error", func(t *testing.T) {
		_, _, err := tree.Split(*tooDeep)
		assert.Error(t, err)
	})

	t.Run("splits jet", func(t *testing.T) {
		okDepth, okPrefix := ok.Depth(), ok.Prefix()

		lExpectedPrefix := make([]byte, len(okPrefix))
		copy(lExpectedPrefix, okPrefix)
		lExpectedPrefix[0] = 0xC0 // 11000000
		rExpectedPrefix := make([]byte, len(okPrefix))
		copy(rExpectedPrefix, okPrefix)
		rExpectedPrefix[0] = 0xE0 // 11100000

		left, right, err := tree.Split(*ok)
		require.NoError(t, err)
		lDepth, lPrefix := left.Depth(), left.Prefix()
		rDepth, rPrefix := right.Depth(), right.Prefix()
		assert.Equal(t, uint8(okDepth+1), lDepth)
		assert.Equal(t, uint8(okDepth+1), rDepth)
		assert.Equal(t, lExpectedPrefix, lPrefix)
		assert.Equal(t, rExpectedPrefix, rPrefix)
	})
}

func TestTree_String(t *testing.T) {
	tree := Tree{
		Head: &Jet{
			Left: &Jet{
				Actual: true,
				Right: &Jet{
					Actual: true,
					Left:   &Jet{Actual: true},
					Right:  &Jet{},
				},
			},
			Right: &Jet{
				Left:  &Jet{},
				Right: &Jet{},
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
		Head: &Jet{},
	}
	assert.Equal(t, "root (level=0 actual=false)\n", emptyTree.String())
}

func TestTree_LeafIDs(t *testing.T) {
	tree := Tree{
		Head: &Jet{
			Left: &Jet{Actual: true},
			Right: &Jet{
				Right: &Jet{
					Left: &Jet{
						Left:  &Jet{Actual: false},
						Right: &Jet{Actual: true},
					},
					Right: &Jet{Actual: true},
				},
			},
		},
	}

	leafIDs := tree.LeafIDs()

	require.Equal(t, len(leafIDs), 3)
	assert.Equal(t, leafIDs[0], NewIDFromString("0"))
	assert.Equal(t, leafIDs[1], NewIDFromString("1101"))
	assert.Equal(t, leafIDs[2], NewIDFromString("111"))
}

func Test_ParsePrefix(t *testing.T) {
	prefix := parsePrefix("1100")
	assert.Equal(t, []byte{0xC0}, prefix)
}
