/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"context"
	"sync"

	"github.com/insolar/insolar/network/servicenetwork"
)

type leaveApproved struct{}

// TerminationHandler handles such node events as graceful stop, abort, etc.
type TerminationHandler interface {
	Leave(context.Context, PulseNumber) chan leaveApproved
	OnLeaveApproved()
	// Abort forces to stop all node components
	Abort()
}

type terminationHandler struct {
	sync.Mutex
	Network     servicenetwork.ServiceNetwork `inject:""`
	done        chan leaveApproved
	terminating bool
}

func NewTerminationHandler() TerminationHandler {
	return &terminationHandler{}
}

func (t terminationHandler) Leave(ctx context.Context, pulseDelta PulseNumber) chan leaveApproved {
	t.Lock()
	defer t.Unlock()

	if !t.terminating {
		t.done = make(chan leaveApproved, 1)
	}

	if pulseDelta == 0 || !t.terminating {
		t.terminating = true
		t.Network.Leave(ctx, pulseDelta)
	}

	return t.done
}

// TODO what if come here few times and second time we try to close closing chanel?
func (t terminationHandler) OnLeaveApproved() {
	t.Lock()
	defer t.Unlock()
	close(t.done)
}

// ci said that log.Fatal causes import cycle
func (t terminationHandler) Abort() {
	panic("Node leave acknowledged by network. Goodbye!")
}
