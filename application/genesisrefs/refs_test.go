// Copyright 2020 Insolar Network Ltd.
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

package genesisrefs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
)

func TestContractPublicKeyShards(t *testing.T) {
	for i, ref := range ContractPublicKeyShards(100) {
		require.Equal(t, genesisrefs.GenesisRef(application.GenesisNamePKShard+strconv.Itoa(i)), ref)
	}
}

func TestContractMigrationAddressShards(t *testing.T) {
	for i, ref := range ContractMigrationAddressShards(100) {
		require.Equal(t, genesisrefs.GenesisRef(application.GenesisNameMigrationShard+strconv.Itoa(i)), ref)
	}
}

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		application.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "insolar:1AAEAAciWtcmPVgAcaIvICkgnSsJmp4Clp650xOHjYks",
		},
		application.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "insolar:1AAEAAWeNhA_NwKaH6E36IJ-2PLvXnJRxiTTNWq1giOg",
		},
		application.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "insolar:1AAEAAVEt_2mipoVG73cbK-v33ne0krJWXkZibayYKJc",
		},
		application.GenesisNameRootAccount: {
			got:    ContractRootAccount,
			expect: "insolar:1AAEAAYUzPb6A9YCwdhstSMjq8g4dV_059cFrscpHemo",
		},
		application.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "insolar:1AAEAAVnfpSe6gLpJptcggYUNNGIu0-_kxjnef5G-nR0",
		},
		application.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "insolar:1AAEAASCuYBHyztkBdO3b5lDXgDsrk12PKOTixEW6kvY",
		},
		application.GenesisNameFeeAccount: {
			got:    ContractFeeAccount,
			expect: "insolar:1AAEAAaN2AfiHUl4HxtCcMV-KhWirOx2MA69ndZVAIpM",
		},
		application.GenesisNameFeeWallet: {
			got:    ContractFeeWallet,
			expect: "insolar:1AAEAAVsLfNjPCXS5hsvt1WHuo0RZIYCs1H3oFC2jxIM",
		},
		application.GenesisNamePKShard: {
			got:    ContractPublicKeyShards(10)[0],
			expect: "insolar:1AAEAAVM1LnFXwPa92gplaRhMroFeWi-gxznLptCtPCc",
		},
		application.GenesisNameMigrationShard: {
			got:    ContractMigrationAddressShards(10)[0],
			expect: "insolar:1AAEAAVvqEYcaRInGWx75iJWNLUSuWXU5XA3L8Qh0fpU",
		},
		application.GenesisNameMigrationAdminAccount: {
			got:    ContractMigrationAccount,
			expect: "insolar:1AAEAAYHatrmIZwsoJ3Fy78F9Q1zZ9bEuxRrZasAcYYo",
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
	ref2 := genesisrefs.GenesisRef(application.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
