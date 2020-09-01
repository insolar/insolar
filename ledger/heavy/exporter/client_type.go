package exporter

const (
	KeyClientType            = "client-type"
	KeyClientVersionHeavy    = "heavy-version"
	KeyClientVersionContract = "contract-version"
)

const AllowedOnHeavyVersion = 2

// Client type indicates the need to check contracts or only the protocol with heavy
type ClientType int

//go:generate stringer -type=ClientType
const (
	Unknown ClientType = iota
	ValidateHeavyVersion
	ValidateContractVersion
)
