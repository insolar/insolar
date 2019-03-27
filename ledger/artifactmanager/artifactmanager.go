/*
 *    Copyright 2019 Insolar Technologies
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
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
)

const (
	getChildrenChunkSize = 10 * 1000
	jetMissRetryCount    = 10
)

// LedgerArtifactManager provides concrete API to storage for processing module.
type LedgerArtifactManager struct {
	DB           storage.DBContext    `inject:""`
	GenesisState storage.GenesisState `inject:""`
	JetStorage   storage.JetStorage   `inject:""`

	DefaultBus                 core.MessageBus                 `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	PulseStorage               core.PulseStorage               `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`

	getChildrenChunkSize int
	senders              *ledgerArtifactSenders
}

// State returns hash state for artifact manager.
func (m *LedgerArtifactManager) State() ([]byte, error) {
	// This is a temporary stab to simulate real hash.
	return m.PlatformCryptographyScheme.IntegrityHasher().Hash([]byte{1, 2, 3}), nil
}

// NewArtifactManger creates new manager instance.
func NewArtifactManger() *LedgerArtifactManager {
	return &LedgerArtifactManager{
		getChildrenChunkSize: getChildrenChunkSize,
		senders:              newLedgerArtifactSenders(),
	}
}

// GenesisRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (m *LedgerArtifactManager) GenesisRef() *core.RecordRef {
	return m.GenesisState.GenesisRef()
}

// RegisterRequest sends message for request registration,
// returns request record Ref if request successfully created or already exists.
func (m *LedgerArtifactManager) RegisterRequest(
	ctx context.Context, obj core.RecordRef, parcel core.Parcel,
) (*core.RecordID, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.RegisterRequest")
	instrumenter := instrument(ctx, "RegisterRequest").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	rec := &record.RequestRecord{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: m.PlatformCryptographyScheme.IntegrityHasher().Hash(message.MustSerializeBytes(parcel.Message())),
		Object:      *obj.Record(),
	}
	recID := record.NewRecordIDFromRecord(
		m.PlatformCryptographyScheme,
		currentPN,
		rec)
	recRef := core.NewRecordRef(*parcel.DefaultTarget().Domain(), *recID)
	id, err := m.setRecord(
		ctx,
		rec,
		*recRef,
		currentPN,
	)
	return id, errors.Wrap(err, "[ RegisterRequest ] ")
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(
	ctx context.Context, code core.RecordRef,
) (core.CodeDescriptor, error) {
	var err error
	instrumenter := instrument(ctx, "GetCode").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetCode")
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(
		bus.Send,
		m.senders.cachedSender(m.PlatformCryptographyScheme),
		followRedirectSender(bus),
		retryJetSender(currentPN, m.JetStorage),
	)

	genericReact, err := sender(ctx, &message.GetCode{Code: code}, nil)

	if err != nil {
		return nil, err
	}

	switch rep := genericReact.(type) {
	case *reply.Code:
		desc := CodeDescriptor{
			ctx:         ctx,
			ref:         code,
			machineType: rep.MachineType,
			code:        rep.Code,
		}
		return &desc, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("GetCode: unexpected reply: %#v", rep)
	}
}

// GetObject returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
// provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetObject(
	ctx context.Context,
	head core.RecordRef,
	state *core.RecordID,
	approved bool,
) (core.ObjectDescriptor, error) {
	var (
		desc *ObjectDescriptor
		err  error
	)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.Getobject")
	instrumenter := instrument(ctx, "GetObject").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		if err != nil && err == ErrObjectDeactivated {
			err = nil // megahack: threat it 2xx
		}
		instrumenter.end()
	}()

	getObjectMsg := &message.GetObject{
		Head:     head,
		State:    state,
		Approved: approved,
	}

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(
		bus.Send,
		followRedirectSender(bus),
		retryJetSender(currentPN, m.JetStorage),
	)

	genericReact, err := sender(ctx, getObjectMsg, nil)
	if err != nil {
		return nil, err
	}

	switch r := genericReact.(type) {
	case *reply.Object:
		desc = &ObjectDescriptor{
			ctx:          ctx,
			am:           m,
			head:         r.Head,
			state:        r.State,
			prototype:    r.Prototype,
			isPrototype:  r.IsPrototype,
			childPointer: r.ChildPointer,
			memory:       r.Memory,
			parent:       r.Parent,
		}
		return desc, err
	case *reply.Error:
		return nil, r.Error()
	default:
		return nil, fmt.Errorf("GetObject: unexpected reply: %#v", genericReact)
	}
}

// GetPendingRequest returns an unclosed pending request
// It takes an id from current LME
// Then goes either to a light node or heavy node
func (m *LedgerArtifactManager) GetPendingRequest(ctx context.Context, objectID core.RecordID) (core.Parcel, error) {
	var err error
	instrumenter := instrument(ctx, "GetRegisterRequest").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetRegisterRequest")
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(
		bus.Send,
		retryJetSender(currentPN, m.JetStorage),
	)

	genericReply, err := sender(ctx, &message.GetPendingRequestID{
		ObjectID: objectID,
	}, nil)
	if err != nil {
		return nil, err
	}

	var requestIDReply *reply.ID
	switch r := genericReply.(type) {
	case *reply.ID:
		requestIDReply = r
	case *reply.Error:
		return nil, r.Error()
	default:
		return nil, fmt.Errorf("GetPendingRequest: unexpected reply: %#v", requestIDReply)
	}

	node, err := m.JetCoordinator.NodeForObject(ctx, objectID, currentPN, requestIDReply.ID.Pulse())

	if err != nil {
		return nil, err
	}

	sender = BuildSender(
		bus.Send,
		retryJetSender(currentPN, m.JetStorage),
	)
	genericReply, err = sender(
		ctx,
		&message.GetRequest{
			Request: requestIDReply.ID,
		}, &core.MessageSendOptions{
			Receiver: node,
		},
	)
	if err != nil {
		return nil, err
	}

	switch r := genericReply.(type) {
	case *reply.Request:
		rec := record.DeserializeRecord(r.Record)
		castedRecord, ok := rec.(*record.RequestRecord)
		if !ok {
			return nil, fmt.Errorf("GetPendingRequest: unexpected message: %#v", r)
		}

		return message.DeserializeParcel(bytes.NewBuffer(castedRecord.Parcel))
	case *reply.Error:
		return nil, r.Error()
	default:
		return nil, fmt.Errorf("GetPendingRequest: unexpected reply: %#v", requestIDReply)
	}
}

// HasPendingRequests returns true if object has unclosed requests.
func (m *LedgerArtifactManager) HasPendingRequests(
	ctx context.Context,
	object core.RecordRef,
) (bool, error) {

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return false, err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(
		bus.Send,
		retryJetSender(currentPN, m.JetStorage),
	)

	genericReact, err := sender(ctx, &message.GetPendingRequests{Object: object}, nil)

	if err != nil {
		return false, err
	}

	switch rep := genericReact.(type) {
	case *reply.HasPendingRequests:
		return rep.Has, nil
	case *reply.Error:
		return false, rep.Error()
	default:
		return false, fmt.Errorf("HasPendingRequests: unexpected reply: %#v", rep)
	}
}

// GetDelegate returns provided object's delegate reference for provided prototype.
//
// Object delegate should be previously created for this object. If object delegate does not exist, an error will
// be returned.
func (m *LedgerArtifactManager) GetDelegate(
	ctx context.Context, head, asType core.RecordRef,
) (*core.RecordRef, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetDelegate")
	instrumenter := instrument(ctx, "GetDelegate").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, followRedirectSender(bus), retryJetSender(currentPN, m.JetStorage))
	genericReact, err := sender(ctx, &message.GetDelegate{
		Head:   head,
		AsType: asType,
	}, nil)
	if err != nil {
		return nil, err
	}

	switch rep := genericReact.(type) {
	case *reply.Delegate:
		return &rep.Head, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("GetDelegate: unexpected reply: %#v", rep)
	}
}

// GetChildren returns children iterator.
//
// During iteration children refs will be fetched from remote source (parent object).
func (m *LedgerArtifactManager) GetChildren(
	ctx context.Context, parent core.RecordRef, pulse *core.PulseNumber,
) (core.RefIterator, error) {
	var err error

	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetChildren")
	instrumenter := instrument(ctx, "GetChildren").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, followRedirectSender(bus), retryJetSender(currentPN, m.JetStorage))
	iter, err := NewChildIterator(ctx, sender, parent, pulse, m.getChildrenChunkSize)
	return iter, err
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	ctx context.Context, domain, request core.RecordRef, typeDec []byte,
) (*core.RecordID, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.DeclareType")
	instrumenter := instrument(ctx, "DeclareType").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	recid, err := m.setRecord(
		ctx,
		&record.TypeRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  domain,
				Request: request,
			},
			TypeDeclaration: typeDec,
		},
		request,
		currentPN,
	)
	return recid, err
}

// DeployCode creates new code record in storage.
//
// CodeRef records are used to activate prototype or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	ctx context.Context,
	domain core.RecordRef,
	request core.RecordRef,
	code []byte,
	machineType core.MachineType,
) (*core.RecordID, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.DeployCode")
	instrumenter := instrument(ctx, "DeployCode").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	codeRec := &record.CodeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domain,
			Request: request,
		},
		Code:        record.CalculateIDForBlob(m.PlatformCryptographyScheme, currentPN, code),
		MachineType: machineType,
	}
	codeID := record.NewRecordIDFromRecord(m.PlatformCryptographyScheme, currentPN, codeRec)
	codeRef := core.NewRecordRef(*domain.Record(), *codeID)

	_, err = m.setBlob(ctx, code, *codeRef, currentPN)
	if err != nil {
		return nil, err
	}
	id, err := m.setRecord(
		ctx,
		codeRec,
		*codeRef,
		currentPN,
	)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivatePrototype(
	ctx context.Context,
	domain, object, parent, code core.RecordRef,
	memory []byte,
) (core.ObjectDescriptor, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.ActivatePrototype")
	instrumenter := instrument(ctx, "ActivatePrototype").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()
	desc, err := m.activateObject(ctx, domain, object, code, true, parent, false, memory)
	return desc, err
}

// ActivateObject creates activate object record in storage. Provided prototype reference will be used as objects prototype
// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObject(
	ctx context.Context,
	domain, object, parent, prototype core.RecordRef,
	asDelegate bool,
	memory []byte,
) (core.ObjectDescriptor, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.ActivateObject")
	instrumenter := instrument(ctx, "ActivateObject").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()
	desc, err := m.activateObject(ctx, domain, object, prototype, false, parent, asDelegate, memory)
	return desc, err
}

// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
// of the object. If object is already deactivated, an error should be returned.
//
// Deactivated object cannot be changed.
func (m *LedgerArtifactManager) DeactivateObject(
	ctx context.Context, domain, request core.RecordRef, object core.ObjectDescriptor,
) (*core.RecordID, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.DeactivateObject")
	instrumenter := instrument(ctx, "DeactivateObject").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)

	desc, err := m.sendUpdateObject(
		ctx,
		&record.DeactivationRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  domain,
				Request: request,
			},
			PrevState: *object.StateID(),
		},
		*object.HeadRef(),
		nil,
		currentPN,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deactivate object")
	}
	return &desc.State, nil
}

// UpdatePrototype creates amend object record in storage. Provided reference should be a reference to the head of the
// prototype. Provided memory well be the new object memory.
//
// Returned reference will be the latest object state (exact) reference.
func (m *LedgerArtifactManager) UpdatePrototype(
	ctx context.Context,
	domain, request core.RecordRef,
	object core.ObjectDescriptor,
	memory []byte,
	code *core.RecordRef,
) (core.ObjectDescriptor, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.UpdatePrototype")
	instrumenter := instrument(ctx, "UpdatePrototype").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	if !object.IsPrototype() {
		err = errors.New("object is not a prototype")
		return nil, err
	}
	desc, err := m.updateObject(ctx, domain, request, object, code, memory)
	return desc, err
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
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.UpdateObject")
	instrumenter := instrument(ctx, "UpdateObject").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	if object.IsPrototype() {
		err = errors.New("object is not an instance")
		return nil, err
	}
	desc, err := m.updateObject(ctx, domain, request, object, nil, memory)
	return desc, err
}

// RegisterValidation marks provided object state as approved or disapproved.
//
// When fetching object, validity can be specified.
func (m *LedgerArtifactManager) RegisterValidation(
	ctx context.Context,
	object core.RecordRef,
	state core.RecordID,
	isValid bool,
	validationMessages []core.Message,
) error {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.RegisterValidation")
	instrumenter := instrument(ctx, "RegisterValidation").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	msg := message.ValidateRecord{
		Object:             object,
		State:              state,
		IsValid:            isValid,
		ValidationMessages: validationMessages,
	}

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return err
	}

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, retryJetSender(currentPN, m.JetStorage))
	_, err = sender(ctx, &msg, nil)

	return err
}

// RegisterResult saves VM method call result.
func (m *LedgerArtifactManager) RegisterResult(
	ctx context.Context, object, request core.RecordRef, payload []byte,
) (*core.RecordID, error) {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.RegisterResult")
	instrumenter := instrument(ctx, "RegisterResult").err(&err)
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	recid, err := m.setRecord(
		ctx,
		&record.ResultRecord{
			Object:  *object.Record(),
			Request: request,
			Payload: payload,
		},
		request,
		currentPN,
	)
	return recid, err
}

// pulse returns current PulseNumber for artifact manager
func (m *LedgerArtifactManager) pulse(ctx context.Context) (pn core.PulseNumber, err error) {
	pulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return
	}

	pn = pulse.PulseNumber
	return
}

func (m *LedgerArtifactManager) activateObject(
	ctx context.Context,
	domain core.RecordRef,
	object core.RecordRef,
	prototype core.RecordRef,
	isPrototype bool,
	parent core.RecordRef,
	asDelegate bool,
	memory []byte,
) (core.ObjectDescriptor, error) {
	parentDesc, err := m.GetObject(ctx, parent, nil, false)
	if err != nil {
		return nil, err
	}
	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	obj, err := m.sendUpdateObject(
		ctx,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  domain,
				Request: object,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory:      record.CalculateIDForBlob(m.PlatformCryptographyScheme, currentPN, memory),
				Image:       prototype,
				IsPrototype: isPrototype,
			},
			Parent:     parent,
			IsDelegate: asDelegate,
		},
		object,
		memory,
		currentPN,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate")
	}

	var (
		prevChild *core.RecordID
		asType    *core.RecordRef
	)
	if parentDesc.ChildPointer() != nil {
		prevChild = parentDesc.ChildPointer()
	}
	if asDelegate {
		asType = &prototype
	}
	_, err = m.registerChild(
		ctx,
		&record.ChildRecord{
			Ref:       object,
			PrevChild: prevChild,
		},
		parent,
		object,
		asType,
		currentPN,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register as child while activating")
	}

	return &ObjectDescriptor{
		ctx:          ctx,
		am:           m,
		head:         obj.Head,
		state:        obj.State,
		prototype:    obj.Prototype,
		childPointer: obj.ChildPointer,
		memory:       memory,
		parent:       obj.Parent,
	}, nil
}

func (m *LedgerArtifactManager) updateObject(
	ctx context.Context,
	domain, request core.RecordRef,
	object core.ObjectDescriptor,
	code *core.RecordRef,
	memory []byte,
) (core.ObjectDescriptor, error) {
	var (
		image *core.RecordRef
		err   error
	)
	if object.IsPrototype() {
		if code != nil {
			image = code
		} else {
			image, err = object.Code()
		}
	} else {
		image, err = object.Prototype()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	obj, err := m.sendUpdateObject(
		ctx,
		&record.ObjectAmendRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  domain,
				Request: request,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Image:       *image,
				IsPrototype: object.IsPrototype(),
			},
			PrevState: *object.StateID(),
		},
		*object.HeadRef(),
		memory,
		currentPN,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	return &ObjectDescriptor{
		ctx:          ctx,
		am:           m,
		head:         obj.Head,
		state:        obj.State,
		prototype:    obj.Prototype,
		childPointer: obj.ChildPointer,
		memory:       memory,
		parent:       obj.Parent,
	}, nil
}

func (m *LedgerArtifactManager) setRecord(
	ctx context.Context,
	rec record.Record,
	target core.RecordRef,
	currentPN core.PulseNumber,
) (*core.RecordID, error) {
	bus := core.MessageBusFromContext(ctx, m.DefaultBus)

	sender := BuildSender(bus.Send, retryJetSender(currentPN, m.JetStorage))
	genericReply, err := sender(ctx, &message.SetRecord{
		Record:    record.SerializeRecord(rec),
		TargetRef: target,
	}, nil)

	if err != nil {
		return nil, err
	}

	switch rep := genericReply.(type) {
	case *reply.ID:
		return &rep.ID, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("setRecord: unexpected reply: %#v", rep)
	}
}

func (m *LedgerArtifactManager) setBlob(
	ctx context.Context,
	blob []byte,
	target core.RecordRef,
	currentPN core.PulseNumber,
) (*core.RecordID, error) {

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, retryJetSender(currentPN, m.JetStorage))
	genericReact, err := sender(ctx, &message.SetBlob{
		Memory:    blob,
		TargetRef: target,
	}, nil)

	if err != nil {
		return nil, err
	}

	switch rep := genericReact.(type) {
	case *reply.ID:
		return &rep.ID, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("setBlob: unexpected reply: %#v", rep)
	}
}

func (m *LedgerArtifactManager) sendUpdateObject(
	ctx context.Context,
	rec record.Record,
	object core.RecordRef,
	memory []byte,
	currentPN core.PulseNumber,
) (*reply.Object, error) {
	// TODO: @andreyromancev. 14.01.19. Uncomment when message streaming or validation is ready.
	// genericRep, err := sendAndRetryJet(ctx, m.bus(ctx), m.db, &message.SetBlob{
	// 	TargetRef: object,
	// 	Memory:    memory,
	// }, currentPulse, jetMissRetryCount, nil)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to save object's memory blob")
	// }
	// if _, ok := genericRep.(*reply.ID); !ok {
	// 	return nil, fmt.Errorf("unexpected reply: %#v\n", genericRep)
	// }

	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, retryJetSender(currentPN, m.JetStorage))
	genericReply, err := sender(
		ctx,
		&message.UpdateObject{
			Record: record.SerializeRecord(rec),
			Object: object,
			Memory: memory,
		}, nil)

	if err != nil {
		return nil, errors.Wrap(err, "UpdateObject message failed")
	}

	switch rep := genericReply.(type) {
	case *reply.Object:
		return rep, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("sendUpdateObject: unexpected reply: %#v", rep)
	}
}

func (m *LedgerArtifactManager) registerChild(
	ctx context.Context,
	rec record.Record,
	parent core.RecordRef,
	child core.RecordRef,
	asType *core.RecordRef,
	currentPN core.PulseNumber,
) (*core.RecordID, error) {
	bus := core.MessageBusFromContext(ctx, m.DefaultBus)
	sender := BuildSender(bus.Send, retryJetSender(currentPN, m.JetStorage))
	genericReact, err := sender(ctx, &message.RegisterChild{
		Record: record.SerializeRecord(rec),
		Parent: parent,
		Child:  child,
		AsType: asType,
	}, nil)

	if err != nil {
		return nil, err
	}

	switch rep := genericReact.(type) {
	case *reply.ID:
		return &rep.ID, nil
	case *reply.Error:
		return nil, rep.Error()
	default:
		return nil, fmt.Errorf("registerChild: unexpected reply: %#v", rep)
	}
}
