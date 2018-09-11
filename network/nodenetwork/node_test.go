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

func TestNewNode(t *testing.T) {
	node := NewNode("nodeID", "hostID", core.String2Ref("domainID"))
	assert.NotNil(t, node)
}

func TestNode_GetDomainID(t *testing.T) {
	expectedDomain := core.String2Ref("domainID")
	node := Node{
		id:        "id",
		reference: core.String2Ref("domainID"),
	}

	assert.Equal(t, expectedDomain, node.GetReference())
}

func TestNode_GetNodeID(t *testing.T) {
	node := Node{
		reference: core.String2Ref("domainID"),
		id:        "id",
		hostID:    "",
	}
	assert.Equal(t, "id", node.GetNodeID())
}

func TestNode_GetNodeRole(t *testing.T) {
	expectedRole := "role"
	node := Node{
		id:        "id",
		reference: core.String2Ref("domainID"),
		hostID:    "",
	}

	node.setRole(expectedRole)
	assert.Equal(t, expectedRole, node.GetNodeRole())
}
