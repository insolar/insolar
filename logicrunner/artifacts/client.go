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

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	insPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
)

const (
	getPendingLimit = 100
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

func (s *localStorage) Prototype(reference insolar.Reference) PrototypeDescriptor {
	objectDescI, ok := s.storage[reference]
	if !ok {
		return nil
	}
	desc, ok := objectDescI.(PrototypeDescriptor)
	if !ok {
		return nil
	}
	return desc
}

func newLocalStorage() *localStorage {
	return &localStorage{
		initialized: false,
		storage:     make(map[insolar.Reference]interface{}),
	}
}

// Client provides concrete API to storage for processing module.
type client struct {
	PCS           insolar.PlatformCryptographyScheme `inject:""`
	PulseAccessor insPulse.Accessor                  `inject:""`

	sender       bus.Sender
	localStorage *localStorage
}

// NewClient creates new client instance.
func NewClient(sender bus.Sender) Client {
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
}

// RegisterIncomingRequest sends message for incoming request registration,
// returns request record Ref if request successfully created or already exists.
func (m *client) RegisterIncomingRequest(ctx context.Context, request *record.IncomingRequest) (*payload.RequestInfo, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "RegisterIncomingRequest", &err)
	defer instrumenter.end()

	incomingRequest := &payload.SetIncomingRequest{Request: record.Wrap(request)}

	// retriesNumber is zero, because we don't retry registering of incoming requests - the caller should
	// re-send the request instead.
	res, err := m.registerRequest(ctx, request, incomingRequest, m.sender)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterIncomingRequest")
	}
	switch {
	case res.Result != nil:
		stats.Record(ctx, metrics.IncomingRequestsClosed.M(1))
	case res.Request != nil:
		stats.Record(ctx, metrics.IncomingRequestsDuplicate.M(1))
	default:
		stats.Record(ctx, metrics.IncomingRequestsNew.M(1))
	}
	return res, err
}

// RegisterOutgoingRequest sends message for outgoing request registration,
// returns request record Ref if request successfully created or already exists.
func (m *client) RegisterOutgoingRequest(ctx context.Context, request *record.OutgoingRequest) (*payload.RequestInfo, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "RegisterOutgoingRequest", &err)
	defer instrumenter.end()

	retrySender := bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1)
	outgoingRequest := &payload.SetOutgoingRequest{Request: record.Wrap(request)}

	res, err := m.registerRequest(ctx, request, outgoingRequest, retrySender)
	if err != nil {
		return nil, errors.Wrap(err, "RegisterOutgoingRequest")
	}

	switch {
	case res.Result != nil:
		stats.Record(ctx, metrics.OutgoingRequestsClosed.M(1))
	case res.Request != nil:
		stats.Record(ctx, metrics.OutgoingRequestsDuplicate.M(1))
	default:
		stats.Record(ctx, metrics.OutgoingRequestsNew.M(1))
	}

	return res, err
}

// GetPulse returns pulse data for pulse number from request.
func (m *client) GetPulse(ctx context.Context, pn insolar.PulseNumber) (insolar.Pulse, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "GetPulse", &err)
	defer instrumenter.end()

	getPulse := &payload.GetPulse{
		PulseNumber: pn,
	}

	pl, err := m.sendToLight(
		ctx, bus.NewRetrySender(
			m.sender,
			m.PulseAccessor, 1,
			1),
		getPulse,
		*insolar.NewReference(*insolar.NewID(pn, nil)),
	)

	if err != nil {
		return insolar.Pulse{}, errors.Wrap(err, "failed to send GetPulse")
	}

	switch p := pl.(type) {
	case *payload.Pulse:
		return *insPulse.FromProto(&p.Pulse), nil
	case *payload.Error:
		err = errors.New(p.Text)
		return insolar.Pulse{}, err
	default:
		err = fmt.Errorf("GetPulse: unexpected reply: %#v", p)
		return insolar.Pulse{}, err
	}
}

// GetCode returns code from code record by provided reference according to provided machine preference.
//
// This method is used by VM to fetch code for execution.
func (m *client) GetCode(
	ctx context.Context, code insolar.Reference,
) (CodeDescriptor, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "GetCode", &err)
	defer instrumenter.end()

	desc := m.localStorage.Code(code)
	if desc != nil {
		return desc, nil
	}

	getCodePl := &payload.GetCode{CodeID: *code.GetLocal()}

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
		err = fmt.Errorf("GetCode: unexpected reply: %#v", p)
		return nil, err
	}
}

// GetObject returns object descriptor with latest state.
func (m *client) GetObject(
	ctx context.Context,
	head insolar.Reference,
	request *insolar.Reference,
) (ObjectDescriptor, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "GetObject", &err)
	defer instrumenter.end()

	if desc := m.localStorage.Object(head); desc != nil {
		return desc, nil
	}
	getObjectRes, err := m.sendGetObject(ctx, head, request)
	if err != nil {
		return nil, err
	}

	if getObjectRes.state.GetIsPrototype() {
		return nil, errors.New("record is prototype, not an object")
	}

	desc := &objectDescriptor{
		head:      head,
		state:     *getObjectRes.index.LatestState,
		prototype: getObjectRes.state.GetImage(),
		memory:    getObjectRes.state.GetMemory(),
		parent:    getObjectRes.index.Parent,
		requestID: getObjectRes.lastRequestID,
	}
	return desc, nil
}

// GetPrototype returns prototype descriptor with latest state.
func (m *client) GetPrototype(
	ctx context.Context,
	head insolar.Reference,
) (PrototypeDescriptor, error) {

	if desc := m.localStorage.Prototype(head); desc != nil {
		return desc, nil
	}

	getObjectRes, err := m.sendGetObject(ctx, head, nil)
	if err != nil {
		return nil, err
	}

	if !getObjectRes.state.GetIsPrototype() {
		return nil, errors.New("record is not a prototype")
	}

	ref := getObjectRes.state.GetImage()
	if ref == nil {
		return nil, errors.New("prototype has no code reference")
	}

	desc := &prototypeDescriptor{
		head:  head,
		state: *getObjectRes.index.LatestState,
		code:  *ref,
	}
	return desc, nil
}

type getObjectRes struct {
	index         *record.Lifeline
	state         record.State
	lastRequestID *insolar.ID
}

func (m *client) sendGetObject(
	ctx context.Context,
	head insolar.Reference,
	request *insolar.Reference,
) (*getObjectRes, error) {
	var (
		err error
		res = &getObjectRes{}
	)
	logger := inslogger.FromContext(ctx).WithField("get_object", head.GetLocal().String())

	pl := payload.GetObject{
		ObjectID: *head.GetLocal(),
	}
	if request != nil {
		pl.RequestID = request.GetLocal()
	}

	msg, err := payload.NewMessage(&pl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}

	r := bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 2)
	reps, done := r.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, head)
	defer done()

	success := func() bool {
		return res.index != nil && res.state != nil
	}

	for rep := range reps {
		replyPayload, err := payload.UnmarshalFromMeta(rep.Payload)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal reply")
		}

		switch p := replyPayload.(type) {
		case *payload.Index:
			logger.Debug("reply index")
			res.index = &record.Lifeline{}
			err := res.index.Unmarshal(p.Index)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal index")
			}
			res.lastRequestID = p.EarliestRequestID
		case *payload.State:
			logger.Debug("reply state")
			rec := record.Material{}
			err = rec.Unmarshal(p.Record)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal state")
			}
			virtual := record.Unwrap(&rec.Virtual)
			s, ok := virtual.(record.State)
			if !ok {
				err = errors.New("wrong state record")
				return nil, err
			}
			res.state = s
		case *payload.Error:
			logger.Debug("reply error: ", p.Text)
			switch p.Code {
			case payload.CodeDeactivated:
				err = insolar.ErrDeactivated
				return nil, err
			default:
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
		err = ErrNoReply
		return nil, err
	}

	return res, nil
}

func (m *client) GetRequest(
	ctx context.Context, object, reqRef insolar.Reference,
) (record.Request, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "GetRequest", &err)
	defer instrumenter.end()

	getRequestPl := &payload.GetRequest{
		ObjectID:  *object.GetLocal(),
		RequestID: *reqRef.GetLocal(),
	}

	pl, err := m.sendToLight(ctx, m.sender, getRequestPl, object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetRequest")
	}

	switch p := pl.(type) {
	case *payload.Request:
		concrete := record.Unwrap(&p.Request)
		var result record.Request

		switch v := concrete.(type) {
		case *record.IncomingRequest:
			result = v
		case *record.OutgoingRequest:
			result = v
		default:
			err = fmt.Errorf("GetRequest: unexpected message: %#v", concrete)
			return nil, err
		}

		return result, nil
	case *payload.Error:
		if p.Code == payload.CodeNotFound {
			return nil, insolar.ErrNotFound
		}
		return nil, errors.New(p.Text)
	default:
		err = errors.Errorf("unexpected reply %T", pl)
		return nil, err
	}
}

// GetPendings returns a list of pending requests
func (m *client) GetPendings(
	ctx context.Context, object insolar.Reference, skip []insolar.ID,
) (
	[]insolar.Reference, error,
) {
	var err error
	ctx, instrumenter := instrument(ctx, "GetPendings", &err)
	defer instrumenter.end()

	getPendingsPl := &payload.GetPendings{
		ObjectID:        *object.GetLocal(),
		Count:           getPendingLimit,
		SkipRequestRefs: skip,
	}

	pl, err := m.sendToLight(ctx, m.sender, getPendingsPl, object)
	if err != nil {
		return []insolar.Reference{}, errors.Wrap(err, "failed to send GetPendings")
	}

	switch concrete := pl.(type) {
	case *payload.IDs:
		res := make([]insolar.Reference, len(concrete.IDs))
		for i := range concrete.IDs {
			res[i] = *insolar.NewRecordReference(concrete.IDs[i])
		}
		return res, nil
	case *payload.Error:
		if concrete.Code == payload.CodeNoPendings {
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
	ctx, instrumenter := instrument(ctx, "HasPendings", &err)
	defer instrumenter.end()

	hasPendingsPl := &payload.HasPendings{
		ObjectID: *object.GetLocal(),
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
func (m *client) DeployCode(ctx context.Context, code []byte, machineType insolar.MachineType) (*insolar.ID, error) {
	var err error
	ctx, instrumenter := instrument(ctx, "DeployCode", &err)
	defer instrumenter.end()

	codeRec := record.Code{
		Code:        code,
		MachineType: machineType,
	}
	virtual := record.Wrap(&codeRec)
	buf, err := virtual.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal record")
	}

	currentPN, err := m.pulse(ctx)
	if err != nil {
		return nil, err
	}

	h := m.PCS.ReferenceHasher()
	recID := *insolar.NewID(currentPN, h.Sum(buf))

	psc := &payload.SetCode{
		Record: buf,
	}

	pl, err := m.sendToLight(
		ctx, bus.NewRetrySender(m.sender, m.PulseAccessor, 1, 1),
		psc, *insolar.NewRecordReference(recID),
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
	ctx, instrumenter := instrument(ctx, "ActivatePrototype", &err)
	defer instrumenter.end()

	err = m.activateObject(ctx, object, code, true, parent, memory)
	return err
}

// pulse returns current PulseNumber for artifact manager
func (m *client) pulse(ctx context.Context) (insolar.PulseNumber, error) {
	pulseObject, err := m.PulseAccessor.Latest(ctx)
	if err != nil {
		return pulse.Unknown, err
	}

	return pulseObject.PulseNumber, nil
}

func (m *client) activateObject(
	ctx context.Context,
	obj insolar.Reference,
	prototype insolar.Reference,
	isPrototype bool,
	parent insolar.Reference,
	memory []byte,
) error {
	_, err := m.GetObject(ctx, parent, nil)
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
		Object:  *obj.GetLocal(),
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

func (m *client) InjectPrototypeDescriptor(reference insolar.Reference, descriptor PrototypeDescriptor) {
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
	var err error
	ctx, instrumenter := instrument(ctx, "RegisterResult", &err)
	defer instrumenter.end()

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
			return &p.ResultID, nil
		case *payload.ErrorResultExists:
			return nil, errors.New("another result already exists")
		case *payload.Error:
			return nil, errors.New(p.Text)
		default:
			return nil, fmt.Errorf("RegisterResult: unexpected reply: %#v", p)
		}
	}

	objReference := result.ObjectReference()
	resultRecord := record.Result{
		Object:  *objReference.GetLocal(),
		Request: request,
		Payload: result.Result(),
	}

	var pl payload.Payload
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
