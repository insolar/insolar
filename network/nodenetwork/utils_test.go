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

package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func Test_diffList(t *testing.T) {
	old := []core.RecordRef{{0}, {1}, {2}, {3}, {4}}
	new := []core.RecordRef{{0}, {2}}
	expected := []core.RecordRef{{1}, {3}, {4}}

	assert.Equal(t, expected, diffList(old, new))

	old = []core.RecordRef{{4}, {0}, {2}, {3}, {1}}
	new = []core.RecordRef{{2}, {0}}
	expected = []core.RecordRef{{1}, {3}, {4}}

	assert.Equal(t, expected, diffList(old, new))

	old = []core.RecordRef{{1}, {2}, {3}, {4}}
	new = []core.RecordRef{{0}, {2}, {3}}
	expected = []core.RecordRef{{1}, {4}}

	assert.Equal(t, expected, diffList(old, new))

	old = []core.RecordRef{{1}, {2}, {3}}
	new = []core.RecordRef{{1}, {2}, {3}}

	assert.Empty(t, diffList(old, new))
}
