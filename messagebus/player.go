package messagebus

import (
	"context"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/localstorage"
)

type player struct {
	sender
	tape tape
	pm   core.PulseManager
}

func NewPlayer(s sender, tape tape, pm core.PulseManager) (*player, error) {
	return &player{sender: s, tape: tape, pm: pm}, nil
}

func (r *player) WriteTape(ctx context.Context, w io.Writer) error {
	panic("can't write the tape from player")
}

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
