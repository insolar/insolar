package proc

import (
	// "github.com/insolar/insolar/testutils"
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChildren_RedirectsToHeavyWhenNoIndex(t *testing.T) {
	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}

	heavyRef := genRandomRef(0)

	jc := jet.NewCoordinatorMock(t)
	jc.HeavyMock.Return(heavyRef, nil)
	jc.IsBeyondLimitMock.Return(true, nil)

	tf := testutils.NewDelegationTokenFactoryMock(t)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)

	childPointer := genRandomID(insolar.FirstPulseNumber)
	gc := GetChildren{
		currentChild: childPointer,
		msg:          &msg,
		parcel: &message.Parcel{
			Msg:         &msg,
			Sender:      *genRandomRef(insolar.FirstPulseNumber),
			PulseNumber: insolar.FirstPulseNumber + 1,
		},
	}
	gc.Dep.Coordinator = jc
	gc.Dep.DelegationTokenFactory = tf

	rep := gc.reply(context.TODO())
	require.NoError(t, rep.Err)
	redirect, ok := rep.Reply.(*reply.GetChildrenRedirectReply)
	require.True(t, ok)
	token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
	assert.Equal(t, []byte{1, 2, 3}, token.Signature)
	assert.Equal(t, heavyRef, redirect.GetReceiver())
}

func TestGetChildren_RedirectToLight(t *testing.T) {
	lightRef := genRandomRef(0)

	jc := jet.NewCoordinatorMock(t)
	jc.IsBeyondLimitMock.Return(false, nil)
	jc.NodeForJetMock.Return(lightRef, nil)

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	childPointer := genRandomID(insolar.FirstPulseNumber)

	indexMemoryStor := object.NewIndexMemory()
	ctx := context.TODO()
	err := indexMemoryStor.Set(ctx, *msg.Parent.Record(), object.Lifeline{
		ChildPointer: genRandomID(insolar.FirstPulseNumber),
		JetID:        insolar.JetID(jetID),
	})
	require.NoError(t, err)

	gc := GetChildren{
		currentChild: childPointer,
		msg:          &msg,
		parcel: &message.Parcel{
			Msg:         &msg,
			Sender:      *genRandomRef(insolar.FirstPulseNumber),
			PulseNumber: insolar.FirstPulseNumber + 1,
		},
	}
	gc.Dep.Coordinator = jc
	gc.Dep.JetStorage = jet.NewStore()
	gc.Dep.JetStorage.Update(ctx, insolar.FirstPulseNumber+1, true)
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.Material, error) {
		return record.Material{}, object.ErrNotFound
	}
	gc.Dep.RecordAccessor = ra

	tf := testutils.NewDelegationTokenFactoryMock(t)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	gc.Dep.DelegationTokenFactory = tf

	rep := gc.reply(ctx)

	require.NoError(t, rep.Err)
	redirect, ok := rep.Reply.(*reply.GetChildrenRedirectReply)
	require.True(t, ok)
	token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
	assert.Equal(t, []byte{1, 2, 3}, token.Signature)
	assert.Equal(t, lightRef, redirect.GetReceiver())
}

func TestGetChildren_RedirectToHeavy(t *testing.T) {
	heavyRef := genRandomRef(0)

	jc := jet.NewCoordinatorMock(t)
	jc.IsBeyondLimitMock.Return(false, nil)
	jc.NodeForJetMock.Return(heavyRef, nil)

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	childPointer := genRandomID(insolar.FirstPulseNumber)

	indexMemoryStor := object.NewIndexMemory()
	ctx := context.TODO()
	err := indexMemoryStor.Set(ctx, *msg.Parent.Record(), object.Lifeline{
		ChildPointer: genRandomID(insolar.FirstPulseNumber),
		JetID:        insolar.JetID(jetID),
	})
	require.NoError(t, err)

	gc := GetChildren{
		currentChild: childPointer,
		msg:          &msg,
		parcel: &message.Parcel{
			Msg:         &msg,
			Sender:      *genRandomRef(insolar.FirstPulseNumber),
			PulseNumber: insolar.FirstPulseNumber + 1,
		},
	}
	gc.Dep.Coordinator = jc
	gc.Dep.JetStorage = jet.NewStore()
	gc.Dep.JetStorage.Update(ctx, insolar.FirstPulseNumber+1, true)
	ra := object.NewRecordAccessorMock(t)
	ra.ForIDFunc = func(ctx context.Context, id insolar.ID) (record.Material, error) {
		return record.Material{}, object.ErrNotFound
	}
	gc.Dep.RecordAccessor = ra

	tf := testutils.NewDelegationTokenFactoryMock(t)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	gc.Dep.DelegationTokenFactory = tf

	rep := gc.reply(ctx)

	require.NoError(t, rep.Err)
	redirect, ok := rep.Reply.(*reply.GetChildrenRedirectReply)
	require.True(t, ok)
	token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
	assert.Equal(t, []byte{1, 2, 3}, token.Signature)
	assert.Equal(t, heavyRef, redirect.GetReceiver())
}
