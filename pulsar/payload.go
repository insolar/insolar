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
	PulseNumber      core.PulseNumber
	EntropySignature []byte
}

func (es *EntropySignaturePayload) Hash(hasher core.Hasher) ([]byte, error) {
	_, err := hasher.Write(es.EntropySignature[:])
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(es.PulseNumber.Bytes())
	if err != nil{
		return nil, err
	}

	return hasher.Sum(nil), err
}

// EntropyPayload is a struct for sending Entropy step
type EntropyPayload struct {
	PulseNumber core.PulseNumber
	Entropy     core.Entropy
}

func (ep *EntropyPayload) Hash(hasher core.Hasher) ([]byte, error) {
	_, err := hasher.Write(ep.Entropy[:])
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(ep.PulseNumber.Bytes())
	if err != nil{
		return nil, err
	}

	return hasher.Sum(nil), err
}

// VectorPayload is a struct for sending vector of Entropy step
type VectorPayload struct {
	PulseNumber core.PulseNumber
	Vector      map[string]*BftCell
}

func (vp *VectorPayload) Hash(hasher core.Hasher) ([]byte, error) {
	var sortedKeys []string
	for key := range  vp.Vector{
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	cborH := &codec.CborHandle{}
	for _, key := range sortedKeys{
		var b bytes.Buffer
		enc := codec.NewEncoder(&b, cborH)
		err := enc.Encode(vp.Vector[key])
		if err != nil {
			return nil, err
		}
		_, err = hasher.Write(b.Bytes())
		if err != nil {
			return nil, err
		}
	}

	_, err := hasher.Write(vp.PulseNumber.Bytes())
	if err != nil{
		return nil, err
	}

	return hasher.Sum(nil), nil
}

// PulsePayload is a struct for sending finished pulse to all pulsars
type PulsePayload struct {
	Pulse core.Pulse
}

func (pp *PulsePayload) Hash(hasher core.Hasher) ([]byte, error) {
	var sortedKeys []string
	for key := range  pp.Pulse.Signs{
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	var b bytes.Buffer
	cborH := &codec.CborHandle{}
	for _, key := range sortedKeys{

		enc := codec.NewEncoder(&b, cborH)
		err := enc.Encode(pp.Pulse.Signs[key])
		if err != nil {
			return nil, err
		}
		_, err = hasher.Write(b.Bytes())
		if err != nil {
			return nil, err
		}
	}

	_, err := hasher.Write(pp.Pulse.Entropy[:])
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(pp.Pulse.PulseNumber.Bytes())
	if err != nil{
		return nil, err
	}

	err = binary.Write(&b, binary.LittleEndian, pp.Pulse.EpochPulseNumber)
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(b.Bytes())
	if err != nil{
		return nil, err
	}

	_, err = hasher.Write(pp.Pulse.OriginID[:])
	if err != nil{
		return nil, err
	}

	return hasher.Sum(nil), nil
}

type PulseSenderConfirmationPayload struct {
	core.PulseSenderConfirmation
}

func (ps *PulseSenderConfirmationPayload) Hash(hasher core.Hasher) ([]byte, error) {
	_, err := hasher.Write(ps.PulseNumber.Bytes())
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write([]byte(ps.ChosenPublicKey))
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(ps.Entropy[:])
	if err != nil{
		return nil, err
	}
	_, err = hasher.Write(ps.Signature)
	if err != nil{
		return nil, err
	}
	return hasher.Sum(nil), nil
}

