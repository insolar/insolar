package handle

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/proc"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/testutils"
)

// redirects when index can't be found
func TestGetChildren_RedirectsWhenNoIndex(t *testing.T) {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	msg := message.GetChildren{
		Parent:    gen.Reference(),
		FromChild: gen.IDPointer(),
	}
	parcel := &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	}

	p := proc.GetChildren{
		Jet:     jetID,
		Message: bus.Message{Parcel: parcel},
	}

	ctx := context.Background()

	recentIndexStorage := recentstorage.NewRecentIndexStorageMock(t)
	recentIndexStorage.AddObjectMock.ExpectOnce(ctx, *msg.Parent.Record())

	recentStorageProvider := recentstorage.NewProviderMock(t)
	recentStorageProvider.GetIndexStorageMock.ExpectOnce(ctx, jetID).Return(recentIndexStorage)
	p.Dep.RecentStorageProvider = recentStorageProvider

	idLocker := storage.NewIDLockerMock(t)
	idLocker.LockMock.ExpectOnce(msg.Parent.Record())
	idLocker.UnlockMock.ExpectOnce(msg.Parent.Record())
	p.Dep.IDLocker = idLocker

	indexAccessor := object.NewIndexAccessorMock(t)
	indexAccessor.ForIDMock.ExpectOnce(ctx, *msg.Parent.Record()).Return(object.Lifeline{}, nil)
	p.Dep.IndexAccessor = indexAccessor

	jetCoordinator := testutils.NewJetCoordinatorMock(t)
	jetCoordinator.IsBeyondLimitMock.ExpectOnce(ctx, parcel.Pulse(), msg.FromChild.Pulse()).Return(true, nil)
	jetCoordinator.HeavyMock.ExpectOnce(ctx, parcel.Pulse()).Return(gen.ReferencePointer(), nil)
	p.Dep.JetCoordinator = jetCoordinator

	token := &delegationtoken.GetChildrenRedirectToken{}
	delegationTokenFactory := testutils.NewDelegationTokenFactoryMock(t)
	delegationTokenFactory.IssueGetChildrenRedirectFunc = func(sender *insolar.Reference, redirectedMessage insolar.Message) (insolar.DelegationToken, error) {
		return token, nil
	}
	p.Dep.DelegationTokenFactory = delegationTokenFactory

	p.Dep.JetStorage = jet.NewStorageMock(t)
	p.Dep.RecordAccessor = object.NewRecordAccessorMock(t)
	p.Dep.TreeUpdater = jet.NewFetcherMock(t)
	p.Dep.IndexSaver = object.NewIndexSaverMock(t)

	err := p.Proceed(ctx)
	require.NoError(t, err)
	require.Equal(t, token, p.Result.Reply.(insolar.RedirectReply).GetToken())
}

// redirect when there is an index but the child can't be found
func TestGetChildren_RedirectWhenFirstChildNotFound(t *testing.T) {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	msg := message.GetChildren{
		Parent:    gen.Reference(),
		FromChild: gen.IDPointer(),
	}
	parcel := &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	}

	p := proc.GetChildren{
		Jet:     jetID,
		Message: bus.Message{Parcel: parcel},
	}

	ctx := context.Background()

	recentIndexStorage := recentstorage.NewRecentIndexStorageMock(t)
	recentIndexStorage.AddObjectMock.ExpectOnce(ctx, *msg.Parent.Record())

	recentStorageProvider := recentstorage.NewProviderMock(t)
	recentStorageProvider.GetIndexStorageMock.ExpectOnce(ctx, jetID).Return(recentIndexStorage)
	p.Dep.RecentStorageProvider = recentStorageProvider

	idLocker := storage.NewIDLockerMock(t)
	idLocker.LockMock.ExpectOnce(msg.Parent.Record())
	idLocker.UnlockMock.ExpectOnce(msg.Parent.Record())
	p.Dep.IDLocker = idLocker

	indexAccessor := object.NewIndexAccessorMock(t)
	indexAccessor.ForIDMock.ExpectOnce(ctx, *msg.Parent.Record()).Return(object.Lifeline{}, nil)
	p.Dep.IndexAccessor = indexAccessor

	jetCoordinator := testutils.NewJetCoordinatorMock(t)
	jetCoordinator.IsBeyondLimitMock.ExpectOnce(ctx, parcel.Pulse(), msg.FromChild.Pulse()).Return(false, nil)
	jetCoordinator.HeavyMock.ExpectOnce(ctx, parcel.Pulse()).Return(gen.ReferencePointer(), nil)

	childJetId := insolar.JetID(gen.ID())
	jetCoordinator.NodeForJetMock.ExpectOnce(ctx, insolar.ID(childJetId), parcel.Pulse(), msg.FromChild.Pulse()).Return(gen.ReferencePointer(), nil)
	p.Dep.JetCoordinator = jetCoordinator

	token := &delegationtoken.GetChildrenRedirectToken{}
	delegationTokenFactory := testutils.NewDelegationTokenFactoryMock(t)
	delegationTokenFactory.IssueGetChildrenRedirectFunc = func(sender *insolar.Reference, redirectedMessage insolar.Message) (insolar.DelegationToken, error) {
		return token, nil
	}
	p.Dep.DelegationTokenFactory = delegationTokenFactory

	jetStorage := jet.NewStorageMock(t)
	jetStorage.ForIDMock.ExpectOnce(ctx, msg.FromChild.Pulse(), *msg.Parent.Record()).Return(childJetId, true)
	p.Dep.JetStorage = jetStorage

	recordAccessor := object.NewRecordAccessorMock(t)
	recordAccessor.ForIDMock.ExpectOnce(ctx, *msg.FromChild).Return(record.MaterialRecord{}, object.ErrNotFound)
	p.Dep.RecordAccessor = recordAccessor

	p.Dep.TreeUpdater = jet.NewFetcherMock(t)
	p.Dep.IndexSaver = object.NewIndexSaverMock(t)

	err := p.Proceed(ctx)
	require.NoError(t, err)
	require.Equal(t, token, p.Result.Reply.(insolar.RedirectReply).GetToken())
}
