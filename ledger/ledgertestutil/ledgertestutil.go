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

package ledgertestutil

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

type messageBusMock struct {
	handler *artifactmanager.MessageHandler
}

func (m *messageBusMock) Start(components core.Components) error {
	panic("implement me")
}

func (m *messageBusMock) Stop() error {
	panic("implement me")
}

func (m *messageBusMock) Send(e core.Message) (core.Reply, error) {
	return m.handler.Handle(e)
}

func (m *messageBusMock) SendAsync(e core.Message) {
	m.handler.Handle(e) // nolint
}

// TmpLedger crteates ledger on top of temporary database.
// Returns *ledger.Ledger andh cleanup function.
func TmpLedger(t testing.TB, dir string) (*ledger.Ledger, func()) {
	var err error
	// Init subcomponents.
	conf := configuration.NewLedger()
	db, dbcancel := storagetest.TmpDB(t, dir)
	handler, err := artifactmanager.NewMessageHandler(db)
	assert.NoError(t, err)
	am, err := artifactmanager.NewArtifactManger(db)
	assert.NoError(t, err)
	jc, err := jetcoordinator.NewJetCoordinator(db, conf.JetCoordinator)
	assert.NoError(t, err)
	pm, err := pulsemanager.NewPulseManager(db, jc)
	assert.NoError(t, err)

	// Bootstrap
	err = db.Bootstrap()
	assert.NoError(t, err)

	// Init components.
	eb := messageBusMock{handler: handler}
	components := core.Components{MessageBus: &eb}

	// Create ledger.
	l := ledger.NewTestLedger(db, am, pm, jc, handler)
	err = l.Start(components)
	assert.NoError(t, err)

	return l, dbcancel
}
