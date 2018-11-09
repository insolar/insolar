package messagebus

import (
	"context"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/localstorage"
)

// Recorder is a MessageBus wrapper that stores received replies to the tape. The tape then can be transferred and
// used by Player to replay those replies.
type recorder struct {
	sender
	tape tape
	pm   core.PulseManager
}

// NewRecorder create new recorder instance.
func NewRecorder(s sender, tape tape, pm core.PulseManager) (*recorder, error) {
	return &recorder{sender: s, tape: tape, pm: pm}, nil
}

// WriteTape writes recorder's tape to the provided writer.
func (r *recorder) WriteTape(ctx context.Context, w io.Writer) error {
	return r.tape.Write(ctx, w)
}

// Send wraps MessageBus Send to save received replies to the tape. This reply is also used to return directly from the
// tape is the message is sent again, thus providing a cash for message replies.
func (r *recorder) Send(ctx context.Context, msg core.Message) (core.Reply, error) {
	var (
		rep core.Reply
		err error
	)
	pulse, err := r.pm.Current(ctx)
	if err != nil {
		return nil, err
	}
	signedMessage, err := r.CreateSignedMessage(ctx, pulse.PulseNumber, msg)
	id := GetMessageHash(signedMessage)

	// Check if Value for this message is already stored.
	rep, err = r.tape.GetReply(ctx, id)
	if err == nil {
		return rep, nil
	}
	if err != localstorage.ErrNotFound {
		return nil, err
	}

	// Actually send message.
	rep, err = r.SendMessage(ctx, pulse, signedMessage)
	if err != nil {
		return nil, err
	}

	// Save the received Value on the storageTape.
	err = r.tape.SetReply(ctx, id, rep)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
