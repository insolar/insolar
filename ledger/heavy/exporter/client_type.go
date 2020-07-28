package exporter

const (
	KeyClientType            = "client_type"
	KeyClientVersionHeavy    = "heavy_version"
	KeyClientVersionContract = "contract_version"
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
