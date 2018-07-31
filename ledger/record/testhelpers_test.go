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

package record

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func str2Bytes(s string) []byte {
	// var b bytes.Buffer
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func str2Hash(s string) Hash {
	// TODO: add check for s length
	var h Hash
	b := str2Bytes(s)
	_ = copy(h[:], b)
	return h
}

func str2ID(s string) ID {
	// TODO: add check for s length
	var id ID
	b := str2Bytes(s)
	_ = copy(id[:], b)
	return id
}

// Test_str2hash is a test for test helper str2Hash.
func Test_str2Hash(t *testing.T) {
	hashStr := "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"
	h := str2Hash(hashStr)
	assert.Equal(t, hashStr, fmt.Sprintf("%x", h))
}

func Test_str2ID(t *testing.T) {
	hashStr := "00001111110000"
	id := str2ID(hashStr)
	assert.Equal(t, hashStr, fmt.Sprintf("%x", id))
}
