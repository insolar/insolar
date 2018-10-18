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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/log"
)

const (
	getChildrenChunkSize = 10 * 1000
)

// LedgerArtifactManager provides concrete API to storage for processing module.
type LedgerArtifactManager struct {
	db         *storage.DB
	messageBus core.MessageBus

	getChildrenChunkSize int
}

// NewArtifactManger creates new manager instance.
func NewArtifactManger(db *storage.DB) (*LedgerArtifactManager, error) {
	return &LedgerArtifactManager{db: db, getChildrenChunkSize: getChildrenChunkSize}, nil
}

// Link links external components.
func (m *LedgerArtifactManager) Link(components core.Components) error {
	m.messageBus = components.MessageBus

	return nil
}

// GenesisRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (m *LedgerArtifactManager) GenesisRef() *core.RecordRef {
	return m.db.GenesisRef().CoreRef()
}

// RegisterRequest sends message for request registration,
// returns request record Ref if request successfully created or already exists.
func (m *LedgerArtifactManager) RegisterRequest(
	ctx core.Context, msg core.Message,
) (*core.RecordID, error) {
	rec := record.CallRequest{
		Payload: message.MustSerializeBytes(msg),
	}
	id, err := m.fetchID(&message.SetRecord{
		Record: record.SerializeRecord(&rec),
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(
	ctx core.Context, code core.RecordRef,
) (core.CodeDescriptor, error) {
	genericReact, err := m.messageBus.Send(&message.GetCode{
		Code: code,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.Code)
	if !ok {
		return nil, ErrUnexpectedReply
	}
	desc := CodeDescriptor{
		ref:         code,
		machineType: react.MachineType,
	}
	desc.cache.code = react.Code

	return &desc, nil
}

// GetClass returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
// provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetClass(
	ctx core.Context, head core.RecordRef, state *core.RecordRef,
) (core.ClassDescriptor, error) {
	genericReact, err := m.messageBus.Send(&message.GetClass{
		Head:  head,
		State: state,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.Class)
	if !ok {
		return nil, ErrUnexpectedReply
	}
	desc := ClassDescriptor{
		am:          m,
		head:        react.Head,
		state:       react.State,
		code:        react.Code,
		machineType: react.MachineType,
	}
	return &desc, nil
}

// GetObject returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
// provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetObject(
	ctx core.Context, head core.RecordRef, state *core.RecordRef,
) (core.ObjectDescriptor, error) {
	genericReact, err := m.messageBus.Send(&message.GetObject{
		Head:  head,
		State: state,
	})

	if err != nil {
		return nil, err
	}

	switch r := genericReact.(type) {
	case *reply.Object:
		desc := ObjectDescriptor{
			am:     m,
			head:   r.Head,
			state:  r.State,
			class:  r.Class,
			memory: r.Memory,
		}
		return &desc, nil
	case *reply.Error:
		return nil, r.Error()
	}

	return nil, ErrUnexpectedReply
}

// GetDelegate returns provided object's delegate reference for provided class.
//
// Object delegate should be previously created for this object. If object delegate does not exist, an error will
// be returned.
func (m *LedgerArtifactManager) GetDelegate(
	ctx core.Context, head, asClass core.RecordRef,
) (*core.RecordRef, error) {
	genericReact, err := m.messageBus.Send(&message.GetDelegate{
		Head:    head,
		AsClass: asClass,
	})

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.Delegate)
	if !ok {
		return nil, ErrUnexpectedReply
	}
	return &react.Head, nil
}

// GetChildren returns children iterator.
//
// During iteration children refs will be fetched from remote source (parent object).
func (m *LedgerArtifactManager) GetChildren(
	ctx core.Context, parent core.RecordRef, pulse *core.PulseNumber,
) (core.RefIterator, error) {
	return NewChildIterator(m.messageBus, parent, pulse, m.getChildrenChunkSize)
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	ctx core.Context, domain, request core.RecordRef, typeDec []byte,
) (*core.RecordRef, error) {
	return m.fetchReference(&message.DeclareType{
		Domain:  domain,
		Request: request,
		TypeDec: typeDec,
	})
}

// DeployCode creates new code record in storage.
//
// Code records are used to activate class or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	ctx core.Context,
	domain core.RecordRef,
	request core.RecordRef,
	code []byte,
	machineType core.MachineType,
) (*core.RecordRef, error) {
	return m.fetchReference(&message.DeployCode{
		Domain:      domain,
		Request:     request,
		Code:        code,
		MachineType: machineType,
	})
}

// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
//
// Request reference will be this class'es identifier and referred as "class head".
func (m *LedgerArtifactManager) ActivateClass(
	ctx core.Context, domain, request, code core.RecordRef,
) (*core.RecordID, error) {
	return m.fetchID(&message.ActivateClass{
		Domain:  domain,
		Request: request,
		Code:    code,
	})
}

// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
// the class. If class is already deactivated, an error should be returned.
//
// Deactivated class cannot be changed or instantiate objects.
func (m *LedgerArtifactManager) DeactivateClass(
	ctx core.Context,
	domain, request, class core.RecordRef,
) (*core.RecordID, error) {
	return m.fetchID(&message.DeactivateClass{
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
	ctx core.Context,
	domain, request, class, code core.RecordRef, migrations []core.RecordRef,
) (*core.RecordID, error) {
	return m.fetchID(&message.UpdateClass{
		Domain:     domain,
		Request:    request,
		Class:      class,
		Code:       code,
		Migrations: migrations,
	})
}

// ActivateObject creates activate object record in storage. Provided class reference will be used as objects class
// memory as memory of created object. If memory is not provided, the class default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObject(
	ctx core.Context,
	domain core.RecordRef,
	object core.RecordRef,
	class core.RecordRef,
	parent core.RecordRef,
	memory []byte,
) (*core.RecordID, error) {
	start := time.Now()
	objID, err := m.fetchID(&message.ActivateObject{
		Domain:  domain,
		Request: object,
		Class:   class,
		Parent:  parent,
		Memory:  memory,
	})

	if err != nil {
		return nil, err
	}

	_, err = m.fetchID(&message.RegisterChild{
		Parent: parent,
		Child:  object,
	})
	if err != nil {
		return nil, err
	}
	log.Debugf("Inside ActivateObject: class - '%s', parent - %s, time - %s", class, parent, time.Since(start))
	return objID, nil
}

// ActivateObjectDelegate is similar to ActivateObj but it created object will be parent's delegate of provided class.
func (m *LedgerArtifactManager) ActivateObjectDelegate(
	ctx core.Context,
	domain core.RecordRef,
	request core.RecordRef,
	class core.RecordRef,
	parent core.RecordRef,
	memory []byte,
) (*core.RecordID, error) {
	return m.fetchID(&message.ActivateObjectDelegate{
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
	ctx core.Context,
	domain, request, object core.RecordRef,
) (*core.RecordID, error) {
	return m.fetchID(&message.DeactivateObject{
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
	ctx core.Context,
	domain core.RecordRef,
	request core.RecordRef,
	object core.RecordRef,
	memory []byte,
) (*core.RecordID, error) {
	return m.fetchID(&message.UpdateObject{
		Domain:  domain,
		Request: request,
		Object:  object,
		Memory:  memory,
	})
}

func (m *LedgerArtifactManager) fetchReference(ev core.Message) (*core.RecordRef, error) {
	genericReact, err := m.messageBus.Send(ev)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.Reference)
	if !ok {
		return nil, ErrUnexpectedReply
	}
	return &react.Ref, nil
}

func (m *LedgerArtifactManager) fetchID(msg core.Message) (*core.RecordID, error) {
	genericReact, err := m.messageBus.Send(msg)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.ID)
	if !ok {
		return nil, ErrUnexpectedReply
	}
	return &react.ID, nil
}
