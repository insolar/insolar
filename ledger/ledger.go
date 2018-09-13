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

package ledger

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger struct {
	db          *storage.DB
	manager     *artifactmanager.LedgerArtifactManager
	coordinator *jetcoordinator.JetCoordinator
}

func (l *Ledger) GetPulseManager() core.PulseManager {
	panic("implement me")
}

func (l *Ledger) GetJetCoordinator() core.JetCoordinator {
	return l.coordinator
}

// GetArtifactManager returns artifact manager to work with.
func (l *Ledger) GetArtifactManager() core.ArtifactManager {
	return l.manager
}

// NewLedger creates new ledger instance.
func NewLedger(conf configuration.Ledger) (*Ledger, error) {
	var err error
	db, err := storage.NewDB(conf, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DB creation failed")
	}
	return NewLedgerWithDB(db)
}

// NewLedgerWithDB creates new ledger with preconfigured storage.DB instance.
func NewLedgerWithDB(db *storage.DB) (*Ledger, error) {
	manager, err := artifactmanager.NewArtifactManger(db)
	if err != nil {
		return nil, errors.Wrap(err, "artifact manager creation failed")
	}
	coordinator, err := jetcoordinator.NewJetCoordinator(db)
	if err != nil {
		return nil, errors.Wrap(err, "jet coordinator creation failed")
	}

	err = db.Bootstrap()
	if err != nil {
		return nil, err
	}

	return &Ledger{
		db:          db,
		manager:     manager,
		coordinator: coordinator,
	}, nil
}

// Start initializes external ledger dependencies.
func (l *Ledger) Start(c core.Components) error {
	// TODO: add links to network and maybe message router
	// mr := c["core.MessageRouter"].(core.MessageRouter)
	// l.messagerouter = mr
	return nil
}

// Stop stops Ledger gracefully.
func (l *Ledger) Stop() error {
	return l.db.Close()
}
