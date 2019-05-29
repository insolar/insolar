//
// Copyright 2019 Insolar Technologies GmbH
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
//

package pulsar

import (
	"bytes"
	"sort"
	"strconv"

	"github.com/insolar/insolar/insolar"
	"github.com/ugorji/go/codec"
)

// Payload is a base struct for pulsar's rpc-message
type Payload struct {
	PublicKey string
	Signature []byte
	Body      PayloadData
}

// PayloadData is a body of Payload
type PayloadData interface {
	Hash(hashProvider insolar.Hasher) ([]byte, error)
}

// HandshakePayload is a struct for handshake step
type HandshakePayload struct {
	Entropy insolar.Entropy
}

// Hash calculates hash of payload
func (hp *HandshakePayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	_, err := hashProvider.Write(hp.Entropy[:])
	if err != nil {
		return nil, err
	}

	return hashProvider.Sum(nil), err
}

// EntropySignaturePayload is a struct for sending Sign of Entropy step
type EntropySignaturePayload struct {
	PulseNumber      insolar.PulseNumber
	EntropySignature []byte
}

// Hash calculates hash of payload
func (es *EntropySignaturePayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	_, err := hashProvider.Write(es.EntropySignature)
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write(es.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}

	return hashProvider.Sum(nil), err
}

// EntropyPayload is a struct for sending Entropy step
type EntropyPayload struct {
	PulseNumber insolar.PulseNumber
	Entropy     insolar.Entropy
}

// Hash calculates hash of payload
func (ep *EntropyPayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	_, err := hashProvider.Write(ep.Entropy[:])
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write(ep.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}

	return hashProvider.Sum(nil), err
}

// VectorPayload is a struct for sending vector of Entropy step
type VectorPayload struct {
	PulseNumber insolar.PulseNumber
	Vector      map[string]*BftCell
}

// Hash calculates hash of payload
func (vp *VectorPayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	sortedKeys := make([]string, 0, len(vp.Vector))
	for key := range vp.Vector {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	cborH := &codec.CborHandle{}
	for _, key := range sortedKeys {
		var b bytes.Buffer
		enc := codec.NewEncoder(&b, cborH)

		threadUnsafeCell := vp.Vector[key]
		threadSaveCell := &BftCell{
			Sign:              threadUnsafeCell.GetSign(),
			Entropy:           threadUnsafeCell.GetEntropy(),
			IsEntropyReceived: threadUnsafeCell.GetIsEntropyReceived(),
		}

		err := enc.Encode(threadSaveCell)
		if err != nil {
			return nil, err
		}
		_, err = hashProvider.Write(b.Bytes())
		if err != nil {
			return nil, err
		}
	}

	_, err := hashProvider.Write(vp.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}

	return hashProvider.Sum(nil), nil
}

// PulsePayload is a struct for sending finished pulse to all pulsars
type PulsePayload struct {
	Pulse insolar.Pulse
}

// Hash calculates hash of payload
func (pp *PulsePayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	sortedKeys := make([]string, 0, len(pp.Pulse.Signs))
	for key := range pp.Pulse.Signs {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var b bytes.Buffer
	cborH := &codec.CborHandle{}
	for _, key := range sortedKeys {

		enc := codec.NewEncoder(&b, cborH)
		err := enc.Encode(pp.Pulse.Signs[key])
		if err != nil {
			return nil, err
		}
		_, err = hashProvider.Write(b.Bytes())
		if err != nil {
			return nil, err
		}
	}

	_, err := hashProvider.Write(pp.Pulse.Entropy[:])
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write(pp.Pulse.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = hashProvider.Write([]byte(strconv.Itoa(pp.Pulse.EpochPulseNumber)))
	if err != nil {
		return nil, err
	}

	_, err = hashProvider.Write(pp.Pulse.OriginID[:])
	if err != nil {
		return nil, err
	}

	return hashProvider.Sum(nil), nil
}

// PulseSenderConfirmationPayload is a struct with info about pulse's confirmations
type PulseSenderConfirmationPayload struct {
	insolar.PulseSenderConfirmation
}

// Hash calculates hash of payload
func (ps *PulseSenderConfirmationPayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	_, err := hashProvider.Write(ps.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write([]byte(ps.ChosenPublicKey))
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write(ps.Entropy[:])
	if err != nil {
		return nil, err
	}
	return hashProvider.Sum(nil), nil
}
