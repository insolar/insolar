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
	// ContractRootAccount is the root account contract reference.
	ContractRootAccount = rootdomain.GenesisRef(insolar.GenesisNameRootAccount)
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = rootdomain.GenesisRef(insolar.GenesisNameMigrationAdminMember)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = rootdomain.GenesisRef(insolar.GenesisNameMigrationAdminWallet)
	// ContractMigrationAccount is the migration account contract reference.
	ContractMigrationAccount = rootdomain.GenesisRef(insolar.GenesisNameMigrationAdminAccount)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = rootdomain.GenesisRef(insolar.GenesisNameDeposit)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = rootdomain.GenesisRef(insolar.GenesisNameCostCenter)
	// ContractFeeWallet is the commission wallet contract reference.
	ContractFeeWallet = rootdomain.GenesisRef(insolar.GenesisNameFeeWallet)
	// ContractFeeAccount is the commission account contract reference.
	ContractFeeAccount = rootdomain.GenesisRef(insolar.GenesisNameFeeAccount)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationDaemonMembers {
			result[i] = rootdomain.GenesisRef(name)
		}
		return
	}()

	// ContractPublicKeyShards is the public key shards contracts references.
	ContractPublicKeyShards = func() (result [insolar.GenesisAmountPublicKeyShards]insolar.Reference) {
		for i, name := range insolar.GenesisNamePublicKeyShards {
			result[i] = rootdomain.GenesisRef(name)
		}
		return
	}()
	// ContractMigrationAddressShards is the migration address shards contracts references.
	ContractMigrationAddressShards = func() (result [insolar.GenesisAmountMigrationAddressShards]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationAddressShards {
			result[i] = rootdomain.GenesisRef(name)
		}
		return
	}()
)
