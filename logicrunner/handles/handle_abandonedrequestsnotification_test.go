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

package handles

import (
	"testing"
)

func TestHandleAbandonedRequestsNotification_Present(t *testing.T) {
	t.Skip("we disabled handling of this notification for now")

	//objectId := testutils.RandomID()
	//objectRef := *insolar.NewReference(objectId)
	//msg := &message.AbandonedRequestsNotification{Object: objectId}
	//parcel := &message.Parcel{Msg: msg}
	//
	//flowMock := flow.NewFlowMock(suite.mc)
	//flowMock.ProcedureMock.Set(func(p context.Context, p1 flow.Procedure, p2 bool) (r error) {
	//	return p1.Proceed(p)
	//})
	//
	//replyChan := mockSender(suite)
	//
	//h := HandleAbandonedRequestsNotification{
	//	Dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
	//	Parcel: parcel,
	//}
	//
	//err := h.Present(suite.ctx, flowMock)
	//suite.Require().NoError(err)
	//
	//_, err = getReply(suite, replyChan)
	//suite.Require().NoError(err)
	//broker := suite.lr.StateStorage.GetExecutionState(objectRef)
	//suite.Equal(true, broker.ledgerHasMoreRequests)
	//_ = suite.lr.Stop(suite.ctx)
	//
	//// LedgerHasMoreRequests false
	//suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{}, suite.pub, suite.sender)
	//
	//broker = suite.lr.StateStorage.UpsertExecutionState(objectRef)
	//broker.ledgerHasMoreRequests = false
	//
	//h = HandleAbandonedRequestsNotification{
	//	Dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
	//	Parcel: parcel,
	//}
	//
	//err = h.Present(suite.ctx, flowMock)
	//suite.Require().NoError(err)
	//_, err = getReply(suite, replyChan)
	//suite.Require().NoError(err)
	//broker = suite.lr.StateStorage.GetExecutionState(objectRef)
	//suite.Equal(true, broker.ledgerHasMoreRequests)
	//_ = suite.lr.Stop(suite.ctx)
	//
	//// LedgerHasMoreRequests already true
	//suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{}, suite.pub, suite.sender)
	//
	//broker = suite.lr.StateStorage.UpsertExecutionState(objectRef)
	//broker.ledgerHasMoreRequests = true
	//
	//h = HandleAbandonedRequestsNotification{
	//	Dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
	//	Parcel: parcel,
	//}
	//
	//err = h.Present(suite.ctx, flowMock)
	//suite.Require().NoError(err)
	//_, err = getReply(suite, replyChan)
	//suite.Require().NoError(err)
	//suite.Require().NoError(err)
	//broker = suite.lr.StateStorage.GetExecutionState(objectRef)
	//suite.Equal(true, broker.ledgerHasMoreRequests)
	//_ = suite.lr.Stop(suite.ctx)
}
