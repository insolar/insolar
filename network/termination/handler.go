package termination

import (
	"context"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/core"
)

type terminationHandler struct {
	sync.Mutex
	done        chan core.LeaveApproved
	terminating bool

	Network      core.Network      `inject:""`
	PulseStorage core.PulseStorage `inject:""`
}

func NewHandler(nw core.Network) core.TerminationHandler {
	return &terminationHandler{Network: nw}
}

// TODO take ETA by role of node
func (t *terminationHandler) Leave(ctx context.Context, leaveAfterPulses core.PulseNumber) chan core.LeaveApproved {
	t.Lock()
	defer t.Unlock()

	if !t.terminating {
		t.terminating = true
		t.done = make(chan core.LeaveApproved, 1)

		if leaveAfterPulses == 0 {
			inslogger.FromContext(ctx).Debug("terminationHandler.Leave() with 0")
			t.Network.Leave(ctx, 0)
		} else {
			pulse, _ := t.PulseStorage.Current(ctx)
			pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber

			inslogger.FromContext(ctx).Debugf("terminationHandler.Leave() with leaveAfterPulses: %+v, in pulse %+v", leaveAfterPulses, pulse.PulseNumber+leaveAfterPulses*pulseDelta)
			t.Network.Leave(ctx, pulse.PulseNumber+leaveAfterPulses*pulseDelta)
		}
	}

	return t.done
}

func (t *terminationHandler) OnLeaveApproved(ctx context.Context) {
	t.Lock()
	defer t.Unlock()
	if t.terminating {
		inslogger.FromContext(ctx).Debug("terminationHandler.OnLeaveApproved() received")
		t.terminating = false
		close(t.done)
	}
}

// ci said that log.Fatal causes import cycle
func (t *terminationHandler) Abort() {
	panic("Node leave acknowledged by network. Goodbye!")
}
