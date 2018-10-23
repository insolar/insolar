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

package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFeature(t *testing.T) {
	feature, err := NewFeature("HELLO", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)

	feature, err = NewFeature("", "v1.1.1", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

	feature, err = NewFeature("START", "", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

	feature, err = NewFeature("HELLO2", "1.2", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	assert.Equal(t, feature.StartVersion.String(), "1.2.0")

	feature, err = NewFeature("START", "abc", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

}
