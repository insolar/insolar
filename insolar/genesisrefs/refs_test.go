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
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/stretchr/testify/require"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJDPJd6QDhsKhhgc5bCCJuDEZyrPpe2EkSCVgMoeQ.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJCCeN3WNGKi6w3YqxHPV7tjxLxsCcookXTe9i6uD.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJDb3zZnEns6R4ChKhE4RFhzbUVxvxUdj58YF22yP.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "1tJCLgYKxM4TABHW8tY3DBxeBZZixWua6iwReJAL4g.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "1tJCUhUMyeumaDA9wPksSjugbQ5uFJ5iYfpsX9yZ7j.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameTariff: {
			got:    ContractStandardTariff,
			expect: "1tJBhczj15YCWdz4AqT2cS2JRs7tYKby8fVogP3GcE.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "1tJDyWCLK4y4JLw7dsCaD9KzEqYTgyGXN5Zp4HuteA.1tJD1hMFxYYt9rHcYuvCMLdCn4AZdPfy4HPaavNWn8",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestRootDomain(t *testing.T) {
	ref1 := rootdomain.RootDomain.Ref()
	ref2 := rootdomain.GenesisRef(insolar.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
