package core

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
	Route(msg Message) (resp Response, err error)
}
