package terminationhandler

import (
	"fmt"

	"github.com/insolar/insolar/core"
)

type testTerminationHandler struct{}

func NewTestHandler() core.TerminationHandler {
	return &testTerminationHandler{}
}

func (t *testTerminationHandler) Abort(reason string) {
	panic(fmt.Sprintf("Node leave acknowledged by network. Goodbye! Reason: %s", reason))
}
