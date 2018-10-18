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

package noderecord

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type RecordInfo struct {
	PublicKey string
	Roles     []core.NodeRole
	IP        string
}

// NodeRecord contains info about node
type NodeRecord struct {
	foundation.BaseContract

	Record RecordInfo
}

// New creates new NodeRecord
func NewNodeRecord(publicKey string, roles []string, ip string) *NodeRecord {
	resultRoles := []core.NodeRole{}
	for _, roleStr := range roles {
		role := core.GetRoleFromString(roleStr)
		if role == core.RoleUnknown {
			// TODO: return error
			panic("Role is not supported: " + roleStr)
		}
		resultRoles = append(resultRoles, role)
	}

	return &NodeRecord{
		Record: RecordInfo{
			PublicKey: publicKey,
			Roles:     resultRoles,
			IP:        ip,
		},
	}
}

// GetPublicKey returns public key
func (nr *NodeRecord) GetPublicKey() string {
	return nr.Record.PublicKey
}

// GetRoleAndPublicKey returns role-pubKey pair
func (nr *NodeRecord) GetRoleAndPublicKey() ([]core.NodeRole, string) {
	return nr.Record.Roles, nr.Record.PublicKey
}

// GetNodeInfo returns RecordInfo
func (nr *NodeRecord) GetNodeInfo() RecordInfo {
	return nr.Record
}

// SelfDestroy makes request to destroy current node record
func (nr *NodeRecord) Destroy() {
	nr.SelfDestruct()
}
