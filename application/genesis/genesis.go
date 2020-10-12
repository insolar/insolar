// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesis

const (
	// GenesisNameRootDomain is the name of root domain contract for genesis record.
	GenesisNameRootDomain = "rootdomain"
	// GenesisNameMember is the name of member contract for genesis record.
	GenesisNameMember     = "member"
	GenesisNameRootMember = "root" + GenesisNameMember
)

// GenesisContractsConfig carries data required for contract object initialization via genesis.
type GenesisContractsConfig struct {
	// RootBalance is a balance of Root Member.
	RootBalance string
	// RootPublicKey is public key of Root Member.
	RootPublicKey string
}
