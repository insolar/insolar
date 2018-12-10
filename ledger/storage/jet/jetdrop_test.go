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

func TestJetDrop_Hash(t *testing.T) {
	drop1 := &JetDrop{
		Pulse:    1,
		PrevHash: []byte{1, 2, 3},
		Hash:     []byte{4, 5, 6},
	}
	drop2 := &JetDrop{
		Pulse:    2,
		PrevHash: []byte{1, 2, 3},
		Hash:     []byte{4, 5, 6},
	}

	b1, err := Encode(drop1)
	assert.NoError(t, err)
	assert.NotNil(t, b1)
	drop1got, err := Decode(b1)
	assert.NoError(t, err)
	assert.Equal(t, drop1, drop1got)

	b2, err := Encode(drop2)
	assert.NoError(t, err)
	assert.NotNil(t, b2)
	drop2got, err := Decode(b2)
	assert.NoError(t, err)
	assert.Equal(t, drop2, drop2got)

	assert.NotEqual(t, drop1got, drop2got)
	assert.NotEqual(t, b1, b2)
}
