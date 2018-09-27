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

package message

import (
	"github.com/insolar/insolar/core"
)

type BaseLedgerMessage struct{}

func (BaseLedgerMessage) GetCaller() *core.RecordRef {
	return nil
}

func (BaseLedgerMessage) TargetRole() core.JetRole {
	return core.RoleHeavyExecutor
}

type GetCode struct {
	BaseLedgerMessage
	Code        core.RecordRef
	MachinePref []core.MachineType
}

func (e *GetCode) Type() core.MessageType {
	return TypeGetCode
}

func (e *GetCode) Target() *core.RecordRef {
	return &e.Code
}

func (e *GetCode) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type GetClass struct {
	BaseLedgerMessage
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

func (e *GetClass) Type() core.MessageType {
	return TypeGetClass
}

func (e *GetClass) Target() *core.RecordRef {
	return &e.Head
}

func (e *GetClass) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type GetObject struct {
	BaseLedgerMessage
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

func (e *GetObject) Type() core.MessageType {
	return TypeGetObject
}

func (e *GetObject) Target() *core.RecordRef {
	return &e.Head
}

func (e *GetObject) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type GetDelegate struct {
	BaseLedgerMessage
	Head    core.RecordRef
	AsClass core.RecordRef
}

func (e *GetDelegate) Type() core.MessageType {
	return TypeGetDelegate
}

func (e *GetDelegate) Target() *core.RecordRef {
	return &e.Head
}

func (e *GetDelegate) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type DeclareType struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	TypeDec []byte
}

func (e *DeclareType) Type() core.MessageType {
	return TypeDeclareType
}

func (e *DeclareType) Target() *core.RecordRef {
	return &e.Request
}

func (e *DeclareType) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type DeployCode struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	CodeMap map[core.MachineType][]byte
}

func (e *DeployCode) Type() core.MessageType {
	return TypeDeployCode
}

func (e *DeployCode) Target() *core.RecordRef {
	return &e.Request
}

func (e *DeployCode) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type ActivateClass struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
}

func (e *ActivateClass) Type() core.MessageType {
	return TypeActivateClass
}

func (e *ActivateClass) Target() *core.RecordRef {
	return &e.Request
}

func (e *ActivateClass) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type DeactivateClass struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
}

func (e *DeactivateClass) Type() core.MessageType {
	return TypeDeactivateClass
}

func (e *DeactivateClass) Target() *core.RecordRef {
	return &e.Class
}

func (e *DeactivateClass) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type UpdateClass struct {
	BaseLedgerMessage
	Domain     core.RecordRef
	Request    core.RecordRef
	Class      core.RecordRef
	Code       core.RecordRef
	Migrations []core.RecordRef
}

func (e *UpdateClass) Type() core.MessageType {
	return TypeUpdateClass
}

func (e *UpdateClass) Target() *core.RecordRef {
	return &e.Class
}

func (e *UpdateClass) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type ActivateObject struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObject) Type() core.MessageType {
	return TypeActivateObject
}

func (e *ActivateObject) Target() *core.RecordRef {
	return &e.Class
}

func (e *ActivateObject) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type ActivateObjectDelegate struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObjectDelegate) Type() core.MessageType {
	return TypeActivateObjectDelegate
}

func (e *ActivateObjectDelegate) Target() *core.RecordRef {
	return &e.Class
}

func (e *ActivateObjectDelegate) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

type DeactivateObject struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
}

func (e *DeactivateObject) Type() core.MessageType {
	return TypeDeactivateObject
}

func (e *DeactivateObject) Target() *core.RecordRef {
	return &e.Object
}

func (e *DeactivateObject) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

// UpdateObject for call of core.ArtifactManager.UpdateObj
type UpdateObject struct {
	BaseLedgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
	Memory  []byte
}

func (e *UpdateObject) Type() core.MessageType {
	return TypeUpdateObject
}

func (e *UpdateObject) Target() *core.RecordRef {
	return &e.Object
}

func (e *UpdateObject) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}
