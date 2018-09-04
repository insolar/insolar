package core

import (
	// TODO: should go away, no imports in TYPES package
	"github.com/insolar/insolar/network/hostnetwork"
)

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte

// Message is a routable packet, ATM just a method call
type Message struct {
	Caller      struct{}
	Constructor bool
	Reference   RecordRef
	Method      string
	Arguments   Arguments
}

// Response to a `Message`
type Response struct {
	Data   []byte
	Result []byte
	Error  error
}

// MessageRouter interface
type MessageRouter interface {
	Component
	Route(ctx hostnetwork.Context, msg Message) (resp Response, err error)
}
