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

	"github.com/ugorji/go/codec"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

// GetEmptyMessage constructs specified message
func getEmptyMessage(mt insolar.MessageType) (insolar.Message, error) {
	switch mt {

	// Logicrunner
	case insolar.TypeCallMethod:
		return &CallMethod{}, nil
	case insolar.TypeCallConstructor:
		return &CallConstructor{}, nil
	case insolar.TypeReturnResults:
		return &ReturnResults{}, nil
	case insolar.TypeExecutorResults:
		return &ExecutorResults{}, nil
	case insolar.TypeValidateCaseBind:
		return &ValidateCaseBind{}, nil
	case insolar.TypeValidationResults:
		return &ValidationResults{}, nil
	case insolar.TypePendingFinished:
		return &PendingFinished{}, nil
	case insolar.TypeStillExecuting:
		return &StillExecuting{}, nil

	// Ledger
	case insolar.TypeGetCode:
		return &GetCode{}, nil
	case insolar.TypeGetObject:
		return &GetObject{}, nil
	case insolar.TypeGetDelegate:
		return &GetDelegate{}, nil
	case insolar.TypeGetChildren:
		return &GetChildren{}, nil
	case insolar.TypeUpdateObject:
		return &UpdateObject{}, nil
	case insolar.TypeRegisterChild:
		return &RegisterChild{}, nil
	case insolar.TypeSetRecord:
		return &SetRecord{}, nil
	case insolar.TypeGetObjectIndex:
		return &GetObjectIndex{}, nil
	case insolar.TypeGetPendingRequests:
		return &GetPendingRequests{}, nil
	case insolar.TypeGetJet:
		return &GetJet{}, nil
	case insolar.TypeAbandonedRequestsNotification:
		return &AbandonedRequestsNotification{}, nil
	case insolar.TypeGetPendingRequestID:
		return &GetPendingRequestID{}, nil
	case insolar.TypeGetRequest:
		return &GetRequest{}, nil

	// heavy sync
	case insolar.TypeHeavyPayload:
		return &HeavyPayload{}, nil
	// Bootstrap
	case insolar.TypeBootstrapRequest:
		return &GenesisRequest{}, nil

	// NodeCert
	case insolar.TypeNodeSignRequest:
		return &NodeSignPayload{}, nil
	default:
		return nil, errors.Errorf("unimplemented message type %d", mt)
	}
}

// MustSerializeBytes returns encoded insolar.Message, panics on error.
func MustSerializeBytes(msg insolar.Message) []byte {
	r, err := Serialize(msg)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

// Serialize returns io.Reader on buffer with encoded insolar.Message.
func Serialize(msg insolar.Message) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(msg.Type())})
	if err != nil {
		return nil, err
	}

	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(msg)
	return buff, err
}

// Deserialize returns decoded message.
func Deserialize(buff io.Reader) (insolar.Message, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize message")
	}

	msg, err := getEmptyMessage(insolar.MessageType(b[0]))
	if err != nil {
		return nil, err
	}
	enc := codec.NewDecoder(buff, &codec.CborHandle{})
	if err = enc.Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// ToBytes serialize a insolar.Message to bytes.
func ToBytes(msg insolar.Message) []byte {
	reqBuff, err := Serialize(msg)
	if err != nil {
		panic(errors.Wrap(err, "failed to serialize message"))
	}
	buff, err := ioutil.ReadAll(reqBuff)
	if err != nil {
		panic(errors.Wrap(err, "failed to serialize message"))
	}
	return buff
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

func init() {
	// Bootstrap
	gob.Register(&NodeSignPayload{})

	// Logicrunner
	gob.Register(&CallConstructor{})
	gob.Register(&CallMethod{})
	gob.Register(&ReturnResults{})
	gob.Register(&ExecutorResults{})
	gob.Register(&ValidateCaseBind{})
	gob.Register(&ValidationResults{})
	gob.Register(&PendingFinished{})
	gob.Register(&StillExecuting{})

	// Ledger
	gob.Register(&GetCode{})
	gob.Register(&GetObject{})
	gob.Register(&GetDelegate{})
	gob.Register(&UpdateObject{})
	gob.Register(&RegisterChild{})
	gob.Register(&SetRecord{})
	gob.Register(&GetObjectIndex{})
	gob.Register(&SetBlob{})
	gob.Register(&ValidateRecord{})
	gob.Register(&GetPendingRequests{})
	gob.Register(&GetJet{})
	gob.Register(&AbandonedRequestsNotification{})
	gob.Register(&HotData{})
	gob.Register(&GetPendingRequestID{})
	gob.Register(&GetRequest{})

	// heavy
	gob.Register(&HeavyPayload{})

	// Bootstrap
	gob.Register(&GenesisRequest{})
	gob.Register(&Parcel{})
	gob.Register(insolar.Reference{})
	gob.Register(&GetChildren{})

	// NodeCert
	gob.Register(&NodeSignPayload{})
}
