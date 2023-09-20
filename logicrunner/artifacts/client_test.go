package artifacts

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type TestRequestResult struct {
	SideEffectType     RequestResultType // every
	RawResult          []byte            // every
	RawObjectReference insolar.Reference // every

	ParentReference insolar.Reference // activate
	ObjectImage     insolar.Reference // amend + activate
	ObjectStateID   insolar.ID        // amend + deactivate
	Memory          []byte            // amend + activate
}

func (s *TestRequestResult) Result() []byte {
	return s.RawResult
}

func (s *TestRequestResult) Activate() (insolar.Reference, insolar.Reference, []byte) {
	return s.ParentReference, s.ObjectImage, s.Memory
}

func (s *TestRequestResult) Amend() (insolar.ID, insolar.Reference, []byte) {
	return s.ObjectStateID, s.ObjectImage, s.Memory
}

func (s *TestRequestResult) Deactivate() insolar.ID {
	return s.ObjectStateID
}

func (s TestRequestResult) Type() RequestResultType {
	return s.SideEffectType
}

func (s *TestRequestResult) ObjectReference() insolar.Reference {
	return s.RawObjectReference
}

func genAPIRequestID() string {
	APIRequestID := utils.RandTraceID()
	if strings.Contains(APIRequestID, "createRandomTraceIDFailed") {
		panic("Failed to generate uuid: " + APIRequestID)
	}
	return APIRequestID
}

func genIncomingRequest() (*record.IncomingRequest, insolar.Reference) {
	baseRef := gen.Reference()
	objectRef := gen.Reference()
	prototypeRef := gen.Reference()

	return &record.IncomingRequest{
		Polymorph:       rand.Int31(),
		CallType:        record.CTMethod,
		Caller:          gen.Reference(),
		CallerPrototype: gen.Reference(),
		Nonce:           0,
		ReturnMode:      record.ReturnSaga,
		Immutable:       false,
		Base:            &baseRef,
		Object:          &objectRef,
		Prototype:       &prototypeRef,
		Method:          "Call",
		Arguments:       []byte{0x80},
		APIRequestID:    genAPIRequestID(),
		Reason:          gen.RecordReference(),
		APINode:         insolar.Reference{},
	}, gen.RecordReference()
}

type ArtifactsMangerClientSuite struct {
	suite.Suite

	mc  *minimock.Controller
	ctx context.Context

	busSender     *bus.SenderMock
	pulseAccessor *pulse.AccessorMock

	amClientOriginal *client
	amClient         Client
}

// Init and run suite
func TestArtifactManager(t *testing.T) {
	suite.Run(t, &ArtifactsMangerClientSuite{Suite: suite.Suite{}})
}

func (s *ArtifactsMangerClientSuite) prepareAMClient() {
	s.mc = minimock.NewController(s.T())

	s.pulseAccessor = pulse.NewAccessorMock(s.mc)
	s.busSender = bus.NewSenderMock(s.mc)

	s.amClientOriginal = &client{
		PCS:           platformpolicy.NewPlatformCryptographyScheme(),
		PulseAccessor: s.pulseAccessor,

		sender:       s.busSender,
		localStorage: newLocalStorage(),
	}

	s.amClient = s.amClientOriginal
}

func (s *ArtifactsMangerClientSuite) prepareContext() {
	s.ctx = inslogger.TestContext(s.T())
}

func (s *ArtifactsMangerClientSuite) BeforeTest(suiteName, testName string) {
	s.prepareContext()
	s.prepareAMClient()
}

func (s *ArtifactsMangerClientSuite) AfterTest(suiteName, testName string) {
	s.mc.Finish()
}

func (s *ArtifactsMangerClientSuite) TestGetRequest() {
	// Arrange
	iRequest, requestRef := genIncomingRequest()
	oRequest := (*record.OutgoingRequest)(iRequest)

	objectRef := iRequest.Object

	for name, test := range map[string]struct {
		response payload.Payload
		check    func(record.Request, error)
	}{
		"incoming success": {
			response: &payload.Request{
				RequestID: *requestRef.GetLocal(),
				Request:   record.Wrap(iRequest),
			},
			check: func(gotRequest record.Request, err error) {
				s.NoError(err)
				s.Equal(*iRequest, *gotRequest.(*record.IncomingRequest))
			},
		},
		"outgoing success": {
			response: &payload.Request{
				RequestID: *requestRef.GetLocal(),
				Request:   record.Wrap(oRequest),
			},
			check: func(gotRequest record.Request, err error) {
				s.NoError(err)
				s.Equal(*oRequest, *gotRequest.(*record.OutgoingRequest))
			},
		},
		"not found": {
			response: &payload.Error{
				Text: "request not found",
				Code: payload.CodeNotFound,
			},
			check: func(_ record.Request, err error) {
				s.Equal(insolar.ErrNotFound, err)
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(_ record.Request, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(_ record.Request, err error) {
				s.Contains(err.Error(), "unexpected reply")
			},
		},
		"unexpected message": {
			response: &payload.Request{
				RequestID: *requestRef.GetLocal(),
				Request:   record.Wrap(&record.Activate{}),
			},
			check: func(_ record.Request, err error) {
				s.Contains(err.Error(), "unexpected message")
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			reqMsg, err := payload.NewMessage(test.response)
			s.Require().NoError(err)

			s.busSender.SendRoleMock.Set(
				func(
					_ context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					s.Equal(insolar.DynamicRoleLightExecutor, role)

					getReq := payload.GetRequest{}
					err := getReq.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					s.Equal(*requestRef.GetLocal(), getReq.RequestID)
					s.Equal(*iRequest.Object, ref)

					meta := payload.Meta{
						Payload: reqMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					reqMsg.Payload = buf

					ch := make(chan *wmMessage.Message, 1)
					ch <- reqMsg
					return ch, func() {}
				},
			)

			// Act
			recordRequest, err := s.amClient.GetRequest(s.ctx, *objectRef, requestRef)
			// Check
			test.check(recordRequest, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestRequest_FailedToSend() {
	// Arrange
	request, requestRef := genIncomingRequest()

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.GetRequest(s.ctx, *request.Object, requestRef)

	// Assert
	s.Error(err)
	s.Contains(err.Error(), "failed to send GetRequest")
}

func (s *ArtifactsMangerClientSuite) TestGetPendings() {
	// Arrange
	objectRef := gen.Reference()
	requestRef1 := gen.RecordReference()
	requestRef2 := gen.RecordReference()
	requestRef3 := gen.RecordReference()

	for name, test := range map[string]struct {
		response payload.Payload
		check    func([]insolar.Reference, error)
	}{
		"success": {
			response: &payload.IDs{
				IDs: []insolar.ID{*requestRef1.GetLocal()},
			},
			check: func(requests []insolar.Reference, err error) {
				s.NoError(err)
				s.Equal([]insolar.Reference{requestRef1}, requests)
			},
		},
		"success_multiple": {
			response: &payload.IDs{
				IDs: []insolar.ID{
					*requestRef1.GetLocal(),
					*requestRef2.GetLocal(),
					*requestRef3.GetLocal(),
				},
			},
			check: func(requests []insolar.Reference, err error) {
				s.NoError(err)
				s.Len(requests, 3)
				s.Contains(requests, requestRef1)
				s.Contains(requests, requestRef2)
				s.Contains(requests, requestRef3)
			},
		},
		"no pendings": {
			response: &payload.Error{
				Text: insolar.ErrNoPendingRequest.Error(),
				Code: payload.CodeNoPendings,
			},
			check: func(requests []insolar.Reference, err error) {
				s.Len(requests, 0)
				s.Error(err)
				s.Equal(err, insolar.ErrNoPendingRequest)
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(requests []insolar.Reference, err error) {
				s.Len(requests, 0)
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(requests []insolar.Reference, err error) {
				s.Len(requests, 0)
				s.Contains(err.Error(), "unexpected reply")
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			resMsg, err := payload.NewMessage(test.response)
			s.Require().NoError(err)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					getPendings := payload.GetPendings{}
					err := getPendings.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					s.Equal(objectRef.GetLocal(), &getPendings.ObjectID)

					meta := payload.Meta{
						Payload: resMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch := make(chan *wmMessage.Message, 1)
					ch <- resMsg
					return ch, func() {}
				},
			)

			// Act
			references, err := s.amClient.GetPendings(s.ctx, objectRef, make([]insolar.ID, 0))

			// Assert
			test.check(references, err)
		})

	}
}

func (s *ArtifactsMangerClientSuite) TestGetPendings_FailedToSend() {
	// Arrange
	request, _ := genIncomingRequest()

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.GetPendings(s.ctx, *request.Object, make([]insolar.ID, 0))

	// Assert
	s.Error(err)
	s.Contains(err.Error(), "failed to send GetPendings")
}

func (s *ArtifactsMangerClientSuite) TestHasPendings() {
	// Arrange
	objectRef := gen.Reference()

	for name, test := range map[string]struct {
		response payload.Payload
		check    func(bool, error)
	}{
		"success ok": {
			response: &payload.PendingsInfo{
				HasPendings: true,
			},
			check: func(hasPendings bool, err error) {
				s.NoError(err)
				s.Equal(true, hasPendings)
			},
		},
		"success not ok": {
			response: &payload.PendingsInfo{
				HasPendings: false,
			},
			check: func(hasPendings bool, err error) {
				s.NoError(err)
				s.Equal(false, hasPendings)
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(_ bool, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(_ bool, err error) {
				s.Contains(err.Error(), "unexpected reply")
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			resMsg, err := payload.NewMessage(test.response)
			s.Require().NoError(err)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					hasPendings := payload.HasPendings{}
					err := hasPendings.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					s.Equal(*objectRef.GetLocal(), hasPendings.ObjectID)

					meta := payload.Meta{
						Payload: resMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch := make(chan *wmMessage.Message, 1)
					ch <- resMsg
					return ch, func() {}
				},
			)

			// Act
			hasPendings, err := s.amClient.HasPendings(s.ctx, objectRef)

			// Assert
			test.check(hasPendings, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestHasPendings_FailedToSend() {
	// Arrange
	request, _ := genIncomingRequest()

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.HasPendings(s.ctx, *request.Object)

	// Assert
	s.Error(err)
	s.Contains(err.Error(), "failed to send HasPendings")
}

func (s *ArtifactsMangerClientSuite) TestDeployCode() {
	// Arrange
	codeID := gen.ID()
	code := []byte(testutils.RandomString())
	machineType := insolar.MachineTypeBuiltin
	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}

	for name, test := range map[string]struct {
		response payload.Payload
		check    func(*insolar.ID, error)
	}{
		"success": {
			response: &payload.ID{
				ID: codeID,
			},
			check: func(id *insolar.ID, err error) {
				s.NoError(err)
				s.Equal(codeID, *id)
			},
		},
		"error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(id *insolar.ID, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
				s.Nil(id)
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(id *insolar.ID, err error) {
				s.Error(err)
				s.Contains(err.Error(), "unexpected reply")
				s.Nil(id)
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()
			s.prepareAMClient()

			resMsg, err := payload.NewMessage(test.response)
			s.Require().NoError(err)

			s.pulseAccessor.LatestMock.Return(pulseObject, nil)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					payloadSetCode := payload.SetCode{}
					err := payloadSetCode.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					virtualRec := &record.Virtual{}
					err = virtualRec.Unmarshal(payloadSetCode.Record)
					s.Require().NoError(err)

					rec := record.Unwrap(virtualRec)
					s.Require().IsType((*record.Code)(nil), rec)

					s.Equal(code, rec.(*record.Code).Code)
					s.Equal(machineType, rec.(*record.Code).MachineType)

					meta := payload.Meta{
						Payload: resMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch := make(chan *wmMessage.Message, 1)
					ch <- resMsg
					return ch, func() {}
				},
			)

			// Act
			deployCodeID, err := s.amClient.DeployCode(s.ctx, code, machineType)

			// Assert
			test.check(deployCodeID, err)

			s.mc.Finish()
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestDeployCode_FailedToSend() {
	// Arrange
	code := []byte(testutils.RandomString())
	machineType := insolar.MachineTypeBuiltin

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	s.pulseAccessor.LatestMock.Return(pulseObject, nil)

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.DeployCode(s.ctx, code, machineType)

	// Assert
	s.Error(err)
	s.Contains(err.Error(), "failed to send SetCode")
}

func (s *ArtifactsMangerClientSuite) TestRegisterIncomingRequest() {
	// Arrange
	incoming, requestRef := genIncomingRequest()
	objectRef := *incoming.Object

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}

	for name, test := range map[string]struct {
		response payload.Payload
		check    func(*payload.RequestInfo, error)
	}{
		"success": {
			response: &payload.RequestInfo{
				ObjectID:  *objectRef.GetLocal(),
				RequestID: *requestRef.GetLocal(),
				Request:   nil,
				Result:    nil,
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Equal(*objectRef.GetLocal(), requestInfo.ObjectID)
				s.Equal(*requestRef.GetLocal(), requestInfo.RequestID)
				s.Nil(requestInfo.Request)
				s.Nil(requestInfo.Result)
			},
		},
		"flow cancelled": {
			response: &payload.Error{
				Code: payload.CodeFlowCanceled,
				Text: "flow cancelled",
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Nil(requestInfo)
				s.Error(err)
				s.Contains(err.Error(), flow.ErrCancelled.Error())
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
				s.Nil(requestInfo)
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Error(err)
				s.Contains(err.Error(), "unexpected reply")
				s.Nil(requestInfo)
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()
			s.prepareAMClient()

			resMsg, err := payload.NewMessage(test.response)
			s.Require().NoError(err)

			s.pulseAccessor.LatestMock.Return(pulseObject, nil)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					payloadSetIncomingRequest := payload.SetIncomingRequest{}
					err := payloadSetIncomingRequest.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					rec := record.Unwrap(&payloadSetIncomingRequest.Request)
					s.Require().IsType((*record.IncomingRequest)(nil), rec)

					meta := payload.Meta{
						Payload: resMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch := make(chan *wmMessage.Message, 1)
					ch <- resMsg
					return ch, func() {}
				},
			)

			reqInfo, err := s.amClient.RegisterIncomingRequest(s.ctx, incoming)

			test.check(reqInfo, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestRegisterOutgoingRequest() {
	// Arrange
	incoming, requestRef := genIncomingRequest()
	objectRef := *incoming.Object
	outgoing := (*record.OutgoingRequest)(incoming)
	initialPulseNumber, pulseOffset := gen.PulseNumber(), insolar.PulseNumber(0)

	for name, test := range map[string]struct {
		response []payload.Payload
		check    func(*payload.RequestInfo, error)
	}{
		"success": {
			response: []payload.Payload{
				&payload.RequestInfo{
					ObjectID:  *objectRef.GetLocal(),
					RequestID: *requestRef.GetLocal(),
					Request:   nil,
					Result:    nil,
				},
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Equal(*objectRef.GetLocal(), requestInfo.ObjectID)
				s.Equal(*requestRef.GetLocal(), requestInfo.RequestID)
				s.Nil(requestInfo.Request)
				s.Nil(requestInfo.Result)
			},
		},
		"flow cancelled": {
			response: []payload.Payload{
				&payload.Error{
					Code: payload.CodeFlowCanceled,
					Text: "flow cancelled",
				},
				&payload.Error{
					Code: payload.CodeFlowCanceled,
					Text: "flow cancelled",
				},
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Nil(requestInfo)
				s.Error(err)
				s.Contains(err.Error(), "timeout while awaiting reply from watermill")
			},
		},
		"other error": {
			response: []payload.Payload{
				&payload.Error{
					Text: "some error",
					Code: payload.CodeUnknown,
				},
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
				s.Nil(requestInfo)
			},
		},
		"unknown payload": {
			response: []payload.Payload{
				&payload.PendingFinished{},
			},
			check: func(requestInfo *payload.RequestInfo, err error) {
				s.Error(err)
				s.Contains(err.Error(), "unexpected reply")
				s.Nil(requestInfo)
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()
			s.prepareAMClient()

			s.pulseAccessor.LatestMock.Set(
				func(ctx context.Context) (insolar.Pulse, error) {
					pulseOffset += 10
					return insolar.Pulse{PulseNumber: initialPulseNumber + pulseOffset}, nil
				},
			)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					payloadSetOutgoingRequest := payload.SetOutgoingRequest{}
					err := payloadSetOutgoingRequest.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					rec := record.Unwrap(&payloadSetOutgoingRequest.Request)
					s.Require().IsType((*record.OutgoingRequest)(nil), rec)

					ch := make(chan *wmMessage.Message, 10)

					for _, pl := range test.response {
						resMsg, err := payload.NewMessage(pl)
						s.Require().NoError(err)

						meta := payload.Meta{
							Payload: resMsg.Payload,
						}
						buf, err := meta.Marshal()
						s.Require().NoError(err)

						resMsg.Payload = buf

						ch <- resMsg
					}
					return ch, func() { close(ch) }
				},
			)

			reqInfo, err := s.amClient.RegisterOutgoingRequest(s.ctx, outgoing)

			test.check(reqInfo, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) packVirtualRecord(objRecord record.Record) []byte {
	virtual := record.Wrap(objRecord)

	byteRecord, err := virtual.Marshal()
	s.Require().NoError(err)
	return byteRecord
}

func (s *ArtifactsMangerClientSuite) packMaterialRecord(objRecord record.Record) []byte {
	material := record.Material{Virtual: record.Wrap(objRecord)}

	byteRecord, err := material.Marshal()
	s.Require().NoError(err)
	return byteRecord
}

func (s *ArtifactsMangerClientSuite) TestGetCode() {
	// Arrange
	codeID := gen.ID()
	codeRef := insolar.NewRecordReference(codeID)

	code := []byte(testutils.RandomString())
	machineType := insolar.MachineTypeBuiltin

	for name, test := range map[string]struct {
		response payload.Payload
		check    func(CodeDescriptor, error)
	}{
		"success": {
			response: &payload.Code{
				Record: s.packMaterialRecord(&record.Code{
					Request:     *insolar.NewRecordReference(codeID),
					Code:        code,
					MachineType: machineType,
				}),
			},
			check: func(descriptor CodeDescriptor, err error) {
				s.NoError(err)
				gotCode, err := descriptor.Code()
				s.NoError(err)
				s.Equal(code, gotCode)
				s.Equal(descriptor.MachineType(), machineType)
			},
		},
		"not found": {
			response: &payload.Error{
				Text: "failed to fetch record",
				Code: payload.CodeNotFound,
			},
			check: func(desc CodeDescriptor, err error) {
				s.Error(err)
				s.Contains(err.Error(), "failed to fetch record")
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(desc CodeDescriptor, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(desc CodeDescriptor, err error) {
				s.Contains(err.Error(), "unexpected reply")
			},
		},
		"unexpected message": {
			response: &payload.Code{
				Record: s.packMaterialRecord(&record.Result{}),
			},
			check: func(desc CodeDescriptor, err error) {
				s.Contains(err.Error(), "unexpected record")
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
			s.pulseAccessor.LatestMock.Return(pulseObject, nil)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					payloadGetCode := payload.GetCode{}
					err := payloadGetCode.Unmarshal(msg.Payload)
					s.Require().NoError(err)

					s.Equal(codeID, payloadGetCode.CodeID)

					ch := make(chan *wmMessage.Message, 10)

					resMsg, err := payload.NewMessage(test.response)
					s.Require().NoError(err)

					meta := payload.Meta{
						Payload: resMsg.Payload,
					}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch <- resMsg
					return ch, func() { close(ch) }
				},
			)

			reqInfo, err := s.amClient.GetCode(s.ctx, *codeRef)

			test.check(reqInfo, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestRegisterResult() {
	requestID := gen.ID()
	requestRef := insolar.NewRecordReference(requestID)
	objectID, resultID := gen.ID(), gen.ID()
	stateID := gen.ID()

	parentRef := gen.Reference()
	imageRef := gen.Reference()

	resultBytes := []byte(testutils.RandomString())
	memoryBytes := []byte(testutils.RandomString())

	for name, test := range map[string]struct {
		response      payload.Payload
		result        RequestResult
		internalCheck func(*wmMessage.Message)
		check         func(error)
	}{
		"success none": {
			response: &payload.ResultInfo{
				ObjectID: objectID,
				ResultID: resultID,
			},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectNone,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),
			},
			internalCheck: func(msg *wmMessage.Message) {
				payloadSetResult := payload.SetResult{}
				unmarshalError := payloadSetResult.Unmarshal(msg.Payload)
				s.Require().NoError(unmarshalError)

				virtualRec := &record.Virtual{}
				unmarshalError = virtualRec.Unmarshal(payloadSetResult.Result)
				s.Require().NoError(unmarshalError)

				rec := record.Unwrap(virtualRec)
				s.Require().IsType((*record.Result)(nil), rec)

				resultRecord := rec.(*record.Result)

				s.Equal(objectID, resultRecord.Object)
				s.Equal(*requestRef, resultRecord.Request)
				s.Equal(resultBytes, resultRecord.Payload)
			},
			check: func(err error) { s.NoError(err) },
		},

		"success activate": {
			response: &payload.ResultInfo{
				ObjectID: objectID,
				ResultID: resultID,
			},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectActivate,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),

				ParentReference: parentRef,
				ObjectImage:     imageRef,
				Memory:          memoryBytes,
			},
			internalCheck: func(msg *wmMessage.Message) {
				// default payload checking
				payloadActivate := payload.Activate{}
				unmarshalError := payloadActivate.Unmarshal(msg.Payload)
				s.Require().NoError(unmarshalError)

				// result record parsing
				virtualRecResult := &record.Virtual{}
				unmarshalError = virtualRecResult.Unmarshal(payloadActivate.Result)
				s.Require().NoError(unmarshalError)

				rec := record.Unwrap(virtualRecResult)
				s.Require().IsType((*record.Result)(nil), rec)

				resultRecord := rec.(*record.Result)
				s.Equal(objectID, resultRecord.Object)
				s.Equal(*requestRef, resultRecord.Request)
				s.Equal(resultBytes, resultRecord.Payload)

				// activate record parsing
				virtualRecActivate := &record.Virtual{}
				unmarshalError = virtualRecActivate.Unmarshal(payloadActivate.Record)
				s.Require().NoError(unmarshalError)

				rec = record.Unwrap(virtualRecActivate)
				s.Require().IsType((*record.Activate)(nil), rec)

				activateRecord := rec.(*record.Activate)
				s.Equal(*requestRef, activateRecord.Request)
				s.Equal(memoryBytes, activateRecord.Memory)
				s.Equal(imageRef, activateRecord.Image)
				s.Equal(false, activateRecord.IsPrototype)
				s.Equal(parentRef, activateRecord.Parent)
			},
			check: func(err error) { s.NoError(err) },
		},

		"success amend": {
			response: &payload.ResultInfo{
				ObjectID: objectID,
				ResultID: resultID,
			},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectAmend,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),

				ObjectImage:   imageRef,
				ObjectStateID: stateID,
				Memory:        memoryBytes,
			},
			internalCheck: func(msg *wmMessage.Message) {
				// default payload checking
				payloadUpdate := payload.Update{}
				unmarshalError := payloadUpdate.Unmarshal(msg.Payload)
				s.Require().NoError(unmarshalError)

				// result record parsing
				virtualRecResult := &record.Virtual{}
				unmarshalError = virtualRecResult.Unmarshal(payloadUpdate.Result)
				s.Require().NoError(unmarshalError)

				rec := record.Unwrap(virtualRecResult)
				s.Require().IsType((*record.Result)(nil), rec)

				resultRecord := rec.(*record.Result)
				s.Equal(objectID, resultRecord.Object)
				s.Equal(*requestRef, resultRecord.Request)
				s.Equal(resultBytes, resultRecord.Payload)

				// amend record parsing
				virtualRecAmend := &record.Virtual{}
				unmarshalError = virtualRecAmend.Unmarshal(payloadUpdate.Record)
				s.Require().NoError(unmarshalError)

				rec = record.Unwrap(virtualRecAmend)
				s.Require().IsType((*record.Amend)(nil), rec)

				amendRecord := rec.(*record.Amend)
				s.Equal(*requestRef, amendRecord.Request)
				s.Equal(memoryBytes, amendRecord.Memory)
				s.Equal(imageRef, amendRecord.Image)
				s.Equal(false, amendRecord.IsPrototype)
			},
			check: func(err error) { s.NoError(err) },
		},

		"success deactivate": {
			response: &payload.ResultInfo{
				ObjectID: objectID,
				ResultID: resultID,
			},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectDeactivate,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),

				ObjectStateID: stateID,
			},
			internalCheck: func(msg *wmMessage.Message) {
				// default payload checking
				payloadDeactivate := payload.Deactivate{}
				unmarshalError := payloadDeactivate.Unmarshal(msg.Payload)
				s.Require().NoError(unmarshalError)

				// result record parsing
				virtualRecResult := &record.Virtual{}
				unmarshalError = virtualRecResult.Unmarshal(payloadDeactivate.Result)
				s.Require().NoError(unmarshalError)

				rec := record.Unwrap(virtualRecResult)
				s.Require().IsType((*record.Result)(nil), rec)

				resultRecord := rec.(*record.Result)
				s.Equal(objectID, resultRecord.Object)
				s.Equal(*requestRef, resultRecord.Request)
				s.Equal(resultBytes, resultRecord.Payload)

				// amend record parsing
				virtualRecDeactivate := &record.Virtual{}
				unmarshalError = virtualRecDeactivate.Unmarshal(payloadDeactivate.Record)
				s.Require().NoError(unmarshalError)

				rec = record.Unwrap(virtualRecDeactivate)
				s.Require().IsType((*record.Deactivate)(nil), rec)

				deactivateRecord := rec.(*record.Deactivate)
				s.Equal(*requestRef, deactivateRecord.Request)
				s.Equal(stateID, deactivateRecord.PrevState)
			},
			check: func(err error) { s.NoError(err) },
		},

		"unknown payload": {
			response: &payload.PendingFinished{},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectNone,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),
			},
			check: func(err error) {
				s.Error(err)
				s.Contains(err.Error(), "unexpected reply")
			},
			internalCheck: func(message *wmMessage.Message) {},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			result: &TestRequestResult{
				SideEffectType:     RequestSideEffectNone,
				RawResult:          resultBytes,
				RawObjectReference: *insolar.NewReference(objectID),
			},
			check: func(err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
			internalCheck: func(message *wmMessage.Message) {},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
			s.pulseAccessor.LatestMock.Return(pulseObject, nil)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					test.internalCheck(msg)

					ch := make(chan *wmMessage.Message, 10)

					resMsg, err := payload.NewMessage(test.response)
					s.Require().NoError(err)

					meta := payload.Meta{Payload: resMsg.Payload}
					buf, err := meta.Marshal()
					s.Require().NoError(err)

					resMsg.Payload = buf

					ch <- resMsg
					return ch, func() { close(ch) }
				},
			)

			err := s.amClient.RegisterResult(s.ctx, *requestRef, test.result)

			test.check(err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) marshalRecord(rec record.Record) []byte {
	result, err := rec.Marshal()
	s.NoError(err)
	return result
}

func (s *ArtifactsMangerClientSuite) TestGetObject() {
	requestID := gen.ID()
	requestRef := insolar.NewRecordReference(requestID)

	objectID := gen.ID()
	stateID := gen.ID()
	objectRef := insolar.NewReference(objectID)

	parentRef := gen.Reference()
	imageRef := gen.Reference()

	codeRef := gen.Reference()

	s.amClientOriginal.localStorage.StoreObject(imageRef, NewPrototypeDescriptor(imageRef, gen.ID(), codeRef))

	memoryBytes := []byte(testutils.RandomString())

	for name, test := range map[string]struct {
		response    []payload.Payload
		checkObject func(ObjectDescriptor, error)
		checkProto  func(PrototypeDescriptor, error)
	}{
		"success_activate": {
			response: []payload.Payload{
				&payload.Index{
					Index: s.marshalRecord(&record.Lifeline{
						LatestState:         &stateID,
						StateID:             record.StateActivation,
						Parent:              parentRef,
						LatestRequest:       nil,
						EarliestOpenRequest: nil,
					}),
				},
				&payload.State{
					Record: s.packMaterialRecord(&record.Activate{
						Request:     *requestRef,
						Memory:      memoryBytes,
						Image:       imageRef,
						IsPrototype: false,
						Parent:      parentRef,
					}),
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().NoError(err)

				proto, err := desc.Prototype()
				s.NoError(err)
				s.Equal(*proto, imageRef)

				s.Equal(desc.Memory(), memoryBytes)
				s.Equal(*desc.HeadRef(), *objectRef)
				s.Equal(*desc.Parent(), parentRef)
				s.Equal(*desc.StateID(), stateID)
			},
			checkProto: func(desc PrototypeDescriptor, err error) {
				s.Require().NoError(err, "GetPrototype error")
				s.Equal(desc.Code(), &codeRef, "prototype code")
			},
		},
		"success_amend": {
			response: []payload.Payload{
				&payload.Index{
					Index: s.marshalRecord(&record.Lifeline{
						LatestState:         &stateID,
						StateID:             record.StateAmend,
						Parent:              parentRef,
						LatestRequest:       nil,
						EarliestOpenRequest: nil,
					}),
				},
				&payload.State{
					Record: s.packMaterialRecord(&record.Amend{
						Request:     *requestRef,
						Memory:      memoryBytes,
						Image:       imageRef,
						IsPrototype: false,
						PrevState:   gen.ID(),
					}),
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().NoError(err, "GetObj error")

				proto, err := desc.Prototype()
				s.NoError(err, "Prototype call error")
				s.Equal(*proto, imageRef, "proto should be image")

				s.Equal(desc.Memory(), memoryBytes, "memory")
				s.Equal(*desc.HeadRef(), *objectRef, "headref")
				s.Equal(*desc.Parent(), parentRef, "parentref")
				s.Equal(*desc.StateID(), stateID, "stateID")
			},
			checkProto: func(desc PrototypeDescriptor, err error) {
				s.Require().NoError(err, "GetPrototype error")
				s.Equal(desc.Code(), &codeRef, "prototype code")
			},
		},
		"success_deactivate_1": {
			response: []payload.Payload{
				&payload.Error{
					Code: payload.CodeDeactivated,
					Text: "object is deactivated",
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err)
				s.Equal(err, insolar.ErrDeactivated)
			},
			checkProto: func(desc PrototypeDescriptor, err error) {
			},
		},
		"success_deactivate_2": {
			response: []payload.Payload{
				&payload.Index{
					Index: s.marshalRecord(&record.Lifeline{
						LatestState:         &stateID,
						StateID:             record.StateAmend,
						Parent:              parentRef,
						LatestRequest:       nil,
						EarliestOpenRequest: nil,
					}),
				},
				&payload.Error{
					Code: payload.CodeDeactivated,
					Text: "object is deactivated",
				},
				&payload.State{
					Record: s.packMaterialRecord(&record.Amend{
						Request:     *requestRef,
						Memory:      memoryBytes,
						Image:       imageRef,
						IsPrototype: false,
						PrevState:   gen.ID(),
					}),
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err, "GetObj error")
				s.Equal(err, insolar.ErrDeactivated, "GetObj error check")
			},
		},
		"unknown error": {
			response: []payload.Payload{
				&payload.Error{
					Code: payload.CodeUnknown,
					Text: "some error",
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Nil(desc, "ObjectDescriptor should be nil")
				s.Require().Error(err, "GetObj error")
				s.Contains(err.Error(), "some error", "GetObj error check")
			},
		},
		"bad lifeline": {
			response: []payload.Payload{
				&payload.Index{Index: s.marshalRecord(&record.Deactivate{})},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err, "GetObj error")
				s.Contains(err.Error(), "failed to unmarshal index", "GetObj error check")
			},
		},
		"bad state 1": {
			response: []payload.Payload{
				&payload.State{Record: s.packMaterialRecord(&record.IncomingRequest{})},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err, "GetObj error")
				s.Contains(err.Error(), "wrong state record", "GetObj error check")
			},
		},
		"bad state 2": {
			response: []payload.Payload{
				&payload.State{Record: s.packVirtualRecord(&record.IncomingRequest{})},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err, "GetObj error")
				s.Contains(err.Error(), "wrong state record", "GetObj error check")
			},
		},
		"unexpected reply": {
			response: []payload.Payload{&payload.GetPendings{}},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Require().Error(err, "GetObj error")
				s.Contains(err.Error(), "GetObject: unexpected reply", "GetObj error check")
			},
		},
		"flow cancelled": {
			response: []payload.Payload{
				&payload.Error{
					Code: payload.CodeFlowCanceled,
					Text: "flow cancelled",
				},
			},
			checkObject: func(desc ObjectDescriptor, err error) {
				s.Error(err, "GetObj error")
				s.Contains(err.Error(), "timeout while awaiting reply from watermill", "GetObj error check")
			},
		},
	} {
		s.Run(name, func() {
			s.prepareContext()

			initialPulseNumber, pulseOffset := gen.PulseNumber(), insolar.PulseNumber(0)
			s.pulseAccessor.LatestMock.Set(
				func(ctx context.Context) (insolar.Pulse, error) {
					pulseOffset += 10
					return insolar.Pulse{PulseNumber: initialPulseNumber + pulseOffset}, nil
				},
			)

			s.busSender.SendRoleMock.Set(
				func(
					p context.Context,
					msg *wmMessage.Message,
					role insolar.DynamicRole,
					ref insolar.Reference,
				) (
					<-chan *wmMessage.Message,
					func(),
				) {
					payloadGetObject := payload.GetObject{}
					err := payloadGetObject.Unmarshal(msg.Payload)
					s.Require().NoError(err, "unmarshal error")

					s.Equal(objectID, payloadGetObject.ObjectID, "payloadGetObject")

					ch := make(chan *wmMessage.Message, 10)
					for _, pl := range test.response {
						resMsg, err := payload.NewMessage(pl)
						s.Require().NoError(err, "marshal error")

						meta := payload.Meta{
							Payload: resMsg.Payload,
						}
						buf, err := meta.Marshal()
						s.Require().NoError(err, "marshal eror")

						resMsg.Payload = buf

						ch <- resMsg
					}
					return ch, func() { close(ch) }
				},
			)

			desc, err := s.amClient.GetObject(s.ctx, *objectRef, nil)

			test.checkObject(desc, err)

			if desc != nil {
				proto, err := desc.Prototype()
				s.Require().NoError(err, "object Prototype call")

				protoDesc, err := s.amClient.GetPrototype(s.ctx, *proto)
				test.checkProto(protoDesc, err)
			}
		})
	}
}

func shouldLoadRef(strRef string) insolar.Reference {
	ref, err := insolar.NewReferenceFromString(strRef)
	if err != nil {
		panic(errors.Wrap(err, "Unexpected error, bailing out"))
	}
	return *ref
}

func (s *ArtifactsMangerClientSuite) TestLocalStorage() {
	codeDesc := NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("insolar:0AAAAyNrxlP_Iiq10drn2FuNMs2VppatXni7MP5Iy47g.record"),
	)
	var protoDesc PrototypeDescriptor
	{ // account
		pRef := shouldLoadRef("insolar:0AAAAyCjqpfzqLqOhivOFDQOK5OO_gW78OzTTniCChIU.record")
		cRef := shouldLoadRef("insolar:0AAAAyNrxlP_Iiq10drn2FuNMs2VppatXni7MP5Iy47g.record")
		protoDesc = NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		)
	}

	s.Run("objectDesc ok", func() {
		s.prepareContext()
		s.prepareAMClient()

		s.amClient.InjectPrototypeDescriptor(*protoDesc.HeadRef(), protoDesc)

		desc, err := s.amClient.GetPrototype(s.ctx, *protoDesc.HeadRef())
		s.NoError(err)
		s.Equal(desc, protoDesc)
	})

	s.Run("codeCache ok", func() {
		s.prepareContext()
		s.prepareAMClient()

		s.amClient.InjectCodeDescriptor(*codeDesc.Ref(), codeDesc)

		desc, err := s.amClient.GetCode(s.ctx, *codeDesc.Ref())
		s.NoError(err)
		s.Equal(desc, codeDesc)
	})

	s.Run("objectDesc code ref miss", func() {
		s.prepareContext()
		s.prepareAMClient()

		response := &payload.Error{
			Text: "failed to fetch record",
			Code: payload.CodeNotFound,
		}

		s.amClient.InjectCodeDescriptor(*codeDesc.Ref(), codeDesc)

		pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
		s.pulseAccessor.LatestMock.Return(pulseObject, nil)

		s.busSender.SendRoleMock.Set(
			func(
				p context.Context,
				msg *wmMessage.Message,
				role insolar.DynamicRole,
				ref insolar.Reference,
			) (
				<-chan *wmMessage.Message,
				func(),
			) {
				ch := make(chan *wmMessage.Message, 10)

				resMsg, err := payload.NewMessage(response)
				s.Require().NoError(err)

				meta := payload.Meta{
					Payload: resMsg.Payload,
				}
				buf, err := meta.Marshal()
				s.Require().NoError(err)

				resMsg.Payload = buf

				ch <- resMsg
				return ch, func() { close(ch) }
			},
		)

		_, err := s.amClient.GetObject(s.ctx, *codeDesc.Ref(), nil)
		s.Error(err)
		s.Contains(err.Error(), "failed to fetch record")
	})

	s.Run("codeDesc object ref miss", func() {
		s.prepareContext()
		s.prepareAMClient()

		response := &payload.Error{
			Text: "failed to fetch record",
			Code: payload.CodeNotFound,
		}

		s.amClient.InjectPrototypeDescriptor(*protoDesc.HeadRef(), protoDesc)

		pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
		s.pulseAccessor.LatestMock.Return(pulseObject, nil)

		s.busSender.SendRoleMock.Set(
			func(
				p context.Context,
				msg *wmMessage.Message,
				role insolar.DynamicRole,
				ref insolar.Reference,
			) (
				<-chan *wmMessage.Message,
				func(),
			) {
				ch := make(chan *wmMessage.Message, 10)

				resMsg, err := payload.NewMessage(response)
				s.Require().NoError(err)

				meta := payload.Meta{
					Payload: resMsg.Payload,
				}
				buf, err := meta.Marshal()
				s.Require().NoError(err)

				resMsg.Payload = buf

				ch <- resMsg
				return ch, func() { close(ch) }
			},
		)

		_, err := s.amClient.GetCode(s.ctx, *protoDesc.HeadRef())
		s.Error(err)
		s.Contains(err.Error(), "failed to fetch record")
	})
}
