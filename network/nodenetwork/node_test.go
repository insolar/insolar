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

package nodenetwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := NewNode("id", nil, "domainID")

	assert.NotNil(t, node)
}

func TestNode_GetDomainID(t *testing.T) {
	expectedDomain := "domain id"
	node := Node{
		id:       "id",
		domainID: expectedDomain,
	}

	assert.Equal(t, expectedDomain, node.GetDomainID())
}

func TestNode_GetNodeID(t *testing.T) {
	expectedID := "id"
	node := Node{
		domainID: "domain id",
		id:       expectedID,
		host:     nil,
	}
	assert.Equal(t, expectedID, node.GetNodeID())
}

func TestNode_GetNodeRole(t *testing.T) {
	expectedRole := "role"
	node := Node{
		id:       "id",
		domainID: "domain id",
		host:     nil,
	}

	node.setRole(expectedRole)
	assert.Equal(t, expectedRole, node.GetNodeRole())
}
