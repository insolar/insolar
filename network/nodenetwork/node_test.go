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

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := Node{id: testutils.RandomRef()}
	assert.NotNil(t, node)
}

func TestNode_GetNodeID(t *testing.T) {
	ref := testutils.RandomRef()
	node := Node{id: ref}
	assert.Equal(t, ref, node.GetID())
}
