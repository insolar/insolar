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
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestPlayer_Send(t *testing.T) {
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	msg := message.GenesisRequest{Name: "test"}
	parcel := message.Parcel{Msg: &msg}
	msgHash := GetMessageHash(pcs, &parcel)
	s := NewsenderMock(mc)
	s.CreateParcelFunc = func(p context.Context, p2 core.Message, p3 core.DelegationToken, p4 core.Pulse) (r core.Parcel, r1 error) {
		return &parcel, nil
	}
	tape := NewtapeMock(mc)
	pulseStorageMock := testutils.NewPulseStorageMock(t)
	pulseStorageMock.CurrentMock.Return(core.GenesisPulse, nil)
	player := newPlayer(s, tape, pcs, pulseStorageMock)

	t.Run("with no reply on the Tape doesn't send the message and returns an error", func(t *testing.T) {
		tape.GetMock.Expect(ctx, msgHash).Return(nil, ErrNoReply)

		_, err := player.Send(ctx, &msg, nil)
		require.Equal(t, ErrNoReply, err)
	})

	t.Run("with reply on the Tape doesn't send the message and returns reply from the storageTape", func(t *testing.T) {
		expectedRep := reply.Object{Memory: []byte{1, 2, 3}}
		item := TapeItem{
			Reply: &expectedRep,
		}
		tape.GetMock.Expect(ctx, msgHash).Return(&item, nil)
		rep, err := player.Send(ctx, &msg, nil)

		require.NoError(t, err)
		require.Equal(t, &expectedRep, rep)
	})
}
