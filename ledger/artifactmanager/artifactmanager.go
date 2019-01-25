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
	"fmt"
	"sync"

	"github.com/insolar/insolar/instrumentation/instracer"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
)

const (
	getChildrenChunkSize = 10 * 1000
	jetMissRetryCount    = 10
)

// LedgerArtifactManager provides concrete API to storage for processing module.
type LedgerArtifactManager struct {
	db                         *storage.DB
	DefaultBus                 core.MessageBus                 `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	PulseStorage               core.PulseStorage               `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`

	codeCacheLock *sync.Mutex
	codeCache     map[core.RecordRef]*cacheEntry

	getChildrenChunkSize int
}

type cacheEntry struct {
	sync.Mutex
	desc core.CodeDescriptor
}

// State returns hash state for artifact manager.
func (m *LedgerArtifactManager) State() ([]byte, error) {
	// This is a temporary stab to simulate real hash.
	return m.PlatformCryptographyScheme.IntegrityHasher().Hash([]byte{1, 2, 3}), nil
}

// NewArtifactManger creates new manager instance.
func NewArtifactManger(db *storage.DB) *LedgerArtifactManager {
	return &LedgerArtifactManager{
		db:                   db,
		getChildrenChunkSize: getChildrenChunkSize,
		codeCacheLock:        &sync.Mutex{},
		codeCache:            make(map[core.RecordRef]*cacheEntry),
	}
}

// GenesisRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (m *LedgerArtifactManager) GenesisRef() *core.RecordRef {
	return m.db.GenesisRef()
}

// RegisterRequest sends message for request registration,
// returns request record Ref if request successfully created or already exists.
func (m *LedgerArtifactManager) RegisterRequest(
	ctx context.Context, obj core.RecordRef, parcel core.Parcel,
) (*core.RecordID, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.RegisterRequest starts ...")
	var err error
	defer instrument(ctx, "RegisterRequest").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	rec := record.RequestRecord{
		Payload: message.MustSerializeBytes(parcel.Message()),
		Object:  *obj.Record(),
	}
	recID := record.NewRecordIDFromRecord(m.PlatformCryptographyScheme, currentPulse.PulseNumber, &rec)
	recRef := core.NewRecordRef(*parcel.DefaultTarget().Domain(), *recID)
	id, err := m.setRecord(
		ctx,
		&rec,
		*recRef,
		*currentPulse,
	)
	return id, errors.Wrap(err, "[ RegisterRequest ] ")
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(
	ctx context.Context, code core.RecordRef,
) (core.CodeDescriptor, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.GetCode starts ...")
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetCode sendAndRetryJet")
	defer span.End()

	var err error
	defer instrument(ctx, "GetCode").err(&err).end()

	m.codeCacheLock.Lock()
	entry, ok := m.codeCache[code]
	if !ok {
		entry = &cacheEntry{}
		m.codeCache[code] = entry
	}
	m.codeCacheLock.Unlock()

	entry.Lock()
	defer entry.Unlock()

	if entry.desc != nil {
		return entry.desc, nil
	}

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	ctx, span = instracer.StartSpan(ctx, "artifactmanager.GetCode sendAndFollowRedirect")
	genericReact, err := sendAndFollowRedirect(
		ctx,
		m.db,
		m.bus(ctx),
		&message.GetCode{Code: code},
		*currentPulse,
	)
	span.End()

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
		entry.desc = &desc
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
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.GetObject starts ...")
	var (
		desc *ObjectDescriptor
		err  error
	)
	defer instrument(ctx, "GetObject").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	getObjectMsg := &message.GetObject{
		Head:     head,
		State:    state,
		Approved: approved,
	}
	rep, err := sendAndFollowRedirect(ctx, m.db, m.bus(ctx), getObjectMsg, *currentPulse)
	if err != nil {
		return nil, err
	}

	switch r := rep.(type) {
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
		return nil, fmt.Errorf("GetObject: unexpected reply: %#v", rep)
	}
}

// HasPendingRequests returns true if object has unclosed requests.
func (m *LedgerArtifactManager) HasPendingRequests(
	ctx context.Context,
	object core.RecordRef,
) (bool, error) {
	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return false, err
	}

	rep, err := sendAndRetryJet(
		ctx,
		m.db,
		m.bus(ctx),
		&message.GetPendingRequests{Object: object},
		*currentPulse,
		jetMissRetryCount,
		nil,
	)
	if err != nil {
		return false, err
	}

	switch rep := rep.(type) {
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
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.GetDelegate starts ...")
	var err error
	defer instrument(ctx, "GetDelegate").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	genericReact, err := sendAndFollowRedirect(
		ctx,
		m.db,
		m.bus(ctx),
		&message.GetDelegate{
			Head:   head,
			AsType: asType,
		},
		*currentPulse,
	)

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
	defer instrument(ctx, "GetChildren").err(&err).end()

	latestPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	iter, err := NewChildIterator(ctx, m.bus(ctx), m.db, parent, pulse, m.getChildrenChunkSize, *latestPulse)
	return iter, err
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	ctx context.Context, domain, request core.RecordRef, typeDec []byte,
) (*core.RecordID, error) {
	var err error
	defer instrument(ctx, "DeclareType").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
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
		*currentPulse,
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
	defer instrument(ctx, "DeployCode").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	codeRec := &record.CodeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domain,
			Request: request,
		},
		Code:        record.CalculateIDForBlob(m.PlatformCryptographyScheme, currentPulse.PulseNumber, code),
		MachineType: machineType,
	}
	codeID := record.NewRecordIDFromRecord(m.PlatformCryptographyScheme, currentPulse.PulseNumber, codeRec)
	codeRef := core.NewRecordRef(*domain.Record(), *codeID)

	_, err = m.setBlob(ctx, code, *codeRef, *currentPulse)
	if err != nil {
		return nil, err
	}
	id, err := m.setRecord(
		ctx,
		codeRec,
		*codeRef,
		*currentPulse,
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
	defer instrument(ctx, "ActivatePrototype").err(&err).end()
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
	defer instrument(ctx, "ActivateObject").err(&err).end()
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
	defer instrument(ctx, "DeactivateObject").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)

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
		*currentPulse,
	)
	if err != nil {
		return nil, err
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
	defer instrument(ctx, "UpdatePrototype").err(&err).end()
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
	defer instrument(ctx, "UpdateObject").err(&err).end()
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
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.RegisterValidation starts ...")
	var err error
	defer instrument(ctx, "RegisterValidation").err(&err).end()

	msg := message.ValidateRecord{
		Object:             object,
		State:              state,
		IsValid:            isValid,
		ValidationMessages: validationMessages,
	}

	currentPulse, err := m.PulseStorage.Current(ctx)
	if err != nil {
		return err
	}

	_, err = sendAndRetryJet(ctx, m.db, m.bus(ctx), &msg, *currentPulse, jetMissRetryCount, nil)
	return err
}

// RegisterResult saves VM method call result.
func (m *LedgerArtifactManager) RegisterResult(
	ctx context.Context, object, request core.RecordRef, payload []byte,
) (*core.RecordID, error) {
	var err error
	defer instrument(ctx, "RegisterResult").err(&err).end()

	currentPulse, err := m.PulseStorage.Current(ctx)
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
		*currentPulse,
	)
	return recid, err
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
	currentPulse, err := m.PulseStorage.Current(ctx)
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
				Memory:      record.CalculateIDForBlob(m.PlatformCryptographyScheme, currentPulse.PulseNumber, memory),
				Image:       prototype,
				IsPrototype: isPrototype,
			},
			Parent:     parent,
			IsDelegate: asDelegate,
		},
		object,
		memory,
		*currentPulse,
	)
	if err != nil {
		return nil, err
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
		*currentPulse,
	)
	if err != nil {
		return nil, err
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
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.updateObject starts ...")
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

	currentPulse, err := m.PulseStorage.Current(ctx)
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
		*currentPulse,
	)
	if err != nil {
		return nil, err
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
	currentPulse core.Pulse,
) (*core.RecordID, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.setRecord starts ...")

	genericReply, err := sendAndRetryJet(
		ctx,
		m.db,
		m.bus(ctx),
		&message.SetRecord{
			Record:    record.SerializeRecord(rec),
			TargetRef: target,
		},
		currentPulse,
		jetMissRetryCount,
		nil,
	)

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
	currentPulse core.Pulse,
) (*core.RecordID, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.setBlob starts ...")
	genericReact, err := sendAndRetryJet(
		ctx,
		m.db,
		m.bus(ctx),
		&message.SetBlob{
			Memory:    blob,
			TargetRef: target,
		},
		currentPulse,
		jetMissRetryCount,
		nil,
	)

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
	currentPulse core.Pulse,
) (*reply.Object, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.sendUpdateObject starts ...")
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

	genericRep, err := sendAndRetryJet(
		ctx,
		m.db,
		m.bus(ctx),
		&message.UpdateObject{
			Record: record.SerializeRecord(rec),
			Object: object,
			Memory: memory,
		},
		currentPulse,
		jetMissRetryCount,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	switch rep := genericRep.(type) {
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
	currentPulse core.Pulse,
) (*core.RecordID, error) {
	inslogger.FromContext(ctx).Debug("LedgerArtifactManager.registerChild starts ...")

	genericReact, err := sendAndRetryJet(
		ctx,
		m.db,
		m.bus(ctx),
		&message.RegisterChild{
			Record: record.SerializeRecord(rec),
			Parent: parent,
			Child:  child,
			AsType: asType,
		}, currentPulse, jetMissRetryCount, nil)

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

func (m *LedgerArtifactManager) bus(ctx context.Context) core.MessageBus {
	return core.MessageBusFromContext(ctx, m.DefaultBus)
}

func sendAndFollowRedirect(
	ctx context.Context,
	db *storage.DB,
	bus core.MessageBus,
	msg core.Message,
	pulse core.Pulse,
) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)
	inslog.Debug("LedgerArtifactManager.sendAndFollowRedirect starts ...")

	rep, err := bus.Send(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := rep.(*reply.JetMiss); ok {
		rep, err = sendAndRetryJet(ctx, db, bus, msg, pulse, jetMissRetryCount, nil)
	}

	if r, ok := rep.(core.RedirectReply); ok {
		stats.Record(ctx, statRedirects.M(1))

		redirected := r.Redirected(msg)
		inslog.Debugf("redirect reciever=%v", r.GetReceiver())

		rep, err = bus.Send(ctx, redirected, &core.MessageSendOptions{
			Token:    r.GetToken(),
			Receiver: r.GetReceiver(),
		})
		if err != nil {
			return nil, err
		}
		if _, ok := rep.(core.RedirectReply); ok {
			return nil, errors.New("double redirects are forbidden")
		}
		return rep, nil
	}

	return rep, err
}

func sendAndRetryJet(
	ctx context.Context,
	db *storage.DB,
	bus core.MessageBus,
	msg core.Message,
	pulse core.Pulse,
	retries int,
	ops *core.MessageSendOptions,
) (core.Reply, error) {
	if retries <= 0 {
		return nil, errors.New("failed to find jet (retry limit exceeded on client)")
	}
	rep, err := bus.Send(ctx, msg, ops)
	if err != nil {
		return nil, err
	}
	if r, ok := rep.(*reply.JetMiss); ok {
		err := db.UpdateJetTree(ctx, pulse.PulseNumber, true, r.JetID)
		if err != nil {
			return nil, err
		}
		return sendAndRetryJet(ctx, db, bus, msg, pulse, retries-1, ops)
	}

	return rep, nil
}
