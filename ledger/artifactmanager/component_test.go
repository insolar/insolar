package artifactmanager

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerArtifactManager_PendingRequest(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	cs := testutils.NewPlatformCryptographyScheme()
	mb := testmessagebus.NewTestMessageBus(t)
	jc := testutils.NewJetCoordinatorMock(mc)
	jc.AmIMock.Return(true, nil)
	am := NewArtifactManger(db)
	am.PlatformCryptographyScheme = cs
	am.DefaultBus = mb
	provider := storage.NewRecentStorageProvider(0)
	handler := NewMessageHandler(db, &configuration.Ledger{})
	handler.Bus = mb
	handler.JetCoordinator = jc
	handler.RecentStorageProvider = provider
	err := handler.Init(ctx)
	require.NoError(t, err)
	objRef := *genRandomRef(0)

	// Register request
	reqID, err := am.RegisterRequest(ctx, objRef, &message.Parcel{Msg: &message.CallMethod{}})
	require.NoError(t, err)

	// Should have pending request.
	requests, err := am.GetPendingRequests(ctx, objRef)
	require.NoError(t, err)
	assert.Equal(t, 1, len(requests))

	// Register result.
	reqRef := *core.NewRecordRef(core.DomainID, *reqID)
	_, err = am.RegisterResult(ctx, objRef, reqRef, nil)
	require.NoError(t, err)

	// Should not have pending request.
	requests, err = am.GetPendingRequests(ctx, objRef)
	require.NoError(t, err)
	assert.Equal(t, 0, len(requests))
}
