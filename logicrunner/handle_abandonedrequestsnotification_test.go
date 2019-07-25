package logicrunner

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
	//	dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
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
	//	dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
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
	//	dep:    &Dependencies{lr: suite.lr, Sender: suite.sender},
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

