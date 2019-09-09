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
	"strconv"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func TestContractPublicKeyShards(t *testing.T) {
	for i, ref := range ContractPublicKeyShards(100) {
		require.Equal(t, GenesisRef(insolar.GenesisNamePKShard+strconv.Itoa(i)), ref)
	}
}

func TestContractMigrationAddressShards(t *testing.T) {
	for i, ref := range ContractMigrationAddressShards(100) {
		require.Equal(t, GenesisRef(insolar.GenesisNameMigrationShard+strconv.Itoa(i)), ref)
	}
}

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "11tJDjojfnTYn2YqF6kxCQimgYhRHuL82ep8NzqrEeE",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "11tJCQmJvVAEzDUartxvGk2t2U2S642nnHAHSCDNdPa",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "11tJBoja1SMYWkw8xHdJJCu5fdjQAJZ6XfMx8YcrBq5",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "11tJCjvL9bzK1HdmaFnvmHGMvNnHYJz2qrN83if4fEf",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "11tJCWaEGnNwk97PS5RbKDErnopfH9wx5r2N1tJnqwc",
		},
		insolar.GenesisNameRootAccount: {
			got:    ContractRootAccount,
			expect: "11tJD3c7peF6Yd7VimufekgnFJg6QvtJBf643JW76L9",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "11tJCbm34yHNdh91AsgmbUBpqAyjeMgy45jZD3kjGa8",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "11tJC1eCWVFJ6digscGgBs2TdrgWntNHCYYAdaAoWEH",
		},
		insolar.GenesisNameFeeAccount: {
			got:    ContractFeeAccount,
			expect: "11tJDMf8Y83BKeEyn9qjjtgAskhRa5mzVbxdVB7Pjez",
		},
		insolar.GenesisNameFeeWallet: {
			got:    ContractFeeWallet,
			expect: "11tJCcTZXZY7zBBsNMtimx1iceLkYCED85Anu1D9R3p",
		},
		insolar.GenesisNameEnterpriseMember: {
			got:    ContractEnterpriseMember,
			expect: "11tJDFu9hDnRvvBKHWcAeDdoYb43nRHHJAm1GdnGaSd",
		},
		insolar.GenesisNameEnterpriseWallet: {
			got:    ContractEnterpriseAccount,
			expect: "11tJEBgzYNSNXj3PVRQu65HHmZmAMYBcDTHezpCWTei",
		},
		insolar.GenesisNameEnterpriseAccount: {
			got:    ContractEnterpriseAccount,
			expect: "11tJEBgzYNSNXj3PVRQu65HHmZmAMYBcDTHezpCWTei",
		},
		insolar.GenesisNamePKShard: {
			got:    ContractPublicKeyShards(10)[0],
			expect: "11tJCXnQ9AAiHGYpSee8jD9AbYu9wTJv8rbbX3kAAza",
		},
		insolar.GenesisNameMigrationShard: {
			got:    ContractMigrationAddressShards(10)[0],
			expect: "11tJCcyeLGqpzYLa3doKn4gmCdtvTSVCGav6sHcbEZ2",
		},
		insolar.GenesisNameMigrationAdminAccount: {
			got:    ContractMigrationAccount,
			expect: "11tJD1cMXNkRUY1yYNtNwse2KBB59nZmUNUHq1vrXLD",
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
