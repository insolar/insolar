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

const AllowedVersionSmartContract = 2

// Get reference RootMember contract.
func GetRootMember() insolar.Reference {
	return genesisrefs.ContractRootMember
}

// Get reference on RootDomain contract.
func GetRootDomain() insolar.Reference {
	return genesisrefs.ContractRootDomain
}

var etheriumAddressRegex = regexp.MustCompile(`^(0x)?[\dA-Fa-f]{40}$`)

// IsEthereumAddress Ethereum address format verifier
func IsEthereumAddress(s string) bool {
	return etheriumAddressRegex.MatchString(s)
}
