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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
)

// TODO: check sender if it was light material in synced pulses:
// sender := genericMsg.GetSender()
// sender.isItWasLMInPulse(pulsenum)
func (h *MessageHandler) handleHeavyPayload(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyPayload)

	if err := h.HeavySync.Store(ctx, insolar.ID(msg.JetID), msg.PulseNum, msg.Records); err != nil {
		return heavyerrreply(err)
	}
	if err := h.HeavySync.StoreDrop(ctx, msg.JetID, msg.Drop); err != nil {
		return heavyerrreply(err)
	}
	if err := h.HeavySync.StoreBlobs(ctx, msg.PulseNum, msg.Blobs); err != nil {
		return heavyerrreply(err)
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) handleHeavyStartStop(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyStartStop)

	// stop
	if msg.Finished {
		if err := h.HeavySync.Stop(ctx, insolar.ID(msg.JetID), msg.PulseNum); err != nil {
			return nil, err
		}
		return &reply.OK{}, nil
	}
	// start
	if err := h.HeavySync.Start(ctx, insolar.ID(msg.JetID), msg.PulseNum); err != nil {
		return heavyerrreply(err)
	}
	return &reply.OK{}, nil
}

func heavyerrreply(err error) (insolar.Reply, error) {
	if herr, ok := err.(*reply.HeavyError); ok {
		return herr, nil
	}
	return nil, err
}
