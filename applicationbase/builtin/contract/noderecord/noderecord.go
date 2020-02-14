// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package noderecord

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// RecordInfo holds record info.
type RecordInfo struct {
	PublicKey string
	Role      insolar.StaticRole
}

// NodeRecord contains info about node.
type NodeRecord struct {
	foundation.BaseContract

	Record RecordInfo
}

// NewNodeRecord creates new NodeRecord.
func NewNodeRecord(publicKey string, roleStr string) (*NodeRecord, error) {
	if len(publicKey) == 0 {
		return nil, fmt.Errorf("public key is required")
	}
	if len(roleStr) == 0 {
		return nil, fmt.Errorf("role is required")
	}

	role := insolar.GetStaticRoleFromString(roleStr)
	if role == insolar.StaticRoleUnknown {
		return nil, fmt.Errorf("role is not supported: %s", roleStr)
	}

	return &NodeRecord{
		Record: RecordInfo{
			PublicKey: publicKey,
			Role:      role,
		},
	}, nil
}

// is needed for proxy
var INSATTR_GetNodeInfo_API = true

// GetNodeInfo returns RecordInfo.
// ins:immutable
func (nr *NodeRecord) GetNodeInfo() (RecordInfo, error) {
	return nr.Record, nil
}

// is needed for proxy
var INSATTR_GetPublicKey_API = true

// GetPublicKey returns public key.
// ins:immutable
func (nr *NodeRecord) GetPublicKey() (string, error) {
	return nr.Record.PublicKey, nil
}

// GetRole returns role.
// ins:immutable
func (nr *NodeRecord) GetRole() (insolar.StaticRole, error) {
	return nr.Record.Role, nil
}
