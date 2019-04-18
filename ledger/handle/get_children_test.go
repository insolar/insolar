package handle

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestGetChildren(t *testing.T) {
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

	p.Dep.RecentStorageProvider = recentstorage.NewProviderMock(t)
	p.Dep.IDLocker = storage.NewIDLockerMock(t)
	p.Dep.IndexAccessor = object.NewIndexAccessorMock(t)
	p.Dep.JetCoordinator = testutils.NewJetCoordinatorMock(t)
	p.Dep.JetStorage = jet.NewStorageMock(t)
	p.Dep.DelegationTokenFactory = testutils.NewDelegationTokenFactoryMock(t)
	p.Dep.RecordAccessor = object.NewRecordAccessorMock(t)
	p.Dep.TreeUpdater = jet.NewTreeUpdaterMock(t)
	p.Dep.IndexSaver = object.NewIndexSaverMock(t)

	assert.NoError(t, nil)
}
