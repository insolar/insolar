//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package termination

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar"
)

type terminationHandler struct {
	sync.Mutex
	done        chan insolar.LeaveApproved
	terminating bool

	Leaver        insolar.Leaver `inject:""`
	PulseAccessor pulse.Accessor `inject:""`
}

func NewHandler(l insolar.Leaver) insolar.TerminationHandler {
	return &terminationHandler{Leaver: l}
}

// TODO take ETA by role of node
func (t *terminationHandler) Leave(ctx context.Context, leaveAfterPulses insolar.PulseNumber) {
	doneChan := t.leave(ctx, leaveAfterPulses)
	<-doneChan
}

func (t *terminationHandler) leave(ctx context.Context, leaveAfterPulses insolar.PulseNumber) chan insolar.LeaveApproved {
	t.Lock()
	defer t.Unlock()

	if !t.terminating {
		t.terminating = true
		t.done = make(chan insolar.LeaveApproved, 1)

		if leaveAfterPulses == 0 {
			inslogger.FromContext(ctx).Debug("terminationHandler.Leave() with 0")
			t.Leaver.Leave(ctx, 0)
		} else {
			pulse, err := t.PulseAccessor.Latest(ctx)
			if err != nil {
				inslogger.FromContext(ctx).Panicf("smth goes wrong. There is no pulse in the storage. err - %v", err)
			}
			pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber

			inslogger.FromContext(ctx).Debugf("terminationHandler.Leave() with leaveAfterPulses: %+v, in pulse %+v", leaveAfterPulses, pulse.PulseNumber+leaveAfterPulses*pulseDelta)
			t.Leaver.Leave(ctx, pulse.PulseNumber+leaveAfterPulses*pulseDelta)
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

func (t *terminationHandler) Abort(reason string) {
	panic(reason)
}

func (t *terminationHandler) Terminating() bool {
	return t.terminating
}
