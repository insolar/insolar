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

package contracts

import (
	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
)

// GenesisContractsStates returns list contract configs for genesis.
//
// Hint: order matters, because of dependency contracts on each other.
func GenesisContractsStates(cfg insolar.GenesisContractsConfig) []insolar.GenesisContractState {
	return []insolar.GenesisContractState{
		rootDomain(),
		nodeDomain(),
		rootMember(cfg.RootPublicKey),
		rootWallet(cfg.RootBalance),
	}
}

func rootDomain() insolar.GenesisContractState {
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			RootMember:    genesisrefs.ContractRootMember,
			NodeDomainRef: genesisrefs.ContractNodeDomain,
		}),
	}
}

func nodeDomain() insolar.GenesisContractState {
	nd, _ := nodedomain.NewNodeDomain()
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameNodeDomain,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(nd),
	}
}

func rootMember(publicKey string) insolar.GenesisContractState {
	m, err := member.New("RootMember", publicKey)
	if err != nil {
		panic("root member constructor failed")
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootMember,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(m),
	}
}

func rootWallet(balance uint) insolar.GenesisContractState {
	w, err := wallet.New(balance)
	if err != nil {
		panic("failed to create wallet instance")
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootWallet,
		ParentName: insolar.GenesisNameRootMember,
		Delegate:   true,
		Memory:     mustGenMemory(w),
	}
}

func mustGenMemory(data interface{}) []byte {
	b, err := insolar.Serialize(data)
	if err != nil {
		panic("failed to serialize contract data")
	}
	return b
}
