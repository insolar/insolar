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

package node

import (
	"testing"

	"github.com/insolar/insolar/network/host/id"
	"github.com/stretchr/testify/assert"
)

func TestNewOrigin_WithIds(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")
	ids, _ := id.NewIDs(10)

	expectedOrigin := &Origin{ids, addr}
	actualOrigin, err := NewOrigin(ids, addr)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrigin, actualOrigin)
}

func TestNewOrigin_WithoutIds(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")

	or, err := NewOrigin(nil, addr)

	assert.NoError(t, err)
	assert.Len(t, or.IDs, 1)
}

func TestOrigin_Contains(t *testing.T) {
	ids, _ := id.NewIDs(20)
	addr, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("10.10.11.11:12345")
	origin, _ := NewOrigin(ids[:10], addr)

	for i := range ids {
		contains := false
		if i < 10 {
			contains = true
		}
		assert.Equal(t, contains, origin.Contains(&Node{ids[i], addr}))
		assert.False(t, origin.Contains(&Node{ids[i], addr2}))
	}
}
