package terminationhandler

import (
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
