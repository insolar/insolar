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
			expect: "1tJDdt2rswYc7W7uKgCMPNJuQsVj98CMEQjMdNiD3P.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJBuVNvfN2aobkUSx2LYfdbrYfCKCuvMTqToBZcPX.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJChm5qnRjXPhD69dBKLQg4YVJ3krQc8TYVYRmte5.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJDMkBFgzhx8MNLujfmwVjm1icuVmwU43gTeH3bfQ.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "1tJDSUusVg63maCCV4JaFFP1qyDj7GV66De88q4A41.11111111111111111111111111111111",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "1tJE2dgG45oigzahQaJxdeqhqHJzsY26aRecDwd3ck.11111111111111111111111111111111",
		},
		insolar.GenesisNameTariff: {
			got:    ContractTariff,
			expect: "1tJCajdax8rhWhdtmhffz3sagQmEzqE81dthtkca12.11111111111111111111111111111111",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "1tJDRFyccUb2bvoxVCgGf4eZcbsFSfT9PY3FArJ8nv.11111111111111111111111111111111",
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
