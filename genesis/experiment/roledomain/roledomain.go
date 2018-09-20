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

package roledomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/experiment/roledomain/rolerecord"
	"github.com/insolar/insolar/toolkit/go/foundation"
)

// RoleDomain holds rolerecords
type RoleDomain struct {
	foundation.BaseContract
}

// NewRoleDomain create new RoleDomain
func NewRoleDomain() *RoleDomain {
	return &RoleDomain{}
}

// RegisterNode registers node in system
func (rd *RoleDomain) RegisterNode(pk string, role core.JetRole) core.RecordRef {
	newRecord := rolerecord.New(pk, role)
	recordHolder := rd.AsChild(newRecord.GetReference())
	return recordHolder.GetReference()
}

// GetNodeRecord get node record by ref
func (rd *RoleDomain) GetNodeRecord(ref core.RecordRef) *rolerecord.RoleRecord {
	return rolerecord.GetObject(ref)
}

// RemoveNode deletes node from registry
func (rd *RoleDomain) RemoveNode(nodeRef core.RecordRef) {
	node := recordHolder.GetObject(nodeRef)
	node.SelfDestructRequest()
}
