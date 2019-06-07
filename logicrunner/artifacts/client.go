//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package artifacts

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
)

const (
	getChildrenChunkSize = 10 * 1000
)

type localStorage struct {
	initialized bool
	storage     map[insolar.Reference]interface{}
}

/*
	objectStorageLock sync.RWMutex
	objectStorage     map[insolar.Reference]ObjectDescriptor

	codeStorageLock sync.RWMutex
	codeStorage     map[insolar.Reference]CodeDescriptor
}*/

func (s *localStorage) Initialized() {
	s.initialized = true
}

func (s *localStorage) StoreObject(reference insolar.Reference, descriptor interface{}) {
	if s.initialized {
		panic("Trying to initialize cache after initialization was finished")
	}
	s.storage[reference] = descriptor
}

func (s *localStorage) Code(reference insolar.Reference) CodeDescriptor {
	codeDescI, ok := s.storage[reference]
	if !ok {
		return nil
	}
	codeDesc, ok := codeDescI.(CodeDescriptor)
	if !ok {
		return nil
	}
	return codeDesc
}

func (s *localStorage) Object(reference insolar.Reference) ObjectDescriptor {
	objectDescI, ok := s.storage[reference]
	if !ok {
		return nil
	}
	objectDesc, ok := objectDescI.(ObjectDescriptor)
	if !ok {
		return nil
	}
	return objectDesc
}

func newLocalStorage() *localStorage {
	return &localStorage{
		initialized: false,
		storage:     make(map[insolar.Reference]interface{}),
	}
}

// Client provides concrete API to storage for processing module.
type client struct {
	JetStorage     jet.Storage                        `inject:""`
	DefaultBus     insolar.MessageBus                 `inject:""`
	PCS            insolar.PlatformCryptographyScheme `inject:""`
	PulseAccessor  pulse.Accessor                     `inject:""`
	JetCoordinator jet.Coordinator                    `inject:""`

	sender               bus.Sender
	getChildrenChunkSize int
	senders              *messagebus.Senders
	localStorage         *localStorage

	codeCache map[insolar.ID]CodeDescriptor
}

// State returns hash state for artifact manager.
func (m *client) State() ([]byte, error) {
	// This is a temporary stab to simulate real hash.
	return m.PCS.IntegrityHasher().Hash([]byte{1, 2, 3}), nil
}

// NewClient creates new client instance.
func NewClient(sender bus.Sender) *client { // nolint
	return &client{
		getChildrenChunkSize: getChildrenChunkSize,
		senders:              messagebus.NewSenders(),
		sender:               sender,
		localStorage:         newLocalStorage(),
		codeCache:            map[insolar.ID]CodeDescriptor{},
	}
}

// RegisterRequest sends message for request registration,
// returns request record Ref if request successfully created or already exists.
func (m *client) RegisterRequest(
	ctx context.Context, request record.Request,
) (*insolar.ID, error) {
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

	virtRec := record.Wrap(request)

	var recRef *insolar.Reference
	switch request.CallType {
	case record.CTMethod:
		recRef = request.Object
	case record.CTSaveAsChild, record.CTSaveAsDelegate, record.CTGenesis:
		hash := record.HashVirtual(m.PCS.ReferenceHasher(), virtRec)
		recID := insolar.NewID(currentPN, hash)
		recRef = insolar.NewReference(*recID)
	default:
		return nil, errors.New("not supported call type " + request.CallType.String())
	}

	id, err := m.setRecord(
		ctx,
		virtRec,
		*recRef,
	)

	return id, errors.Wrap(err, "[ RegisterRequest ] ")
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *client) GetCode(
	ctx context.Context, code insolar.Reference,
) (CodeDescriptor, error) {
	var (
		desc CodeDescriptor
		err  error
	)

	desc = m.localStorage.Code(code)
	if desc != nil {
		return desc, nil
	}

	instrumenter := instrument(ctx, "GetCode").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetCode")
	defer func() {
		if err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
		}
		span.End()
		instrumenter.end()
	}()

	if cd, ok := m.codeCache[*code.Record()]; ok {
		return cd, nil
	}

	msg, err := payload.NewMessage(&payload.GetCode{
		CodeID: *code.Record(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}

	reps, done := m.sender.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, code)
	defer done()

	rep, ok := <-reps
	if !ok {
		return nil, errors.New("no reply")
	}
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch p := pl.(type) {
	case *payload.Code:
		rec := record.Material{}
		err := rec.Unmarshal(p.Record)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal record")
		}
		virtual := record.Unwrap(rec.Virtual)
		codeRecord, ok := virtual.(*record.Code)
		if !ok {
			return nil, errors.Wrapf(err, "unexpected record %T", virtual)
		}
		desc = &codeDescriptor{
			ref:         code,
			machineType: codeRecord.MachineType,
			code:        p.Code,
		}
		m.codeCache[*code.Record()] = desc
		return desc, nil
	case *payload.Error:
		return nil, errors.New(p.Text)
	default:
		return nil, fmt.Errorf("GetObject: unexpected reply: %#v", p)
	}
}

// GetObject returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
// provide methods for fetching all related data.
func (m *client) GetObject(
	ctx context.Context,
	head insolar.Reference,
) (ObjectDescriptor, error) {
	var (
		err error
	)

	if desc := m.localStorage.Object(head); desc != nil {
		return desc, nil
	}

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

	msg, err := payload.NewMessage(&payload.GetObject{
		ObjectID: *head.Record(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}
	reps, done := m.sender.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, head)
	defer done()

	var (
		index        *object.Lifeline
		statePayload *payload.State
	)
	success := func() bool {
		return index != nil && statePayload != nil
	}

	logger := inslogger.FromContext(ctx)
	for rep := range reps {
		replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal reply")
		}

		switch p := replyPayload.(type) {
		case *payload.Index:
			logger.Info("rep index")
			index = &object.Lifeline{}
			err := index.Unmarshal(p.Index)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal index")
			}
		case *payload.State:
			logger.Info("rep state")
			statePayload = p
		case *payload.Error:
			logger.Info("rep error: ", p.Text)
			switch p.Code {
			case payload.CodeDeactivated:
				return nil, insolar.ErrDeactivated
			default:
				return nil, errors.New(p.Text)
			}
		default:
			return nil, fmt.Errorf("GetObject: unexpected reply: %#v", p)
		}

		if success() {
			break
		}
	}
	if !success() {
		logger.WithField("correlation_id", middleware.MessageCorrelationID(msg)).Error("no reply")
		return nil, errors.New("no reply")
	}

	rec := record.Material{}
	err = rec.Unmarshal(statePayload.Record)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal state")
	}
	virtual := record.Unwrap(rec.Virtual)
	s, ok := virtual.(record.State)
	if !ok {
		return nil, errors.New("wrong state record")
	}
	state := s

	desc := &objectDescriptor{
		head:         head,
		state:        *index.LatestState,
		prototype:    state.GetImage(),
		isPrototype:  state.GetIsPrototype(),
		childPointer: index.ChildPointer,
		memory:       statePayload.Memory,
		parent:       index.Parent,
	}
	return desc, err
}

// GetPendingRequest returns an unclosed pending request
// It takes an id from current LME
// Then goes either to a light node or heavy node
func (m *client) GetPendingRequest(ctx context.Context, objectID insolar.ID) (*insolar.Reference, insolar.Parcel, error) {
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

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.RetryJetSender(m.JetStorage),
	)

	genericReply, err := sender(ctx, &message.GetPendingRequestID{
		ObjectID: objectID,
	}, nil)
	if err != nil {
		return nil, nil, err
	}

	var requestID insolar.ID
	switch r := genericReply.(type) {
	case *reply.ID:
		requestID = r.ID
	case *reply.Error:
		return nil, nil, r.Error()
	default:
		return nil, nil, fmt.Errorf("GetPendingRequest: unexpected reply: %#v", genericReply)
	}

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, nil, err
	}

	node, err := m.JetCoordinator.NodeForObject(ctx, objectID, currentPN, requestID.Pulse())
	if err != nil {
		return nil, nil, err
	}

	sender = messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryJetSender(m.JetStorage),
	)
	genericReply, err = sender(
		ctx,
		&message.GetRequest{
			Request: requestID,
		}, &insolar.MessageSendOptions{
			Receiver: node,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	switch r := genericReply.(type) {
	case *reply.Request:
		rec := record.Virtual{}
		err = rec.Unmarshal(r.Record)
		if err != nil {
			return nil, nil, errors.Wrap(err, "GetPendingRequest: can't deserialize record")
		}
		concrete := record.Unwrap(&rec)
		castedRecord, ok := concrete.(*record.Request)
		if !ok {
			return nil, nil, fmt.Errorf("GetPendingRequest: unexpected message: %#v", r)
		}

		serviceData := message.ServiceData{
			LogTraceID:    inslogger.TraceID(ctx),
			LogLevel:      inslogger.GetLoggerLevel(ctx),
			TraceSpanData: instracer.MustSerialize(ctx),
		}
		return insolar.NewReference(requestID), &message.Parcel{Msg: &message.CallMethod{
			Request: *castedRecord}, ServiceData: serviceData}, nil
	case *reply.Error:
		return nil, nil, r.Error()
	default:
		return nil, nil, fmt.Errorf("GetPendingRequest: unexpected reply: %#v", genericReply)
	}
}

// HasPendingRequests returns true if object has unclosed requests.
func (m *client) HasPendingRequests(
	ctx context.Context,
	object insolar.Reference,
) (bool, error) {
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.HasPendingRequests")
	defer span.End()

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.RetryJetSender(m.JetStorage),
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
func (m *client) GetDelegate(
	ctx context.Context, head, asType insolar.Reference,
) (*insolar.Reference, error) {
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

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.FollowRedirectSender(m.DefaultBus),
		messagebus.RetryJetSender(m.JetStorage),
	)
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
func (m *client) GetChildren(
	ctx context.Context, parent insolar.Reference, pulse *insolar.PulseNumber,
) (RefIterator, error) {
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

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.FollowRedirectSender(m.DefaultBus),
		messagebus.RetryJetSender(m.JetStorage),
	)
	iter, err := NewChildIterator(ctx, sender, parent, pulse, m.getChildrenChunkSize)
	return iter, err
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *client) DeclareType(
	ctx context.Context, domain, request insolar.Reference, typeDec []byte,
) (*insolar.ID, error) {
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

	typeRec := record.Type{
		Domain:          domain,
		Request:         request,
		TypeDeclaration: typeDec,
	}
	virtRec := record.Wrap(typeRec)

	recid, err := m.setRecord(
		ctx,
		virtRec,
		request,
	)

	return recid, err
}

// DeployCode creates new code record in storage.
//
// CodeRef records are used to activate prototype or as migration code for an object.
func (m *client) DeployCode(
	ctx context.Context,
	domain insolar.Reference,
	request insolar.Reference,
	code []byte,
	machineType insolar.MachineType,
) (*insolar.ID, error) {
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

	h := m.PCS.ReferenceHasher()
	_, err = h.Write(code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate hash")
	}
	blobID := *insolar.NewID(currentPN, h.Sum(nil))

	codeRec := record.Code{
		Domain:      domain,
		Request:     request,
		Code:        blobID,
		MachineType: machineType,
	}
	buf, err := codeRec.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal record")
	}

	h = m.PCS.ReferenceHasher()
	_, err = h.Write(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate hash")
	}
	recID := *insolar.NewID(currentPN, h.Sum(nil))

	msg, err := payload.NewMessage(&payload.SetCode{
		Record: buf,
		Code:   code,
	})

	reps, done := m.sender.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, *insolar.NewReference(recID))
	defer done()

	rep, ok := <-reps
	if !ok {
		return nil, errors.New("no reply")
	}
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch p := pl.(type) {
	case *payload.ID:
		return &p.ID, nil
	case *payload.Error:
		return nil, errors.New(p.Text)
	default:
		return nil, fmt.Errorf("GetObject: unexpected reply: %#v", p)
	}
}

// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *client) ActivatePrototype(
	ctx context.Context,
	domain, object, parent, code insolar.Reference,
	memory []byte,
) (ObjectDescriptor, error) {
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
func (m *client) ActivateObject(
	ctx context.Context,
	domain, object, parent, prototype insolar.Reference,
	asDelegate bool,
	memory []byte,
) (ObjectDescriptor, error) {
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
func (m *client) DeactivateObject(
	ctx context.Context, domain, request insolar.Reference, obj ObjectDescriptor,
) (*insolar.ID, error) {
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

	deactivate := record.Deactivate{
		Domain:    domain,
		Request:   request,
		PrevState: *obj.StateID(),
	}
	virtRec := record.Wrap(deactivate)

	desc, err := m.sendUpdateObject(
		ctx,
		virtRec,
		*obj.HeadRef(),
		nil,
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
func (m *client) UpdatePrototype(
	ctx context.Context,
	domain, request insolar.Reference,
	object ObjectDescriptor,
	memory []byte,
	code *insolar.Reference,
) (ObjectDescriptor, error) {
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
func (m *client) UpdateObject(
	ctx context.Context,
	domain, request insolar.Reference,
	object ObjectDescriptor,
	memory []byte,
) (ObjectDescriptor, error) {
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
func (m *client) RegisterValidation(
	ctx context.Context,
	object insolar.Reference,
	state insolar.ID,
	isValid bool,
	validationMessages []insolar.Message,
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

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryJetSender(m.JetStorage),
	)
	_, err = sender(ctx, &msg, nil)

	return err
}

// RegisterResult saves VM method call result.
func (m *client) RegisterResult(
	ctx context.Context, obj, request insolar.Reference, payload []byte,
) (*insolar.ID, error) {
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

	res := record.Result{
		Object:  *obj.Record(),
		Request: request,
		Payload: payload,
	}
	virtRec := record.Wrap(res)

	recid, err := m.setRecord(
		ctx,
		virtRec,
		obj,
	)
	return recid, err
}

// pulse returns current PulseNumber for artifact manager
func (m *client) pulse(ctx context.Context) (pn insolar.PulseNumber, err error) {
	pulse, err := m.PulseAccessor.Latest(ctx)
	if err != nil {
		return
	}

	pn = pulse.PulseNumber
	return
}

func (m *client) activateObject(
	ctx context.Context,
	domain insolar.Reference,
	obj insolar.Reference,
	prototype insolar.Reference,
	isPrototype bool,
	parent insolar.Reference,
	asDelegate bool,
	memory []byte,
) (ObjectDescriptor, error) {
	parentDesc, err := m.GetObject(ctx, parent)
	if err != nil {
		return nil, err
	}
	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	activate := record.Activate{
		Domain:      domain,
		Request:     obj,
		Memory:      *object.CalculateIDForBlob(m.PCS, currentPN, memory),
		Image:       prototype,
		IsPrototype: isPrototype,
		Parent:      parent,
		IsDelegate:  asDelegate,
	}
	virtRec := record.Wrap(activate)

	o, err := m.sendUpdateObject(
		ctx,
		virtRec,
		obj,
		memory,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate")
	}

	var (
		asType *insolar.Reference
	)
	child := record.Child{Ref: obj}
	if parentDesc.ChildPointer() != nil {
		child.PrevChild = *parentDesc.ChildPointer()
	}
	if asDelegate {
		asType = &prototype
	}
	virtChild := record.Wrap(child)

	_, err = m.registerChild(
		ctx,
		virtChild,
		parent,
		obj,
		asType,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register as child while activating")
	}

	return &objectDescriptor{
		head:         o.Head,
		state:        o.State,
		prototype:    o.Prototype,
		childPointer: o.ChildPointer,
		memory:       memory,
		parent:       o.Parent,
	}, nil
}

func (m *client) updateObject(
	ctx context.Context,
	domain, request insolar.Reference,
	obj ObjectDescriptor,
	code *insolar.Reference,
	memory []byte,
) (ObjectDescriptor, error) {
	var (
		image *insolar.Reference
		err   error
	)
	if obj.IsPrototype() {
		if code != nil {
			image = code
		} else {
			image, err = obj.Code()
		}
	} else {
		image, err = obj.Prototype()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	amend := record.Amend{
		Domain:      domain,
		Request:     request,
		Image:       *image,
		IsPrototype: obj.IsPrototype(),
		PrevState:   *obj.StateID(),
	}
	virtRec := record.Wrap(amend)

	o, err := m.sendUpdateObject(
		ctx,
		virtRec,
		*obj.HeadRef(),
		memory,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	return &objectDescriptor{
		head:         o.Head,
		state:        o.State,
		prototype:    o.Prototype,
		childPointer: o.ChildPointer,
		memory:       memory,
		parent:       o.Parent,
	}, nil
}

func (m *client) setRecord(
	ctx context.Context,
	rec record.Virtual,
	target insolar.Reference,
) (*insolar.ID, error) {
	data, err := rec.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "setRecord: can't serialize record")
	}
	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.RetryJetSender(m.JetStorage),
		messagebus.RetryFlowCancelled(m.PulseAccessor),
	)
	genericReply, err := sender(ctx, &message.SetRecord{
		Record:    data,
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

func (m *client) setBlob(
	ctx context.Context,
	blob []byte,
	target insolar.Reference,
) (*insolar.ID, error) {

	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryJetSender(m.JetStorage),
		messagebus.RetryFlowCancelled(m.PulseAccessor),
	)
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

func (m *client) sendUpdateObject(
	ctx context.Context,
	rec record.Virtual,
	obj insolar.Reference,
	memory []byte,
) (*reply.Object, error) {
	data, err := rec.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "setRecord: can't serialize record")
	}
	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.RetryJetSender(m.JetStorage),
		messagebus.RetryFlowCancelled(m.PulseAccessor),
	)
	genericReply, err := sender(
		ctx,
		&message.UpdateObject{
			Record: data,
			Object: obj,
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

func (m *client) registerChild(
	ctx context.Context,
	rec record.Virtual,
	parent insolar.Reference,
	child insolar.Reference,
	asType *insolar.Reference,
) (*insolar.ID, error) {
	data, err := rec.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "setRecord: can't serialize record")
	}
	sender := messagebus.BuildSender(
		m.DefaultBus.Send,
		messagebus.RetryIncorrectPulse(m.PulseAccessor),
		messagebus.RetryJetSender(m.JetStorage),
		messagebus.RetryFlowCancelled(m.PulseAccessor),
	)
	genericReact, err := sender(ctx, &message.RegisterChild{
		Record: data,
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

func (m *client) InjectCodeDescriptor(reference insolar.Reference, descriptor CodeDescriptor) {
	m.localStorage.StoreObject(reference, descriptor)
}

func (m *client) InjectObjectDescriptor(reference insolar.Reference, descriptor ObjectDescriptor) {
	m.localStorage.StoreObject(reference, descriptor)
}

func (m *client) InjectFinish() {
	m.localStorage.Initialized()
}
