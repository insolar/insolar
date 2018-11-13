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
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/localstorage"
)

// Player is a MessageBus wrapper that replays replies from provided tape. The tape can be created and by Recorder
// and transferred to player.
type player struct {
	sender
	tape   tape
	pm     core.PulseManager
	scheme core.PlatformCryptographyScheme
}

// newPlayer creates player instance. It will replay replies from provided tape.
func newPlayer(s sender, tape tape, pm core.PulseManager, scheme core.PlatformCryptographyScheme) *player {
	return &player{sender: s, tape: tape, pm: pm, scheme: scheme}
}

// WriteTape for player is not available.
func (r *player) WriteTape(ctx context.Context, w io.Writer) error {
	panic("can't write the tape from player")
}

// Send wraps MessageBus Send to reply replies from the tape. If reply for this message is not on the tape, an error
// will be returned.
func (r *player) Send(ctx context.Context, msg core.Message) (core.Reply, error) {
	var (
		rep core.Reply
		err error
	)

	parcel, err := r.CreateParcel(ctx, msg)
	id := GetMessageHash(r.scheme, parcel)

	// Value from storageTape.
	rep, err = r.tape.GetReply(ctx, id)
	if err == nil {
		return rep, nil
	}
	if err == localstorage.ErrNotFound {
		return nil, ErrNoReply
	} else {
		return nil, err
	}

	return rep, nil
}
