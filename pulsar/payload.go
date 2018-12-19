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
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// Payload is a base struct for pulsar's rpc-message
type Payload struct {
	PublicKey string
	Signature []byte
	Body      PayloadData
}

type PayloadData interface {
	Hash(hasher core.Hasher) ([]byte, error)
}

// HandshakePayload is a struct for handshake step
type HandshakePayload struct {
	Entropy core.Entropy
}

func (hp *HandshakePayload) Hash(hasher core.Hasher) ([]byte, error) {
	_, err := hasher.Write(hp.Entropy[:])
	if err != nil{
		return nil, err
	}

	return hasher.Sum(nil), err
}
// EntropySignaturePayload is a struct for sending Sign of Entropy step
type EntropySignaturePayload struct {
	PulseNumber core.PulseNumber
	Signature   []byte
}

// EntropyPayload is a struct for sending Entropy step
type EntropyPayload struct {
	PulseNumber core.PulseNumber
	Entropy     core.Entropy
}

// VectorPayload is a struct for sending vector of Entropy step
type VectorPayload struct {
	PulseNumber core.PulseNumber
	Vector      map[string]*BftCell
}

// PulsePayload is a struct for sending finished pulse to all pulsars
type PulsePayload struct {
	Pulse core.Pulse
}
