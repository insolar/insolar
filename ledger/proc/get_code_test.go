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
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGetCode_handle_RecordAccessorErrNotFound_Error(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{}, object.ErrNotFound
	}
	coo := testutils.NewJetCoordinatorMock(t)
	coo.NodeForJetFunc = func(ctx context.Context, id insolar.ID, rootPN insolar.PulseNumber, targetPN insolar.PulseNumber) (*insolar.Reference, error) {
		require.Equal(t, insolar.ID(jetID), id)
		require.Equal(t, insolar.PulseNumber(42), rootPN)
		require.Equal(t, gc.Code.Record().Pulse(), targetPN)
		return nil, errors.New("test error")
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Coordinator = coo
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "test error")
}

func TestGetCode_handle_RecordAccessorErrNotFound(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	parcel.MessageFunc = func() insolar.Message {
		return &message.GetCode{}

	}
	senderRef := gen.Reference()
	parcel.GetSenderFunc = func() insolar.Reference {
		return senderRef
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{}, object.ErrNotFound
	}
	nodeRef := gen.Reference()
	coo := testutils.NewJetCoordinatorMock(t)
	coo.NodeForJetFunc = func(ctx context.Context, id insolar.ID, rootPN insolar.PulseNumber, targetPN insolar.PulseNumber) (*insolar.Reference, error) {
		require.Equal(t, insolar.ID(jetID), id)
		require.Equal(t, insolar.PulseNumber(42), rootPN)
		require.Equal(t, gc.Code.Record().Pulse(), targetPN)
		return &nodeRef, nil
	}
	token := &delegationtoken.GetCodeRedirectToken{}
	factory := testutils.NewDelegationTokenFactoryMock(t)
	factory.IssueGetCodeRedirectFunc = func(ref *insolar.Reference, msg insolar.Message) (insolar.DelegationToken, error) {
		require.Equal(t, senderRef, *ref)
		return token, nil
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Coordinator = coo
	gc.Dep.DelegationTokenFactory = factory
	ctx := context.Background()
	result, err := gc.handle(ctx)
	require.NoError(t, err)
	reply, ok := result.(*reply.GetCodeRedirectReply)
	require.True(t, ok)
	require.Equal(t, token, reply.GetToken())
	require.Equal(t, nodeRef, *reply.Receiver)
}

func TestGetCode_handle_RecordAccessorError(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{}, errors.New("test error")
	}
	gc.Dep.RecordAccessor = ra
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "test error")
}

func TestGetCode_handle_CastError(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{}, nil
	}
	gc.Dep.RecordAccessor = ra
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "failed to retrieve code record: invalid reference")
}

func TestGetCode_handle_AccessorErrNotFound_Error(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	codeID := gen.ID()
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{
			Record: &object.CodeRecord{
				Code: &codeID,
			},
		}, nil
	}
	acc := blob.NewAccessorMock(t)
	acc.ForIDFunc = func(ctx context.Context, id insolar.ID) (blob.Blob, error) {
		require.Equal(t, id, codeID)
		return blob.Blob{}, blob.ErrNotFound
	}
	coo := testutils.NewJetCoordinatorMock(t)
	coo.HeavyFunc = func(ctx context.Context, pulse insolar.PulseNumber) (*insolar.Reference, error) {
		require.Equal(t, insolar.PulseNumber(42), pulse)
		return nil, errors.New("test error")
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Accessor = acc
	gc.Dep.Coordinator = coo
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "test error")
}

func TestGetCode_handle_AccessorErrNotFound(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	codeID := gen.ID()
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{
			Record: &object.CodeRecord{
				Code: &codeID,
			},
		}, nil
	}
	acc := blob.NewAccessorMock(t)
	acc.ForIDFunc = func(ctx context.Context, id insolar.ID) (blob.Blob, error) {
		require.Equal(t, id, codeID)
		return blob.Blob{}, blob.ErrNotFound
	}
	coo := testutils.NewJetCoordinatorMock(t)
	heavyRef := gen.Reference()
	coo.HeavyFunc = func(ctx context.Context, pulse insolar.PulseNumber) (*insolar.Reference, error) {
		require.Equal(t, insolar.PulseNumber(42), pulse)
		return &heavyRef, nil
	}
	bus := testutils.NewMessageBusMock(t)
	bus.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		return nil, errors.New("test error")
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Accessor = acc
	gc.Dep.Coordinator = coo
	gc.Dep.Bus = bus
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "failed to send: test error")
}

func TestGetCode_handle_AccessorError(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	codeID := gen.ID()
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{
			Record: &object.CodeRecord{
				Code: &codeID,
			},
		}, nil
	}
	acc := blob.NewAccessorMock(t)
	acc.ForIDFunc = func(ctx context.Context, id insolar.ID) (blob.Blob, error) {
		require.Equal(t, id, codeID)
		return blob.Blob{}, errors.New("test error")
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Accessor = acc
	ctx := context.Background()
	_, err := gc.handle(ctx)
	require.EqualError(t, err, "test error")
}

func TestGetCode_handle(t *testing.T) {
	jetID := gen.JetID()
	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	gc := GetCode{
		JetID: jetID,
		Code:  gen.Reference(),
		Message: bus.Message{
			Parcel: parcel,
		},
	}
	codeID := gen.ID()
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
		require.Equal(t, *gc.Code.Record(), id)
		return record.MaterialRecord{
			Record: &object.CodeRecord{
				Code:        &codeID,
				MachineType: 42,
			},
		}, nil
	}
	acc := blob.NewAccessorMock(t)
	acc.ForIDFunc = func(ctx context.Context, id insolar.ID) (blob.Blob, error) {
		require.Equal(t, id, codeID)
		return blob.Blob{
			Value: []byte("test blob"),
		}, nil
	}
	gc.Dep.RecordAccessor = ra
	gc.Dep.Accessor = acc
	ctx := context.Background()
	res, err := gc.handle(ctx)
	require.NoError(t, err)
	require.Equal(t, reply.TypeCode, res.Type())
	result, ok := res.(*reply.Code)
	require.True(t, ok)
	require.Equal(t, []byte("test blob"), result.Code)
	require.Equal(t, insolar.MachineType(42), result.MachineType)
}

func TestGetCode_saveCodeFromHeavy_SendFailed(t *testing.T) {
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

func TestGetCode_saveCodeFromHeavy_WrongAnswer(t *testing.T) {
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

func TestGetCode_saveCodeFromHeavy_SaveFailed(t *testing.T) {
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

func TestGetCode_saveCodeFromHeavy(t *testing.T) {
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
