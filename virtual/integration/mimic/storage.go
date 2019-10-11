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
)

type Storage interface {
	store.DB

	// Object
	GetObject(object insolar.ID) (record.State, *record.Index, *insolar.ID, error)
	SetObject(object insolar.ID, requestID insolar.ID, newState record.State) error
	// Request
	GetRequest(request insolar.ID) (record.Request, error)
	SetRequest(request record.Request) (*insolar.ID, []byte, []byte, error)
	// Result
	GetResult(requestID insolar.ID) (*insolar.ID, []byte, error)
	SetResult(result *record.Result) (*insolar.ID, error)
	RollbackSetResult(result *record.Result)
	// Code
	SetCode(code record.Code) (insolar.ID, error)
	GetCode(codeID insolar.ID) ([]byte, error)
	// Pendings
	GetPendings(object insolar.ID, limit uint32) ([]insolar.ID, error)
	HasPendings(object insolar.ID) (bool, error)
}

type storage struct {
	// components
	pcs insolar.PlatformCryptographyScheme
	pa  pulse.Accessor

	Code        map[insolar.ID]*CodeEntity
	Objects     map[insolar.ID]*ObjectEntity
	Requests    map[insolar.ID]*RequestEntity
	Results     map[insolar.ID]*RequestEntity
	MiscStorage map[string][]byte
}

var (
	// object related errors
	ErrNotFound           = errors.New("object not found")
	ErrAlreadyActivated   = errors.New("object already activated")
	ErrAlreadyDeactivated = errors.New("object already deactivated")
	ErrDeactivated        = errors.New("object is activated")
	ErrNotActivated       = errors.New("object isn't activated")

	// request related errors
	ErrRequestParentNotFound = errors.New("parent request not found")
	ErrRequestNotFound       = errors.New("request not found")
	ErrRequestExists         = errors.New("request already exists")
	ErrResultExists          = errors.New("request result exists")
	ErrResultNotFound        = errors.New("request result not found")

	// code related errors
	ErrCodeExists   = errors.New("code already exists")
	ErrCodeNotFound = errors.New("code not found")
)

func (s *storage) calculateVirtualID(pulse insolar.PulseNumber, virtual record.Virtual) insolar.ID {
	hash := record.HashVirtual(s.pcs.ReferenceHasher(), virtual)
	return *insolar.NewID(pulse, hash)
}

func (s *storage) calculateRecordID(pulse insolar.PulseNumber, rec record.Record) insolar.ID {
	virtual := record.Wrap(rec)
	return s.calculateVirtualID(pulse, virtual)
}

func (s *storage) calculateSideEffectID(state *ObjectState) insolar.ID {
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

func (s *storage) GetObject(objectID insolar.ID) (record.State, *record.Index, *insolar.ID, error) {
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

func (s *storage) SetObject(objectID insolar.ID, requestID insolar.ID, newState record.State) error {
	switch newState.(type) {
	case *record.Activate:
		// TODO[bigbes]: take Objects lock
		if s.Objects[objectID] != nil {
			return ErrAlreadyActivated
		}

		s.Objects[objectID] = &ObjectEntity{
			ObjectChanges: []ObjectState{
				{
					State:     newState,
					RequestID: requestID,
				},
			},
			RequestsMap:  make(map[insolar.ID]*RequestEntity),
			RequestsList: nil,
		}
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

func (s *storage) GetPendings(object insolar.ID, limit uint32) ([]insolar.ID, error) {
	state := s.Objects[object]
	if state == nil {
		return nil, ErrNotFound
	}

	rv := make([]insolar.ID, 0)
	for _, request := range state.RequestsList {
		if limit == 0 {
			break
		} else if request.Status == RequestFinished {
			continue
		}

		rv = append(rv, request.ID)
		limit--
	}

	return rv, nil
}

func (s *storage) HasPendings(object insolar.ID) (bool, error) {
	pendings, err := s.GetPendings(object, 1)
	if err != nil {
		return false, err
	}
	return len(pendings) != 0, nil
}

func (s *storage) GetRequest(request insolar.ID) (record.Request, error) {
	if s.Requests[request] == nil {
		return nil, ErrRequestNotFound
	}
	return s.Requests[request].Request, nil
}

func (s *storage) SetRequest(request record.Request) (*insolar.ID, []byte, []byte, error) {
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
	if request.GetCallType() != record.CTGenesis {
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
		if !ok {
			return nil, nil, nil, ErrNotFound
		} else if len(state.ObjectChanges) == 0 {
			return nil, nil, nil, ErrNotActivated
		} else if state.isDeactivated() {
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

	return &requestID, nil, nil, nil
}

func (s *storage) GetResult(requestID insolar.ID) (*insolar.ID, []byte, error) {
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

func (s *storage) SetResult(result *record.Result) (*insolar.ID, error) {
	if !result.Object.IsEmpty() {
		state, ok := s.Objects[result.Object]
		if !ok {
			return nil, ErrNotFound
		} else if len(state.ObjectChanges) == 0 {
			return nil, ErrNotActivated
		} else if state.isDeactivated() {
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

func (s *storage) RollbackSetResult(result *record.Result) {
	request := s.Requests[*result.Request.GetLocal()]

	request.Status = RequestRegistered
	request.Result = nil
	delete(s.Results, request.ResultID)

	request.ResultID = insolar.ID{}
}

func (s *storage) GetCode(codeID insolar.ID) ([]byte, error) {
	code, ok := s.Code[codeID]
	if !ok {
		return nil, ErrCodeNotFound
	}

	return code.Code, nil
}

func (s *storage) SetCode(code record.Code) (insolar.ID, error) {
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

func (s *storage) Get(key store.Key) ([]byte, error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	value, ok := s.MiscStorage[string(fullKey)]
	if !ok {
		return nil, store.ErrNotFound
	}
	return value, nil
}

func (s *storage) Set(key store.Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	s.MiscStorage[string(fullKey)] = value
	return nil
}

func (s *storage) Delete(key store.Key) error                               { panic("not implemented") } // NOT NEEDED
func (s *storage) Backend() *badger.DB                                      { panic("not implemented") } // NOT NEEDED
func (s *storage) NewIterator(pivot store.Key, reverse bool) store.Iterator { panic("not implemented") } // NOT NEEDED

func NewStorage(
	pcs insolar.PlatformCryptographyScheme,
	pa pulse.Accessor,
) Storage {
	return &storage{
		pcs: pcs,
		pa:  pa,

		Code:        make(map[insolar.ID]*CodeEntity),
		Objects:     make(map[insolar.ID]*ObjectEntity),
		Requests:    make(map[insolar.ID]*RequestEntity),
		Results:     make(map[insolar.ID]*RequestEntity),
		MiscStorage: make(map[string][]byte),
	}
}
