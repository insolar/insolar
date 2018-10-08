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

package record

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func str2Bytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func str2Hash(s string) [core.RecordHashSize]byte {
	// TODO: add check for s length
	var h [core.RecordHashSize]byte
	b := str2Bytes(s)
	_ = copy(h[:], b)
	return h
}

func str2ID(s string) ID {
	k := str2Hash(s)
	return Bytes2ID(k[:])
}

// Test_str2hash is a test for test helper str2Hash.
func Test_str2Hash(t *testing.T) {
	hashStr := "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"
	h := str2Hash(hashStr)
	assert.Equal(t, hashStr, fmt.Sprintf("%x", h))
}

func Test_str2Bytes(t *testing.T) {
	idStr := "00001111" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"
	id := str2Bytes(idStr)
	assert.Equal(t, idStr, fmt.Sprintf("%x", id))
}
