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

package rootdomain

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationdaemon"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationshard"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/pkshard"
)

const (
	GenesisPrototypeSuffix = "_proto"
)

func init() {
	for _, el := range insolar.GenesisNameMigrationDaemonMembers {
		genesisrefs.PredefinedPrototypes[el+GenesisPrototypeSuffix] = *member.PrototypeReference
	}

	for _, el := range insolar.GenesisNameMigrationDaemons {
		genesisrefs.PredefinedPrototypes[el+GenesisPrototypeSuffix] = *migrationdaemon.PrototypeReference
	}

	for _, el := range insolar.GenesisNamePublicKeyShards {
		genesisrefs.PredefinedPrototypes[el+GenesisPrototypeSuffix] = *pkshard.PrototypeReference
	}

	for _, el := range insolar.GenesisNameMigrationAddressShards {
		genesisrefs.PredefinedPrototypes[el+GenesisPrototypeSuffix] = *migrationshard.PrototypeReference
	}
}
