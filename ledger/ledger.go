/*
 *    Copyright 2019 Insolar Technologies
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

	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/exporter"
	"github.com/insolar/insolar/ledger/heavyserver"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/log"
)

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger struct {
	db              storage.DBContext
	ArtifactManager core.ArtifactManager `inject:""`
	PulseManager    core.PulseManager    `inject:""`
	JetCoordinator  core.JetCoordinator  `inject:""`
}

// Deprecated: remove after deleting TmpLedger
// GetPulseManager returns PulseManager.
func (l *Ledger) GetPulseManager() core.PulseManager {
	log.Warn("GetPulseManager is deprecated. Use component injection.")
	return l.PulseManager
}

// Deprecated: remove after deleting TmpLedger
// GetJetCoordinator returns JetCoordinator.
func (l *Ledger) GetJetCoordinator() core.JetCoordinator {
	log.Warn("GetJetCoordinator is deprecated. Use component injection.")
	return l.JetCoordinator
}

// Deprecated: remove after deleting TmpLedger
// GetArtifactManager returns artifact manager to work with.
func (l *Ledger) GetArtifactManager() core.ArtifactManager {
	log.Warn("GetArtifactManager is deprecated. Use component injection.")
	return l.ArtifactManager
}

// NewTestLedger is the util function for creation of Ledger with provided
// private members (suitable for tests).
func NewTestLedger(
	db storage.DBContext,
	am *artifactmanager.LedgerArtifactManager,
	pm *pulsemanager.PulseManager,
	jc core.JetCoordinator,
) *Ledger {
	return &Ledger{
		db:              db,
		ArtifactManager: am,
		PulseManager:    pm,
		JetCoordinator:  jc,
	}
}

// GetLedgerComponents returns ledger components.
func GetLedgerComponents(conf configuration.Ledger, certificate core.Certificate) []interface{} {
	db, err := storage.NewDB(conf, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	var pulseTracker storage.PulseTracker
	// TODO: @imarkin 18.02.18 - Comparision with core.StaticRoleUnknown is a hack for genesis pulse (INS-1537)
	switch certificate.GetRole() {
	case core.StaticRoleUnknown, core.StaticRoleHeavyMaterial:
		pulseTracker = storage.NewPulseTracker()
	default:
		pulseTracker = storage.NewPulseTrackerMemory()
	}

	return []interface{}{
		db,
		storage.NewCleaner(),
		pulseTracker,
		storage.NewPulseStorage(),
		storage.NewJetStorage(),
		storage.NewDropStorage(conf.JetSizesHistoryDepth),
		nodes.NewStorage(),
		storage.NewObjectStorage(),
		storage.NewReplicaStorage(),
		storage.NewGenesisInitializer(),
		recentstorage.NewRecentStorageProvider(conf.RecentStorage.DefaultTTL),
		artifactmanager.NewHotDataWaiterConcrete(),
		artifactmanager.NewArtifactManger(),
		jetcoordinator.NewJetCoordinator(conf.LightChainLimit),
		pulsemanager.NewPulseManager(conf),
		artifactmanager.NewMessageHandler(&conf, certificate),
		heavyserver.NewSync(db),
		exporter.NewExporter(conf.Exporter),
	}
}

// Start stub.
func (l *Ledger) Start(ctx context.Context) error {
	return nil
}

// Stop stops Ledger gracefully.
func (l *Ledger) Stop(ctx context.Context) error {
	return nil
}
