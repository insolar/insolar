package core

// RoutingToken is the base interface for the routing token
type Token interface {
	To() *RecordRef
	From() *RecordRef
	Pulse() PulseNumber
	MsgHash() []byte
	Sign() []byte
}
