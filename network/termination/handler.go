package termination

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/servicenetwork"
)

type terminationHandler struct {
	sync.Mutex
	done        chan core.LeaveApproved
	terminating bool

	Network servicenetwork.ServiceNetwork `inject:""`
}

func NewHandler() core.TerminationHandler {
	return &terminationHandler{}
}

func (t *terminationHandler) Leave(ctx context.Context, pulseDelta core.PulseNumber) chan core.LeaveApproved {
	t.Lock()
	defer t.Unlock()

	if !t.terminating {
		t.done = make(chan core.LeaveApproved, 1)
	}

	if pulseDelta == 0 || !t.terminating {
		t.terminating = true
		t.Network.Leave(ctx, pulseDelta)
	}

	return t.done
}

// TODO what if come here few times and second time we try to close closing chanel?
func (t *terminationHandler) OnLeaveApproved() {
	t.Lock()
	defer t.Unlock()
	close(t.done)
}

// ci said that log.Fatal causes import cycle
func (t *terminationHandler) Abort() {
	panic("Node leave acknowledged by network. Goodbye!")
}
