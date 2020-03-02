// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package appfoundation

import (
	"regexp"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
)

// Get reference CostCenter contract.
func GetCostCenter() insolar.Reference {
	return genesisrefs.ContractCostCenter
}

// Get reference MigrationAdminMember contract.
func GetMigrationAdminMember() insolar.Reference {
	return genesisrefs.ContractMigrationAdminMember
}

// Get reference RootMember contract.
func GetRootMember() insolar.Reference {
	return genesisrefs.ContractRootMember
}

// Get reference on MigrationAdmin contract.
func GetMigrationAdmin() insolar.Reference {
	return genesisrefs.ContractMigrationAdmin
}

// Get reference on RootDomain contract.
func GetRootDomain() insolar.Reference {
	return genesisrefs.ContractRootDomain
}

// Get reference on  migrationdaemon contract by  migration member.
func GetMigrationDaemon(migrationMember insolar.Reference) (insolar.Reference, error) {
	return genesisrefs.ContractMigrationMap[migrationMember], nil
}

// Check member is migration daemon member or not
func IsMigrationDaemonMember(member insolar.Reference) bool {
	for _, mDaemonMember := range genesisrefs.ContractMigrationDaemonMembers {
		if mDaemonMember.Equal(member) {
			return true
		}
	}
	return false
}

var etheriumAddressRegex = regexp.MustCompile(`^(0x)?[\dA-Fa-f]{40}$`)

// IsEthereumAddress Ethereum address format verifier
func IsEthereumAddress(s string) bool {
	return etheriumAddressRegex.MatchString(s)
}
