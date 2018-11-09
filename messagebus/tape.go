package messagebus

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"

	"github.com/satori/go.uuid"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
)

type tape interface {
	Write(ctx context.Context, writer io.Writer) error
	GetReply(ctx context.Context, msgHash []byte) (core.Reply, error)
	SetReply(ctx context.Context, msgHash []byte, rep core.Reply) error
}

// storagetape saves and fetches message replies to/from local storage.
//
// It uses <storagetape id> + <message hash> for Value keys.
type storagetape struct {
	ls    core.LocalStorage
	pulse core.PulseNumber
	id    uuid.UUID
}

type couple struct {
	Key   []byte
	Value []byte
}

// NewTape creates new storagetape with random id.
func NewTape(ls core.LocalStorage, pulse core.PulseNumber) (*storagetape, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &storagetape{ls: ls, pulse: pulse, id: id}, nil
}

// NewTapeFromReader creates and fills a new storagetape from a stream.
//
// This is a very long operation, as it saves replies in storage until the stream is exhausted.
func NewTapeFromReader(ctx context.Context, ls core.LocalStorage, reader io.Reader) (*storagetape, error) {
	var err error
	tape := storagetape{ls: ls}

	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&tape.pulse)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&tape.id)
	if err != nil {
		return nil, err
	}
	for {
		var rep couple
		err = decoder.Decode(&rep)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		err = tape.setReplyBinary(ctx, rep.Key, rep.Value)
		if err != nil {
			return nil, err
		}
	}

	return &tape, nil
}
