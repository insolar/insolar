package messagebus

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestPlayer_Send(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	msg := message.GenesisRequest{Name: "test"}
	signedMessage := message.SignedMessage{Msg: &msg}
	msgHash := GetMessageHash(&signedMessage)
	pm := testutils.NewPulseManagerMock(mc)
	pm.CurrentMock.Return(&core.Pulse{PulseNumber: 42}, nil)
	s := NewsenderMock(mc)
	s.CreateSignedMessageFunc = func(c context.Context, p core.PulseNumber, m core.Message) (core.SignedMessage, error) {
		return &signedMessage, nil
	}
	tape := NewtapeMock(mc)
	player, err := NewPlayer(s, tape, pm)

	t.Run("with no reply on the storagetape doesn't send the message and returns an error", func(t *testing.T) {
		tape.GetReplyMock.Expect(ctx, msgHash).Return(nil, ErrNoReply)

		_, err = player.Send(ctx, &msg)
		assert.Equal(t, ErrNoReply, err)
	})

	t.Run("with reply on the storagetape doesn't send the message and returns reply from the storagetape", func(t *testing.T) {
		expectedRep := reply.Object{Memory: []byte{1, 2, 3}}
		tape.GetReplyMock.Expect(ctx, msgHash).Return(&expectedRep, nil)
		rep, err := player.Send(ctx, &msg)

		assert.NoError(t, err)
		assert.Equal(t, &expectedRep, rep)
	})
}
