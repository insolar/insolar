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

	id, err := NewID()

	assert.NoError(t, err)
	assert.Len(t, id, 20)
	assert.Equal(t, ID("1234567890abcdefghij"), id)
}

func TestNewIDs(t *testing.T) {
	random = newMockReader()

	ids, err := NewIDs(123)

	assert.NoError(t, err)
	assert.Len(t, ids, 123)
	assert.Len(t, ids[0], 20)
	assert.Equal(t, ID("1234567890abcdefghij"), ids[0])
	assert.Equal(t, ID("klmnopqrstuvwxyzABCD"), ids[1])
}

func TestID_Equal(t *testing.T) {
	tests := []struct {
		id1, id2 ID
		equal    bool
		name     string
	}{
		{ID("1234567890abcdefghij"), ID("1234567890abcdefghij"), true, "same ids"},
		{ID("1234567890abcdefghij"), ID("klmnopqrstuvwxyzABCD"), false, "different ids"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, test.id1.Equal(test.id2))
		})
	}
}

func TestID_String(t *testing.T) {
	random = newMockReader()
	id, _ := NewID()

	assert.Equal(t, "gkdhQDvLi23xxjXjhpMWaTt5byb", id.String())
}
