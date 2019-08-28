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

// Package reply represents responses to messages of the messagebus
package reply

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"
)

const (
	// Generic

	// TypeError is reply with error.
	TypeError = insolar.ReplyType(iota + 1)
	// TypeOK is a generic reply for signaling a positive result.
	TypeOK
	// TypeNotOK is a generic reply for signaling a negative result.
	TypeNotOK

	// Logicrunner

	// TypeCallMethod - two binary fields: data and results.
	TypeCallMethod
	// TypeRegisterRequest - request for execution was registered
	TypeRegisterRequest
)

// ErrType is used to determine and compare reply errors.
type ErrType int

const (
	// ErrDeactivated returned when requested object is deactivated.
	ErrDeactivated = iota + 1
	// ErrStateNotAvailable is returned when requested object is deactivated.
	ErrStateNotAvailable
	// ErrHotDataTimeout is returned when no hot data received for a specific jet
	ErrHotDataTimeout
	// ErrNoPendingRequests is returned when there are no pending requests on current LME
	ErrNoPendingRequests
	// FlowCancelled is returned when a new pulse happened in the process of message execution
	FlowCancelled
)

func getEmptyReply(t insolar.ReplyType) (insolar.Reply, error) {
	switch t {
	case TypeCallMethod:
		return &CallMethod{}, nil
	case TypeRegisterRequest:
		return &RegisterRequest{}, nil
	case TypeError:
		return &Error{}, nil
	case TypeOK:
		return &OK{}, nil

	default:
		return nil, errors.Errorf("unimplemented reply type: '%d'", t)
	}
}

// Serialize returns encoded reply.
func Serialize(reply insolar.Reply) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(reply.Type())})
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(buff)
	err = enc.Encode(reply)
	return buff, err
}

// Deserialize returns decoded reply.
func Deserialize(buff io.Reader) (insolar.Reply, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short input to deserialize a message reply")
	}

	reply, err := getEmptyReply(insolar.ReplyType(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(reply)
	return reply, err
}

// ToBytes deserializes reply to bytes.
func ToBytes(rep insolar.Reply) []byte {
	repBuff, err := Serialize(rep)
	if err != nil {
		panic("failed to serialize reply: " + err.Error())
	}
	return repBuff.(*bytes.Buffer).Bytes()
}

func init() {
	gob.Register(&CallMethod{})
	gob.Register(&RegisterRequest{})
	gob.Register(&Error{})
	gob.Register(&OK{})
}

// UnmarshalFromMeta reads only payload skipping meta decoding. Use this instead of regular Unmarshal if you don't need
// Meta data.
func UnmarshalFromMeta(meta []byte) (insolar.Reply, error) {
	m := payload.Meta{}
	// Can be optimized by using proto.NewBuffer.
	err := m.Unmarshal(meta)
	if err != nil {
		return nil, err
	}

	rep, err := Deserialize(bytes.NewBuffer(m.Payload))
	if err != nil {
		return nil, errors.Wrap(err, "can't deserialize payload to reply")
	}
	return rep, nil
}
