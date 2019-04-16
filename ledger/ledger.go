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
	"github.com/insolar/insolar/ledger/replication/light"
	"github.com/insolar/insolar/ledger/storage/pulse"
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
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
)

// GetLedgerComponents returns ledger components.
func GetLedgerComponents(conf configuration.Ledger, msgBus insolar.MessageBus, certificate insolar.Certificate) []interface{} {
	idLocker := storage.NewIDLocker()

	legacyDB, err := storage.NewDB(conf, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	db, err := store.NewBadgerDB(conf.Storage.DataDirectoryNewDB)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	jetStorage := jet.NewStore()
	nodeStorage := node.NewStorage()
	rsProvide := recentstorage.NewRecentStorageProvider()
	jetCoordinator := jetcoordinator.NewJetCoordinator(conf.LightChainLimit)

	var dropModifier drop.Modifier
	var dropAccessor drop.Accessor
	var dropCleaner drop.Cleaner

	var blobCleaner blob.Cleaner
	var blobModifier blob.Modifier
	var blobAccessor blob.Accessor
	var blobCollectionAccessor blob.CollectionAccessor

	var pulseAccessor pulse.Accessor
	var pulseAppender pulse.Appender
	var pulseCalculator pulse.Calculator
	var pulseShifter pulse.Shifter

	var recordModifier object.RecordModifier
	var recordAccessor object.RecordAccessor
	var recSyncAccessor object.RecordCollectionAccessor
	var recordCleaner object.RecordCleaner

	var indexStorage object.IndexStorage
	var collectionIndexAccessor object.IndexCollectionAccessor
	var indexCleaner object.IndexCleaner
	var lightIndexModifier object.ExtendedIndexModifier

	// Comparision with insolar.StaticRoleUnknown is a hack for genesis pulse (INS-1537)
	switch certificate.GetRole() {
	case insolar.StaticRoleHeavyMaterial:
		ps := pulse.NewDB(db)
		pulseAccessor = ps
		pulseAppender = ps
		pulseCalculator = ps

		dropDB := drop.NewDB(db)
		dropModifier = dropDB
		dropAccessor = dropDB

		// should be replaced with db
		blobDB := blob.NewDB(db)
		blobModifier = blobDB
		blobAccessor = blobDB

		records := object.NewRecordDB(db)
		recordModifier = records
		recordAccessor = records

		indexDB := object.NewIndexDB(db)
		indexStorage = indexDB
	default:
		ps := pulse.NewStorageMem()
		pulseAccessor = ps
		pulseAppender = ps
		pulseCalculator = ps
		pulseShifter = ps

		dropDB := drop.NewStorageMemory()
		dropModifier = dropDB
		dropAccessor = dropDB
		dropCleaner = dropDB

		blobDB := blob.NewStorageMemory()
		blobModifier = blobDB
		blobAccessor = blobDB
		blobCleaner = blobDB
		blobCollectionAccessor = blobDB

		records := object.NewRecordMemory()
		recordModifier = records
		recordAccessor = records
		recSyncAccessor = records
		recordCleaner = records

		indexDB := object.NewIndexMemory()
		indexStorage = indexDB
		indexCleaner = indexDB
		lightIndexModifier = indexDB
		collectionIndexAccessor = indexDB
	}

	dataGatherer := light.NewDataGatherer(dropAccessor, blobCollectionAccessor, recSyncAccessor, collectionIndexAccessor)
	lightCleaner := light.NewCleaner(jetStorage, nodeStorage, dropCleaner, blobCleaner, recordCleaner, indexCleaner, rsProvide, pulseShifter, jet.NewCalculator(jetCoordinator, jetStorage), pulseCalculator, conf.LightChainLimit)

	lSyncer := light.NewToHeavySyncer(jet.NewCalculator(jetCoordinator, jetStorage), dataGatherer, lightCleaner, msgBus, conf.LightToHeavySync, pulseCalculator)

	pm := pulsemanager.NewPulseManager(conf, dropCleaner, blobCleaner, blobCollectionAccessor, pulseShifter, recordCleaner, recSyncAccessor, collectionIndexAccessor, indexCleaner, lSyncer)

	components := []interface{}{
		legacyDB,
		db,
		idLocker,
		dropModifier,
		dropAccessor,
		blobModifier,
		blobAccessor,
		pulseAccessor,
		pulseAppender,
		pulseCalculator,
		recordModifier,
		recordAccessor,
		indexStorage,
		jetStorage,
		nodeStorage,
		storage.NewReplicaStorage(),
		rsProvide,
		artifactmanager.NewHotDataWaiterConcrete(),
		jetCoordinator,
		heavyserver.NewSync(legacyDB, recordModifier),
	}

	switch certificate.GetRole() {
	case insolar.StaticRoleUnknown, insolar.StaticRoleLightMaterial:
		h := artifactmanager.NewMessageHandler(indexStorage, lightIndexModifier, &conf)
		pm.MessageHandler = h
		components = append(components, h)
		components = append(components, pm)
	case insolar.StaticRoleHeavyMaterial:
		components = append(components, pm)
		components = append(components, heavy.Components()...)
	}

	return components
}
