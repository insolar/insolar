///
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///

package proc

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_saveCodeFromHeavy_SendFailed(t *testing.T) {
	jetID := gen.JetID()
	codeRef := gen.Reference()
	blobID := gen.ID()
	heavyRef := gen.Reference()

	bus := testutils.NewMessageBusMock(t)
	bus.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.Equal(t, heavyRef, *opt.Receiver)
		require.Equal(t, codeRef, *msg.DefaultTarget())
		return nil, errors.New("test error")
	}
	gc := GetCode{}
	gc.Dep.Bus = bus
	ctx := context.Background()

	_, err := gc.saveCodeFromHeavy(ctx, jetID, codeRef, blobID, &heavyRef)
	require.EqualError(t, err, "failed to send: test error")
}

func Test_saveCodeFromHeavy_WrongAnswer(t *testing.T) {
	jetID := gen.JetID()
	codeRef := gen.Reference()
	blobID := gen.ID()
	heavyRef := gen.Reference()

	bus := testutils.NewMessageBusMock(t)
	bus.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.Equal(t, heavyRef, *opt.Receiver)
		require.Equal(t, codeRef, *msg.DefaultTarget())
		return &reply.NotOK{}, nil
	}
	gc := GetCode{}
	gc.Dep.Bus = bus
	ctx := context.Background()

	_, err := gc.saveCodeFromHeavy(ctx, jetID, codeRef, blobID, &heavyRef)
	require.EqualError(t, err, "failed to fetch code: unexpected reply type *reply.NotOK (reply=&{})")
}

func Test_saveCodeFromHeavy_SaveFailed(t *testing.T) {
	jetID := gen.JetID()
	codeRef := gen.Reference()
	blobID := gen.ID()
	heavyRef := gen.Reference()

	bus := testutils.NewMessageBusMock(t)
	bus.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.Equal(t, heavyRef, *opt.Receiver)
		require.Equal(t, codeRef, *msg.DefaultTarget())
		return &reply.Code{
			Code: []byte("test data"),
		}, nil
	}
	modifier := blob.NewModifierMock(t)
	modifier.SetFunc = func(ctx context.Context, ID insolar.ID, blob blob.Blob) error {
		require.Equal(t, blobID, ID)
		require.Equal(t, jetID, blob.JetID)
		require.Equal(t, []byte("test data"), blob.Value)
		return errors.New("test error")
	}

	gc := GetCode{}
	gc.Dep.Bus = bus
	gc.Dep.BlobModifier = modifier
	ctx := context.Background()

	_, err := gc.saveCodeFromHeavy(ctx, jetID, codeRef, blobID, &heavyRef)
	require.EqualError(t, err, "failed to save: test error")
}

func Test_saveCodeFromHeavy(t *testing.T) {
	jetID := gen.JetID()
	codeRef := gen.Reference()
	blobID := gen.ID()
	heavyRef := gen.Reference()

	bus := testutils.NewMessageBusMock(t)
	bus.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.Equal(t, heavyRef, *opt.Receiver)
		require.Equal(t, codeRef, *msg.DefaultTarget())
		return &reply.Code{
			Code: []byte("test data"),
		}, nil
	}
	modifier := blob.NewModifierMock(t)
	modifier.SetFunc = func(ctx context.Context, ID insolar.ID, blob blob.Blob) error {
		require.Equal(t, blobID, ID)
		require.Equal(t, jetID, blob.JetID)
		require.Equal(t, []byte("test data"), blob.Value)
		return errors.New("test error")
	}

	gc := GetCode{}
	gc.Dep.Bus = bus
	gc.Dep.BlobModifier = modifier
	ctx := context.Background()

	_, err := gc.saveCodeFromHeavy(ctx, jetID, codeRef, blobID, &heavyRef)
	require.EqualError(t, err, "failed to save: test error")
}
