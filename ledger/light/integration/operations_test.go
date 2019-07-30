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

func callSetCode(ctx context.Context, s *Server) (payload.Payload, record.Virtual) {
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

func callGetCode(ctx context.Context, s *Server, id insolar.ID) payload.Payload {
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

func callSetIncomingRequest(
	ctx context.Context, s *Server, objectID, reasonID insolar.ID, isCreation, isAPI bool,
) (payload.Payload, record.Virtual) {
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
		req.Caller = gen.Reference()
	}
	rec := record.Wrap(&req)
	reps, done := s.Send(ctx, &payload.SetIncomingRequest{
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

func callSetOutgoingRequest(
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

func callSetResult(
	ctx context.Context, s *Server, objectID, requestID insolar.ID,
) (payload.Payload, record.Virtual) {
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
	reps, done := s.Send(ctx, &payload.SetResult{
		Result: buf,
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

	return nil, record.Virtual{}
}

func sendMessage(
	ctx context.Context, s *Server, msg payload.Payload,
) payload.Payload {
	reps, done := s.Send(ctx, msg)
	defer done()

	rep := <-reps
	pl, err := payload.UnmarshalFromMeta(rep.Payload)
	panicIfErr(err)

	return pl
}
func callGetRequest(ctx context.Context, s *Server, requestID insolar.ID) payload.Payload {
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

func callActivateObject(ctx context.Context, s *Server, objectID insolar.ID) (payload.Payload, record.Virtual) {
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

func callAmendObject(ctx context.Context, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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

func callDeactivateObject(ctx context.Context, s *Server, objectID, requestID insolar.ID) (payload.Payload, record.Virtual) {
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

func callGetObject(ctx context.Context, s *Server, objectID insolar.ID) (payload.Payload, payload.Payload) {
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

func callGetPendings(ctx context.Context, s *Server, objectID insolar.ID) payload.Payload {
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
