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

import (
	"errors"
	"net"
	"net/rpc"

	"github.com/insolar/insolar/core"
)

type RequestType string

const (
	HealthCheck RequestType = "Pulsar.HealthCheck"
	Handshake   RequestType = "Pulsar.MakeHandshake"
)

func (state RequestType) String() string {
	return string(state)
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}

type Handler struct {
	pulsar *Pulsar
}

func (handler *Handler) HealthCheck(request *Payload, response *Payload) error {
	return nil
}

func (handler *Handler) MakeHandshake(request *Payload, response *Payload) error {
	neighbour, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		return err
	}

	result, err := checkSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Signature check failed")
	}

	generator := StandardEntropyGenerator{}
	convertedKey, err := ExportPublicKey(&handler.pulsar.PrivateKey.PublicKey)
	if err != nil {
		return err
	}
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = singData(handler.pulsar.PrivateKey, message.Body)
	if err != nil {
		return err
	}
	*response = message

	if neighbour.OutgoingClient == nil {
		conn, err := net.Dial(neighbour.ConnectionType.String(), neighbour.ConnectionAddress)
		if err != nil {
			return err
		}
		neighbour.OutgoingClient = &RpcConnection{RpcClient: rpc.NewClient(conn)}
	}

	return nil
}
