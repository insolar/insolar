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

package messagebus

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/localstorage"
)

// Player is a MessageBus wrapper that replays replies from provided tape. The tape can be created and by Recorder
// and transferred to player.
type player struct {
	sender
	tape         tape
	scheme       core.PlatformCryptographyScheme
	pulseStorage core.PulseStorage
}

// newPlayer creates player instance. It will replay replies from provided tape.
func newPlayer(s sender, tape tape, scheme core.PlatformCryptographyScheme, pulseStorage core.PulseStorage) *player {
	return &player{
		sender:       s,
		tape:         tape,
		scheme:       scheme,
		pulseStorage: pulseStorage,
	}
}

// Send wraps MessageBus Send to reply replies from the tape. If reply for this message is not on the tape, an error
// will be returned.
func (p *player) Send(ctx context.Context, msg core.Message, ops *core.MessageSendOptions) (core.Reply, error) {
	currentPulse, err := p.pulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := p.CreateParcel(ctx, msg, ops.Safe().Token, *currentPulse)
	if err != nil {
		return nil, err
	}
	id := GetMessageHash(p.scheme, parcel)

	item, err := p.tape.Get(ctx, id)
	if err != nil {
		if err == localstorage.ErrNotFound {
			return nil, ErrNoReply
		}
		return nil, err
	}

	return item.Reply, item.Error
}

func (p *player) OnPulse(context.Context, core.Pulse) error {
	panic("This method must not be called")
}
