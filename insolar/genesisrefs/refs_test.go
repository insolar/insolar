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
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "11tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "11tJDPJd6QDhsKhhgc5bCCJuDEZyrPpe2EkSCVgMoeQ",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "11tJCCeN3WNGKi6w3YqxHPV7tjxLxsCcookXTe9i6uD",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "11tJDb3zZnEns6R4ChKhE4RFhzbUVxvxUdj58YF22yP",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "11tJCLgYKxM4TABHW8tY3DBxeBZZixWua6iwReJAL4g",
		},
		insolar.GenesisNameRootAccount: {
			got:    ContractRootAccount,
			expect: "11tJC166eWveiffkDiF3zyqc3zhHyqULzkNUMn1VmSc",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "11tJCUhUMyeumaDA9wPksSjugbQ5uFJ5iYfpsX9yZ7j",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "11tJDyWCLK4y4JLw7dsCaD9KzEqYTgyGXN5Zp4HuteA",
		},
		insolar.GenesisNameFeeAccount: {
			got:    ContractFeeAccount,
			expect: "11tJEHA5P78QZLoboHnrVKNAxGys78xmBGVXDVca2HE",
		},
		insolar.GenesisNamePKShard: {
			got:    ContractPublicKeyShards[0],
			expect: "11tJCPZRjHWbFXQT5xNMzhm33ZWMQMSw2f5s39hYkNM",
		},
		insolar.GenesisNameMigrationShard: {
			got:    ContractMigrationAddressShards[0],
			expect: "11tJDMfQ7GmZ2AU4efkVyPYjQ9ExkpN9uMqpqBieYwA",
		},
		insolar.GenesisNameMigrationAdminAccount: {
			got:    ContractMigrationAccount,
			expect: "11tJEFK9k1NtRjxCVbTAshtt8CDuZjV6W8m6f6RYXxD",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestRootDomain(t *testing.T) {
	ref1 := ContractRootDomain
	ref2 := GenesisRef(insolar.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
