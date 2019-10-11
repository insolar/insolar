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
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
)

// TODO[bigbes]: check for oldest mutable

type Ledger interface {
	ProcessMessage(meta payload.Meta, pl payload.Payload) []payload.Payload
}

type mimicLedger struct {
	lock sync.Mutex

	// components
	pcs insolar.PlatformCryptographyScheme
	pa  pulse.Accessor

	ctx     context.Context
	storage Storage
}

func NewMimicLedger(
	pcs insolar.PlatformCryptographyScheme,
	pa pulse.Accessor,
) Ledger {
	return &mimicLedger{
		pcs:     pcs,
		pa:      pa,
		storage: NewStorage(pcs, pa),
	}
}

func (p *mimicLedger) processGetPendings(pl *payload.GetPendings) []payload.Payload {
	requests, err := p.storage.GetPendings(pl.ObjectID, pl.Count)
	if err == ErrNotFound {
		return []payload.Payload{
			&payload.Error{
				Text: insolar.ErrNoPendingRequest.Error(),
				Code: payload.CodeNoPendings,
			},
		}
	} else if err != nil {
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return []payload.Payload{
		&payload.IDs{
			IDs: requests,
		},
	}
}

func (p *mimicLedger) processGetRequest(pl *payload.GetRequest) []payload.Payload {
	request, err := p.storage.GetRequest(pl.RequestID)
	if err == ErrRequestNotFound {
		return []payload.Payload{
			&payload.Error{
				Text: err.Error(),
				Code: payload.CodeNotFound,
			},
		}
	}

	// TODO[bigbes]: may throw if getRequest for Outgoing. Possible?
	virtReqRecord := record.Wrap(request.(*record.IncomingRequest))
	return []payload.Payload{
		&payload.Request{
			RequestID: pl.RequestID,
			Request:   virtReqRecord,
		},
	}
}

func (p *mimicLedger) setRequestCommon(request record.Request) []payload.Payload {
	requestID, reqBuf, resBuf, err := p.storage.SetRequest(request)
	switch err {
	case nil, ErrRequestExists:
		break
	case ErrRequestParentNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNonActivated,
				Text: err.Error(),
			},
		}
	case ErrNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNotFound,
				Text: err.Error(),
			},
		}
	case ErrNotActivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNonActivated,
				Text: err.Error(),
			},
		}
	case ErrAlreadyDeactivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeDeactivated,
				Text: err.Error(),
			},
		}
	default:
		panic("unexpected error: " + err.Error())
	}

	var objectID insolar.ID
	if objectRef := request.AffinityRef(); objectRef != nil {
		objectID = *objectRef.GetLocal()
	}

	if requestID == nil {
		panic("requestID is nil, shouldn't be")
	}

	return []payload.Payload{
		&payload.RequestInfo{
			ObjectID:  objectID,
			RequestID: *requestID,
			Request:   reqBuf,
			Result:    resBuf,
		},
	}
}

func (p *mimicLedger) setIncomingRequest(pl *payload.SetIncomingRequest) []payload.Payload {
	rec := record.Unwrap(&pl.Request)
	request, ok := rec.(*record.IncomingRequest)
	if !ok {
		panic(fmt.Sprintf("wrong request type, expected Incoming: %T", rec))
	}

	return p.setRequestCommon(request)
}

func (p *mimicLedger) setOutgoingRequest(pl *payload.SetOutgoingRequest) []payload.Payload {
	rec := record.Unwrap(&pl.Request)
	request, ok := rec.(*record.OutgoingRequest)
	if !ok {
		panic(fmt.Sprintf("wrong request type, expected Outgoing: %T", rec))
	}

	return p.setRequestCommon(request)
}

func (p *mimicLedger) setResultCommon(result *record.Result) ([]payload.Payload, bool) {
	resultID, err := p.storage.SetResult(result)
	switch err {
	case nil:
		break
	case ErrResultExists:
		id, resultBuf, err := p.storage.GetResult(*result.Request.GetLocal())
		if err != nil {
			panic("unexpected error: " + err.Error())
		}

		materialDuplicatedRec := record.Material{}
		if err := materialDuplicatedRec.Unmarshal(resultBuf); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal Material Result record").Error())
		}

		storedPayload := record.Unwrap(&materialDuplicatedRec.Virtual).(*record.Result).Payload
		if bytes.Compare(storedPayload, result.Payload) != 0 {
			return []payload.Payload{
				&payload.ErrorResultExists{
					ObjectID: result.Object,
					ResultID: *id,
					Result:   resultBuf,
				},
			}, true
		}

		return []payload.Payload{
			&payload.ResultInfo{
				ObjectID: result.Object,
				ResultID: *id,
			},
		}, true
	case ErrNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNotFound,
				Text: err.Error(),
			},
		}, false
	case ErrNotActivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNonActivated,
				Text: err.Error(),
			},
		}, false
	case ErrAlreadyDeactivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeDeactivated,
				Text: err.Error(),
			},
		}, false
	case ErrRequestNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeRequestNotFound,
				Text: err.Error(),
			},
		}, false
	default:
		panic("unexpected error: " + err.Error())
	}

	return []payload.Payload{
		&payload.ResultInfo{
			ObjectID: result.Object,
			ResultID: *resultID,
		},
	}, false
}

// TODO[bigbes]: check outgoings
func (p *mimicLedger) setResult(pl *payload.SetResult) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResult, _ := p.setResultCommon(result)
	return setResult
}

func (p *mimicLedger) activate(pl *payload.Activate) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}
	// resultID := setResultResult[0].(*payload.ResultInfo).ResultID

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Result); err != nil {
		p.storage.RollbackSetResult(result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	activate, ok := record.Unwrap(&virtualActivateRec).(*record.Activate)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	objectID := result.Object
	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(objectID, requestID, activate)
	if err != nil {
		p.storage.RollbackSetResult(result)

		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return setResultResult
}

func (p *mimicLedger) update(pl *payload.Update) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Result); err != nil {
		p.storage.RollbackSetResult(result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	amend, ok := record.Unwrap(&virtualActivateRec).(*record.Amend)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	objectID := result.Object
	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(objectID, requestID, amend)
	if err != nil {
		p.storage.RollbackSetResult(result)

		if err == ErrNotFound {
			return []payload.Payload{
				&payload.Error{
					Code: payload.CodeNotFound,
					Text: err.Error(),
				},
			}
		} else if err == ErrAlreadyDeactivated {
			return []payload.Payload{
				&payload.Error{
					Code: payload.CodeDeactivated,
					Text: err.Error(),
				},
			}
		}

		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return setResultResult
}
func (p *mimicLedger) deactivate(pl *payload.Deactivate) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Result); err != nil {
		p.storage.RollbackSetResult(result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	deactivate, ok := record.Unwrap(&virtualActivateRec).(*record.Deactivate)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	objectID := result.Object
	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(objectID, requestID, deactivate)
	if err != nil {
		p.storage.RollbackSetResult(result)

		if err == ErrNotFound {
			return []payload.Payload{
				&payload.Error{
					Code: payload.CodeNotFound,
					Text: err.Error(),
				},
			}
		} else if err == ErrAlreadyDeactivated {
			return []payload.Payload{
				&payload.Error{
					Code: payload.CodeDeactivated,
					Text: err.Error(),
				},
			}
		}

		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return setResultResult
}

func (p *mimicLedger) hasPendings(pl *payload.HasPendings) []payload.Payload {
	hasPendings, err := p.storage.HasPendings(pl.ObjectID)
	if err == ErrNotFound {
		return []payload.Payload{
			&payload.Error{
				Text: insolar.ErrNoPendingRequest.Error(),
				Code: payload.CodeNoPendings,
			},
		}
	} else if err != nil {
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return []payload.Payload{
		&payload.PendingsInfo{
			HasPendings: hasPendings,
		},
	}
}

func (p *mimicLedger) getObject(pl *payload.GetObject) []payload.Payload {
	state, index, firstRequestID, err := p.storage.GetObject(pl.ObjectID)
	switch err {
	case nil:
		break
	case ErrNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNotFound,
				Text: err.Error(),
			},
		}
	case ErrNotActivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNonActivated,
				Text: err.Error(),
			},
		}
	case ErrAlreadyDeactivated:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeDeactivated,
				Text: err.Error(),
			},
		}
	default:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	material := record.Material{
		Virtual:  record.Wrap(state),
		ID:       pl.ObjectID,
		ObjectID: pl.ObjectID,
		JetID:    gen.JetID(),
	}
	stateBuf, err := material.Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal Material State record").Error())
	}

	indexBuf, err := index.Lifeline.Marshal()
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal Lifeline record").Error())
	}

	return []payload.Payload{
		&payload.Index{
			Index:             indexBuf,
			EarliestRequestID: firstRequestID,
		},
		&payload.State{
			Record: stateBuf,
		},
	}
}

func (p *mimicLedger) getCode(pl *payload.GetCode) []payload.Payload { panic("implement me") }
func (p *mimicLedger) setCode(pl *payload.SetCode) []payload.Payload { panic("implement me") }

func (p *mimicLedger) ProcessMessage(meta payload.Meta, pl payload.Payload) []payload.Payload {
	p.lock.Lock()
	defer p.lock.Unlock()

	switch data := pl.(type) {
	case *payload.GetPendings:
		return p.processGetPendings(data)
	case *payload.GetRequest:
		return p.processGetRequest(data)
	case *payload.SetIncomingRequest:
		return p.setIncomingRequest(data)
	case *payload.SetOutgoingRequest:
		return p.setOutgoingRequest(data)
	case *payload.SetResult:
		return p.setResult(data)
	case *payload.Activate:
		return p.activate(data)
	case *payload.Update:
		return p.update(data)
	case *payload.Deactivate:
		return p.deactivate(data)
	case *payload.HasPendings:
		return p.hasPendings(data)
	case *payload.GetObject:
		return p.getObject(data)
	case *payload.GetCode:
		return p.getCode(data)
	case *payload.SetCode:
		return p.setCode(data)
	default:
		panic(fmt.Sprintf("unexpected message to light %T", pl))
	}
}
