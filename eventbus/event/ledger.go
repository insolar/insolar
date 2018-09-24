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

package event

import (
	"io"

	"github.com/insolar/insolar/core"
)

type GetCode struct {
	Code        core.RecordRef
	MachinePref []core.MachineType
}

func (e *GetCode) Serialize() (io.Reader, error) {
	return serialize(e, TypeGetCode)
}

func (e *GetCode) GetReference() core.RecordRef {
	return e.Code
}

func (e *GetCode) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetCode) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type GetClass struct {
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

func (e *GetClass) Serialize() (io.Reader, error) {
	return serialize(e, TypeGetClass)
}

func (e *GetClass) GetReference() core.RecordRef {
	return e.Head
}

func (e *GetClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetClass) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type GetObject struct {
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

func (e *GetObject) Serialize() (io.Reader, error) {
	return serialize(e, TypeGetObject)
}

func (e *GetObject) GetReference() core.RecordRef {
	return e.Head
}

func (e *GetObject) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetObject) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type GetDelegate struct {
	Head    core.RecordRef
	AsClass core.RecordRef
}

func (e *GetDelegate) Serialize() (io.Reader, error) {
	return serialize(e, TypeGetDelegate)
}

func (e *GetDelegate) GetReference() core.RecordRef {
	return e.Head
}

func (e *GetDelegate) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetDelegate) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type DeclareType struct {
	Domain  core.RecordRef
	Request core.RecordRef
	TypeDec []byte
}

func (e *DeclareType) Serialize() (io.Reader, error) {
	return serialize(e, TypeDeclareType)
}

func (e *DeclareType) GetReference() core.RecordRef {
	return e.Request
}

func (e *DeclareType) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeclareType) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type DeployCode struct {
	Domain  core.RecordRef
	Request core.RecordRef
	CodeMap map[core.MachineType][]byte
}

func (e *DeployCode) Serialize() (io.Reader, error) {
	return serialize(e, TypeDeployCode)
}

func (e *DeployCode) GetReference() core.RecordRef {
	// XXX: ?
	return e.Request
}

func (e *DeployCode) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeployCode) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type ActivateClass struct {
	Domain  core.RecordRef
	Request core.RecordRef
}

func (e *ActivateClass) Serialize() (io.Reader, error) {
	return serialize(e, TypeActivateClass)
}

func (e *ActivateClass) GetReference() core.RecordRef {
	// XXX: ?
	return e.Request
}

func (e *ActivateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateClass) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type DeactivateClass struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
}

func (e *DeactivateClass) Serialize() (io.Reader, error) {
	return serialize(e, TypeDeactivateClass)
}

func (e *DeactivateClass) GetReference() core.RecordRef {
	return e.Class
}

func (e *DeactivateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeactivateClass) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type UpdateClass struct {
	Domain     core.RecordRef
	Request    core.RecordRef
	Class      core.RecordRef
	Code       core.RecordRef
	Migrations []core.RecordRef
}

func (e *UpdateClass) Serialize() (io.Reader, error) {
	return serialize(e, TypeUpdateClass)
}

func (e *UpdateClass) GetReference() core.RecordRef {
	return e.Class
}

func (e *UpdateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *UpdateClass) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type ActivateObject struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObject) Serialize() (io.Reader, error) {
	return serialize(e, TypeActivateObject)
}

func (e *ActivateObject) GetReference() core.RecordRef {
	return e.Class
}

func (e *ActivateObject) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateObject) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type ActivateObjectDelegate struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObjectDelegate) Serialize() (io.Reader, error) {
	return serialize(e, TypeActivateObjectDelegate)
}

func (e *ActivateObjectDelegate) GetReference() core.RecordRef {
	return e.Class
}

func (e *ActivateObjectDelegate) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateObjectDelegate) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

type DeactivateObject struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
}

func (e *DeactivateObject) Serialize() (io.Reader, error) {
	return serialize(e, TypeDeactivateObject)
}

func (e *DeactivateObject) GetReference() core.RecordRef {
	return e.Object
}

func (e *DeactivateObject) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeactivateObject) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}

// UpdateObject event for call of core.ArtifactManager.UpdateObj
type UpdateObject struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
	Memory  []byte
}

func (e *UpdateObject) Serialize() (io.Reader, error) {
	return serialize(e, TypeUpdateObject)
}

func (e *UpdateObject) GetReference() core.RecordRef {
	return e.Object
}

func (e *UpdateObject) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *UpdateObject) React(c core.Components) (core.Reaction, error) {
	return c.Ledger.HandleEvent(e)
}
