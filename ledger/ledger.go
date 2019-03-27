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
	"github.com/insolar/insolar/ledger/artifactmanager"
	"github.com/insolar/insolar/ledger/heavy"
	"github.com/insolar/insolar/ledger/heavyserver"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/genesis"
	"github.com/insolar/insolar/ledger/storage/node"
)

// GetLedgerComponents returns ledger components.
func GetLedgerComponents(conf configuration.Ledger, certificate insolar.Certificate) []interface{} {
	idLocker := storage.NewIDLocker()

	store, err := storage.NewDB(conf, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	dbBadger, err := db.NewBadgerDB(conf)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	var pulseTracker storage.PulseTracker
	var dropModifier drop.Modifier
	var dropAccessor drop.Accessor
	// TODO: @imarkin 18.02.18 - Comparision with insolar.StaticRoleUnknown is a hack for genesis pulse (INS-1537)
	switch certificate.GetRole() {
	case insolar.StaticRoleUnknown, insolar.StaticRoleHeavyMaterial:
		pulseTracker = storage.NewPulseTracker()

		dropDB := drop.NewStorageDB(dbBadger)
		dropModifier = dropDB
		dropAccessor = dropDB
	default:
		pulseTracker = storage.NewPulseTrackerMemory()

		dropDB := drop.NewStorageMemory()
		dropModifier = dropDB
		dropAccessor = dropDB
	}

	components := []interface{}{
		store,
		dbBadger,
		idLocker,
		dropModifier,
		dropAccessor,
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
		pulsemanager.NewPulseManager(conf),
		heavyserver.NewSync(store),
	}

	switch certificate.GetRole() {
	case insolar.StaticRoleUnknown, insolar.StaticRoleLightMaterial:
		components = append(components, artifactmanager.NewMessageHandler(&conf))
	case insolar.StaticRoleHeavyMaterial:
		components = append(components, heavy.Components()...)
	}

	return components
}
