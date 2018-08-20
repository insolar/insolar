/*
 *    Copyright 2018 INS Ecosystem
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
