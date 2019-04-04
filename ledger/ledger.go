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

package ledger

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/heavy"
	"github.com/insolar/insolar/ledger/heavyserver"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/genesis"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
)

// GetLedgerComponents returns ledger components.
func GetLedgerComponents(conf configuration.Ledger, certificate insolar.Certificate) []interface{} {
	idLocker := storage.NewIDLocker()

	legacyDB, err := storage.NewDB(conf, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	db, err := store.NewBadgerDB(conf.Storage.DataDirectoryNewDB)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	var pulseTracker storage.PulseTracker

	var dropModifier drop.Modifier
	var dropAccessor drop.Accessor
	var dropCleaner drop.Cleaner

	var blobCleaner blob.Cleaner
	var blobModifier blob.Modifier
	var blobAccessor blob.Accessor
	var blobSyncAccessor blob.CollectionAccessor

	var recordModifier object.RecordModifier
	var recordAccessor object.RecordAccessor
	var recSyncAccessor object.RecordCollectionAccessor
	var recordCleaner object.RecordCleaner
	// Comparision with insolar.StaticRoleUnknown is a hack for genesis pulse (INS-1537)
	switch certificate.GetRole() {
	case insolar.StaticRoleUnknown, insolar.StaticRoleHeavyMaterial:
		pulseTracker = storage.NewPulseTracker()

		dropDB := drop.NewStorageDB(db)
		dropModifier = dropDB
		dropAccessor = dropDB

		// should be replaced with db
		blobDB := blob.NewStorageDB(db)
		blobModifier = blobDB
		blobAccessor = blobDB

		records := object.NewRecordDB(db)
		recordModifier = records
		recordAccessor = records
	default:
		pulseTracker = storage.NewPulseTrackerMemory()

		dropDB := drop.NewStorageMemory()
		dropModifier = dropDB
		dropAccessor = dropDB
		dropCleaner = dropDB

		blobDB := blob.NewStorageMemory()
		blobModifier = blobDB
		blobAccessor = blobDB
		blobCleaner = blobDB
		blobSyncAccessor = blobDB

		records := object.NewRecordMemory()
		recordModifier = records
		recordAccessor = records
		recSyncAccessor = records
		recordCleaner = records
	}

	components := []interface{}{
		legacyDB,
		db,
		idLocker,
		dropModifier,
		dropAccessor,
		blobModifier,
		blobAccessor,
		recordModifier,
		recordAccessor,
		storage.NewCleaner(),
		pulseTracker,
		storage.NewPulseStorage(),
		jet.NewStore(),
		node.NewStorage(),
		storage.NewObjectStorage(),
		storage.NewReplicaStorage(),
		genesis.NewGenesisInitializer(),
		recentstorage.NewRecentStorageProvider(conf.RecentStorage.DefaultTTL),
		artifactmanager.NewHotDataWaiterConcrete(),
		jetcoordinator.NewJetCoordinator(conf.LightChainLimit),
		heavyserver.NewSync(legacyDB, recordModifier),
	}

	switch certificate.GetRole() {
	case insolar.StaticRoleUnknown, insolar.StaticRoleLightMaterial:
		components = append(components, artifactmanager.NewMessageHandler(&conf))
		components = append(components, pulsemanager.NewPulseManager(conf, dropCleaner, blobCleaner, blobSyncAccessor, recordCleaner, recSyncAccessor))
	case insolar.StaticRoleHeavyMaterial:
		components = append(components, pulsemanager.NewPulseManager(conf, dropCleaner, blobCleaner, blobSyncAccessor, recordCleaner, recSyncAccessor))
		components = append(components, heavy.Components()...)
	}

	return components
}
