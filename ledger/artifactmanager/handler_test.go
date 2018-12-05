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
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/record"
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

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)

	msg := message.GetObject{
		Head: *genRandomRef(0),
	}
	objIndex := index.ObjectLifeline{LatestState: genRandomID(0)}

	tf.IssueGetObjectRedirectMock.Return(&delegationtoken.GetObjectRedirect{Signature: []byte{1, 2, 3}}, nil)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Head, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(t, err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.Recent = recentStorageMock
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb

	t.Run("fetches index from heavy when no index", func(t *testing.T) {
		lightRef := genRandomRef(0)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleHeavyExecutor, &msg.Head, 0,
		).Return(
			[]core.RecordRef{*lightRef}, nil,
		)
		rep, err := h.handleGetObject(ctx, 0, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirect)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
		assert.Nil(t, redirect.StateID)

		idx, err := db.GetObjectIndex(ctx, msg.Head.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	t.Run("redirect to light when has index and state later than limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		stateID := genRandomID(0)
		err := db.SetObjectIndex(ctx, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleLightExecutor, &msg.Head, 0,
		).Return(
			[]core.RecordRef{*lightRef}, nil,
		)
		rep, err := h.handleGetObject(ctx, 1, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirect)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
		assert.Equal(t, stateID, redirect.StateID)
	})

	t.Run("redirect to heavy when has index and state earlier than limit", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		stateID := genRandomID(0)

		err := db.SetObjectIndex(ctx, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleHeavyExecutor, &msg.Head, 5,
		).Return(
			[]core.RecordRef{*heavyRef}, nil,
		)
		rep, err := h.handleGetObject(ctx, 5, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirect)
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

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirect{Signature: []byte{1, 2, 3}}, nil)
	mb := testutils.NewMessageBusMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	objIndex := index.ObjectLifeline{LatestState: genRandomID(0)}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Parent, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(t, err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	h.Recent = recentStorageMock

	t.Run("redirects to heavy when no index", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleHeavyExecutor, &msg.Parent, 0,
		).Return(
			[]core.RecordRef{*heavyRef}, nil,
		)
		rep, err := h.handleGetChildren(ctx, 0, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirect)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())

		idx, err := db.GetObjectIndex(ctx, msg.Parent.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	t.Run("redirect to light when has index and child later than limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		err := db.SetObjectIndex(ctx, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(0),
		})
		require.NoError(t, err)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleLightExecutor, &msg.Parent, 0,
		).Return(
			[]core.RecordRef{*lightRef}, nil,
		)
		rep, err := h.handleGetChildren(ctx, 1, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirect)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	t.Run("redirect to heavy when has index and child earlier than limit", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		err := db.SetObjectIndex(ctx, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(0),
		})
		require.NoError(t, err)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleHeavyExecutor, &msg.Parent, 5,
		).Return(
			[]core.RecordRef{*heavyRef}, nil,
		)
		rep, err := h.handleGetChildren(ctx, 5, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirect)
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

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.Recent = recentStorageMock
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
	heavyRef := genRandomRef(0)
	jc.QueryRoleMock.Expect(
		ctx, core.DynamicRoleHeavyExecutor, &msg.Head, 0,
	).Return(
		[]core.RecordRef{*heavyRef}, nil,
	)
	rep, err := h.handleGetDelegate(ctx, 0, &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	delegateRep, ok := rep.(*reply.Delegate)
	require.True(t, ok)
	assert.Equal(t, delegate, delegateRep.Head)

	idx, err := db.GetObjectIndex(ctx, msg.Head.Record(), false)
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

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.Recent = recentStorageMock

	objIndex := index.ObjectLifeline{LatestState: genRandomID(0), State: record.StateActivation}
	amendRecord := record.ObjectAmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(t, err)
	amendID := core.NewRecordID(0, amendHash.Sum(nil))

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
	heavyRef := genRandomRef(0)
	jc.QueryRoleMock.Expect(
		ctx, core.DynamicRoleHeavyExecutor, &msg.Object, 0,
	).Return(
		[]core.RecordRef{*heavyRef}, nil,
	)
	rep, err := h.handleUpdateObject(ctx, 0, &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	objRep, ok := rep.(*reply.Object)
	require.True(t, ok)
	assert.Equal(t, *amendID, objRep.State)

	idx, err := db.GetObjectIndex(ctx, msg.Object.Record(), false)
	require.NoError(t, err)
	assert.Equal(t, amendID, idx.LatestState)
}

func TestMessageHandler_HandleGetObjectIndex(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.Recent = recentStorageMock

	msg := message.GetObjectIndex{
		Object: *genRandomRef(0),
	}
	objectIndex := index.ObjectLifeline{LatestState: genRandomID(0)}
	err := db.SetObjectIndex(ctx, msg.Object.Record(), &objectIndex)
	require.NoError(t, err)

	rep, err := h.handleGetObjectIndex(ctx, &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	indexRep, ok := rep.(*reply.ObjectIndex)
	require.True(t, ok)
	decodedIndex, err := index.DecodeObjectLifeline(indexRep.Index)
	require.NoError(t, err)
	assert.Equal(t, objectIndex, *decodedIndex)
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

	msg := message.GetCode{
		Code: *genRandomRef(0),
	}
	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	tf.IssueGetCodeRedirectMock.Return(&delegationtoken.GetCodeRedirect{Signature: []byte{1, 2, 3}}, nil)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})

	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	h.Recent = recentStorageMock

	t.Run("redirects to light when created after limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleLightExecutor, &msg.Code, 0,
		).Return(
			[]core.RecordRef{*lightRef}, nil,
		)
		rep, err := h.handleGetCode(ctx, 0, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirect)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	t.Run("redirects to heavy when created before limit", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		jc.QueryRoleMock.Expect(
			ctx, core.DynamicRoleHeavyExecutor, &msg.Code, 5,
		).Return(
			[]core.RecordRef{*heavyRef}, nil,
		)
		rep, err := h.handleGetCode(ctx, 5, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirect)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirect)
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

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.Recent = recentStorageMock

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
	heavyRef := genRandomRef(0)
	jc.QueryRoleMock.Expect(
		ctx, core.DynamicRoleHeavyExecutor, &msg.Parent, 0,
	).Return(
		[]core.RecordRef{*heavyRef}, nil,
	)
	rep, err := h.handleRegisterChild(ctx, 0, &message.Parcel{
		Msg: &msg,
	})
	require.NoError(t, err)
	objRep, ok := rep.(*reply.ID)
	require.True(t, ok)
	assert.Equal(t, *childID, objRep.ID)

	idx, err := db.GetObjectIndex(ctx, msg.Parent.Record(), false)
	require.NoError(t, err)
	assert.Equal(t, childID, idx.ChildPointer)
}

func TestMessageHandler_HandleHotRecords(t *testing.T) {
	ctx := inslogger.TestContext(t)

	idCreator, idCreatorCleaner := storagetest.TmpDB(ctx, t)
	defer idCreatorCleaner()
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	err := db.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(t, err)

	firstID := core.NewRecordID(core.FirstPulseNumber, []byte{1, 2, 3})
	secondId, _ := idCreator.SetRecord(ctx, core.FirstPulseNumber, &record.CodeRecord{})

	firstIndex, _ := index.EncodeObjectLifeline(&index.ObjectLifeline{
		LatestState: firstID,
	})
	err = db.SetObjectIndex(ctx, firstID, &index.ObjectLifeline{
		LatestState: firstID,
	})
	require.NoError(t, err)
	hotIndexes := &message.HotData{
		PulseNumber: core.FirstPulseNumber,
		RecentObjects: map[core.RecordID]*message.HotIndex{
			*firstID: {
				Index: firstIndex,
				TTL:   321,
			},
		},
		PendingRequests: map[core.RecordID][]byte{
			*secondId: record.SerializeRecord(&record.CodeRecord{}),
		},
		Drop: jetdrop.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}},
	}

	recentMock := recentstorage.NewRecentStorageMock(t)
	recentMock.AddPendingRequestFunc = func(p core.RecordID) {
		require.Equal(t, p, *secondId)
	}
	recentMock.AddObjectWithTLLFunc = func(p core.RecordID, ttl int, isMine bool) {
		require.Equal(t, p, *firstID)
		require.Equal(t, 320, ttl)
		require.Equal(t, true, isMine)
	}

	h := NewMessageHandler(db, &configuration.Ledger{})
	h.Recent = recentMock

	res, err := h.handleHotRecords(ctx, &message.Parcel{Msg: hotIndexes})

	require.NoError(t, err)
	require.Equal(t, res, &reply.OK{})

	savedDrop, err := h.db.GetDrop(ctx, core.FirstPulseNumber)
	require.NoError(t, err)
	require.Equal(t, &jetdrop.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}}, savedDrop)

	recentMock.MinimockFinish()
}

func TestMessageHandler_HandleValidationCheck(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	h := NewMessageHandler(db, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.Recent = recentStorageMock

	t.Run("returns not ok when not valid", func(t *testing.T) {
		validatedStateID, err := db.SetRecord(ctx, 0, &record.ObjectAmendRecord{})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: genRandomID(0),
		}

		rep, err := h.handleValidationCheck(ctx, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.NotOK)
		assert.True(t, ok)
	})

	t.Run("returns ok when valid", func(t *testing.T) {
		approvedStateID := *genRandomID(0)
		validatedStateID, err := db.SetRecord(ctx, 0, &record.ObjectAmendRecord{
			PrevState: approvedStateID,
		})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: &approvedStateID,
		}

		rep, err := h.handleValidationCheck(ctx, &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.OK)
		assert.True(t, ok)
	})
}
