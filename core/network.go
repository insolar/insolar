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

// Cascade contains routing data for cascade sending
type Cascade struct {
	// NodeIds contains the slice of node identifiers that will receive the message
	NodeIds []RecordRef
	// Entropy is used for pseudorandom cascade building
	Entropy Entropy
	// Replication factor is the number of children nodes of the each node of the cascade
	ReplicationFactor uint
}

// RemoteProcedure is remote procedure call function.
type RemoteProcedure func(args [][]byte) ([]byte, error)

// Network is interface for network modules facade.
type Network interface {
	// SendMessage sends a message.
	SendMessage(nodeID RecordRef, method string, msg Message) ([]byte, error)
	// SendMessage sends a message.
	SendCascadeMessage(data Cascade, method string, msg Message) error
	// GetAddress returns an origin address.
	GetAddress() string
	// RemoteProcedureRegister is remote procedure register func.
	RemoteProcedureRegister(name string, method RemoteProcedure)
}
