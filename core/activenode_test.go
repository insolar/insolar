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

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJetRoleMask_Set(t *testing.T) {
	var mask JetRoleMask
	assert.Equal(t, uint64(0), uint64(mask))
	mask.Set(RoleVirtualValidator)
	assert.True(t, mask.IsSet(RoleVirtualValidator))
	assert.False(t, mask.IsSet(RoleVirtualExecutor))
	mask.Set(RoleVirtualExecutor)
	assert.True(t, mask.IsSet(RoleVirtualExecutor))
	mask.Unset(RoleVirtualValidator)
	assert.False(t, mask.IsSet(RoleVirtualValidator))
	assert.True(t, mask.IsSet(RoleVirtualExecutor))
}
