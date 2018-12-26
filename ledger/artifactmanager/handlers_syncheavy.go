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
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// TODO: check sender if it was light material in synced pulses:
// sender := genericMsg.GetSender()
// sender.isItWasLMInPulse(pulsenum)
func (h *MessageHandler) handleHeavyPayload(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyPayload)

	inslog := inslogger.FromContext(ctx).WithField("pulseNum", msg.PulseNum)
	inslog = inslog.WithField("jetID", msg.JetID)
	inslog.Debugf("Heavy sync: get payload message with %v records", len(msg.Records))

	if err := h.HeavySync.Store(ctx, msg.JetID, msg.PulseNum, msg.Records); err != nil {
		inslog.Error("Heavy store failed", err)
		return heavyerrreply(err), err
	}
	inslog.Debugf("Heavy sync: stores %v records", len(msg.Records))
	return &reply.OK{}, nil
}

func (h *MessageHandler) handleHeavyStartStop(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyStartStop)

	inslog := inslogger.FromContext(ctx).WithField("pulseNum", msg.PulseNum)
	inslog = inslog.WithField("jetID", fmt.Sprintf("%+v", msg.JetID))

	// stop
	if msg.Finished {
		inslog.Debug("Heavy sync: get stop message")
		if err := h.HeavySync.Stop(ctx, msg.JetID, msg.PulseNum); err != nil {
			return nil, err
		}
		return &reply.OK{}, nil
	}
	// start
	inslog.Debug("Heavy sync: get start message")
	if err := h.HeavySync.Start(ctx, msg.JetID, msg.PulseNum); err != nil {
		return heavyerrreply(err), err
	}
	return &reply.OK{}, nil
}

func (h *MessageHandler) handleHeavyReset(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyReset)

	inslog := inslogger.FromContext(ctx).WithField("pulseNum", msg.PulseNum)
	inslog = inslog.WithField("jetID", fmt.Sprintf("%+v", msg.JetID))

	inslog.Debug("Heavy sync: get reset message")
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
