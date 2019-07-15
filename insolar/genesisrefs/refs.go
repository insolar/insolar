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
	ContractNodeDomain = rootdomain.GenesisRef(insolar.GenesisNameNodeDomain)
	// ContractNodeRecord is the node contract reference.
	ContractNodeRecord = rootdomain.GenesisRef(insolar.GenesisNameNodeRecord)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = rootdomain.GenesisRef(insolar.GenesisNameRootMember)
	// ContractRootWallet is the root wallet contract reference.
	ContractRootWallet = rootdomain.GenesisRef(insolar.GenesisNameRootWallet)
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = rootdomain.GenesisRef(insolar.GenesisNameMigrationAdminMember)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = rootdomain.GenesisRef(insolar.GenesisNameMigrationWallet)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = rootdomain.GenesisRef(insolar.GenesisNameDeposit)
	// ContractTariff is the tariff contract reference.
	ContractTariff = rootdomain.GenesisRef(insolar.GenesisNameTariff)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = rootdomain.GenesisRef(insolar.GenesisNameCostCenter)
	// ContractFeeWallet is the commission wallet contract reference.
	ContractFeeWallet = rootdomain.GenesisRef(insolar.GenesisNameCommissionWallet)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result []insolar.Reference) {
		for _, name := range insolar.GenesisNameMigrationDaemonMembers {
			result = append(result, rootdomain.GenesisRef(name))
		}
		return
	}()
)
