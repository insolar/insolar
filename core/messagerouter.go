/*
 *    Copyright 2018 Insolar
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

package core

import "io"

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte

// Message is a routable packet, ATM just a method call
type Message interface {
	// Serialize serializes message.
	Serialize() (io.Reader, error)
	// Get reference returns referenced object.
	GetReference() RecordRef
	// GetOperatingRole returns operating jet role for given message type.
	GetOperatingRole() JetRole
}

// Response to a `Message`
type Response interface {
	// Serialize serializes message.
	Serialize() (io.Reader, error)
}

// MessageRouter interface
type MessageRouter interface {
	Route(msg Message) (resp Response, err error)
}
