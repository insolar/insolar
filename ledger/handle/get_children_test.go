package handle

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/ledger/recentstorage"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/proc"
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

func createProc(t *testing.T, ctx context.Context) *proc.GetChildren {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	parcel := &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	}

	p := proc.GetChildren{
		Jet:     jetID,
		Message: bus.Message{Parcel: parcel /* TODO replyTo */},
	}

	recentIndexStorage := recentstorage.NewRecentIndexStorageMock(t)
	recentIndexStorage.AddObjectMock.ExpectOnce(ctx, *msg.Parent.Record())

	recentStorageProvider := recentstorage.NewProviderMock(t)
	recentStorageProvider.GetIndexStorageMock.ExpectOnce(ctx, jetID).Return(recentIndexStorage)
	p.Dep.RecentStorageProvider = recentStorageProvider

	idLocker := storage.NewIDLockerMock(t)
	idLocker.LockMock.ExpectOnce(msg.Parent.Record())
	idLocker.UnlockMock.ExpectOnce(msg.Parent.Record())

	p.Dep.IDLocker = idLocker
	p.Dep.IndexAccessor = object.NewIndexAccessorMock(t)
	p.Dep.JetCoordinator = testutils.NewJetCoordinatorMock(t)
	p.Dep.JetStorage = jet.NewStorageMock(t)
	p.Dep.DelegationTokenFactory = testutils.NewDelegationTokenFactoryMock(t)
	p.Dep.RecordAccessor = object.NewRecordAccessorMock(t)
	p.Dep.TreeUpdater = jet.NewTreeUpdaterMock(t)
	p.Dep.IndexSaver = object.NewIndexSaverMock(t)
	return &p
}

// redirects to heavy when no index
func TestGetChildren_RedirectsToHaveWhenNoIndex(t *testing.T) {
	ctx := context.Background()
	p := createProc(t, ctx)
	/*
		fl := flow.NewFlowMock(t)
		fl.ProcedureFunc = func(ctx context.Context, pr flow.Procedure) error {
			if fetchJet, ok := pr.(*proc.FetchJet); ok {
				require.Equal(t, "TypeGetChildren", fetchJet.Parcel.Message().Type().String())
				fetchJet.Result.Jet = insolar.JetID(jetID)
				return nil
			} else if getChildren, ok := pr.(*proc.GetChildren); ok {
				require.Equal(t, getChildren.Jet, jetID)
				//require.Equal(t, getChildren.Code, codeRef)
				//require.Equal(t, getChildren.Message, msg)
				return nil
			}
			t.Fatal("you shouldn't be here")
			return nil
		}
	*/

	err := p.Proceed(ctx)
	//err := p.Present(ctx, fl)
	require.NoError(t, err)
}

/*
// redirect to light when has index and child later than limit
func TestGetChildren_RedirectsToLightChildLaterThanLimit(t *testing.T) {
	// p := createProc(t)
	assert.NoError(t, nil)
}

// redirect to heavy when has index and child earlier than limit
func TestGetChildren_RedirectsToHeavyChildEarlierThanLimit(t *testing.T) {
	// p := createProc(t)
	assert.NoError(t, nil)
}
*/
