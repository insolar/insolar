package meta

const (
	// TypeReply is Type for messages with insolar.Reply in Payload
	TypeReply = "reply"

	// TypeReturnResults is Type of messages with *payload.ReturnResults in Payload
	// should be handled by contractrequester
	TypeReturnResults = "returnresults"
)

const (
	// Pulse is key for Pulse
	Pulse = "pulse"

	// Type is key for Type
	Type = "type"

	// Sender is key for Sender
	Sender = "sender"

	// Receiver is key for Receiver
	Receiver = "receiver"

	// TraceID is key for traceID
	TraceID = "TraceID"

	// SpanData is key for a span data
	SpanData = "SpanData"
)
