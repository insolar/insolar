package pulsar

// RequestType is a enum-like strings
// It identifies the type of the rpc-call
type RequestType string

const (
	// HealthCheck is a method for checking connection between pulsars
	HealthCheck RequestType = "Pulsar.HealthCheck"

	// Handshake is a method for creating connection between pulsars
	Handshake RequestType = "Pulsar.MakeHandshake"

	// ReceiveSignatureForEntropy is a method for receiving signs from peers
	ReceiveSignatureForEntropy RequestType = "Pulsar.ReceiveSignatureForEntropy"

	// ReceiveEntropy is a method for receiving entropy from peers
	ReceiveEntropy RequestType = "Pulsar.ReceiveEntropy"

	// ReceiveVector is a method for receiving vectors from peers
	ReceiveVector RequestType = "Pulsar.ReceiveVector"

	// ReceiveChosenSignature is a method for receiving signature for sending from peers
	ReceiveChosenSignature RequestType = "Pulsar.ReceiveChosenSignature"

	// ReceivePulse is a method for receiving pulse from the sender
	ReceivePulse RequestType = "Pulsar.ReceivePulse"
)

func (state RequestType) String() string {
	return string(state)
}
