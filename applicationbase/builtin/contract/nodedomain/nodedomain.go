// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package nodedomain

import (
	"fmt"

	"github.com/insolar/insolar/applicationbase/builtin/proxy/noderecord"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// NodeDomain holds node records.
type NodeDomain struct {
	foundation.BaseContract

	NodeIndexPublicKey foundation.StableMap
}

// NewNodeDomain create new NodeDomain.
func NewNodeDomain() (*NodeDomain, error) {
	return &NodeDomain{
		NodeIndexPublicKey: make(foundation.StableMap),
	}, nil
}

// RegisterNode registers node in system.
func (nd *NodeDomain) RegisterNode(publicKey string, role string) (string, error) {
	canonicalKey, err := foundation.ExtractCanonicalPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("extracting canonical pk failed, current value %v", publicKey)
	}

	_, ok := nd.NodeIndexPublicKey[canonicalKey]
	if ok {
		return "", fmt.Errorf("node already exist with this public key: %s", publicKey)
	}

	newNode := noderecord.NewNodeRecord(publicKey, role)
	node, err := newNode.AsChild(nd.GetReference())
	if err != nil {
		return "", fmt.Errorf("failed to save as child: %s", err.Error())
	}

	newNodeRef := node.GetReference().String()

	nd.NodeIndexPublicKey[canonicalKey] = newNodeRef

	return newNodeRef, err
}

// GetNodeRefByPublicKey returns node reference.
// ins:immutable
func (nd *NodeDomain) GetNodeRefByPublicKey(publicKey string) (string, error) {
	canonicalKey, err := foundation.ExtractCanonicalPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("extracting canonical pk failed, current value %v", publicKey)
	}
	nodeRef, ok := nd.NodeIndexPublicKey[canonicalKey]
	if !ok {
		return "", fmt.Errorf("network node not found by public key: %s", publicKey)
	}
	return nodeRef, nil
}
