/*
 *    Copyright 2018 INS Ecosystem
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

package ledger

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

// Store is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger struct {
	store       *storage.DB
	manager     *artifactmanager.LedgerArtifactManager
	coordinator *jetcoordinator.JetCoordinator
}

// GetManager returns artifact manager to work with.
func (l *Ledger) GetManager() core.ArtifactManager {
	return l.manager
}

// NewLedger creates new ledger instance.
func NewLedger(conf configuration.Ledger) (core.Ledger, error) {
	store, err := storage.NewStore(conf.DataDirectory, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DB creation failed")
	}
	manager, err := artifactmanager.NewArtifactManger(store)
	if err != nil {
		return nil, errors.Wrap(err, "artifact manager creation failed")
	}
	coordinator, err := jetcoordinator.NewJetCoordinator(store)
	if err != nil {
		return nil, errors.Wrap(err, "jet coordinator creation failed")
	}
	ledger := &Ledger{
		store:       store,
		manager:     manager,
		coordinator: coordinator,
	}
	return ledger, nil
}

// Start initializes external ledger dependencies.
func (l *Ledger) Start(c core.Components) error {
	// TODO: add links to network and maybe message router
	// mr := c["core.MessageRouter"].(core.MessageRouter)
	// l.messagerouter = mr
	return nil
}

// Stop stops Store gracefully.
func (l *Ledger) Stop() error {
	return l.store.Close()
}
