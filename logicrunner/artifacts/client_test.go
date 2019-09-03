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

package artifacts

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
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
		ReturnMode:      record.ReturnNoWait,
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

func (s *ArtifactsMangerClientSuite) TestGetAbandonedRequest() {
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
			recordRequest, err := s.amClient.GetAbandonedRequest(s.ctx, *objectRef, requestRef)
			// Check
			test.check(recordRequest, err)
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestGetAbandonedRequest_FailedToSend() {
	// Arrange
	request, requestRef := genIncomingRequest()

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.GetAbandonedRequest(s.ctx, *request.Object, requestRef)

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
			references, err := s.amClient.GetPendings(s.ctx, objectRef)

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
	_, err := s.amClient.GetPendings(s.ctx, *request.Object)

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
	machineType := insolar.MachineTypeGoPlugin
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

			emptyRef := *insolar.NewEmptyReference()

			// Act
			deployCodeID, err := s.amClient.DeployCode(s.ctx, emptyRef, emptyRef, code, machineType)

			// Assert
			test.check(deployCodeID, err)

			s.mc.Finish()
		})
	}
}

func (s *ArtifactsMangerClientSuite) TestDeployCode_FailedToSend() {
	// Arrange
	emptyRef := *insolar.NewEmptyReference()
	code := []byte(testutils.RandomString())
	machineType := insolar.MachineTypeGoPlugin

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	s.pulseAccessor.LatestMock.Return(pulseObject, nil)

	ch := make(chan *wmMessage.Message)
	close(ch)
	s.busSender.SendRoleMock.Return(ch, func() {})

	// Act
	_, err := s.amClient.DeployCode(s.ctx, emptyRef, emptyRef, code, machineType)

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

	_ = &payload.RequestInfo{
		ObjectID:  *objectRef.GetLocal(),
		RequestID: *requestRef.GetLocal(),
		Request:   nil,
		Result:    nil,
	}

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

func (s *ArtifactsMangerClientSuite) TestGetCode() {

}

func (s *ArtifactsMangerClientSuite) TestGetObject() {

}

func (s *ArtifactsMangerClientSuite) TestActivatePrototype() {

}

func (s *ArtifactsMangerClientSuite) TestRegisterResult() {

}
