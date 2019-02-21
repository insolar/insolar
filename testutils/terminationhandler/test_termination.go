package terminationhandler

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

type testTerminationHandler struct{}

func NewTestHandler() core.TerminationHandler {
	return &testTerminationHandler{}
}

func (t *testTerminationHandler) Abort() {
	log.Error("Node leave acknowledged by network. Goodbye!")
}

func (t testTerminationHandler) Leave(ctx context.Context, pulseDelta core.PulseNumber) chan core.LeaveApproved {
	panic("implement me")
}

func (t testTerminationHandler) OnLeaveApproved() {
	panic("implement me")
}
