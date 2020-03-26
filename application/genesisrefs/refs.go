// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesisrefs

import (
	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
)

const (
	FundsDepositName = "genesis_deposit"
)

var applicationPrototypes = map[string]insolar.Reference{
	application.GenesisNameRootDomain + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameRootDomain, 0),
	application.GenesisNameRootMember + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameMember + genesisrefs.PrototypeSuffix:     *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
}

func init() {
	for i, val := range applicationPrototypes {
		genesisrefs.PredefinedPrototypes[i] = val
	}
}

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = genesisrefs.GenesisRef(application.GenesisNameRootDomain)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = genesisrefs.GenesisRef(application.GenesisNameRootMember)
)
