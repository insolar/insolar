/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package ledgertestutils

import (
	"context"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/localstorage"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

// TmpLedger crteates ledger on top of temporary database.
// Returns *ledger.Ledger andh cleanup function.
// FIXME: THIS METHOD IS DEPRECATED. USE MOCKS.
func TmpLedger(t *testing.T, dir string, c core.Components) (*ledger.Ledger, func()) {
	log.Warn("TmpLedger is deprecated. Use mocks.")

	pcs := platformpolicy.NewPlatformCryptographyScheme()

	// Init subcomponents.
	ctx := inslogger.TestContext(t)
	conf := configuration.NewLedger()
	db, dbcancel := storagetest.TmpDB(ctx, t, storagetest.Dir(dir))

	handler := artifactmanager.NewMessageHandler(db, nil)
	handler.PlatformCryptographyScheme = pcs

	am := artifactmanager.NewArtifactManger(db)
	am.PlatformCryptographyScheme = pcs
	jc := jetcoordinator.NewJetCoordinator(db, conf.JetCoordinator)
	jc.PlatformCryptographyScheme = pcs
	conf.PulseManager.HeavySyncEnabled = false
	pm := pulsemanager.NewPulseManager(db, conf)
	ls := localstorage.NewLocalStorage(db)

	// Init components.
	if c.MessageBus == nil {
		c.MessageBus = testmessagebus.NewTestMessageBus(t)
	}
	if c.NodeNetwork == nil {
		c.NodeNetwork = nodenetwork.NewNodeKeeper(nodenetwork.NewNode(core.RecordRef{}, core.StaticRoleUnknown, nil, "", ""))
	}

	gilMock := testutils.NewGlobalInsolarLockMock(t)
	gilMock.AcquireFunc = func(context.Context) {}
	gilMock.ReleaseFunc = func(context.Context) {}

	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func() {}

	handler.Bus = c.MessageBus
	am.DefaultBus = c.MessageBus
	jc.NodeNet = c.NodeNetwork
	pm.NodeNet = c.NodeNetwork
	pm.GIL = gilMock
	pm.Bus = c.MessageBus
	pm.LR = c.LogicRunner
	pm.ActiveListSwapper = alsMock

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()
	handler.Recent = recentStorageMock

	err := handler.Init(ctx)
	if err != nil {
		panic(err)
	}

	// Create ledger.
	l := ledger.NewTestLedger(db, am, pm, jc, ls)

	return l, dbcancel
}
