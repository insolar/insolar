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
)

// Recorder is a MessageBus wrapper that stores received replies to the tape. The tape then can be transferred and
// used by Player to replay those replies.
type recorder struct {
	sender
	tape   tape
	scheme core.PlatformCryptographyScheme
}

// newRecorder create new recorder instance.
func newRecorder(s sender, tape tape, scheme core.PlatformCryptographyScheme) *recorder {
	return &recorder{sender: s, tape: tape, scheme: scheme}
}

// WriteTape writes recorder's tape to the provided writer.
func (r *recorder) WriteTape(ctx context.Context, w io.Writer) error {
	return r.tape.Write(ctx, w)
}

// Send wraps MessageBus Send to save received replies to the tape. This reply is also used to return directly from the
// tape is the message is sent again, thus providing a cash for message replies.
func (r *recorder) Send(ctx context.Context, msg core.Message, currentPulse core.Pulse, ops *core.MessageSendOptions) (core.Reply, error) {
	var (
		rep core.Reply
		err error
	)

	parcel, err := r.CreateParcel(ctx, msg, ops.Safe().Token, currentPulse)
	if err != nil {
		return nil, err
	}

	// Actually send message.
	rep, err = r.SendParcel(ctx, parcel, currentPulse, ops)
	if err != nil {
		return nil, err
	}

	// Save the received Value on the tape.
	id := GetMessageHash(r.scheme, parcel)
	err = r.tape.SetReply(ctx, id, rep)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
