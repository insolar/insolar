// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
