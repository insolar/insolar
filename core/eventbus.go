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

// Event is a routable packet, ATM just a method call
type Event interface {
	// Serialize serializes event.
	Serialize() (io.Reader, error)
	// Get reference returns referenced object.
	GetReference() RecordRef
	// GetOperatingRole returns operating jet role for given event type.
	GetOperatingRole() JetRole
	// React handles event and returns associated response.
	React(Components) (Reaction, error)
}

// Reaction for an `Event`
type Reaction interface {
	// Serialize serializes event.
	Serialize() (io.Reader, error)
}

// EventBus interface
type EventBus interface {
	Route(event Event) (resp Reaction, err error)
}
