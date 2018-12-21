/*
 *    Copyright 2018 Insolar
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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/hack"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// TODO: check sender if it was light material in synced pulses:
// sender := genericMsg.GetSender()
// sender.isItWasLMInPulse(pulsenum)
func (h *MessageHandler) handleHeavyPayload(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)
	if hack.SkipValidation(ctx) {
		return &reply.OK{}, nil
	}
	msg := genericMsg.Message().(*message.HeavyPayload)
	inslog.Debugf("Heavy sync: get start payload message with %v records", len(msg.Records))
	if err := h.HeavySync.Store(ctx, msg.JetID, msg.PulseNum, msg.Records); err != nil {
		return heavyerrreply(err), err
	}
	return &reply.OK{}, nil
}

func (h *MessageHandler) handleHeavyStartStop(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)
	if hack.SkipValidation(ctx) {
		return &reply.OK{}, nil
	}

	msg := genericMsg.Message().(*message.HeavyStartStop)
	// stop branch
	if msg.Finished {
		inslog.Debugf("Heavy sync: get stop message for pulse %v", msg.PulseNum)
		if err := h.HeavySync.Stop(ctx, msg.JetID, msg.PulseNum); err != nil {
			return nil, err
		}
		return &reply.OK{}, nil
	}
	// start
	inslog.Debugf("Heavy sync: get start message for pulse %v", msg.PulseNum)
	if err := h.HeavySync.Start(ctx, msg.JetID, msg.PulseNum); err != nil {
		return heavyerrreply(err), err
	}
	return &reply.OK{}, nil
}

func (h *MessageHandler) handleHeavyReset(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)
	if hack.SkipValidation(ctx) {
		return &reply.OK{}, nil
	}

	msg := genericMsg.Message().(*message.HeavyReset)
	inslog.Debugf("Heavy sync: get reset message for pulse %v", msg.PulseNum)
	if err := h.HeavySync.Reset(ctx, msg.JetID, msg.PulseNum); err != nil {
		return heavyerrreply(err), err
	}
	return &reply.OK{}, nil
}

func heavyerrreply(err error) core.Reply {
	if herr, ok := err.(*reply.HeavyError); ok {
		return herr
	}
	return nil
}
