package handle

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/proc"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/testutils"
)

var (
	domainID = *genRandomID(0)
)

func genRandomID(pulse insolar.PulseNumber) *insolar.ID {
	buff := [insolar.RecordIDSize - insolar.PulseNumberSize]byte{}
	_, err := rand.Read(buff[:])
	if err != nil {
		panic(err)
	}
	return insolar.NewID(pulse, buff[:])
}

func genRefWithID(id *insolar.ID) *insolar.Reference {
	return insolar.NewReference(domainID, *id)
}

func genRandomRef(pulse insolar.PulseNumber) *insolar.Reference {
	return genRefWithID(genRandomID(pulse))
}

// redirects to heavy when no index
func TestGetChildren_RedirectsToHaveWhenNoIndex(t *testing.T) {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	msg := message.GetChildren{
		Parent:    *genRandomRef(0),
		FromChild: genRandomID(0),
	}
	parcel := &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	}

	p := proc.GetChildren{
		Jet:     jetID,
		Message: bus.Message{Parcel: parcel /* TODO replyTo */},
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
	jetCoordinator.HeavyMock.ExpectOnce(ctx, parcel.Pulse()).Return(genRandomRef(0), nil)
	p.Dep.JetCoordinator = jetCoordinator

	token := &delegationtoken.GetChildrenRedirectToken{}
	delegationTokenFactory := testutils.NewDelegationTokenFactoryMock(t)
	delegationTokenFactory.IssueGetChildrenRedirectFunc = func(sender *insolar.Reference, redirectedMessage insolar.Message) (insolar.DelegationToken, error) {
		return token, nil
	}
	p.Dep.DelegationTokenFactory = delegationTokenFactory

	p.Dep.JetStorage = jet.NewStorageMock(t)
	p.Dep.RecordAccessor = object.NewRecordAccessorMock(t)
	p.Dep.TreeUpdater = jet.NewTreeUpdaterMock(t)
	p.Dep.IndexSaver = object.NewIndexSaverMock(t)

	err := p.Proceed(ctx)
	require.NoError(t, err)
	require.Equal(t, token, p.Result.Reply.(insolar.RedirectReply).GetToken())
}

// redirect to light when has index and child later than limit
func TestGetChildren_RedirectsToLightChildLaterThanLimit(t *testing.T) {
	// p := createProc(t)
	assert.NoError(t, nil)
}

/*
// redirect to heavy when has index and child earlier than limit
func TestGetChildren_RedirectsToHeavyChildEarlierThanLimit(t *testing.T) {
	// p := createProc(t)
	assert.NoError(t, nil)
}
*/
