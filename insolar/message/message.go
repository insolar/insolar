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

// Package message represents message that messagebus can route
package message

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

// GetEmptyMessage constructs specified message
func getEmptyMessage(mt insolar.MessageType) (insolar.Message, error) {
	switch mt {

	// Logicrunner
	case insolar.TypeCallMethod:
		return &CallMethod{}, nil
	case insolar.TypeReturnResults:
		return &ReturnResults{}, nil
	case insolar.TypeExecutorResults:
		return &ExecutorResults{}, nil
	case insolar.TypeValidationResults:
		return &ValidationResults{}, nil
	case insolar.TypePendingFinished:
		return &PendingFinished{}, nil
	case insolar.TypeAdditionalCallFromPreviousExecutor:
		return &AdditionalCallFromPreviousExecutor{}, nil
	case insolar.TypeStillExecuting:
		return &StillExecuting{}, nil

	// Ledger
	case insolar.TypeGetObjectIndex:
		return &GetObjectIndex{}, nil

	// heavy sync
	case insolar.TypeHeavyPayload:
		return &HeavyPayload{}, nil
	// Genesis
	case insolar.TypeGenesisRequest:
		return &GenesisRequest{}, nil
	default:
		return nil, errors.Errorf("unimplemented message type %d", mt)
	}
}

// Deserialize returns decoded message.
func Deserialize(buf []byte) (insolar.Message, error) {
	msg, err := getEmptyMessage(insolar.MessageType(buf[0]))
	if err != nil {
		return nil, err
	}
	buf = buf[1:]
	err = insolar.Deserialize(buf, &msg)
	return msg, err
}

// MustSerialize serialize a insolar.Message to bytes.
func MustSerialize(msg insolar.Message) []byte {
	r := insolar.MustSerialize(msg)
	r = append([]byte{byte(msg.Type())}, r...)
	return r
}

// SerializeParcel returns io.Reader on buffer with encoded insolar.Parcel.
func SerializeParcel(parcel insolar.Parcel) (io.Reader, error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err := enc.Encode(parcel)
	return buff, err
}

// DeserializeParcel returns decoded signed message.
func DeserializeParcel(buff io.Reader) (insolar.Parcel, error) {
	var signed Parcel
	enc := gob.NewDecoder(buff)
	err := enc.Decode(&signed)
	return &signed, err
}

// ParcelToBytes deserialize a insolar.Parcel to bytes.
func ParcelToBytes(msg insolar.Parcel) []byte {
	reqBuff, err := SerializeParcel(msg)
	if err != nil {
		panic("failed to serialize message: " + err.Error())
	}
	buf, err := ioutil.ReadAll(reqBuff)
	if err != nil {
		panic("failed to serialize message: " + err.Error())
	}
	return buf
}

// ParcelMessageHash returns hash of parcel's message calculated with provided cryptography scheme.
func ParcelMessageHash(pcs insolar.PlatformCryptographyScheme, parcel insolar.Parcel) []byte {
	return pcs.IntegrityHasher().Hash(MustSerialize(parcel.Message()))
}

func init() {
	// Logicrunner
	gob.Register(&CallMethod{})
	gob.Register(&ReturnResults{})
	gob.Register(&ExecutorResults{})
	gob.Register(&AdditionalCallFromPreviousExecutor{})
	gob.Register(&ValidationResults{})
	gob.Register(&PendingFinished{})
	gob.Register(&StillExecuting{})

	// Ledger
	gob.Register(&GetObjectIndex{})

	// heavy
	gob.Register(&HeavyPayload{})

	// Bootstrap
	gob.Register(&GenesisRequest{})
	gob.Register(&Parcel{})
	gob.Register(insolar.Reference{})
}
