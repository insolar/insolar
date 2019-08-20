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

package integration_test

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
)

func CallSetCode(ctx context.Context, s *Server) (payload.Payload, record.Virtual) {
	code := make([]byte, 100)
	_, err := rand.Read(code)
	panicIfErr(err)
	rec := record.Wrap(&record.Code{Code: code})
	buf, err := rec.Marshal()
	panicIfErr(err)
	reps, done := s.Send(ctx, &payload.SetCode{
		Record: buf,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.ID:
		return pl, rec
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}

	return nil, rec
}

func CallGetCode(ctx context.Context, s *Server, id insolar.ID) payload.Payload {
	reps, done := s.Send(ctx, &payload.GetCode{
		CodeID: id,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl
	case *payload.Code:
		return pl
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}
	return nil
}

// func MakeSetIncomingRequestFromAPI(objectID, reasonID insolar.ID, isCreation bool) (payload.SetIncomingRequest, record.Virtual) {
// 	return MakeSetIncomingRequest(objectID, reasonID, insolar.ID{}, isCreation, true)
// }

func MakeSetIncomingRequest(objectID, reasonID insolar.ID, reasonObjectID insolar.ID, isCreation, isAPI bool) (payload.SetIncomingRequest, record.Virtual) {
	args := make([]byte, 100)
	_, err := rand.Read(args)
	panicIfErr(err)

	req := record.IncomingRequest{
		Arguments: args,
		Reason:    *insolar.NewReference(reasonID),
	}

	if isCreation {
		req.CallType = record.CTSaveAsChild
	} else {
		req.Object = insolar.NewReference(objectID)
	}

	if isAPI {
		req.APINode = gen.Reference()
	} else {
		req.Caller = *insolar.NewReference(reasonObjectID)
	}

	rec := record.Wrap(&req)
	pl := payload.SetIncomingRequest{
		Request: rec,
	}
	return pl, rec
}

func MakeSetOutgoingRequest(
	objectID, reasonID insolar.ID, detached bool,
) (payload.SetOutgoingRequest, record.Virtual) {
	args := make([]byte, 100)
	_, err := rand.Read(args)
	panicIfErr(err)

	rm := record.ReturnResult
	if detached {
		rm = record.ReturnSaga
	}
	req := record.OutgoingRequest{
		Caller:     *insolar.NewReference(objectID),
		Arguments:  args,
		ReturnMode: rm,
		Reason:     *insolar.NewReference(reasonID),
		APINode:    gen.Reference(),
	}

	rec := record.Wrap(&req)

	pl := payload.SetOutgoingRequest{
		Request: rec,
	}
	return pl, rec
}

func CallSetOutgoingRequest(
	ctx context.Context, s *Server, objectID, reasonID insolar.ID, detached bool,
) (payload.Payload, record.Virtual) {
	args := make([]byte, 100)
	_, err := rand.Read(args)
	panicIfErr(err)
	rm := record.ReturnResult
	if detached {
		rm = record.ReturnSaga
	}
	rec := record.Wrap(&record.OutgoingRequest{
		Caller:     *insolar.NewReference(objectID),
		Arguments:  args,
		ReturnMode: rm,
		Reason:     *insolar.NewReference(reasonID),
		APINode:    gen.Reference(),
	})
	reps, done := s.Send(ctx, &payload.SetOutgoingRequest{
		Request: rec,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.RequestInfo:
		return pl, rec
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}

	return nil, record.Virtual{}
}

func MakeSetResult(objectID, requestID insolar.ID) (payload.SetResult, record.Virtual) {
	data := make([]byte, 100)
	_, err := rand.Read(data)
	panicIfErr(err)
	rec := record.Wrap(&record.Result{
		Object:  objectID,
		Request: *insolar.NewReference(requestID),
		Payload: data,
	})
	buf, err := rec.Marshal()
	panicIfErr(err)
	pl := payload.SetResult{
		Result: buf,
	}
	return pl, rec
}

func SendMessage(
	ctx context.Context, s *Server, msg payload.Payload,
) payload.Payload {
	reps, done := s.Send(ctx, msg)
	defer done()

	rep, ok := <-reps
	if !ok {
		panic("no reply")
	}
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	return pl
}

func CallGetRequest(ctx context.Context, s *Server, requestID insolar.ID) payload.Payload {
	reps, done := s.Send(ctx, &payload.GetRequest{
		RequestID: requestID,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl
	case *payload.Request:
		return pl
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}

	return nil
}

func CallActivateObject(ctx context.Context, s *Server, objectID insolar.ID) (payload.Payload, record.Virtual) {
	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	panicIfErr(err)
	rec := record.Wrap(&record.Activate{
		Request: *insolar.NewReference(objectID),
		Memory:  mem,
	})
	buf, err := rec.Marshal()
	panicIfErr(err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	panicIfErr(err)
	resultRecord := record.Wrap(&record.Result{
		Request: *insolar.NewReference(objectID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	panicIfErr(err)

	reps, done := s.Send(ctx, &payload.Activate{
		Record: buf,
		Result: resBuf,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl, rec
	case *payload.ResultInfo:
		return pl, rec
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}
	return nil, rec
}

func CallAmendObject(ctx context.Context, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	panicIfErr(err)
	rec := record.Wrap(&record.Amend{
		Memory: mem,
	})
	buf, err := rec.Marshal()
	panicIfErr(err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	panicIfErr(err)
	resultRecord := record.Wrap(&record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	panicIfErr(err)

	reps, done := s.Send(ctx, &payload.Update{
		Record: buf,
		Result: resBuf,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.ResultInfo:
		return pl, rec
	case *payload.Error:
		return pl, rec
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}
	return nil, rec
}

func CallDeactivateObject(ctx context.Context, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
	mem := make([]byte, 100)
	_, err := rand.Read(mem)
	panicIfErr(err)
	rec := record.Wrap(&record.Deactivate{
		Request: *insolar.NewReference(objectID),
	})
	buf, err := rec.Marshal()
	panicIfErr(err)
	res := make([]byte, 100)
	_, err = rand.Read(res)
	panicIfErr(err)
	resultRecord := record.Wrap(&record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
		Payload: res,
	})
	resBuf, err := resultRecord.Marshal()
	panicIfErr(err)

	reps, done := s.Send(ctx, &payload.Deactivate{
		Record: buf,
		Result: resBuf,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.ResultInfo:
		return pl, rec
	case *payload.Error:
		return pl, rec
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}
	return pl, rec
}

func CallGetObject(ctx context.Context, s *Server, objectID insolar.ID) (payload.Payload, payload.Payload) {
	reps, d := s.Send(ctx, &payload.GetObject{
		ObjectID: objectID,
	})
	defer d()

	var (
		lifeline, state payload.Payload
	)
	done := func() bool {
		return lifeline != nil && state != nil
	}
	for rep := range reps {
		pl, err := payload.UnmarshalFromMeta(rep.Payload)
		panicIfErr(err)
		switch pl.(type) {
		case *payload.Index:
			lifeline = pl
		case *payload.State:
			state = pl
		case *payload.Error:
			return pl, nil
		default:
			panic(fmt.Sprintf("received unexpected reply %T", pl))
		}

		if done() {
			break
		}
	}
	if !done() {
		panic("no reply from GetObject")
	}
	return lifeline, state
}

func CallGetPendings(ctx context.Context, s *Server, objectID insolar.ID) payload.Payload {
	reps, done := s.Send(ctx, &payload.GetPendings{
		ObjectID: objectID,
	})
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)
	switch pl.(type) {
	case *payload.Error:
		return pl
	case *payload.IDs:
		return pl
	default:
		panic(fmt.Sprintf("received unexpected reply %T", pl))
	}

	return nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RequireNotError(pl payload.Payload) {
	if err, ok := pl.(*payload.Error); ok {
		panic(err.Text)
	}
}

func RequireError(pl payload.Payload) {
	if _, ok := pl.(*payload.Error); !ok {
		panic("expected error")
	}
}

func RequireErrorCode(pl payload.Payload, expectedCode uint32) {
	RequireError(pl)
	err := pl.(*payload.Error)
	if err.Code != expectedCode {
		panic(fmt.Sprintf("expected error code %d, got %d (%s)", expectedCode, err.Code, err.Text))
	}
}
