package adapterhelper

import (
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/messagebus"
	"github.com/stretchr/testify/require"
)

type mockReply struct {
	data string
}

func (mr *mockReply) Type() insolar.ReplyType {
	return 0
}

func TestSendResponseHelper(t *testing.T) {
	f := messagebus.NewFuture()
	event := insolar.ConveyorPendingMessage{Future: f}
	testReply := &mockReply{data: "Put-in"}

	slotElementHelperMock := fsm.NewSlotElementHelperMock(t)
	slotElementHelperMock.GetInputEventFunc = func() (r interface{}) {
		return event
	}
	slotElementHelperMock.SendTaskFunc = func(p adapterid.ID, response interface{}, p2 uint32) (r error) {
		f := response.(adapter.SendResponseTask).Future
		f.SetResult(testReply)
		return nil
	}

	adapterCatalog := NewCatalog()
	err := adapterCatalog.SendResponseHelper.SendResponse(slotElementHelperMock, testReply, 42)
	require.NoError(t, err)

	gotReply, err := f.GetResult(time.Second)
	require.NoError(t, err)
	require.Equal(t, testReply, gotReply)
}

func TestSendResponseHelper_BadInput(t *testing.T) {
	slotElementHelperMock := fsm.NewSlotElementHelperMock(t)
	slotElementHelperMock.GetInputEventFunc = func() (r interface{}) {
		return 33
	}
	adapterCatalog := NewCatalog()
	err := adapterCatalog.SendResponseHelper.SendResponse(slotElementHelperMock, &mockReply{}, 44)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Input event is not insolar.ConveyorPendingMessage")
}
