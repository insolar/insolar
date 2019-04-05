//
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
//

package artifacts

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/messagebus"
	networknode "github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

// TMPLedger
// DEPRECATED
type TMPLedger struct {
	db              storage.DBContext
	ArtifactManager Client
	PulseManager    insolar.PulseManager   `inject:""`
	JetCoordinator  insolar.JetCoordinator `inject:""`
}

// Deprecated: remove after deleting TmpLedger
// GetPulseManager returns PulseManager.
func (l *TMPLedger) GetPulseManager() insolar.PulseManager {
	log.Warn("GetPulseManager is deprecated. Use component injection.")
	return l.PulseManager
}

// Deprecated: remove after deleting TmpLedger
// GetJetCoordinator returns JetCoordinator.
func (l *TMPLedger) GetJetCoordinator() insolar.JetCoordinator {
	log.Warn("GetJetCoordinator is deprecated. Use component injection.")
	return l.JetCoordinator
}

// Deprecated: remove after deleting TmpLedger
// GetArtifactManager returns artifact manager to work with.
func (l *TMPLedger) GetArtifactManager() Client {
	log.Warn("GetArtifactManager is deprecated. Use component injection.")
	return l.ArtifactManager
}

// NewTestLedger is the util function for creation of Ledger with provided
// private members (suitable for tests).
func NewTestLedger(
	db storage.DBContext,
	am Client,
	pm *pulsemanager.PulseManager,
	jc insolar.JetCoordinator,
) *TMPLedger {
	return &TMPLedger{
		db:              db,
		ArtifactManager: am,
		PulseManager:    pm,
		JetCoordinator:  jc,
	}
}

// TmpLedger creates ledger on top of temporary database.
// Returns *ledger.Ledger and cleanup function.
// DEPRECATED
func TmpLedger(t *testing.T, dir string, c insolar.Components) (*TMPLedger, storage.DBContext, func()) {
	log.Warn("TmpLedger is deprecated. Use mocks.")

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	mc := minimock.NewController(t)

	// Init subcomponents.
	ctx := inslogger.TestContext(t)
	conf := configuration.NewLedger()
	tmpDB, recMem, dbcancel := storagetest.TmpDB(ctx, t, storagetest.Dir(dir))
	memoryMockDB := store.NewMemoryMockDB()

	cm := &component.Manager{}
	pt := storage.NewPulseTracker()
	ps := storage.NewPulseStorage()
	js := jet.NewStore()
	os := storage.NewObjectStorage()
	ns := node.NewStorage()
	ds := drop.NewDB(memoryMockDB)
	bs := blob.NewDB(memoryMockDB)
	rs := storage.NewReplicaStorage()
	cl := storage.NewCleaner()

	recordAccessor := recMem
	recordModifier := recMem

	am := NewClient()
	am.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

	conf.PulseManager.HeavySyncEnabled = false
	pm := pulsemanager.NewPulseManager()
	jc := testutils.NewJetCoordinatorMock(mc)
	jc.IsAuthorizedMock.Return(true, nil)
	jc.LightExecutorForJetMock.Return(&insolar.Reference{}, nil)
	jc.HeavyMock.Return(&insolar.Reference{}, nil)
	jc.MeMock.Return(insolar.Reference{})
	jc.IsBeyondLimitMock.Return(false, nil)

	// Init components.
	if c.MessageBus == nil {
		mb := testmessagebus.NewTestMessageBus(t)
		mb.PulseStorage = ps
		c.MessageBus = mb
	} else {
		switch mb := c.MessageBus.(type) {
		case *messagebus.MessageBus:
			mb.PulseStorage = ps
		case *testmessagebus.TestMessageBus:
			mb.PulseStorage = ps
		default:
			panic("unknown message bus")
		}
	}
	if c.NodeNetwork == nil {
		c.NodeNetwork = nodenetwork.NewNodeKeeper(networknode.NewNode(insolar.Reference{}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:5432", ""))
	}

	handler := artifactmanager.NewMessageHandler(&conf)
	handler.PulseTracker = pt
	handler.JetStorage = js
	handler.Nodes = ns
	handler.DBContext = tmpDB
	handler.ObjectStorage = os
	handler.DropModifier = ds
	handler.BlobModifier = bs
	handler.BlobAccessor = bs
	handler.RecordModifier = recordModifier
	handler.RecordAccessor = recordAccessor

	idLockerMock := storage.NewIDLockerMock(t)
	idLockerMock.LockMock.Return()
	idLockerMock.UnlockMock.Return()

	handler.IDLocker = idLockerMock

	handler.PlatformCryptographyScheme = pcs
	handler.JetCoordinator = jc

	am.DefaultBus = c.MessageBus
	am.JetCoordinator = jc

	cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		tmpDB,
		memoryMockDB,
		js,
		os,
		ns,
		pt,
		ps,
		ds,
		am,
		rs,
		cl,
		recordAccessor,
		recordModifier,
	)

	err := cm.Init(ctx)
	if err != nil {
		t.Error("ComponentManager init failed", err)
	}
	err = cm.Start(ctx)
	if err != nil {
		t.Error("ComponentManager start failed", err)
	}

	pulse, err := pt.GetLatestPulse(ctx)
	require.NoError(t, err)
	ps.Set(&pulse.Pulse)

	gilMock := testutils.NewGlobalInsolarLockMock(t)
	gilMock.AcquireFunc = func(context.Context) {}
	gilMock.ReleaseFunc = func(context.Context) {}

	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func(context.Context) error { return nil }

	handler.Bus = c.MessageBus

	pm.NodeNet = c.NodeNetwork
	pm.GIL = gilMock
	pm.Bus = c.MessageBus
	pm.LR = c.LogicRunner
	pm.ActiveListSwapper = alsMock
	pm.PulseStorage = ps
	pm.Nodes = ns
	pm.NodeSetter = ns
	pm.PulseTracker = pt
	pm.JetModifier = js

	hdw := artifactmanager.NewHotDataWaiterConcrete()

	handler.HotDataWaiter = hdw

	indexMock := recentstorage.NewRecentIndexStorageMock(t)
	pendingMock := recentstorage.NewPendingStorageMock(t)

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)
	provideMock.CountMock.Return(0)

	handler.RecentStorageProvider = provideMock

	err = handler.Init(ctx)
	if err != nil {
		panic(err)
	}

	// Create ledger.
	l := NewTestLedger(tmpDB, am, pm, jc)

	return l, tmpDB, dbcancel
}
