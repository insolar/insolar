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
	"github.com/pkg/errors"
)

type GetCode struct {
	Code        core.RecordRef
	MachinePref []core.MachineType
}

func (e *GetCode) Serialize() (io.Reader, error) {
	return serialize(e, GetCodeType)
}

func (e *GetCode) GetReference() core.RecordRef {
	return e.Code
}

func (e *GetCode) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetCode) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}
	return ledger.GetArtifactManager().HandleEvent(e)
}

type GetLatestClass struct {
	Head core.RecordRef
}

func (e *GetLatestClass) Serialize() (io.Reader, error) {
	return serialize(e, GetLatestClassType)
}

func (e *GetLatestClass) GetReference() core.RecordRef {
	return e.Head
}

func (e *GetLatestClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetLatestClass) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}
	return ledger.GetArtifactManager().HandleEvent(e)
}

type GetLatestObj struct {
	Head core.RecordRef
}

func (e *GetLatestObj) Serialize() (io.Reader, error) {
	return serialize(e, GetLatestObjType)
}

func (e *GetLatestObj) GetReference() core.RecordRef {
	return e.Head
}

func (e *GetLatestObj) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetLatestObj) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}
	return ledger.GetArtifactManager().HandleEvent(e)
}

type DeclareType struct {
	Domain  core.RecordRef
	Request core.RecordRef
	TypeDec []byte
}

func (e *DeclareType) Serialize() (io.Reader, error) {
	return serialize(e, DeclareTypeType)
}

func (e *DeclareType) GetReference() core.RecordRef {
	// XXX: ?
	return e.Request
}

func (e *DeclareType) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeclareType) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}
	return ledger.GetArtifactManager().HandleEvent(e)
}

type DeployCode struct {
	Domain  core.RecordRef
	Request core.RecordRef
	CodeMap map[core.MachineType][]byte
}

func (e *DeployCode) Serialize() (io.Reader, error) {
	return serialize(e, DeployCodeType)
}

func (e *DeployCode) GetReference() core.RecordRef {
	// XXX: ?
	return e.Request
}

func (e *DeployCode) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeployCode) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type ActivateClass struct {
	Domain  core.RecordRef
	Request core.RecordRef
}

func (e *ActivateClass) Serialize() (io.Reader, error) {
	return serialize(e, ActivateClassType)
}

func (e *ActivateClass) GetReference() core.RecordRef {
	// XXX: ?
	return e.Request
}

func (e *ActivateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateClass) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type DeactivateClass struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
}

func (e *DeactivateClass) Serialize() (io.Reader, error) {
	return serialize(e, DeactivateClassType)
}

func (e *DeactivateClass) GetReference() core.RecordRef {
	return e.Class
}

func (e *DeactivateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeactivateClass) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type UpdateClass struct {
	Domain        core.RecordRef
	Request       core.RecordRef
	Class         core.RecordRef
	Code          core.RecordRef
	MigrationRefs []core.RecordRef
}

func (e *UpdateClass) Serialize() (io.Reader, error) {
	return serialize(e, UpdateClassType)
}

func (e *UpdateClass) GetReference() core.RecordRef {
	// XXX: or Code ?
	return e.Class
}

func (e *UpdateClass) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *UpdateClass) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type ActivateObj struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObj) Serialize() (io.Reader, error) {
	return serialize(e, ActivateObjType)
}

func (e *ActivateObj) GetReference() core.RecordRef {
	return e.Class
}

func (e *ActivateObj) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateObj) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type ActivateObjDelegate struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

func (e *ActivateObjDelegate) Serialize() (io.Reader, error) {
	return serialize(e, ActivateObjDelegateType)
}

func (e *ActivateObjDelegate) GetReference() core.RecordRef {
	return e.Class
}

func (e *ActivateObjDelegate) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *ActivateObjDelegate) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

type DeactivateObj struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Obj     core.RecordRef
}

func (e *DeactivateObj) Serialize() (io.Reader, error) {
	return serialize(e, DeactivateObjType)
}

func (e *DeactivateObj) GetReference() core.RecordRef {
	return e.Obj
}

func (e *DeactivateObj) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *DeactivateObj) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}

// UpdateObj event for call of core.ArtifactManager.UpdateObj
type UpdateObj struct {
	Domain  core.RecordRef
	Request core.RecordRef
	Obj     core.RecordRef
	Memory  []byte
}

func (e *UpdateObj) Serialize() (io.Reader, error) {
	return serialize(e, UpdateObjType)
}

func (e *UpdateObj) GetReference() core.RecordRef {
	return e.Obj
}

func (e *UpdateObj) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *UpdateObj) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}
