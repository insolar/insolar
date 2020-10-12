// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesis

import (
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
)

const (
	FundsDepositName = "genesis_deposit"
)

var applicationPrototypes = map[string]insolar.Reference{
	GenesisNameRootDomain + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, GenesisNameRootDomain, 0),
	GenesisNameRootMember + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, GenesisNameMember, 0),
	GenesisNameMember + genesisrefs.PrototypeSuffix:     *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, GenesisNameMember, 0),
}

func init() {
	for i, val := range applicationPrototypes {
		genesisrefs.PredefinedPrototypes[i] = val
	}
}

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = genesisrefs.GenesisRef(GenesisNameRootDomain)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = genesisrefs.GenesisRef(GenesisNameRootMember)
)

// Get reference RootMember contract.
func GetRootMember() insolar.Reference {
	return ContractRootMember
}

// Get reference on RootDomain contract.
func GetRootDomain() insolar.Reference {
	return ContractRootDomain
}
