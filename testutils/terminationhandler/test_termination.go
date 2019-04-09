package terminationhandler

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

type testTerminationHandler struct{}

func NewTestHandler() core.TerminationHandler {
	return &testTerminationHandler{}
}

func (t *testTerminationHandler) Abort(reason string) {
	log.Error(reason)
}
