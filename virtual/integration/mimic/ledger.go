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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/testutils"
)

// TODO[bigbes]: check for oldest mutable

type DebugLedger interface {
	AddCode(ctx context.Context, code []byte) (*insolar.ID, error)
	AddObject(ctx context.Context, image insolar.ID, isPrototype bool, memory []byte) (*insolar.ID, error)
	LoadGenesis(ctx context.Context, genesisDirectory string) error
}

type Ledger interface {
	DebugLedger

	ProcessMessage(meta payload.Meta, pl payload.Payload) []payload.Payload
}

type mimicLedger struct {
	lock sync.Mutex

	// components
	pcs       insolar.PlatformCryptographyScheme
	pAccessor pulse.Accessor
	pAppender pulse.Appender

	ctx     context.Context
	storage Storage
}

func NewMimicLedger(
	ctx context.Context,
	pcs insolar.PlatformCryptographyScheme,
	pAccessor pulse.Accessor,
	pAppender pulse.Appender,
) Ledger {
	ctx, _ = inslogger.WithField(ctx, "component", "mimic")
	return &mimicLedger{
		pcs:       pcs,
		pAppender: pAppender,
		pAccessor: pAccessor,

		ctx:     ctx,
		storage: NewStorage(pcs, pAccessor),
	}
}

func (p *mimicLedger) processGetPendings(ctx context.Context, pl *payload.GetPendings) []payload.Payload {
	requests, err := p.storage.GetPendings(ctx, pl.ObjectID, pl.Count, pl.SkipRequestRefs)
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
	default:
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

func (p *mimicLedger) processGetRequest(ctx context.Context, pl *payload.GetRequest) []payload.Payload {
	request, err := p.storage.GetRequest(ctx, pl.RequestID)
	switch err {
	case nil:
		break
	case ErrRequestNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNotFound,
				Text: err.Error(),
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

func (p *mimicLedger) setRequestCommon(ctx context.Context, request record.Request) []payload.Payload {
	requestID, reqBuf, resBuf, err := p.storage.SetRequest(ctx, request)
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

	pulseObject, err := p.pAccessor.Latest(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "failed to obtained latest pulse"))
	}
	objectID := p.storage.CalculateRequestAffinityRef(request, pulseObject.PulseNumber)

	var method string
	if req, ok := request.(*record.IncomingRequest); ok {
		method = req.Method
	} else if req, ok := request.(*record.OutgoingRequest); ok {
		method = req.Method
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"type":           request.GetCallType().String(),
		"resultObjectID": objectID.String(),
		"method":         method,
	}).Info("Registering request")

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

func (p *mimicLedger) processSetIncomingRequest(ctx context.Context, pl *payload.SetIncomingRequest) []payload.Payload {
	rec := record.Unwrap(&pl.Request)
	request, ok := rec.(*record.IncomingRequest)
	if !ok {
		panic(fmt.Sprintf("wrong request type, expected Incoming: %T", rec))
	}

	return p.setRequestCommon(ctx, request)
}

func (p *mimicLedger) processSetOutgoingRequest(ctx context.Context, pl *payload.SetOutgoingRequest) []payload.Payload {
	rec := record.Unwrap(&pl.Request)
	request, ok := rec.(*record.OutgoingRequest)
	if !ok {
		panic(fmt.Sprintf("wrong request type, expected Outgoing: %T", rec))
	}

	return p.setRequestCommon(ctx, request)
}

func (p *mimicLedger) setResultCommon(ctx context.Context, result *record.Result) ([]payload.Payload, bool) {
	resultID, err := p.storage.SetResult(ctx, result)
	switch err {
	case nil:
		break
	case ErrResultExists: // duplicate result already exists
		id, resultBuf, err := p.storage.GetResult(ctx, *result.Request.GetLocal())
		if err != nil {
			panic("unexpected error: " + err.Error())
		}

		materialDuplicatedRec := record.Material{}
		if err := materialDuplicatedRec.Unmarshal(resultBuf); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal Material Result record").Error())
		}

		storedPayload := record.Unwrap(&materialDuplicatedRec.Virtual).(*record.Result).Payload
		if !bytes.Equal(storedPayload, result.Payload) {
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
func (p *mimicLedger) processSetResult(ctx context.Context, pl *payload.SetResult) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResult, _ := p.setResultCommon(ctx, result)
	return setResult
}

func (p *mimicLedger) processActivate(ctx context.Context, pl *payload.Activate) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(ctx, result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}
	// resultID := setResultResult[0].(*payload.ResultInfo).ResultID

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Record); err != nil {
		p.storage.RollbackSetResult(ctx, result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	activate, ok := record.Unwrap(&virtualActivateRec).(*record.Activate)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(ctx, requestID, activate, insolar.ID{})
	if err != nil {
		p.storage.RollbackSetResult(ctx, result)

		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return setResultResult
}

func (p *mimicLedger) processUpdate(ctx context.Context, pl *payload.Update) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(ctx, result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Record); err != nil {
		p.storage.RollbackSetResult(ctx, result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	amend, ok := record.Unwrap(&virtualActivateRec).(*record.Amend)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	objectID := result.Object
	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(ctx, requestID, amend, objectID)
	if err != nil {
		p.storage.RollbackSetResult(ctx, result)

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
func (p *mimicLedger) processDeactivate(ctx context.Context, pl *payload.Deactivate) []payload.Payload {
	virtualRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualRec.Unmarshal(pl.Result); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	rec := record.Unwrap(&virtualRec) // record.Result
	result, ok := rec.(*record.Result)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	setResultResult, isDuplicate := p.setResultCommon(ctx, result)
	if _, ok := setResultResult[0].(*payload.ResultInfo); !ok || isDuplicate {
		return setResultResult
	}

	virtualActivateRec := record.Virtual{} // wrapped virtual record.Result
	if err := virtualActivateRec.Unmarshal(pl.Record); err != nil {
		p.storage.RollbackSetResult(ctx, result)
		panic(errors.Wrap(err, "failed to unmarshal Result record").Error())
	}

	deactivate, ok := record.Unwrap(&virtualActivateRec).(*record.Deactivate)
	if !ok {
		panic(fmt.Errorf("wrong result type: %T", rec))
	}

	objectID := result.Object
	requestID := *result.Request.GetLocal()

	err := p.storage.SetObject(ctx, requestID, deactivate, objectID)
	if err != nil {
		p.storage.RollbackSetResult(ctx, result)

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

func (p *mimicLedger) processHasPendings(ctx context.Context, pl *payload.HasPendings) []payload.Payload {
	_, err := p.storage.HasPendings(ctx, pl.ObjectID)
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
	default:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeUnknown,
				Text: err.Error(),
			},
		}
	}

	return []payload.Payload{
		&payload.PendingsInfo{
			HasPendings: false,
		},
	}
}

func (p *mimicLedger) processGetObject(ctx context.Context, pl *payload.GetObject) []payload.Payload {
	state, index, firstRequestID, err := p.storage.GetObject(ctx, pl.ObjectID)
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

func (p *mimicLedger) processGetCode(ctx context.Context, pl *payload.GetCode) []payload.Payload {
	codeBuf, err := p.storage.GetCode(ctx, pl.CodeID)
	switch err {
	case nil:
		break
	case ErrCodeNotFound:
		return []payload.Payload{
			&payload.Error{
				Code: payload.CodeNotFound,
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

	return []payload.Payload{
		&payload.Code{
			Record: codeBuf,
		},
	}
}

func (p *mimicLedger) processSetCode(ctx context.Context, pl *payload.SetCode) []payload.Payload {
	panic("implement me")
}

func (p *mimicLedger) ProcessMessage(meta payload.Meta, pl payload.Payload) []payload.Payload {
	p.lock.Lock()
	defer p.lock.Unlock()

	msgType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		panic(errors.Wrap(err, "unknown payload type"))
	}

	ctx, logger := inslogger.WithFields(p.ctx, map[string]interface{}{
		"sender":      meta.Sender.String(),
		"receiver":    meta.Receiver.String(),
		"senderPulse": meta.Pulse,
		"msgType":     msgType.String(),
	})
	logger.Info("Processing message")

	var result []payload.Payload

	switch data := pl.(type) {
	case *payload.GetPendings:
		result = p.processGetPendings(ctx, data)
	case *payload.GetRequest:
		result = p.processGetRequest(ctx, data)
	case *payload.SetIncomingRequest:
		result = p.processSetIncomingRequest(ctx, data)
	case *payload.SetOutgoingRequest:
		result = p.processSetOutgoingRequest(ctx, data)
	case *payload.SetResult:
		result = p.processSetResult(ctx, data)
	case *payload.Activate:
		result = p.processActivate(ctx, data)
	case *payload.Update:
		result = p.processUpdate(ctx, data)
	case *payload.Deactivate:
		result = p.processDeactivate(ctx, data)
	case *payload.HasPendings:
		result = p.processHasPendings(ctx, data)
	case *payload.GetObject:
		result = p.processGetObject(ctx, data)
	case *payload.GetCode:
		result = p.processGetCode(ctx, data)
	case *payload.SetCode:
		result = p.processSetCode(ctx, data)
	default:
		panic(fmt.Sprintf("unexpected message to light %T", pl))
	}

	if err, ok := result[0].(*payload.Error); ok {
		logger.WithField("error", err.Text).Error("Failed to process message")
	}

	return result
}

func (p *mimicLedger) AddObject(ctx context.Context, image insolar.ID, isPrototype bool, memory []byte) (*insolar.ID, error) {
	id, _, _, err := p.storage.SetRequest(ctx, &record.IncomingRequest{
		CallType: record.CTGenesis,
		Method:   testutils.RandomString(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to set request")
	}

	requestRef := *insolar.NewRecordReference(*id)

	result := &record.Result{
		Object:  insolar.ID{},
		Request: requestRef,
		Payload: []byte{},
	}
	_, err = p.storage.SetResult(ctx, result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set result")
	}

	err = p.storage.SetObject(ctx, *id, &record.Activate{
		Request:     requestRef,
		Memory:      memory,
		Image:       *insolar.NewReference(image),
		IsPrototype: isPrototype,
	}, insolar.ID{})
	if err != nil {
		p.storage.RollbackSetResult(ctx, result)
		return nil, errors.Wrap(err, "failed to activate object")
	}

	return id, nil
}

func (p *mimicLedger) AddCode(ctx context.Context, code []byte) (*insolar.ID, error) {
	id, err := p.storage.SetCode(ctx, record.Code{
		Code:        code,
		MachineType: insolar.MachineTypeGoPlugin,
	})
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *mimicLedger) LoadGenesis(ctx context.Context, dirPath string) error {
	genesisContractsConfig, err := ReadGenesisContractsConfig(dirPath)
	if err != nil {
		return errors.Wrap(err, "failed to load genesis config")
	}

	genesisObject := &genesis.Genesis{
		ArtifactManager: NewClient(p.storage),
		BaseRecord: &genesis.BaseRecord{
			DB:             p.storage,
			DropModifier:   &dropModifierMock{},
			PulseAppender:  p.pAppender,
			PulseAccessor:  p.pAccessor,
			RecordModifier: &recordModifierMock{},
			IndexModifier:  &indexModifierMock{},
		},
		DiscoveryNodes:  []insolar.DiscoveryNodeRegister{},
		ContractsConfig: *genesisContractsConfig,
	}

	if err := genesisObject.Start(ctx); err != nil {
		return errors.Wrap(err, "failed to load genesis")
	}

	return nil
}
