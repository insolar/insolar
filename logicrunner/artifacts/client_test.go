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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
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

	busSender      *bus.SenderMock
	jetStorage     *jet.StorageMock
	pulseAccessor  *pulse.AccessorMock
	jetCoordinator *jet.CoordinatorMock

	amClientOriginal *client
	amClient         Client
}

// Init and run suite
func TestArtifactManager(t *testing.T) {
	suite.Run(t, &ArtifactsMangerClientSuite{Suite: suite.Suite{}})
}

func (s *ArtifactsMangerClientSuite) createAMClient() *client {
	s.mc = minimock.NewController(s.T())

	s.jetStorage = jet.NewStorageMock(s.mc)
	s.pulseAccessor = pulse.NewAccessorMock(s.mc)
	s.jetCoordinator = jet.NewCoordinatorMock(s.mc)

	s.busSender = bus.NewSenderMock(s.mc)

	return &client{
		JetStorage:     s.jetStorage,
		PCS:            platformpolicy.NewPlatformCryptographyScheme(),
		PulseAccessor:  s.pulseAccessor,
		JetCoordinator: s.jetCoordinator,

		sender:       s.busSender,
		localStorage: newLocalStorage(),
	}
}

func (s *ArtifactsMangerClientSuite) BeforeTest(suiteName, testName string) {
	s.ctx = inslogger.TestContext(s.T())

	s.amClientOriginal = s.createAMClient()
	s.amClient = s.amClientOriginal
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
				s.Assert().NoError(err)
				s.Equal(*iRequest, *gotRequest.(*record.IncomingRequest))
			},
		},
		"outgoing success": {
			response: &payload.Request{
				RequestID: *requestRef.GetLocal(),
				Request:   record.Wrap(oRequest),
			},
			check: func(gotRequest record.Request, err error) {
				s.Assert().NoError(err)
				s.Equal(*oRequest, *gotRequest.(*record.OutgoingRequest))
			},
		},
		"not found": {
			response: &payload.Error{
				Text: "request not found",
				Code: payload.CodeNotFound,
			},
			check: func(gotRequest record.Request, err error) {
				s.Equal(insolar.ErrNotFound, err)
			},
		},
		"other error": {
			response: &payload.Error{
				Text: "some error",
				Code: payload.CodeUnknown,
			},
			check: func(gotRequest record.Request, err error) {
				s.Error(err)
				s.Contains(err.Error(), "some error")
			},
		},
		"unknown payload": {
			response: &payload.PendingFinished{},
			check: func(gotRequest record.Request, err error) {
				s.Assert().Contains(err.Error(), "unexpected reply")
			},
		},
		"unexpected message": {
			response: &payload.Request{
				RequestID: *requestRef.GetLocal(),
				Request:   record.Wrap(&record.Activate{}),
			},
			check: func(gotRequest record.Request, err error) {
				s.Assert().Contains(err.Error(), "unexpected message")
			},
		},
	} {
		s.Run(name, func() {
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

func (s *ArtifactsMangerClientSuite) TestGetPendings_Success() {
	// Arrange
	request, requestRef := genIncomingRequest()

	resultIDs := &payload.IDs{
		IDs: []insolar.ID{*requestRef.GetLocal()},
	}
	resMsg, err := payload.NewMessage(resultIDs)
	require.NoError(s.T(), err)

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

			s.Equal(*request.Object.GetLocal(), getPendings.ObjectID)

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
	res, err := s.amClient.GetPendings(s.ctx, *request.Object)

	// Assert
	s.NoError(err)
	s.Equal([]insolar.Reference{requestRef}, res)
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
	s.Contains(err.Error(), "failed to send GetPendings")
}

func (s *ArtifactsMangerClientSuite) TestHasPendings_Success() {
	// Arrange
	objectRef := gen.Reference()

	resultHas := &payload.PendingsInfo{
		HasPendings: true,
	}
	resMsg, err := payload.NewMessage(resultHas)
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
	res, err := s.amClient.HasPendings(s.ctx, objectRef)

	// Assert
	s.NoError(err)
	s.Equal(true, res)
}
