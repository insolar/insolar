package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type Transcript struct {
	ObjectDescriptor artifacts.ObjectDescriptor
	Context          context.Context
	LogicContext     *insolar.LogicCallContext
	Request          *record.IncomingRequest
	RequestRef       insolar.Reference
	Nonce            uint64
	Deactivate       bool
	OutgoingRequests []OutgoingRequest
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
	}
}

func convertRecordReferenceToSelfReference(recordRef insolar.Reference) *insolar.Reference {
	if !recordRef.IsRecordScope() {
		panic("recordRef is not record reference, ref=" + recordRef.String())
	}
	recordID := recordRef.GetLocal()
	return insolar.NewReference(*recordID)
}

// NewTranscriptCloneContext creates a transcript with fresh context created from
// contextSource which can be either other Context or ServiceData. In general
// transcript shouldn't be created with context as execution can take minutes.
func NewTranscriptCloneContext(
	ctxSource interface{},
	requestRef insolar.Reference,
	request record.IncomingRequest,
) *Transcript {
	if request.CallType != record.CTMethod {
		request.Object = convertRecordReferenceToSelfReference(requestRef)
	}

	var prevCtx context.Context

	switch sourceTyped := ctxSource.(type) {
	case context.Context:
		prevCtx = freshContextFromContext(sourceTyped, request.APIRequestID)
	case *payload.ServiceData:
		prevCtx = contextFromServiceData(sourceTyped)
	default:
		panic(fmt.Errorf("unexpected type of context source: %T", ctxSource))
	}

	newCtx, _ := inslogger.WithFields(
		context.Background(),
		map[string]interface{}{
			"request": requestRef.String(),
			"object":  request.Object.String(),
			"method":  request.Method,
		},
	)
	newCtx, _ = inslogger.WithTraceField(newCtx, inslogger.TraceID(prevCtx))

	return NewTranscript(newCtx, requestRef, request)
}

func (t *Transcript) AddOutgoingRequest(
	ctx context.Context, request record.IncomingRequest, result []byte, err error,
) {
	rec := OutgoingRequest{
		Request:  request,
		Response: result,
		Error:    err,
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

func contextFromServiceData(data *payload.ServiceData) context.Context {
	ctx := inslogger.ContextWithTrace(context.Background(), data.LogTraceID)
	ctx = inslogger.WithLoggerLevel(ctx, data.LogLevel)
	if data.TraceSpanData != nil {
		parentSpan := instracer.MustDeserialize(data.TraceSpanData)
		return instracer.WithParentSpan(ctx, parentSpan)
	}
	return ctx
}

func freshContextFromContext(ctx context.Context, reqID string) context.Context {
	res := context.Background()

	logLevel := inslogger.GetLoggerLevel(ctx)
	if logLevel != insolar.NoLevel {
		res = inslogger.WithLoggerLevel(res, logLevel)
	}

	// we know that trace id is equal to APIRequestID, in a few cases we
	// call this function and don't have correct context and trace id around
	// except in request record as APIRequestID field
	res = inslogger.ContextWithTrace(res, reqID)

	parentSpan, ok := instracer.ParentSpan(ctx)
	if ok {
		res = instracer.WithParentSpan(res, parentSpan)
	}

	return res
}
