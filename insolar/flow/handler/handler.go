//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const handleTimeout = 10 * time.Second

type Handler struct {
	handles struct {
		present flow.MakeHandle
	}
	controller *thread.Controller
}

func NewHandler(present flow.MakeHandle) *Handler {
	h := &Handler{
		controller: thread.NewController(),
	}
	h.handles.present = present
	return h
}

// ChangePulse is a handle for pulse change vent.
func (h *Handler) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	h.controller.Pulse()
}

func (h *Handler) WrapBusHandle(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply),
		Parcel:  parcel,
	}
	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", parcel.Pulse()))
	ctx = pulse.ContextWith(ctx, parcel.Pulse())
	go func() {
		f := thread.NewThread(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		if err != nil {
			select {
			case msg.ReplyTo <- bus.Reply{Err: err}:
			default:
			}
			logger.Error("Handling failed", err)
		}
	}()
	rep := <-msg.ReplyTo
	return rep.Reply, rep.Err
}
