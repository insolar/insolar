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
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/recentstorage"
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
	ArtifactManager Client
	PulseManager    insolar.PulseManager `inject:""`
	JetCoordinator  jet.Coordinator      `inject:""`
}

// Deprecated: remove after deleting TmpLedger
// GetPulseManager returns PulseManager.
func (l *TMPLedger) GetPulseManager() insolar.PulseManager {
	log.Warn("GetPulseManager is deprecated. Use component injection.")
	return l.PulseManager
}

// Deprecated: remove after deleting TmpLedger
// GetJetCoordinator returns Coordinator.
func (l *TMPLedger) GetJetCoordinator() jet.Coordinator {
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
	am Client,
	pm *pulsemanager.PulseManager,
	jc jet.Coordinator,
) *TMPLedger {
	return &TMPLedger{
		ArtifactManager: am,
		PulseManager:    pm,
		JetCoordinator:  jc,
	}
}

// TmpLedger creates ledger on top of temporary database.
// Returns *ledger.Ledger and cleanup function.
// DEPRECATED
func TmpLedger(t *testing.T, dir string, c insolar.Components) *TMPLedger {
	log.Warn("TmpLedger is deprecated. Use mocks.")

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	mc := minimock.NewController(t)
	ps := pulse.NewStorageMem()
	is := object.NewIndexMemory()

	// Init subcomponents.
	ctx := inslogger.TestContext(t)
	conf := configuration.NewLedger()
	recordStorage := object.NewRecordMemory()
	memoryMockDB := store.NewMemoryMockDB()

	cm := &component.Manager{}
	js := jet.NewStore()
	ns := node.NewStorage()
	ds := drop.NewDB(memoryMockDB)
	bs := blob.NewDB(memoryMockDB)

	genesisBaseRecord := &genesis.BaseRecord{
		DB:             memoryMockDB,
		DropModifier:   ds,
		PulseAppender:  ps,
		PulseAccessor:  ps,
		RecordModifier: recordStorage,
		IndexModifier:  is,
	}
	_, err := genesisBaseRecord.CreateIfNeeded(ctx)
	if err != nil {
		t.Error(err, "failed to create base genesis record")
	}

	recordAccessor := recordStorage
	recordModifier := recordStorage

	am := NewClient()
	am.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

	pm := pulsemanager.NewPulseManager()
	jc := jet.NewCoordinatorMock(mc)
	jc.IsAuthorizedMock.Return(true, nil)
	jc.LightExecutorForJetMock.Return(&insolar.Reference{}, nil)
	jc.HeavyMock.Return(&insolar.Reference{}, nil)
	jc.MeMock.Return(insolar.Reference{})
	jc.IsBeyondLimitMock.Return(false, nil)

	// Init components.
	if c.MessageBus == nil {
		mb := testmessagebus.NewTestMessageBus(t)
		mb.PulseAccessor = ps
		c.MessageBus = mb
	} else {
		switch mb := c.MessageBus.(type) {
		case *messagebus.MessageBus:
			mb.PulseAccessor = ps
		case *testmessagebus.TestMessageBus:
			mb.PulseAccessor = ps
		default:
			panic("unknown message bus")
		}
	}
	if c.NodeNetwork == nil {
		c.NodeNetwork = nodenetwork.NewNodeKeeper(networknode.NewNode(insolar.Reference{}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:5432", ""))
	}

	handler := artifactmanager.NewMessageHandler(is, is, &conf)
	handler.JetStorage = js
	handler.Nodes = ns
	handler.IndexStorage = is
	handler.DropModifier = ds
	handler.BlobModifier = bs
	handler.BlobAccessor = bs
	handler.Blobs = bs
	handler.RecordModifier = recordModifier
	handler.RecordAccessor = recordAccessor

	idLockerMock := object.NewIDLockerMock(t)
	idLockerMock.LockMock.Return()
	idLockerMock.UnlockMock.Return()

	handler.IDLocker = idLockerMock

	handler.PlatformCryptographyScheme = pcs
	handler.JetCoordinator = jc

	am.DefaultBus = c.MessageBus
	am.JetCoordinator = jc

	cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		memoryMockDB,
		js,
		ns,
		is,
		ps,
		ps,
		ds,
		am,
		recordAccessor,
		recordModifier,
	)

	err = cm.Init(ctx)
	if err != nil {
		t.Error("ComponentManager init failed", err)
	}
	err = cm.Start(ctx)
	if err != nil {
		t.Error("ComponentManager start failed", err)
	}

	gilMock := testutils.NewGlobalInsolarLockMock(t)
	gilMock.AcquireFunc = func(context.Context) {}
	gilMock.ReleaseFunc = func(context.Context) {}

	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func(context.Context, insolar.PulseNumber) error { return nil }

	handler.Bus = c.MessageBus

	pm.NodeNet = c.NodeNetwork
	pm.GIL = gilMock
	pm.Bus = c.MessageBus
	pm.LR = c.LogicRunner
	pm.ActiveListSwapper = alsMock
	// pm.PulseStorage = ps
	pm.Nodes = ns
	pm.NodeSetter = ns
	pm.JetModifier = js

	pm.PulseAccessor = ps
	pm.PulseAppender = ps

	hdw := hot.NewChannelWaiter()

	handler.HotDataWaiter = hdw
	handler.JetReleaser = hdw

	pendingMock := recentstorage.NewPendingStorageMock(t)

	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetPendingStorageMock.Return(pendingMock)
	provideMock.CountMock.Return(0)

	handler.RecentStorageProvider = provideMock

	err = handler.Init(ctx)
	if err != nil {
		panic(err)
	}

	// Create ledger.
	l := NewTestLedger(am, pm, jc)

	return l
}
