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

package artifactmanager

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
	"github.com/insolar/insolar/ledger/storage"
)

// LedgerArtifactManager provides concrete API to storage for processing module.
type LedgerArtifactManager struct {
	db       *storage.DB
	eventBus core.EventBus
}

// NewArtifactManger creates new manager instance.
func NewArtifactManger(db *storage.DB) (*LedgerArtifactManager, error) {
	return &LedgerArtifactManager{db: db}, nil
}

// Link links external components components.
func (m *LedgerArtifactManager) Link(components core.Components) error {
	m.eventBus = components.EventBus

	return nil
}

// RootRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (m *LedgerArtifactManager) RootRef() *core.RecordRef {
	return m.db.RootRef().CoreRef()
}

// GetCode returns code from code record by provided reference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(
	code core.RecordRef, machinePref []core.MachineType,
) (core.CodeDescriptor, error) {
	genericReact, err := m.eventBus.Dispatch(&event.GetCode{
		Code:        code,
		MachinePref: machinePref,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reaction.Code)
	if !ok {
		return nil, ErrUnexpectedReaction
	}
	desc := CodeDescriptor{
		machinePref: machinePref,
		ref:         code,

		machineType: react.MachineType,
		code:        react.Code,
	}

	return &desc, nil
}

// GetClass returns descriptor for latest state of the class known to storage.
// If the class is deactivated, an error should be returned.
//
// Returned descriptor will provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetClass(head core.RecordRef, state *core.RecordRef) (core.ClassDescriptor, error) {
	genericReact, err := m.eventBus.Dispatch(&event.GetClass{
		Head:  head,
		State: state,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reaction.Class)
	if !ok {
		return nil, ErrUnexpectedReaction
	}
	desc := ClassDescriptor{
		am:    m,
		head:  react.Head,
		state: react.State,
		code:  react.Code,
	}
	return &desc, nil
}

// GetObject returns descriptor for latest state of the object known to storage.
// If the object or the class is deactivated, an error should be returned.
//
// Returned descriptor will provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetObject(head core.RecordRef, state *core.RecordRef) (core.ObjectDescriptor, error) {
	genericReact, err := m.eventBus.Dispatch(&event.GetObject{
		Head:  head,
		State: state,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reaction.Object)
	if !ok {
		return nil, ErrUnexpectedReaction
	}
	desc := ObjectDescriptor{
		am:       m,
		head:     react.Head,
		state:    react.State,
		class:    react.Class,
		memory:   react.Memory,
		children: react.Children,
	}
	return &desc, nil
}

// GetDelegate returns provided object's delegate reference for provided class.
//
// Object delegate should be previously created for this object. If object delegate does not exist, an error will
// be returned.
func (m *LedgerArtifactManager) GetDelegate(head, asClass core.RecordRef) (*core.RecordRef, error) {
	genericReact, err := m.eventBus.Dispatch(&event.GetDelegate{
		Head:    head,
		AsClass: asClass,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reaction.Delegate)
	if !ok {
		return nil, ErrUnexpectedReaction
	}
	return &react.Head, nil
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	domain, request core.RecordRef, typeDec []byte,
) (*core.RecordRef, error) {
	return m.sendReference(&event.DeclareType{
		Domain:  domain,
		Request: request,
		TypeDec: typeDec,
	})
}

// DeployCode creates new code record in storage.
//
// Code records are used to activate class or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	domain, request core.RecordRef, codeMap map[core.MachineType][]byte,
) (*core.RecordRef, error) {
	return m.sendReference(&event.DeployCode{
		Domain:  domain,
		Request: request,
		CodeMap: codeMap,
	})
}

// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
//
// Activation reference will be this class'es identifier and referred as "class head".
func (m *LedgerArtifactManager) ActivateClass(
	domain, request core.RecordRef,
) (*core.RecordRef, error) {
	return m.sendReference(&event.ActivateClass{
		Domain:  domain,
		Request: request,
	})
}

// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
// the class. If class is already deactivated, an error should be returned.
//
// Deactivated class cannot be changed or instantiate objects.
func (m *LedgerArtifactManager) DeactivateClass(
	domain, request, class core.RecordRef,
) (*core.RecordRef, error) {
	return m.sendReference(&event.DeactivateClass{
		Domain:  domain,
		Request: request,
		Class:   class,
	})
}

// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
// the class. Migrations are references to code records.
//
// Returned reference will be the latest class state (exact) reference. Migration code will be executed by VM to
// migrate objects memory in the order they appear in provided slice.
func (m *LedgerArtifactManager) UpdateClass(
	domain, request, class, code core.RecordRef, migrations []core.RecordRef,
) (*core.RecordRef, error) {
	return m.sendReference(&event.UpdateClass{
		Domain:     domain,
		Request:    request,
		Class:      class,
		Code:       code,
		Migrations: migrations,
	})
}

// ActivateObject creates activate object record in storage. Provided class reference will be used as objects class
// memory as memory of crated object. If memory is not provided, the class default memory will be used.
//
// Activation reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObject(
	domain, request, class, parent core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	return m.sendReference(&event.ActivateObject{
		Domain:  domain,
		Request: request,
		Class:   class,
		Parent:  parent,
		Memory:  memory,
	})
}

// ActivateObjectDelegate is similar to ActivateObj but it created object will be parent's delegate of provided class.
func (m *LedgerArtifactManager) ActivateObjectDelegate(
	domain, request, class, parent core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	return m.sendReference(&event.ActivateObjectDelegate{
		Domain:  domain,
		Request: request,
		Class:   class,
		Parent:  parent,
		Memory:  memory,
	})
}

// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
// of the object. If object is already deactivated, an error should be returned.
//
// Deactivated object cannot be changed.
func (m *LedgerArtifactManager) DeactivateObject(
	domain, request, object core.RecordRef,
) (*core.RecordRef, error) {
	return m.sendReference(&event.DeactivateObject{
		Domain:  domain,
		Request: request,
		Object:  object,
	})
}

// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
// object. Provided memory well be the new object memory.
//
// Returned reference will be the latest object state (exact) reference.
func (m *LedgerArtifactManager) UpdateObject(
	domain, request, object core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	return m.sendReference(&event.UpdateObject{
		Domain:  domain,
		Request: request,
		Object:  object,
		Memory:  memory,
	})
}

func (m *LedgerArtifactManager) sendReference(ev core.Event) (*core.RecordRef, error) {
	genericReact, err := m.eventBus.Dispatch(ev)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reaction.Reference)
	if !ok {
		return nil, ErrUnexpectedReaction
	}
	return &react.Ref, nil
}
