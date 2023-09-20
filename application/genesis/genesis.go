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
