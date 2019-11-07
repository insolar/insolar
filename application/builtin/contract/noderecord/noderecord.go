//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
