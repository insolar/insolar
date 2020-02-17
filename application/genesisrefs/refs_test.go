// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesisrefs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/insolar"
)

func TestContractPublicKeyShards(t *testing.T) {
	for i, ref := range ContractPublicKeyShards(100) {
		require.Equal(t, GenesisRef(application.GenesisNamePKShard+strconv.Itoa(i)), ref)
	}
}

func TestContractMigrationAddressShards(t *testing.T) {
	for i, ref := range ContractMigrationAddressShards(100) {
		require.Equal(t, GenesisRef(application.GenesisNameMigrationShard+strconv.Itoa(i)), ref)
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
		GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "insolar:1AAEAAUdxJyWoY-IjQatMYpOk51MZx9tEThkqd1dSB1U",
		},
		GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "insolar:1AAEAAQy4dc1JKDJGNd5YfU7ow3DFrW_9j7v772siVMQ",
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
	ref2 := GenesisRef(application.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}

func TestGenesisRef(t *testing.T) {
	var (
		pubKey    = "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf+vsMVU75xH8uj5WRcOqYdHXtaHH\nN0na2RVQ1xbhsVybYPae3ujNHeQCPj+RaJyMVhb6Aj/AOsTTOPFswwIDAQ==\n-----END PUBLIC KEY-----\n"
		pubKeyRef = "insolar:1AAEAAcEp7HwQByGOr6rZwkyiRA3wR2POYCrDIhqBJyY"
	)
	genesisRef := GenesisRef(pubKey)
	require.Equal(t, pubKeyRef, genesisRef.String(), "reference by name always the same")
}
