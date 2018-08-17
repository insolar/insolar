package jetdrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJetDrop_Hash(t *testing.T) {
	drop1 := JetDrop{
		PrevHash:     []byte{1, 2, 3},
		RecordHashes: [][]byte{{4}, {5}, {6}},
	}
	drop2 := JetDrop{
		PrevHash:     []byte{1, 2, 3},
		RecordHashes: [][]byte{{4}, {5}, {6}},
	}
	drop3 := JetDrop{
		PrevHash:     []byte{1, 2, 3},
		RecordHashes: [][]byte{{5}, {4}, {6}},
	}

	h1, err := drop1.Hash()
	assert.NoError(t, err)
	h2, err := drop2.Hash()
	assert.NoError(t, err)
	h3, err := drop3.Hash()
	assert.NoError(t, err)
	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h1, h3)
}
