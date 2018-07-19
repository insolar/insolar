/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package transport

import (
	"github.com/insolar/insolar/network/host/message"
)

// Transport is an interface for insolar transport.
type Transport interface {

	// SendRequest sends message to destination. Sequence number is generated automatically.
	SendRequest(*message.Message) (Future, error)

	// SendResponse sends message for request with passed request id.
	SendResponse(message.RequestID, *message.Message) error

	// Start starts thread to listen incoming messages.
	Start() error

	// Stop gracefully stops listening.
	Stop()

	// Close disposing all transport underlying structures after stop are called.
	Close()

	// Messages returns channel to listen incoming messages.
	Messages() <-chan *message.Message

	// Stopped returns signal channel to support graceful shutdown.
	Stopped() <-chan bool
}
