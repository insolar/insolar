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

package artifactmanager

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/delegationtoken"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHandler_HandleGetObject_Redirects(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	msg := message.GetObject{
		Head: *genRandomRef(core.FirstPulseNumber),
	}
	objIndex := index.ObjectLifeline{LatestState: genRandomID(0)}

	tf.IssueGetObjectRedirectMock.Return(&delegationtoken.GetObjectRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 2,
	})
	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	t.Run("fetches index from heavy when no index", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if m, ok := gm.(*message.GetObjectIndex); ok {
				assert.Equal(t, msg.Head, m.Object)
				buf, err := index.EncodeObjectLifeline(&objIndex)
				require.NoError(t, err)
				return &reply.ObjectIndex{Index: buf}, nil
			}

			panic("unexpected call")
		}
		jc.HeavyMock.Return(heavyRef, nil)
		rep, err := h.handleGetObject(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
		assert.Nil(t, redirect.StateID)

		idx, err := db.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	t.Run("redirect to light when has index and state later than limit", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			switch p {
			case core.FirstPulseNumber:
				return lightRef, nil
			case core.FirstPulseNumber + 1:
				return &core.RecordRef{}, nil
			}
			panic("unexpected call")
		}
		stateID := genRandomID(core.FirstPulseNumber)
		err = db.SetObjectIndex(ctx, jetID, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		rep, err := h.handleGetObject(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
		assert.Equal(t, stateID, redirect.StateID)
	})

	t.Run("redirect to heavy when has index and state earlier than limit", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		err = db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 2})
		require.NoError(t, err)
		heavyRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		jc.HeavyMock.Return(heavyRef, nil)
		stateID := genRandomID(core.FirstPulseNumber)

		err = db.SetObjectIndex(ctx, jetID, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		rep, err := h.handleGetObject(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
		assert.Equal(t, stateID, redirect.StateID)
	})
}

func TestMessageHandler_HandleGetChildren_Redirects(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	objIndex := index.ObjectLifeline{LatestState: genRandomID(0)}
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)
	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	t.Run("redirects to heavy when no index", func(t *testing.T) {
		mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if m, ok := gm.(*message.GetObjectIndex); ok {
				assert.Equal(t, msg.Parent, m.Object)
				buf, err := index.EncodeObjectLifeline(&objIndex)
				require.NoError(t, err)
				return &reply.ObjectIndex{Index: buf}, nil
			}

			panic("unexpected call")
		}
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		heavyRef := genRandomRef(0)
		jc.HeavyMock.Return(heavyRef, nil)
		rep, err := h.handleGetChildren(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())

		idx, err := db.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	t.Run("redirect to light when has index and child later than limit", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			switch p {
			case core.FirstPulseNumber:
				return lightRef, nil
			case core.FirstPulseNumber + 1:
				return &core.RecordRef{}, nil
			}
			panic("unexpected call")
		}
		err = db.SetObjectIndex(ctx, jetID, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(core.FirstPulseNumber),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	t.Run("redirect to heavy when has index and child earlier than limit", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		err = db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 2})
		require.NoError(t, err)
		heavyRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		jc.HeavyMock.Return(heavyRef, nil)
		err = db.SetObjectIndex(ctx, jetID, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(core.FirstPulseNumber),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
	})
}

func TestMessageHandler_HandleGetDelegate_FetchesIndexFromHeavy(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	delegateType := *genRandomRef(0)
	delegate := *genRandomRef(0)
	objIndex := index.ObjectLifeline{Delegates: map[core.RecordRef]core.RecordRef{delegateType: delegate}}
	msg := message.GetDelegate{
		Head:   *genRandomRef(0),
		AsType: delegateType,
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Head, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(t, err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleGetDelegate(contextWithJet(ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	delegateRep, ok := rep.(*reply.Delegate)
	require.True(t, ok)
	assert.Equal(t, delegate, delegateRep.Head)

	idx, err := db.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
	require.NoError(t, err)
	assert.Equal(t, objIndex.Delegates, idx.Delegates)
}

func TestMessageHandler_HandleUpdateObject_FetchesIndexFromHeavy(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{LatestState: genRandomID(0), State: record.StateActivation}
	amendRecord := record.ObjectAmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(t, err)

	msg := message.UpdateObject{
		Record: record.SerializeRecord(&amendRecord),
		Object: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Object, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(t, err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err = h.Init(ctx)
	require.NoError(t, err)
	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleUpdateObject(contextWithJet(ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber,
	})
	require.NoError(t, err)
	objRep, ok := rep.(*reply.Object)
	require.True(t, ok)

	idx, err := db.GetObjectIndex(ctx, jetID, msg.Object.Record(), false)
	require.NoError(t, err)
	assert.Equal(t, objRep.State, *idx.LatestState)
}

func TestMessageHandler_HandleGetObjectIndex(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)
	msg := message.GetObjectIndex{
		Object: *genRandomRef(0),
	}
	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetCoordinator = jc
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objectIndex := index.ObjectLifeline{LatestState: genRandomID(0)}
	err = db.SetObjectIndex(ctx, jetID, msg.Object.Record(), &objectIndex)
	require.NoError(t, err)

	rep, err := h.handleGetObjectIndex(contextWithJet(ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	indexRep, ok := rep.(*reply.ObjectIndex)
	require.True(t, ok)
	decodedIndex, err := index.DecodeObjectLifeline(indexRep.Index)
	require.NoError(t, err)
	assert.Equal(t, objectIndex, *decodedIndex)
}

func TestMessageHandler_HandleHasPendingRequests(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	msg := message.GetPendingRequests{
		Object: *genRandomRef(0),
	}
	pendingRequests := []core.RecordID{
		*genRandomID(core.FirstPulseNumber),
		*genRandomID(core.FirstPulseNumber),
	}

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.GetRequestsForObjectMock.Return(pendingRequests)

	jetID := *jet.NewID(0, nil)
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(db, &configuration.Ledger{})
	h.JetCoordinator = jc
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	rep, err := h.handleHasPendingRequests(contextWithJet(ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber + 1,
	})
	require.NoError(t, err)
	has, ok := rep.(*reply.HasPendingRequests)
	require.True(t, ok)
	assert.True(t, has.Has)
}

func TestMessageHandler_HandleGetCode_Redirects(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	tf.IssueGetCodeRedirectMock.Return(&delegationtoken.GetCodeRedirectToken{Signature: []byte{1, 2, 3}}, nil)

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	jetID := *jet.NewID(0, nil)
	msg := message.GetCode{
		Code: *genRandomRef(core.FirstPulseNumber),
	}

	t.Run("redirects to light before limit threshold", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			return lightRef, nil
		}
		rep, err := h.handleGetCode(ctx, &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	t.Run("redirects to heavy after limit threshold", func(t *testing.T) {
		err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		err = db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 2})
		require.NoError(t, err)
		heavyRef := genRandomRef(0)
		jc.HeavyMock.Return(heavyRef, nil)
		rep, err := h.handleGetCode(contextWithJet(ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
	})
}

func TestMessageHandler_HandleRegisterChild_FetchesIndexFromHeavy(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 2,
	})

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{LatestState: genRandomID(0), State: record.StateActivation}
	childRecord := record.ChildRecord{
		Ref:       *genRandomRef(0),
		PrevChild: nil,
	}
	amendHash := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err := childRecord.WriteHashData(amendHash)
	require.NoError(t, err)
	childID := core.NewRecordID(0, amendHash.Sum(nil))

	msg := message.RegisterChild{
		Record: record.SerializeRecord(&childRecord),
		Parent: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Parent, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(t, err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err = h.Init(ctx)
	require.NoError(t, err)
	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleRegisterChild(contextWithJet(ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	objRep, ok := rep.(*reply.ID)
	require.True(t, ok)
	assert.Equal(t, *childID, objRep.ID)

	idx, err := db.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
	require.NoError(t, err)
	assert.Equal(t, childID, idx.ChildPointer)
}

const testDropSize uint64 = 100

func addDropSizeToDB(ctx context.Context, t *testing.T, db *storage.DB, jetID core.RecordID) {
	dropSizeData := &jet.DropSize{
		JetID:    jetID,
		PulseNo:  core.FirstPulseNumber,
		DropSize: testDropSize,
	}

	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}

	hasher := testutils.NewPlatformCryptographyScheme().IntegrityHasher()
	_, err := dropSizeData.WriteHashData(hasher)
	require.NoError(t, err)

	signature, err := cryptoServiceMock.Sign(hasher.Sum(nil))
	require.NoError(t, err)

	dropSizeData.Signature = signature.Bytes()

	err = db.AddDropSize(ctx, dropSizeData)
	require.NoError(t, err)
}

func TestMessageHandler_HandleHotRecords(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	jetID := testutils.RandomJet()

	cs := testutils.NewPlatformCryptographyScheme()
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	db.PlatformCryptographyScheme = cs
	err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(t, err)

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	firstID := core.NewRecordID(core.FirstPulseNumber, []byte{1, 2, 3})
	secondId := record.NewRecordIDFromRecord(cs, core.FirstPulseNumber, &record.CodeRecord{})

	firstIndex, _ := index.EncodeObjectLifeline(&index.ObjectLifeline{
		LatestState: firstID,
	})
	err = db.SetObjectIndex(ctx, jetID, firstID, &index.ObjectLifeline{
		LatestState: firstID,
	})

	dropSizeHistory, err := db.GetDropSizeHistory(ctx, jetID)
	require.NoError(t, err)
	require.Equal(t, jet.DropSizeHistory{}, dropSizeHistory)
	addDropSizeToDB(ctx, t, db, jetID)

	dropSizeHistory, err = db.GetDropSizeHistory(ctx, jetID)
	require.NoError(t, err)

	obj := core.RecordID{}
	hotIndexes := &message.HotData{
		Jet:         *core.NewRecordRef(core.DomainID, jetID),
		PulseNumber: core.FirstPulseNumber,
		RecentObjects: map[core.RecordID]*message.HotIndex{
			*firstID: {
				Index: firstIndex,
				TTL:   321,
			},
		},
		PendingRequests: map[core.RecordID]map[core.RecordID][]byte{
			obj: {
				*secondId: record.SerializeRecord(&record.CodeRecord{}),
			},
		},
		Drop:               jet.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}},
		DropJet:            jetID,
		JetDropSizeHistory: dropSizeHistory,
	}

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestFunc = func(o, p core.RecordID) {
		require.Equal(t, o, obj)
		require.Equal(t, p, *secondId)
	}
	recentStorageMock.AddObjectWithTLLFunc = func(p core.RecordID, ttl int) {
		require.Equal(t, p, *firstID)
		require.Equal(t, 320, ttl)
	}
	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h := NewMessageHandler(db, &configuration.Ledger{})
	h.JetCoordinator = jc
	h.RecentStorageProvider = provideMock
	h.Bus = mb
	err = h.Init(ctx)
	require.NoError(t, err)

	res, err := h.handleHotRecords(ctx, &message.Parcel{Msg: hotIndexes})

	require.NoError(t, err)
	require.Equal(t, res, &reply.OK{})

	savedDrop, err := h.db.GetDrop(ctx, jetID, core.FirstPulseNumber)
	require.NoError(t, err)
	require.Equal(t, &jet.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}}, savedDrop)

	// check drop size list
	dropSizeHistory, err = db.GetDropSizeHistory(ctx, jetID)
	require.NoError(t, err)
	require.Equal(t, testDropSize, dropSizeHistory[0].DropSize)
	require.Equal(t, jetID, dropSizeHistory[0].JetID)
	require.Equal(t, core.FirstPulseNumber, int(dropSizeHistory[0].PulseNo))

	recentStorageMock.MinimockFinish()

}

func TestMessageHandler_HandleValidationCheck(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetCoordinator = jc
	h.Bus = mb
	err := h.Init(ctx)
	require.NoError(t, err)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	t.Run("returns not ok when not valid", func(t *testing.T) {
		validatedStateID, err := db.SetRecord(ctx, jetID, 0, &record.ObjectAmendRecord{})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: genRandomID(0),
		}

		rep, err := h.handleValidationCheck(contextWithJet(ctx, jetID), &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.NotOK)
		assert.True(t, ok)
	})

	t.Run("returns ok when valid", func(t *testing.T) {
		approvedStateID := *genRandomID(0)
		validatedStateID, err := db.SetRecord(ctx, jetID, 0, &record.ObjectAmendRecord{
			PrevState: approvedStateID,
		})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: &approvedStateID,
		}

		rep, err := h.handleValidationCheck(contextWithJet(ctx, jetID), &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.OK)
		assert.True(t, ok)
	})
}

func TestMessageHandler_HandleJetDrop_SaveJet(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer func() {
		cleaner()
		mc.Finish()
	}()

	jetID := jet.NewID(0, []byte{2})
	msg := message.JetDrop{
		JetID: *jetID,
	}
	expectedSetId := jet.IDSet{
		*jetID: struct{}{},
	}

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	// Act
	response, err := h.handleJetDrop(ctx, &message.Parcel{Msg: &msg})
	require.NoError(t, err)

	idSet, err := db.GetJets(ctx)
	require.NoError(t, err)
	require.NotNil(t, idSet)

	// Assert
	require.Equal(t, &reply.OK{}, response)
	for id := range expectedSetId {
		require.True(t, idSet.Has(id))
	}
}

func TestMessageHandler_HandleJetDrop_SaveJet_ExistingMap(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer func() {
		cleaner()
		mc.Finish()
	}()

	jetID := jet.NewID(0, []byte{2})
	secondJetID := jet.NewID(0, []byte{3})
	msg := message.JetDrop{
		JetID: *jetID,
	}
	secondMsg := message.JetDrop{
		JetID: *secondJetID,
	}
	expectedSetId := jet.IDSet{
		*jetID:       struct{}{},
		*secondJetID: struct{}{},
	}

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	// Act
	response, err := h.handleJetDrop(ctx, &message.Parcel{Msg: &msg})
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, response)

	secondResponse, err := h.handleJetDrop(ctx, &message.Parcel{Msg: &secondMsg})
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, secondResponse)

	idSet, err := db.GetJets(ctx)
	require.NoError(t, err)
	require.NotNil(t, idSet)

	// Assert
	for id := range expectedSetId {
		require.True(t, idSet.Has(id))
	}
}
