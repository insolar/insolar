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

package ledgertestutils

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/internal/jet"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/genesis"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/require"
)

// TmpLedger crteates ledger on top of temporary database.
// Returns *ledger.Ledger and cleanup function.
// FIXME: THIS METHOD IS DEPRECATED. USE MOCKS.
func TmpLedger(t *testing.T, dir string, handlersRole insolar.StaticRole, c insolar.Components, closeJets bool) (*ledger.Ledger, storage.DBContext, func()) {
	log.Warn("TmpLedger is deprecated. Use mocks.")

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	mc := minimock.NewController(t)

	// Init subcomponents.
	ctx := inslogger.TestContext(t)
	conf := configuration.NewLedger()
	tmpDB, dbcancel := storagetest.TmpDB(ctx, t, storagetest.Dir(dir))

	cm := &component.Manager{}
	gi := genesis.NewGenesisInitializer()
	pt := storage.NewPulseTracker()
	ps := storage.NewPulseStorage()
	js := jet.NewStore()
	os := storage.NewObjectStorage()
	ns := node.NewStorage()
	ds := drop.NewStorageDB()
	rs := storage.NewReplicaStorage()
	cl := storage.NewCleaner()

	am := artifactmanager.NewArtifactManger()
	am.PlatformCryptographyScheme = pcs

	conf.PulseManager.HeavySyncEnabled = false
	pm := pulsemanager.NewPulseManager(conf)
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
		c.NodeNetwork = nodenetwork.NewNodeKeeper(nodenetwork.NewNode(insolar.Reference{}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:5432", ""))
	}

	certificate := testutils.NewCertificateMock(t)
	certificate.GetRoleMock.Return(handlersRole)

	handler := artifactmanager.NewMessageHandler(&conf, certificate)
	handler.PulseTracker = pt
	handler.JetStorage = js
	handler.Nodes = ns
	handler.DBContext = tmpDB
	handler.ObjectStorage = os
	handler.DropModifier = ds

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
		db.NewMemoryMockDB(),
		js,
		os,
		ns,
		pt,
		ps,
		ds,
		gi,
		am,
		rs,
		cl,
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
	pm.JetAccessor = js
	pm.JetModifier = js
	pm.DropModifier = ds
	pm.DropAccessor = ds
	pm.DropCleaner = ds
	pm.ObjectStorage = os
	pm.Nodes = ns
	pm.NodeSetter = ns
	pm.PulseTracker = pt
	pm.ReplicaStorage = rs
	pm.StorageCleaner = cl

	hdw := artifactmanager.NewHotDataWaiterConcrete()

	pm.HotDataWaiter = hdw
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

	if closeJets {
		err := pm.HotDataWaiter.Unlock(ctx, insolar.ID(*insolar.NewJetID(0, nil)))
		require.NoError(t, err)
	}

	// Create ledger.
	l := ledger.NewTestLedger(tmpDB, am, pm, jc)

	return l, tmpDB, dbcancel
}
