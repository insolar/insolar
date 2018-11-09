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
	tape tape
	pm   core.PulseManager
}

// NewPlayer creates player instance. It will replay replies from provided tape.
func NewPlayer(s sender, tape tape, pm core.PulseManager) (*player, error) {
	return &player{sender: s, tape: tape, pm: pm}, nil
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
	pulse, err := r.pm.Current(ctx)
	if err != nil {
		return nil, err
	}
	signedMessage, err := r.CreateSignedMessage(ctx, pulse.PulseNumber, msg)
	id := GetMessageHash(signedMessage)

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
