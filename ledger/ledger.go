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
	"context"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/blockexplorer"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/localstorage"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger struct {
	db      *storage.DB
	am      *artifactmanager.LedgerArtifactManager
	pm      *pulsemanager.PulseManager
	jc      *jetcoordinator.JetCoordinator
	handler *artifactmanager.MessageHandler
	ls      *localstorage.LocalStorage
	be      *blockexplorer.BlockExplorerManager
}

// GetPulseManager returns PulseManager.
func (l *Ledger) GetPulseManager() core.PulseManager {
	return l.pm
}

// GetJetCoordinator returns JetCoordinator.
func (l *Ledger) GetJetCoordinator() core.JetCoordinator {
	return l.jc
}

// GetArtifactManager returns artifact manager to work with.
func (l *Ledger) GetArtifactManager() core.ArtifactManager {
	return l.am
}

// GetLocalStorage returns local storage to work with.
func (l *Ledger) GetLocalStorage() core.LocalStorage {
	return l.ls
}

// GetBlockExplorer returns block explorer to work with.
func (l *Ledger) GetBlockExplorer() core.BlockExplorer {
	return l.be
}

// NewLedger creates new ledger instance.
func NewLedger(ctx context.Context, conf configuration.Ledger) (*Ledger, error) {
	var err error
	db, err := storage.NewDB(conf, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DB creation failed")
	}
	am, err := artifactmanager.NewArtifactManger(db)
	if err != nil {
		return nil, errors.Wrap(err, "artifact manager creation failed")
	}
	jc, err := jetcoordinator.NewJetCoordinator(db, conf.JetCoordinator)
	if err != nil {
		return nil, errors.Wrap(err, "jet coordinator creation failed")
	}
	pm, err := pulsemanager.NewPulseManager(db)
	if err != nil {
		return nil, errors.Wrap(err, "pulse manager creation failed")
	}
	handler, err := artifactmanager.NewMessageHandler(db)
	if err != nil {
		return nil, err
	}
	ls, err := localstorage.NewLocalStorage(db)
	if err != nil {
		return nil, err
	}
	be, err := blockexplorer.NewBlockExplorer(db)
	if err != nil {
		return nil, errors.Wrap(err, "block explorer creation failed")
	}

	err = db.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	ledger := Ledger{
		db:      db,
		am:      am,
		pm:      pm,
		jc:      jc,
		handler: handler,
		ls:      ls,
		be:      be,
	}

	return &ledger, nil
}

// NewTestLedger is the util function for creation of Ledger with provided
// private members (suitable for tests).
func NewTestLedger(
	db *storage.DB,
	am *artifactmanager.LedgerArtifactManager,
	pm *pulsemanager.PulseManager,
	jc *jetcoordinator.JetCoordinator,
	amh *artifactmanager.MessageHandler,
	ls *localstorage.LocalStorage,
	be *blockexplorer.BlockExplorerManager,
) *Ledger {
	return &Ledger{
		db:      db,
		am:      am,
		pm:      pm,
		jc:      jc,
		handler: amh,
		ls:      ls,
		be:      be,
	}
}

// Start initializes external ledger dependencies.
func (l *Ledger) Start(ctx context.Context, c core.Components) error {
	var err error
	if err = l.am.Link(c); err != nil {
		return err
	}
	if err = l.be.Link(c); err != nil {
		return err
	}
	if err = l.pm.Link(c); err != nil {
		return err
	}
	if err = l.handler.Link(c); err != nil {
		return err
	}
	if err = l.jc.Link(c); err != nil {
		return err
	}

	return nil
}

// Stop stops Ledger gracefully.
func (l *Ledger) Stop(ctx context.Context) error {
	return l.db.Close()
}
