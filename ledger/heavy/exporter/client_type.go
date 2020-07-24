package exporter

const KeyClientType = "client_type"
const KeyClientVersionHeavy = "heavy_version"
const KeyClientVersionContract = "contract_version"

// Client type indicates the need to check contracts or only the protocol with heavy
type ClientType int

//go:generate stringer -type=ClientType
const (
	Unknown ClientType = iota
	ValidateHeavyVersion
	ValidateContractVersion
)
