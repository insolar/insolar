package core

type JetRole int

const (
	NoRole = JetRole(iota)
	VirtualExecutor
	VirtualValidator
	HeavyExecutor
	LightExecutor
	LightValidator
)

type NetworkAddress string

type PulseNumber uint32

type JetID RecordRef

type JetCoordinator interface {
	Component
	// AmI Checks Me for role on concrete pulse for this address
	AmI(role JetRole, ref RecordRef, number PulseNumber) bool
	IsIt(role JetRole, ref RecordRef, number PulseNumber) bool

	GetVirtualExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
	GetVirtualValidators(pulse PulseNumber, ref RecordRef) []NetworkAddress

	// TODO: depends on JetTree
	//GetJetID(ref RecordRef) JetID

	// TODO: calc JetID from RecordRef inside
	GetLightExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
	GetLightValidators(pulse PulseNumber, ref RecordRef) []NetworkAddress

	// TODO: calc JetID from RecordRef inside
	GetHeavyExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
}
