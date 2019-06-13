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

package genesisrefs

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/rootdomain"
)

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = rootdomain.RootDomain.Ref()
	// ContractNodeDomain is the node domain contract reference.
	ContractNodeDomain = rootdomain.GenesisRef(insolar.GetGenesisNameNodeDomain())
	// ContractNodeRecord is the node contract reference.
	ContractNodeRecord = rootdomain.GenesisRef(insolar.GetGenesisNameNodeRecord())
	// ContractRootMember is the root member contract reference.
	ContractRootMember = rootdomain.GenesisRef(insolar.GetGenesisNameRootMember())
	// ContractWallet is the root wallet contract reference.
	ContractRootWallet = rootdomain.GenesisRef(insolar.GetGenesisNameRootWallet())
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = rootdomain.GenesisRef(insolar.GetGenesisNameMigrationAdminMember())
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = rootdomain.GenesisRef(insolar.GetGenesisNameMigrationWallet())
	// ContractMigrationAdminMembers is the migration damon members contracts references.
	ContractMigrationDamonMembers = [10]insolar.Reference{}
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = rootdomain.GenesisRef(insolar.GetGenesisNameDeposit())
	// ContractTariff is the tariff contract reference.
	ContractTariff = rootdomain.GenesisRef(insolar.GetGenesisNameTariff())
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = rootdomain.GenesisRef(insolar.GetGenesisNameCostCenter())
)
