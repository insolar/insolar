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
			expect: "1tJBs4NHBSTZKqGET49Se31ken7i6oVhEfsnVyu6VK.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJBxdoc3hAM5aLStE4AWqVhx4DexNp85WdAWkZgQ3.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJBtBaT1r27eYNfFFkeKWpcp39ahVXaEfdaVDVT7K.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJBqpFDFVRnHghz4bFtZx5Cidnf3U5vvVqgxMBKKX.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "1tJBjgjkRDcjiqCt14hLN5bNCCUnxg9PcH3naR5vbL.11111111111111111111111111111111",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "1tJDXFp5aTKYSX4u1k9JY1bhEDhPqyggagprNtdvir.11111111111111111111111111111111",
		},
		insolar.GenesisNameTariff: {
			got:    ContractTariff,
			expect: "1tJDwWoTjy1WArGJ8vWWyooAXGmtKuDxTyEqvNwTk7.11111111111111111111111111111111",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "1tJCwPq32u3rFaTx7akLTcYLMa9FLTUrz2ykV1Md8t.11111111111111111111111111111111",
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
