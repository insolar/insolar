// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
