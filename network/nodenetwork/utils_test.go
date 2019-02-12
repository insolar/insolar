/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func Test_removeFromList(t *testing.T) {
	// table driven tests
	tests := map[string]struct {
		list          []core.RecordRef
		nodesToRemove []core.RecordRef
		expected      []core.RecordRef
	}{
		"simpleDiff": {
			list:          []core.RecordRef{{0}, {1}, {2}, {3}, {4}},
			nodesToRemove: []core.RecordRef{{0}, {2}},
			expected:      []core.RecordRef{{1}, {3}, {4}},
		},
		"equals": {
			list:          []core.RecordRef{{1}, {2}, {3}},
			nodesToRemove: []core.RecordRef{{1}, {2}, {3}},
			expected:      []core.RecordRef{},
		},
		"emptyRemoveList": {
			list:          []core.RecordRef{{1}, {2}, {3}},
			nodesToRemove: []core.RecordRef{},
			expected:      []core.RecordRef{{1}, {2}, {3}},
		},
		"emptyList": {
			list:          []core.RecordRef{},
			nodesToRemove: []core.RecordRef{{1}, {2}, {3}},
			expected:      []core.RecordRef{},
		},
		"allEmpty": {
			list:          []core.RecordRef{},
			nodesToRemove: []core.RecordRef{},
			expected:      []core.RecordRef{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, removeFromList(test.list, test.nodesToRemove))
		})
	}
}
