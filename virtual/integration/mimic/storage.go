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

package mimic

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Storage interface {
	store.DB

	// Object
	GetObject(ctx context.Context, objectID insolar.ID) (record.State, *record.Index, *insolar.ID, error)
	SetObject(ctx context.Context, requestID insolar.ID, newState record.State, object insolar.ID) error
	// Request
	GetRequest(ctx context.Context, requestID insolar.ID) (record.Request, error)
	SetRequest(ctx context.Context, request record.Request) (*insolar.ID, []byte, []byte, error)
	// Result
	GetResult(ctx context.Context, requestID insolar.ID) (*insolar.ID, []byte, error)
	SetResult(ctx context.Context, result *record.Result) (*insolar.ID, error)
	RollbackSetResult(ctx context.Context, result *record.Result)
	// Code
	SetCode(ctx context.Context, code record.Code) (insolar.ID, error)
	GetCode(ctx context.Context, codeID insolar.ID) ([]byte, error)
	// Pendings
	GetPendings(ctx context.Context, objectID insolar.ID, limit uint32, skipRequestRefs []insolar.ID) ([]insolar.ID, error)
	HasPendings(ctx context.Context, objectID insolar.ID) (bool, error)
	GetOutgoingSagas(ctx context.Context, requestID insolar.ID) ([]*outgoingInfo, error)

	ObjectRequestsAreClosed(ctx context.Context, objectID insolar.ID) (bool, error)

	CalculateRequestAffinityRef(request record.Request, pulseNumber insolar.PulseNumber) insolar.ID
}

type mimicStorage struct {
	// components
	pcs insolar.PlatformCryptographyScheme
	pa  pulse.Accessor

	Code        map[insolar.ID]*CodeEntity
	Objects     map[insolar.ID]*ObjectEntity
	Requests    map[insolar.ID]*RequestEntity
	Results     map[insolar.ID]*RequestEntity
	MiscStorage map[string][]byte
}

func (s *mimicStorage) GetOutgoingSagas(ctx context.Context, requestID insolar.ID) ([]*outgoingInfo, error) {
	return s.Requests[requestID].getSagaOutgoingRequestIDs(), nil
}

var (
	// object related errors

	ErrNotFound           = errors.New("object not found")
	ErrAlreadyActivated   = errors.New("object already activated")
	ErrAlreadyDeactivated = errors.New("object already deactivated")
	ErrDeactivated        = errors.New("object is activated")
	ErrNotActivated       = errors.New("object isn't activated")

	// request/result related errors

	ErrRequestParentNotFound     = errors.New("parent request not found")
	ErrRequestNotFound           = errors.New("request not found")
	ErrRequestExists             = errors.New("request already exists")
	ErrRequestHasOpenedOutgoings = errors.New("request is reason for non closed outgoing request")
	ErrResultExists              = errors.New("request result exists")
	ErrResultNotFound            = errors.New("request result not found")
	ErrNoPendings                = errors.New("no pending requests are available")

	// code related errors

	ErrCodeExists   = errors.New("code already exists")
	ErrCodeNotFound = errors.New("code not found")
)

func (s *mimicStorage) calculateVirtualID(pulse insolar.PulseNumber, virtual record.Virtual) insolar.ID {
	hash := record.HashVirtual(s.pcs.ReferenceHasher(), virtual)
	return *insolar.NewID(pulse, hash)
}

func (s *mimicStorage) calculateRecordID(pulse insolar.PulseNumber, rec record.Record) insolar.ID {
	virtual := record.Wrap(rec)
	return s.calculateVirtualID(pulse, virtual)
}

func (s *mimicStorage) calculateSideEffectID(state *ObjectState) insolar.ID {
	request, ok := s.Requests[state.RequestID]
	if !ok {
		panic("failed to find request")
	}

	resultID := request.ResultID
	if resultID.IsEmpty() {
		panic("result is empty (should be closed)")
	}

	return s.calculateRecordID(resultID.Pulse(), state.State)
}

func (s *mimicStorage) CalculateRequestAffinityRef(request record.Request, pulseNumber insolar.PulseNumber) insolar.ID {
	return *record.CalculateRequestAffinityRef(request, pulseNumber, s.pcs).GetLocal()
}

func (s *mimicStorage) GetObject(_ context.Context, objectID insolar.ID) (record.State, *record.Index, *insolar.ID, error) {
	object, ok := s.Objects[objectID]
	if !ok {
		return nil, nil, nil, ErrNotFound
	}

	if len(object.ObjectChanges) == 0 {
		return nil, nil, nil, ErrNotActivated
	}

	if object.isDeactivated() {
		return nil, nil, nil, ErrDeactivated
	}

	latestObjectState := object.ObjectChanges[len(object.ObjectChanges)-1]
	latestObjectStateID := s.calculateSideEffectID(&latestObjectState)

	openRequestsCount, latestRequest, earliestRequest := object.getRequestsInfo()

	index := &record.Index{
		ObjID: objectID,
		Lifeline: record.Lifeline{
			LatestState:         &latestObjectStateID,
			StateID:             object.getLatestStateID(),
			Parent:              insolar.Reference{},
			LatestRequest:       nil,
			EarliestOpenRequest: nil,
			OpenRequestsCount:   openRequestsCount,
		},
		LifelineLastUsed: 0,
		PendingRecords:   nil,
	}

	if latestRequest != nil {
		index.Lifeline.LatestRequest = &latestRequest.ID
	}
	if earliestRequest != nil {
		pulseNumber := earliestRequest.getPulse()
		index.Lifeline.EarliestOpenRequest = &pulseNumber
	}

	return latestObjectState.State, index, nil, nil
}

func (s *mimicStorage) SetObject(ctx context.Context, requestID insolar.ID, newState record.State, objectID insolar.ID) error {
	switch newState.(type) {
	case *record.Activate:
		// TODO[bigbes]: take Objects lock
		if s.Objects[requestID] != nil && len(s.Objects[requestID].ObjectChanges) > 0 {
			return ErrAlreadyActivated
		}

		if _, ok := s.Objects[requestID]; !ok {
			inslogger.FromContext(ctx).Error("object is empty")
			s.Objects[requestID] = &ObjectEntity{
				ObjectChanges: nil,
				RequestsMap:   make(map[insolar.ID]*RequestEntity),
				RequestsList:  nil,
			}
		}
		objectEntity := s.Objects[requestID]
		objectEntity.ObjectChanges = append(objectEntity.ObjectChanges, ObjectState{
			State:     newState,
			RequestID: requestID,
		})

	case *record.Amend:
		if s.Objects[objectID] == nil {
			return ErrNotFound
		}

		state := s.Objects[objectID]
		if state.isDeactivated() {
			return ErrAlreadyDeactivated
		}
		state.ObjectChanges = append(state.ObjectChanges,
			ObjectState{
				State:     newState,
				RequestID: requestID,
			},
		)
	case *record.Deactivate:
		if s.Objects[objectID] == nil {
			return ErrNotFound
		}

		state := s.Objects[objectID]
		if state.isDeactivated() {
			return ErrAlreadyDeactivated
		}
		state.ObjectChanges = append(state.ObjectChanges,
			ObjectState{
				State:     newState,
				RequestID: requestID,
			},
		)
	default:
		panic(fmt.Sprintf("unexpected type %T", newState))
	}

	return nil
}

func contains(a []insolar.ID, x insolar.ID) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (s *mimicStorage) GetPendings(ctx context.Context, objectID insolar.ID, limit uint32, skipRequestRefs []insolar.ID) ([]insolar.ID, error) {
	state := s.Objects[objectID]
	if state == nil {
		return nil, ErrNotFound
	}

	// TODO[bigbes]: this logic should be rethought, probably we should add request to this not created object, anyway
	// not yet activated object, return Request
	if len(state.ObjectChanges) == 0 {
		return []insolar.ID{objectID}, nil
	}

	rv := make([]insolar.ID, 0)
	for _, request := range state.RequestsList {
		switch {
		case limit == 0:
			// we're done with our limit
			break
		case request.Status == RequestFinished, contains(skipRequestRefs, request.ID):
			// we're ignoring finished requests or one that we were asked to ignore
			continue
		}

		rv = append(rv, request.ID)
		limit--
	}

	return rv, nil
}

func (s *mimicStorage) HasPendings(ctx context.Context, objectID insolar.ID) (bool, error) {
	pendings, err := s.GetPendings(ctx, objectID, 1, nil)
	if err != nil {
		return false, err
	}
	return len(pendings) != 0, nil
}

func (s *mimicStorage) GetRequest(_ context.Context, requestID insolar.ID) (record.Request, error) {
	requestEntity, ok := s.Requests[requestID]
	if !ok {
		return nil, ErrRequestNotFound
	}
	return requestEntity.Request, nil
}

func (s *mimicStorage) SetRequest(_ context.Context, request record.Request) (*insolar.ID, []byte, []byte, error) {
	isOutgoingRequest := false
	if _, ok := request.(*record.OutgoingRequest); ok {
		isOutgoingRequest = true
	} else if _, ok := request.(*record.IncomingRequest); !ok {
		panic(fmt.Sprintf("Unknown request %T", request))
	}

	latest, err := s.pa.Latest(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed to obtained latest pulse"))
	}
	requestID := s.calculateRecordID(latest.PulseNumber, request)

	var objectID *insolar.ID
	if objectRef := request.AffinityRef(); objectRef != nil {
		objectID = objectRef.GetLocal()
	}

	// TODO[bigbes]: find duplicates
	if request, ok := s.Requests[requestID]; ok {
		material := record.Material{
			Virtual: record.Wrap(request.Request),
			ID:      requestID,
			JetID:   gen.JetID(),
		}
		if objectID != nil {
			material.ObjectID = *objectID
		}
		reqBuf, err := material.Marshal()
		if err != nil {
			panic(errors.Wrap(err, "failed to marshal Material Result record").Error())
		}
		return &requestID, reqBuf, request.Result, ErrRequestExists
	}

	parentRequestID := *request.ReasonRef().GetLocal()
	if request.GetCallType() != record.CTGenesis && !request.IsAPIRequest() {
		if parentRequestID.IsEmpty() {
			return nil, nil, nil, errors.New("bad request: reason is empty")
		}
		if s.Requests[parentRequestID] == nil {
			return nil, nil, nil, ErrRequestParentNotFound
		}
	}

	var re *RequestEntity
	if isOutgoingRequest {
		re = NewOutgoingRequestEntity(requestID, request)
	} else {
		re = NewIncomingRequestEntity(requestID, request)
	}

	if objectID != nil {
		state, ok := s.Objects[*objectID]
		switch {
		case !ok:
			return nil, nil, nil, ErrNotFound
		case len(state.ObjectChanges) == 0:
			return nil, nil, nil, ErrNotActivated
		case state.isDeactivated():
			return nil, nil, nil, ErrAlreadyDeactivated
		}

		if !isOutgoingRequest {
			state.addIncomingRequest(re)
		}
	}

	if isOutgoingRequest {
		s.Requests[parentRequestID].appendOutgoing(re)
	}

	s.Requests[requestID] = re

	if !isOutgoingRequest && (request.GetCallType() == record.CTSaveAsChild || request.GetCallType() == record.CTGenesis) {
		s.Objects[requestID] = &ObjectEntity{
			ObjectChanges: nil,
			RequestsMap:   make(map[insolar.ID]*RequestEntity),
			RequestsList:  nil,
		}
	}

	return &requestID, nil, nil, nil
}

func (s *mimicStorage) GetResult(_ context.Context, requestID insolar.ID) (*insolar.ID, []byte, error) {
	request, ok := s.Requests[requestID]
	if !ok {
		return nil, nil, ErrRequestNotFound
	}

	if request.Status != RequestFinished {
		return nil, nil, ErrResultNotFound
	}

	materialRec := record.Material{}
	if err := materialRec.Unmarshal(request.Result); !ok {
		panic(errors.Wrap(err, "failed to unmarshal Material Result record").Error())
	}

	return &materialRec.ID, request.Result, nil
}

// TODO[bigbes]: check for non-closed outgoings
func (s *mimicStorage) SetResult(_ context.Context, result *record.Result) (*insolar.ID, error) {
	if !result.Object.IsEmpty() {
		state, ok := s.Objects[result.Object]
		switch {
		case !ok:
			return nil, ErrNotFound
		// TODO[bigbes]: this check needs to be returned (not in case when we're activating object)
		// TODO[bigbes]: this logic should be rethought
		// case len(state.ObjectChanges) == 0:
		// 	return nil, ErrNotActivated
		case state.isDeactivated():
			return nil, ErrAlreadyDeactivated
		}
	}

	request, ok := s.Requests[*result.Request.GetLocal()]
	if !ok {
		return nil, ErrRequestNotFound
	}

	latest, err := s.pa.Latest(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed to obtained latest pulse"))
	}
	virtual := record.Wrap(result)
	resultID := s.calculateVirtualID(latest.PulseNumber, virtual)

	if request.Status == RequestFinished {
		return nil, ErrResultExists
	}

	if request.hasOpenedOutgoings() {
		return nil, ErrRequestHasOpenedOutgoings
	}

	materialRec := record.Material{
		Virtual:  virtual,
		ID:       resultID,
		ObjectID: result.Object,
		JetID:    gen.JetID(),
	}

	resultBytes, err := materialRec.Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal Material Result record").Error())
	}

	request.Status = RequestFinished
	request.Result = resultBytes
	request.ResultID = resultID
	s.Results[resultID] = request

	return &resultID, nil
}

func (s *mimicStorage) RollbackSetResult(_ context.Context, result *record.Result) {
	request := s.Requests[*result.Request.GetLocal()]

	request.Status = RequestRegistered
	request.Result = nil
	delete(s.Results, request.ResultID)

	request.ResultID = insolar.ID{}
}

func (s *mimicStorage) GetCode(_ context.Context, codeID insolar.ID) ([]byte, error) {
	code, ok := s.Code[codeID]
	if !ok {
		return nil, ErrCodeNotFound
	}

	return code.Code, nil
}

func (s *mimicStorage) SetCode(_ context.Context, code record.Code) (insolar.ID, error) {
	virtual := record.Wrap(&code)

	latest, err := s.pa.Latest(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed to obtained latest pulse"))
	}
	codeID := s.calculateVirtualID(latest.PulseNumber, virtual)

	if _, ok := s.Code[codeID]; ok {
		return insolar.ID{}, ErrCodeExists
	}

	material := record.Material{
		Virtual: virtual,
		ID:      codeID,
		JetID:   gen.JetID(),
	}
	resultBytes, err := material.Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal record"))
	}

	s.Code[codeID] = &CodeEntity{
		Code: resultBytes,
	}

	return codeID, nil
}

func (s *mimicStorage) ObjectRequestsAreClosed(ctx context.Context, objectID insolar.ID) (bool, error) {
	pendings, err := s.GetPendings(ctx, objectID, 1, []insolar.ID{})
	if err != nil {
		return false, err
	}
	return len(pendings) == 0, nil
}

func (s *mimicStorage) Get(key store.Key) ([]byte, error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	value, ok := s.MiscStorage[string(fullKey)]
	if !ok {
		return nil, store.ErrNotFound
	}
	return value, nil
}

func (s *mimicStorage) Set(key store.Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	s.MiscStorage[string(fullKey)] = value
	return nil
}

func (s *mimicStorage) Delete(key store.Key) error { panic("not implemented") } // NOT NEEDED
func (s *mimicStorage) Backend() *badger.DB        { panic("not implemented") } // NOT NEEDED
func (s *mimicStorage) NewIterator(pivot store.Key, reverse bool) store.Iterator {
	panic("not implemented")
} // NOT NEEDED

func NewStorage(
	pcs insolar.PlatformCryptographyScheme,
	pa pulse.Accessor,
) Storage {
	return &mimicStorage{
		pcs: pcs,
		pa:  pa,

		Code:        make(map[insolar.ID]*CodeEntity),
		Objects:     make(map[insolar.ID]*ObjectEntity),
		Requests:    make(map[insolar.ID]*RequestEntity),
		Results:     make(map[insolar.ID]*RequestEntity),
		MiscStorage: make(map[string][]byte),
	}
}
