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

package request

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

// Just to make Goland happy
func TestStubCurrentVer(t *testing.T) {
	newVer := NewVersion("v1.2.3")
	assert.Equal(t, newVer.Major, 1, "Major verify passed")
	assert.Equal(t, newVer.Minor, 2, "Minor verify passed")
	assert.Equal(t, newVer.Revision, 3, "Revision verify passed")
}
