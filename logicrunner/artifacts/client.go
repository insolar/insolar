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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type localStorage struct {
	initialized bool
	storage     map[insolar.Reference]interface{}
}

func (s *localStorage) Initialized() {
	s.initialized = true
}

func (s *localStorage) StoreObject(reference insolar.Reference, descriptor interface{}) {
	if s.initialized {
		panic("Trying to initialize singleFlightCache after initialization was finished")
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
	PCS            insolar.PlatformCryptographyScheme `inject:""`
	PulseAccessor  pulse.Accessor                     `inject:""`
	JetCoordinator jet.Coordinator                    `inject:""`

	sender       bus.Sender
	localStorage *localStorage
}

// State returns hash state for artifact manager.
func (m *client) State() []byte {
	// This is a temporary stab to simulate real hash.
	return m.PCS.IntegrityHasher().Hash([]byte{1, 2, 3})
}

// NewClient creates new client instance.
func NewClient(sender bus.Sender) *client { // nolint
	return &client{
		sender:       sender,
		localStorage: newLocalStorage(),
	}
}

// registerRequest registers incoming or outgoing request.
func (m *client) registerRequest(
	ctx context.Context, req record.Request, msgPayload payload.Payload, sender bus.Sender,
) (*payload.RequestInfo, error) {
	affinityRef, err := m.calculateAffinityReference(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "registerRequest: failed to calculate affinity reference")
	}

	ctx, span := instracer.StartSpan(ctx, "artifactmanager.registerRequest")
	instrumenter := instrument(ctx, "registerRequest").err(&err)
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	pl, err := m.sendToLight(ctx, sender, msgPayload, *affinityRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't register request")
	}

	switch p := pl.(type) {
	case *payload.RequestInfo:
		return p, nil
	case *payload.Error:
		if p.Code == payload.CodeFlowCanceled {
			err = flow.ErrCancelled
		} else {
			err = &payload.CodedError{Code: p.Code, Text: p.Text}
		}
		return nil, err
	default:
		err = fmt.Errorf("registerRequest: unexpected reply: %#v", p)
		return nil, err
	}
}

func (m *client) calculateAffinityReference(ctx context.Context, requestRecord record.Request) (*insolar.Reference, error) {
	pulseNumber, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}
	return record.CalculateRequestAffinityRef(requestRecord, pulseNumber, m.PCS), nil
	// recID := *insolar.NewID(currentPN, h.Sum(nil))
	// return insolar.NewReference(recID), nil
}

// RegisterIncomingRequest sends message for incoming request registration,
// returns request record Ref if request successfully created or already exists.
func (m *client) RegisterIncomingRequest(ctx context.Context, request *record.IncomingRequest) (*payload.RequestInfo, error) {
	incomingRequest := &payload.SetIncomingRequest{Request: record.Wrap(request)}

	// retriesNumber is zero, because we don't retry registering of incoming requests - the caller should
	// re-send the request instead.
	res, err := m.registerRequest(ctx, request, incomingRequest, m.sender)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterIncomingRequest")
	}
	return res, err
}

// RegisterOutgoingRequest sends message for outgoing request registration,
// returns request record Ref if request successfully created or already exists.
func (m *client) RegisterOutgoingRequest(ctx context.Context, request *record.OutgoingRequest) (*payload.RequestInfo, error) {
	outgoingRequest := &payload.SetOutgoingRequest{Request: record.Wrap(request)}
	res, err := m.registerRequest(
		ctx, request, outgoingRequest, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1),
	)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterOutgoingRequest")
	}
	return res, err
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
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	getCodePl := &payload.GetCode{CodeID: *code.Record()}

	pl, err := m.sendToLight(
		ctx, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1), getCodePl, code,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetCode")
	}

	switch p := pl.(type) {
	case *payload.Code:
		rec := record.Material{}
		err := rec.Unmarshal(p.Record)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal record")
		}
		virtual := record.Unwrap(&rec.Virtual)
		codeRecord, ok := virtual.(*record.Code)
		if !ok {
			return nil, fmt.Errorf("unexpected record %T", virtual)
		}
		desc = &codeDescriptor{
			ref:         code,
			machineType: codeRecord.MachineType,
			code:        codeRecord.Code,
		}
		return desc, nil
	case *payload.Error:
		err = errors.New(p.Text)
		return nil, err
	default:
		err = fmt.Errorf("GetObject: unexpected reply: %#v", p)
		return nil, err
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
			instracer.AddError(span, err)
		}
		span.End()
		if err != nil && err == ErrObjectDeactivated {
			err = nil // megahack: threat it 2xx
		}
		instrumenter.end()
	}()

	logger := inslogger.FromContext(ctx).WithField("object", head.Record().DebugString())

	msg, err := payload.NewMessage(&payload.GetObject{
		ObjectID: *head.Record(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}

	r := bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 2)
	reps, done := r.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, head)
	defer done()

	var (
		index        *record.Lifeline
		statePayload *payload.State
	)
	success := func() bool {
		return index != nil && statePayload != nil
	}

	for rep := range reps {
		replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal reply")
		}

		switch p := replyPayload.(type) {
		case *payload.Index:
			logger.Debug("reply index")
			index = &record.Lifeline{}
			err := index.Unmarshal(p.Index)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal index")
			}
		case *payload.State:
			logger.Debug("reply state")
			statePayload = p
		case *payload.Error:
			logger.Debug("reply error: ", p.Text)
			switch p.Code {
			case payload.CodeDeactivated:
				err = insolar.ErrDeactivated
				return nil, err
			default:
				logger.Errorf("reply error: %v, objectID: %v", p.Text, head.Record().DebugString())
				err = errors.New(p.Text)
				return nil, err
			}
		default:
			err = fmt.Errorf("GetObject: unexpected reply: %#v", p)
			return nil, err
		}

		if success() {
			break
		}
	}
	if !success() {
		logger.Error(ErrNoReply)
		err = ErrNoReply
		return nil, ErrNoReply
	}

	rec := record.Material{}
	err = rec.Unmarshal(statePayload.Record)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal state")
	}
	virtual := record.Unwrap(&rec.Virtual)
	s, ok := virtual.(record.State)
	if !ok {
		err = errors.New("wrong state record")
		return nil, err
	}
	state := s

	desc := &objectDescriptor{
		head:        head,
		state:       *index.LatestState,
		prototype:   state.GetImage(),
		isPrototype: state.GetIsPrototype(),
		memory:      statePayload.Memory,
		parent:      index.Parent,
	}
	return desc, nil
}

func (m *client) GetAbandonedRequest(
	ctx context.Context, object, reqRef insolar.Reference,
) (record.Request, error) {
	var err error
	instrumenter := instrument(ctx, "GetAbandonedRequest").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifacts.GetAbandonedRequest")
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	getRequestPl := &payload.GetRequest{
		ObjectID:  *object.Record(),
		RequestID: *reqRef.Record(),
	}

	pl, err := m.sendToLight(ctx, m.sender, getRequestPl, object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetRequest")
	}
	req, ok := pl.(*payload.Request)
	if !ok {
		err = fmt.Errorf("unexpected reply %T", pl)
		return nil, err
	}

	concrete := record.Unwrap(&req.Request)
	var result record.Request

	switch v := concrete.(type) {
	case *record.IncomingRequest:
		result = v
	case *record.OutgoingRequest:
		result = v
	default:
		err = fmt.Errorf("GetAbandonedRequest: unexpected message: %#v", concrete)
		return nil, err
	}

	return result, nil
}

// GetPendings returns a list of pending requests
func (m *client) GetPendings(ctx context.Context, object insolar.Reference) ([]insolar.Reference, error) {
	var err error
	instrumenter := instrument(ctx, "GetPendings").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.GetPendings")
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	getPendingsPl := &payload.GetPendings{
		ObjectID: *object.Record(),
	}

	pl, err := m.sendToLight(ctx, m.sender, getPendingsPl, object)
	if err != nil {
		return []insolar.Reference{}, errors.Wrap(err, "failed to send GetPendings")
	}

	switch concrete := pl.(type) {
	case *payload.IDs:
		res := make([]insolar.Reference, len(concrete.IDs))
		for i := range concrete.IDs {
			res[i] = *insolar.NewReference(concrete.IDs[i])
		}
		return res, nil
	case *payload.Error:
		if concrete.Code == payload.CodeNoPendings {
			err = insolar.ErrNoPendingRequest
			return []insolar.Reference{}, insolar.ErrNoPendingRequest
		}
		err = errors.New(concrete.Text)
		return []insolar.Reference{}, err
	default:
		return []insolar.Reference{}, fmt.Errorf("unexpected reply %T", pl)
	}
}

// HasPendings returns true if object has unclosed requests.
func (m *client) HasPendings(
	ctx context.Context,
	object insolar.Reference,
) (bool, error) {
	var err error
	instrumenter := instrument(ctx, "HasPendings").err(&err)
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.HasPendings")
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	hasPendingsPl := &payload.HasPendings{
		ObjectID: *object.Record(),
	}

	pl, err := m.sendToLight(ctx, m.sender, hasPendingsPl, object)
	if err != nil {
		err = errors.Wrap(err, "failed to send HasPendings")
		return false, err
	}

	switch concrete := pl.(type) {
	case *payload.PendingsInfo:
		return concrete.HasPendings, nil
	case *payload.Error:
		err = errors.New(concrete.Text)
		return false, err
	default:
		err = fmt.Errorf("HasPendings: unexpected reply %T", pl)
		return false, err
	}
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
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	codeRec := record.Code{
		Domain:      domain,
		Code:        code,
		MachineType: machineType,
	}
	virtual := record.Wrap(&codeRec)
	buf, err := virtual.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal record")
	}

	h := m.PCS.ReferenceHasher()
	_, err = h.Write(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate hash")
	}
	recID := *insolar.NewID(currentPN, h.Sum(nil))

	psc := &payload.SetCode{
		Record: buf,
	}

	pl, err := m.sendToLight(
		ctx, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1),
		psc, *insolar.NewReference(recID),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send SetCode")
	}

	switch p := pl.(type) {
	case *payload.ID:
		return &p.ID, nil
	case *payload.Error:
		err = errors.New(p.Text)
		return nil, err
	default:
		err = fmt.Errorf("DeployCode: unexpected reply: %#v", p)
		return nil, err
	}
}

// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *client) ActivatePrototype(
	ctx context.Context,
	object, parent, code insolar.Reference,
	memory []byte,
) error {
	var err error
	ctx, span := instracer.StartSpan(ctx, "artifactmanager.ActivatePrototype")
	instrumenter := instrument(ctx, "ActivatePrototype").err(&err)
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()
	err = m.activateObject(ctx, object, code, true, parent, memory)
	return err
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
	obj insolar.Reference,
	prototype insolar.Reference,
	isPrototype bool,
	parent insolar.Reference,
	memory []byte,
) error {
	_, err := m.GetObject(ctx, parent)
	if err != nil {
		return errors.Wrap(err, "wrong parent")
	}

	activate := record.Activate{
		Request:     obj,
		Memory:      memory,
		Image:       prototype,
		IsPrototype: isPrototype,
		Parent:      parent,
	}

	result := record.Result{
		Object:  *obj.Record(),
		Request: obj,
	}

	virtActivate := record.Wrap(&activate)
	virtResult := record.Wrap(&result)

	activateBuf, err := virtActivate.Marshal()
	if err != nil {
		return errors.Wrap(err, "ActivateObject: can't serialize record")
	}
	resultBuf, err := virtResult.Marshal()
	if err != nil {
		return errors.Wrap(err, "ActivateObject: can't serialize record")
	}

	pa := &payload.Activate{
		Record: activateBuf,
		Result: resultBuf,
	}

	pl, err := m.sendToLight(ctx, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1), pa, obj)
	if err != nil {
		return errors.Wrap(err, "can't send activation and result records")
	}

	switch p := pl.(type) {
	case *payload.ResultInfo:
		return nil
	case *payload.Error:
		return errors.New(p.Text)
	default:
		return fmt.Errorf("ActivateObject: unexpected reply: %#v", p)
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

// RegisterResult saves VM method call result with it's side effects
func (m *client) RegisterResult(
	ctx context.Context,
	request insolar.Reference,
	result RequestResult,
) error {

	sendResult := func(
		payloadInput payload.Payload,
		obj insolar.Reference,
	) (*insolar.ID, error) {

		payloadOutput, err := m.sendToLight(
			ctx, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1), payloadInput, obj,
		)
		if err != nil {
			return nil, err
		}

		switch p := payloadOutput.(type) {
		case *payload.ResultInfo:
			return &payloadOutput.(*payload.ResultInfo).ResultID, nil
		case *payload.Error:
			return nil, errors.New(p.Text)
		default:
			return nil, fmt.Errorf("RegisterResult: unexpected reply: %#v", p)
		}
	}

	var (
		pl  payload.Payload
		err error
	)

	ctx, span := instracer.StartSpan(ctx, "artifactmanager.RegisterResult")
	instrumenter := instrument(ctx, "RegisterResult").err(&err)
	defer func() {
		if err != nil {
			instracer.AddError(span, err)
		}
		span.End()
		instrumenter.end()
	}()
	span.AddAttributes(trace.StringAttribute("SideEffect", result.Type().String()))

	objReference := result.ObjectReference()
	resultRecord := record.Result{
		Object:  *objReference.Record(),
		Request: request,
		Payload: result.Result(),
	}

	switch result.Type() {
	// ActivateObject creates activate object record in storage. Provided prototype reference will be used as objects prototype
	// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	case RequestSideEffectActivate:
		parentRef, imageRef, memory := result.Activate()

		vResultRecord := record.Wrap(&resultRecord)
		vActivateRecord := record.Wrap(&record.Activate{
			Request:     request,
			Memory:      memory,
			Image:       imageRef,
			IsPrototype: false,
			Parent:      parentRef,
		})

		plTyped := payload.Activate{}
		plTyped.Record, err = vActivateRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Activate record")
		}
		plTyped.Result, err = vResultRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Result record")
		}
		pl = &plTyped

	// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	case RequestSideEffectAmend:
		objectStateID, objectImage, memory := result.Amend()

		vResultRecord := record.Wrap(&resultRecord)
		vAmendRecord := record.Wrap(&record.Amend{
			Request:     request,
			Memory:      memory,
			Image:       objectImage,
			IsPrototype: false,
			PrevState:   objectStateID,
		})

		plTyped := payload.Update{}
		plTyped.Record, err = vAmendRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Amend record")
		}
		plTyped.Result, err = vResultRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Result record")
		}
		pl = &plTyped

	// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	case RequestSideEffectDeactivate:
		objectStateID := result.Deactivate()

		vResultRecord := record.Wrap(&resultRecord)
		vDeactivateRecord := record.Wrap(&record.Deactivate{
			Request:   request,
			PrevState: objectStateID,
		})

		plTyped := payload.Deactivate{}
		plTyped.Record, err = vDeactivateRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Deactivate record")
		}
		plTyped.Result, err = vResultRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Result record")
		}
		pl = &plTyped

	case RequestSideEffectNone:
		vResultRecord := record.Wrap(&resultRecord)

		plTyped := payload.SetResult{}
		plTyped.Result, err = vResultRecord.Marshal()
		if err != nil {
			return errors.Wrap(err, "RegisterResult: can't serialize Result record")
		}
		pl = &plTyped

	default:
		err = errors.Errorf("RegisterResult: Unknown side effect %d", result.Type())
		return err
	}

	_, err = sendResult(pl, result.ObjectReference())
	if err != nil {
		return errors.Wrapf(err, "RegisterResult: Failed to send results: %s", result.Type().String())
	}

	return nil
}

func (m *client) sendToLight(
	ctx context.Context,
	sender bus.Sender,
	ppl payload.Payload,
	ref insolar.Reference,
) (payload.Payload, error) {

	msg, err := payload.NewMessage(ppl)
	if err != nil {
		return nil, err
	}

	reps, done := sender.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, ref)
	defer done()

	rep, ok := <-reps
	if !ok {
		return nil, ErrNoReply
	}

	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	return pl, nil
}
