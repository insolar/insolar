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
// +build slowtest

package integration

import (
	"github.com/insolar/insolar/insolar/payload"
)

// call -> API -> Contract 1 -> Contract 2 -> <Save Result> Error
//
// call -> API -> Contract 1 <- Contract 2

type PseudoLedger interface {
	ProcessMessage(meta payload.Meta, pl payload.Payload) []payload.Payload
}

// func (p *pseudoLedger) ProcessMessage1(meta payload.Meta, pl payload.Payload) []payload.Payload {
// 	lifeTime := &record.Lifeline{
// 		LatestState: objID,
// 	}
// 	indexData, err := lifeTime.Marshal()
// 	require.NoError(t, err)
//
// 	codeRecord := &record.Material{
// 		Virtual: record.Wrap(&record.Code{
// 			Request: reqRef,
// 			Code:    []byte("no code"),
// 		}),
// 	}
// 	codeData, err := codeRecord.Marshal()
// 	require.NoError(t, err)
//
// 	call := 0
//
// 	switch data := pl.(type) {
// 	case *payload.GetPendings:
// 		if call > 0 {
// 			return []payload.Payload{
// 				&payload.Error{
// 					Code: payload.CodeNoPendings,
// 					Text: "No pendings",
// 				},
// 			}
// 		}
// 		call++
// 		return []payload.Payload{
// 			&payload.IDs{IDs: []insolar.ID{requestID}},
// 		}
// 	case *payload.GetRequest:
// 		req := &record.IncomingRequest{
// 			Object:       &objRef,
// 			Method:       "good.CallMethod",
// 			Arguments:    nil,
// 			APIRequestID: utils.TraceID(ctx),
// 			APINode:      virtual.ref,
// 			Reason:       api.MakeReason(insolar.GenesisPulse.PulseNumber, nil),
// 			Immutable:    true,
// 		}
//
// 		virtReqRecord := record.Wrap(req)
//
// 		return []payload.Payload{
// 			&payload.Request{
// 				RequestID: data.RequestID,
// 				Request:   virtReqRecord,
// 			},
// 		}
// 		// getters
// 	case *payload.SetIncomingRequest:
// 		buf, err := data.Request.Marshal()
// 		if err != nil {
// 			panic(errors.Wrap(err, "failed to marshal request"))
// 		}
// 		requestID = *insolar.NewID(gen.PulseNumber(), hasher.Hash(buf))
// 		return []payload.Payload{&payload.RequestInfo{
// 			ObjectID:  *objID,
// 			RequestID: requestID,
// 		}}
// 	case *payload.SetOutgoingRequest:
// 		return []payload.Payload{&payload.RequestInfo{
// 			ObjectID:  *objID,
// 			RequestID: requestID,
// 		}}
// 		// setters
// 	case *payload.SetResult:
// 		return []payload.Payload{&payload.ResultInfo{
// 			ResultID: *insolar.NewID(gen.PulseNumber(), hasher.Hash(data.Result)),
// 		}}
// 	case *payload.Activate:
// 		return []payload.Payload{&payload.ResultInfo{}}
// 	case *payload.HasPendings:
// 		return []payload.Payload{&payload.PendingsInfo{HasPendings: false}}
// 	case *payload.GetObject:
// 		resRecord := &record.Activate{
// 			Request: *insolar.NewReference(requestID),
// 			Image:   walletRef,
// 		}
//
// 		if data.ObjectID == *walletRef.GetLocal() {
// 			resRecord.IsPrototype = true
// 		}
//
// 		virtResRecord := record.Wrap(resRecord)
// 		recMaterial := &record.Material{
// 			Virtual:  virtResRecord,
// 			ID:       gen.ID(),
// 			ObjectID: gen.ID(),
// 		}
// 		recordData, err := recMaterial.Marshal()
// 		require.NoError(t, err)
//
// 		return []payload.Payload{
// 			&payload.Index{
// 				Index: indexData,
// 			},
// 			&payload.State{Record: recordData},
// 		}
// 	case *payload.GetCode:
// 		return []payload.Payload{
// 			&payload.Code{
// 				Record: codeData,
// 			},
// 		}
// 	case *payload.Update:
// 		return []payload.Payload{
// 			&payload.ResultInfo{
// 				ResultID: *insolar.NewID(gen.PulseNumber(), hasher.Hash(data.Result)),
// 			},
// 		}
// 	}
//
// 	panic(fmt.Sprintf("unexpected message to light %T", pl))
// }
//
// func TestVirtual_UnwantedResponse(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
// 	cfg := DefaultVMConfig()
//
// 	t.Run("after save error", func(t *testing.T) {
// 		objRef := gen.Reference()
// 		reqRef := gen.RecordReference()
// 		objID := objRef.GetLocal()
//
// 		var requestID insolar.ID
//
// 		expectedRes := struct {
// 			blip string
// 		}{
// 			blip: "blop",
// 		}
//
// 		hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
//
// 		leMock := testutils.NewMachineLogicExecutorMock(t).CallMethodMock.Set(
// 			func(_ context.Context, _ *insolar.LogicCallContext, _ insolar.Reference, _ []byte, _ string, _ insolar.Arguments) ([]byte, insolar.Arguments, error) {
// 				return insolar.MustSerialize(expectedRes), insolar.MustSerialize(expectedRes), nil
// 			},
// 		)
//
// 		mmMock := machinesmanager.NewMachinesManagerMock(t).GetExecutorMock.Set(
// 			func(_ insolar.MachineType) (insolar.MachineLogicExecutor, error) {
// 				return leMock, nil
// 			},
// 		).RegisterExecutorMock.Set(
// 			func(_ insolar.MachineType, _ insolar.MachineLogicExecutor) error {
// 				return nil
// 			},
// 		)
//
// 		s, err := NewServer(t, ctx, cfg, flowRouter, mmMock)
//
// 		require.NoError(t, err)
// 		defer s.Stop(ctx)
//
// 		// First pulse goes in storage then interrupts.
// 		s.SetPulse(ctx)
//
// 		res, requestRef, err := CallContract(
// 			s, &objRef, "good.CallMethod", nil, s.pulse.PulseNumber,
// 		)
//
// 		require.NoError(t, err)
// 		require.NotEmpty(t, requestRef)
// 		require.Equal(t, &reply.CallMethod{
// 			Object: &objRef,
// 			Result: insolar.MustSerialize(expectedRes),
// 		}, res)
// 	})
// }
//
// var walletRef = shouldLoadRef("0111A5e49cJW6GKGegWBhtgrJs7nFh1kSWhBtT2VgK4t.record")
//
// func shouldLoadRef(strRef string) insolar.Reference {
// 	ref, err := insolar.NewReferenceFromBase58(strRef)
// 	if err != nil {
// 		panic(errors.Wrap(err, "Unexpected error, bailing out"))
// 	}
// 	return *ref
// }
