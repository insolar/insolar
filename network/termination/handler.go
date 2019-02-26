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

	Network *servicenetwork.ServiceNetwork
}

func NewHandler(nw *servicenetwork.ServiceNetwork) core.TerminationHandler {
	return &terminationHandler{Network: nw}
}

// TODO take ETA by role of node
func (t *terminationHandler) Leave(ctx context.Context) chan core.LeaveApproved {
	t.Lock()
	defer t.Unlock()

	if !t.terminating {
		t.terminating = true
		t.done = make(chan core.LeaveApproved, 1)

		pulse, _ := t.Network.PulseStorage.Current(ctx)
		pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber
		t.Network.Leave(ctx, pulse.PulseNumber+10*pulseDelta)
	}

	return t.done
}

func (t *terminationHandler) OnLeaveApproved() {
	t.Lock()
	defer t.Unlock()
	if t.terminating {
		close(t.done)
	}
}

// ci said that log.Fatal causes import cycle
func (t *terminationHandler) Abort() {
	panic("Node leave acknowledged by network. Goodbye!")
}
