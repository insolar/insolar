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
	"bytes"
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
)

// Tape is an abstraction for saving replies for messages and restoring them.
//
// There can be many active tapes simultaneously and they do not share saved replies.
//go:generate minimock -i github.com/insolar/insolar/messagebus.tape -o .
type tape interface {
	Write(ctx context.Context, writer io.Writer) error
	GetReply(ctx context.Context, msgHash []byte) (core.Reply, error)
	SetReply(ctx context.Context, msgHash []byte, rep core.Reply) error
}

// memoryTape saves and fetches message replies to/from memory array.
//
// It uses <storageTape id> + <message hash> for Value keys.
type memoryTape struct {
	pulse   core.PulseNumber
	storage []memoryTapeMessage
}

type memoryTapeMessage struct {
	msgHash []byte
	reply   core.Reply
}

func newMemoryTape(pulse core.PulseNumber) *memoryTape {
	return &memoryTape{
		pulse: pulse,
	}
}

func newMemoryTapeFromReader(ctx context.Context, r io.Reader) (*memoryTape, error) {
	t := memoryTape{}
	ch := new(codec.CborHandle)
	decoder := codec.NewDecoder(r, ch)
	err := decoder.Decode(&t.pulse)
	if err != nil {
		return nil, errors.Wrap(err, "[ MemoryTape ] can't read pulse")
	}
	err = decoder.Decode(&t.storage)
	if err != nil {
		return nil, errors.Wrap(err, "[ MemoryTape ] can't read storage")
	}
	return &t, nil
}

func (t *memoryTape) Write(ctx context.Context, w io.Writer) error {
	encoder := codec.NewEncoder(w, new(codec.CborHandle))

	err := encoder.Encode(t.pulse)
	if err != nil {
		return errors.Wrap(err, "[ MemoryTape ] can't write pulse")
	}

	err = encoder.Encode(t.storage)
	if err != nil {
		return errors.Wrap(err, "[ MemoryTape ] can't write storage")
	}
	return nil
}

func (t *memoryTape) GetReply(ctx context.Context, msgHash []byte) (core.Reply, error) {
	if len(t.storage) == 0 {
		return nil, errors.New("Validation error. Message is not expected")
	}
	if !bytes.Equal(msgHash, t.storage[0].msgHash) {
		return nil, errors.New("Validation error. Message mismatch")
	}
	ret := t.storage[0]
	t.storage = t.storage[1:]
	return ret.reply, nil
}

func (t *memoryTape) SetReply(ctx context.Context, msgHash []byte, rep core.Reply) error {
	t.storage = append(t.storage, memoryTapeMessage{
		msgHash: msgHash,
		reply:   rep,
	})
	return nil
}
