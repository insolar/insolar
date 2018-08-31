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
