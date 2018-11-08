package core

// Token is the base interface for the routing token
type Token interface {
	GetTo() *RecordRef
	GetFrom() *RecordRef
	GetPulse() PulseNumber
	GetMsgHash() []byte
	GetSign() []byte
}
