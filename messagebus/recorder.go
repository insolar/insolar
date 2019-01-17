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
	tape         tape
	scheme       core.PlatformCryptographyScheme
	pulseStorage core.PulseStorage
}

// newRecorder create new recorder instance.
func newRecorder(s sender, tape tape, scheme core.PlatformCryptographyScheme, pulseStorage core.PulseStorage) *recorder {
	return &recorder{
		sender:       s,
		tape:         tape,
		scheme:       scheme,
		pulseStorage: pulseStorage,
	}
}

// WriteTape writes recorder's tape to the provided writer.
func (r *recorder) WriteTape(ctx context.Context, w io.Writer) error {
	return r.tape.Write(ctx, w)
}

// Send wraps MessageBus Send to save received replies to the tape. This reply is also used to return directly from the
// tape is the message is sent again, thus providing a cash for message replies.
func (r *recorder) Send(ctx context.Context, msg core.Message, ops *core.MessageSendOptions) (core.Reply, error) {
	currentPulse, err := r.pulseStorage.Current(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := r.CreateParcel(ctx, msg, ops.Safe().Token, *currentPulse)
	if err != nil {
		return nil, err
	}

	// Actually send message.
	rep, sendErr := r.SendParcel(ctx, parcel, *currentPulse, ops)

	// Save the received Value on the tape.
	id := GetMessageHash(r.scheme, parcel)
	err = r.tape.Set(ctx, id, rep, sendErr)
	if err != nil {
		return nil, err
	}

	if sendErr != nil {
		return nil, sendErr
	}
	return rep, nil
}

func (r *recorder) OnPulse(context.Context, core.Pulse) error {
	panic("This method must not be called")
}
