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

package pulsar

// RequestType is a enum-like strings
// It identifies the type of the rpc-call
type RequestType string

const (
	// HealthCheck is a method for checking connection between pulsars
	HealthCheck RequestType = "Pulsar.HealthCheck"

	// Handshake is a method for creating connection between pulsars
	Handshake RequestType = "Pulsar.MakeHandshake"

	// ReceiveSignatureForEntropy is a method for receiving signs from peers
	ReceiveSignatureForEntropy RequestType = "Pulsar.ReceiveSignatureForEntropy"

	// ReceiveEntropy is a method for receiving Entropy from peers
	ReceiveEntropy RequestType = "Pulsar.ReceiveEntropy"

	// ReceiveVector is a method for receiving vectors from peers
	ReceiveVector RequestType = "Pulsar.ReceiveVector"

	// ReceiveChosenSignature is a method for receiving signature for sending from peers
	ReceiveChosenSignature RequestType = "Pulsar.ReceiveChosenSignature"

	// ReceivePulse is a method for receiving pulse from the sender
	ReceivePulse RequestType = "Pulsar.ReceivePulse"
)

func (state RequestType) String() string {
	return string(state)
}
