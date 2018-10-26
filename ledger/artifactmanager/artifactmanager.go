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
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
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
	ctx context.Context, msg core.Message,
) (*core.RecordID, error) {
	return m.setRecord(
		&record.CallRequest{
			Payload: message.MustSerializeBytes(msg),
		},
		*msg.Target(),
	)
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(
	ctx context.Context, code core.RecordRef,
) (core.CodeDescriptor, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.GetCode{Code: code},
	)

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
	ctx context.Context, head core.RecordRef, state *core.RecordID,
) (core.ClassDescriptor, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.GetClass{
			Head:  head,
			State: state,
		},
	)

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
	ctx context.Context, head core.RecordRef, state *core.RecordID, approved bool,
) (core.ObjectDescriptor, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.GetObject{
			Head:     head,
			State:    state,
			Approved: approved,
		},
	)

	if err != nil {
		return nil, err
	}

	switch r := genericReact.(type) {
	case *reply.Object:
		desc := ObjectDescriptor{
			am:           m,
			head:         r.Head,
			state:        r.State,
			class:        r.Class,
			childPointer: r.ChildPointer,
			memory:       r.Memory,
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
	ctx context.Context, head, asClass core.RecordRef,
) (*core.RecordRef, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.GetDelegate{
			Head:    head,
			AsClass: asClass,
		},
	)

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
	ctx context.Context, parent core.RecordRef, pulse *core.PulseNumber,
) (core.RefIterator, error) {
	return NewChildIterator(m.messageBus, parent, pulse, m.getChildrenChunkSize)
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	ctx context.Context, domain, request core.RecordRef, typeDec []byte,
) (*core.RecordID, error) {
	return m.setRecord(
		&record.TypeRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			TypeDeclaration: typeDec,
		},
		request,
	)
}

// DeployCode creates new code record in storage.
//
// Code records are used to activate class or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	ctx context.Context,
	domain core.RecordRef,
	request core.RecordRef,
	code []byte,
	machineType core.MachineType,
) (*core.RecordID, error) {
	return m.setRecord(
		&record.CodeRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			Code:        code,
			MachineType: machineType,
		},
		request,
	)
}

// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
//
// Request reference will be this class'es identifier and referred as "class head".
func (m *LedgerArtifactManager) ActivateClass(
	ctx context.Context, domain, request, code core.RecordRef, machineType core.MachineType,
) (*core.RecordID, error) {
	return m.updateClass(
		&record.ClassActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			ClassStateRecord: record.ClassStateRecord{
				MachineType: machineType,
				Code:        record.Core2Reference(code),
			},
		},
		request,
	)
}

// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
// the class. If class is already deactivated, an error should be returned.
//
// Deactivated class cannot be changed or instantiate objects.
func (m *LedgerArtifactManager) DeactivateClass(
	ctx context.Context,
	domain, request, class core.RecordRef, state core.RecordID,
) (*core.RecordID, error) {
	return m.updateClass(
		&record.DeactivationRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			PrevState: record.Bytes2ID(state[:]),
		},
		class,
	)
}

// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
// the class. Migrations are references to code records.
//
// Returned reference will be the latest class state (exact) reference. Migration code will be executed by VM to
// migrate objects memory in the order they appear in provided slice.
func (m *LedgerArtifactManager) UpdateClass(
	ctx context.Context,
	domain, request, class, code core.RecordRef, machineType core.MachineType, state core.RecordID,
) (*core.RecordID, error) {
	return m.updateClass(
		&record.ClassAmendRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			ClassStateRecord: record.ClassStateRecord{
				Code:        record.Core2Reference(code),
				MachineType: machineType,
			},
			PrevState: record.Bytes2ID(state[:]),
		},
		class,
	)
}

// ActivateObject creates activate object record in storage. Provided class reference will be used as objects class
// memory as memory of created object. If memory is not provided, the class default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObject(
	ctx context.Context,
	domain core.RecordRef,
	object core.RecordRef,
	class core.RecordRef,
	parent core.RecordRef,
	asDelegate bool,
	memory []byte,
) (core.ObjectDescriptor, error) {
	objectRef := record.Core2Reference(object)

	parendDesc, err := m.GetObject(ctx, parent, nil, false)
	if err != nil {
		return nil, err
	}

	obj, err := m.updateObject(
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: objectRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: memory,
			},
			Class:    record.Core2Reference(class),
			Parent:   record.Core2Reference(parent),
			Delegate: asDelegate,
		},
		object,
		&class,
	)
	if err != nil {
		return nil, err
	}

	var (
		prevChild *record.ID
		asClass   *core.RecordRef
	)
	if parendDesc.ChildPointer() != nil {
		c := record.Bytes2ID(parendDesc.ChildPointer()[:])
		prevChild = &c
	}
	if asDelegate {
		asClass = &class
	}
	_, err = m.registerChild(
		&record.ChildRecord{
			Ref:       objectRef,
			PrevChild: prevChild,
		},
		parent,
		object,
		asClass,
	)
	if err != nil {
		return nil, err
	}

	return &ObjectDescriptor{
		am:           m,
		head:         obj.Head,
		state:        obj.State,
		class:        obj.Class,
		childPointer: obj.ChildPointer,
		memory:       memory,
	}, nil
}

// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
// of the object. If object is already deactivated, an error should be returned.
//
// Deactivated object cannot be changed.
func (m *LedgerArtifactManager) DeactivateObject(
	ctx context.Context, domain, request core.RecordRef, object core.ObjectDescriptor,
) (*core.RecordID, error) {
	desc, err := m.updateObject(
		&record.DeactivationRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			PrevState: record.Bytes2ID(object.StateID()[:]),
		},
		*object.HeadRef(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &desc.State, nil
}

// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
// object. Provided memory well be the new object memory.
//
// Returned reference will be the latest object state (exact) reference.
func (m *LedgerArtifactManager) UpdateObject(
	ctx context.Context,
	domain, request core.RecordRef,
	object core.ObjectDescriptor,
	memory []byte,
) (core.ObjectDescriptor, error) {
	obj, err := m.updateObject(
		&record.ObjectAmendRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  record.Core2Reference(domain),
				Request: record.Core2Reference(request),
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: memory,
			},
			PrevState: record.Bytes2ID(object.StateID()[:]),
		},
		*object.HeadRef(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ObjectDescriptor{
		am:           m,
		head:         obj.Head,
		state:        obj.State,
		class:        obj.Class,
		childPointer: obj.ChildPointer,
		memory:       memory,
	}, nil
}

// RegisterValidation marks provided object state as approved or disapproved.
//
// When fetching object, validity can be specified.
func (m *LedgerArtifactManager) RegisterValidation(
	ctx context.Context, object core.RecordRef, state core.RecordID, isValid bool, validationMessages []core.Message,
) error {
	msg := message.ValidateRecord{
		Object:             object,
		State:              state,
		IsValid:            isValid,
		ValidationMessages: validationMessages,
	}
	_, err := m.messageBus.Send(ctx, &msg)
	if err != nil {
		return err
	}

	return nil
}

// RegisterResult saves VM method call result.
func (m *LedgerArtifactManager) RegisterResult(
	ctx core.Context, request core.RecordRef, payload []byte,
) (*core.RecordID, error) {
	return m.setRecord(
		&record.ResultRecord{
			Request: record.Core2Reference(request),
			Payload: payload,
		},
		request,
	)
}

func (m *LedgerArtifactManager) setRecord(rec record.Record, target core.RecordRef) (*core.RecordID, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.SetRecord{
			Record:    record.SerializeRecord(rec),
			TargetRef: target,
		},
	)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.ID)
	if !ok {
		return nil, ErrUnexpectedReply
	}

	return &react.ID, nil
}

func (m *LedgerArtifactManager) updateClass(rec record.Record, class core.RecordRef) (*core.RecordID, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.UpdateClass{
			Record: record.SerializeRecord(rec),
			Class:  class,
		},
	)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.ID)
	if !ok {
		return nil, ErrUnexpectedReply
	}

	return &react.ID, nil
}

func (m *LedgerArtifactManager) updateObject(
	rec record.Record, object core.RecordRef, class *core.RecordRef,
) (*reply.Object, error) {
	genericReact, err := m.messageBus.Send(
		context.TODO(),
		&message.UpdateObject{
			Record: record.SerializeRecord(rec),
			Object: object,
			Class:  class,
		},
	)

	if err != nil {
		return nil, err
	}

	rep, ok := genericReact.(*reply.Object)
	if !ok {
		return nil, ErrUnexpectedReply
	}

	return rep, nil
}

func (m *LedgerArtifactManager) registerChild(
	rec record.Record, parent, child core.RecordRef, asClass *core.RecordRef,
) (*core.RecordID, error) {
	genericReact, err := m.messageBus.Send(context.TODO(),
		&message.RegisterChild{
			Record:  record.SerializeRecord(rec),
			Parent:  parent,
			Child:   child,
			AsClass: asClass,
		},
	)

	if err != nil {
		return nil, err
	}

	react, ok := genericReact.(*reply.ID)
	if !ok {
		return nil, ErrUnexpectedReply
	}

	return &react.ID, nil
}
