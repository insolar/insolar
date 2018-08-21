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

package id

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	bytes []byte
	ptr   int
}

func newMockReader() *mockReader {
	return &mockReader{
		bytes: []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		ptr:   0,
	}
}

func (mr *mockReader) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = mr.bytes[mr.ptr]
		mr.ptr++
		if mr.ptr >= len(mr.bytes) {
			mr.ptr = 0
		}
	}
	return len(b), nil
}

func TestNewID(t *testing.T) {
	random = newMockReader()

	id, err := NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})

	assert.NoError(t, err)
	assert.Len(t, id.GetHash(), 20)
	id1, _ := NewID(nil)
	id1.SetHash([]byte("1234567890abcdefghij"))
	assert.Equal(t, id1.GetHash(), id.GetHash())
}

func TestID_Equal(t *testing.T) {
	id1, _ := NewID(nil)
	id1.SetHash([]byte("1234567890abcdefghij"))
	id2, _ := NewID(nil)
	id2.SetHash([]byte("klmnopqrstuvwxyzABCD"))
	tests := []struct {
		id1, id2 ID
		equal    bool
		name     string
	}{
		{id1, id1, true, "same ids"},
		{id1, id2, false, "different ids"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, test.id1.HashEqual(test.id2.GetHash()))
		})
	}
}

func TestID_String(t *testing.T) {
	random = newMockReader()
	id, _ := NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})

	assert.Equal(t, "gkdhQDvLi23xxjXjhpMWaTt5byb", id.HashString())
}
