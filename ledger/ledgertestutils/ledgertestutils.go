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
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/localstorage"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

// TmpLedger crteates ledger on top of temporary database.
// Returns *ledger.Ledger andh cleanup function.
func TmpLedger(t testing.TB, dir string, c core.Components) (*ledger.Ledger, func()) {
	var err error
	// Init subcomponents.
	ctx := inslogger.TestContext(t.(*testing.T))
	conf := configuration.NewLedger()
	db, dbcancel := storagetest.TmpDB(ctx, t, dir)

	handler, err := artifactmanager.NewMessageHandler(db, storage.NewRecentObjectsIndex(0))
	assert.NoError(t, err)
	am, err := artifactmanager.NewArtifactManger(db)
	assert.NoError(t, err)
	jc, err := jetcoordinator.NewJetCoordinator(db, conf.JetCoordinator)
	assert.NoError(t, err)
	pm, err := pulsemanager.NewPulseManager(db)
	assert.NoError(t, err)
	ls, err := localstorage.NewLocalStorage(db)
	assert.NoError(t, err)

	// Init components.
	if c.MessageBus == nil {
		c.MessageBus = testmessagebus.NewTestMessageBus()
	}
	if c.NodeNetwork == nil {
		c.NodeNetwork = nodenetwork.NewNodeKeeper(nodenetwork.NewNode(core.RecordRef{}, nil, nil, 0, "", ""))
	}

	// Create ledger.
	l := ledger.NewTestLedger(db, am, pm, jc, handler, ls)
	err = l.Start(ctx, c)
	assert.NoError(t, err)

	return l, dbcancel
}
