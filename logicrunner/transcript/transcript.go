package transcript

import (
	"context"
	"reflect"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type OutgoingRequest struct {
	Request   record.IncomingRequest
	NewObject *insolar.Reference
	Response  []byte
	Error     error
}

type Transcript struct {
	ObjectDescriptor artifacts.ObjectDescriptor
	Context          context.Context
	LogicContext     *insolar.LogicCallContext
	Request          *record.IncomingRequest
	RequestRef       insolar.Reference
	RequesterNode    *insolar.Reference
	Nonce            uint64
	Deactivate       bool
	OutgoingRequests []OutgoingRequest
	FromLedger       bool
}

func NewTranscript(
	ctx context.Context,
	requestRef insolar.Reference,
	request record.IncomingRequest,
) *Transcript {

	return &Transcript{
		Context:    ctx,
		Request:    &request,
		RequestRef: requestRef,
		Nonce:      0,
		Deactivate: false,

		FromLedger: false,
	}
}

func (t *Transcript) AddOutgoingRequest(
	ctx context.Context, request record.IncomingRequest, result []byte, newObject *insolar.Reference, err error,
) {
	rec := OutgoingRequest{
		Request:   request,
		Response:  result,
		NewObject: newObject,
		Error:     err,
	}
	t.OutgoingRequests = append(t.OutgoingRequests, rec)
}

func (t *Transcript) HasOutgoingRequest(
	ctx context.Context, request record.IncomingRequest,
) *OutgoingRequest {
	for i := range t.OutgoingRequests {
		if reflect.DeepEqual(t.OutgoingRequests[i].Request, request) {
			return &t.OutgoingRequests[i]
		}
	}
	return nil
}
