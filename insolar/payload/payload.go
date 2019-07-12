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

package payload

import (
	"github.com/gogo/protobuf/proto"
	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type Type uint32

//go:generate stringer -type=Type

const (
	TypeUnknown Type = iota
	TypeMeta
	TypeError
	TypeID
	TypeState
	TypeGetObject
	TypePassState
	TypeObjIndex
	TypeObjState
	TypeIndex
	TypePass
	TypeGetCode
	TypeCode
	TypeSetCode
	TypeSetIncomingRequest
	TypeSetOutgoingRequest
	TypeGetFilament
	TypeFilamentSegment
	TypeSetResult
	TypeActivate
	TypeRequestInfo
	TypeGotHotConfirmation
	TypeDeactivate
	TypeUpdate

	// should be the last (required by TypesMap)
	_latestType
)

// TypesMap contains Type name (gen by stringer) to type mapping.
var TypesMap = func() map[string]Type {
	m := map[string]Type{}
	for i := TypeUnknown; i < _latestType; i++ {
		m[i.String()] = i
	}
	return m
}()

// Payload represents any kind of data that can be encoded in consistent manner.
type Payload interface {
	Marshal() ([]byte, error)
}

const (
	MessageHashSize = 28
)

type MessageHash [MessageHashSize]byte

func (h *MessageHash) MarshalTo(data []byte) (int, error) {
	if len(data) < len(h) {
		return 0, errors.New("Not enough bytes to marshal PulseNumber")
	}
	copy(data, h[:])
	return len(h), nil
}

func (h *MessageHash) Unmarshal(data []byte) error {
	if len(data) < MessageHashSize {
		return errors.New("not enough bytes")
	}
	copy(h[:], data)
	return nil
}

func (h MessageHash) Equal(other MessageHash) bool {
	return h == other
}

func (h MessageHash) Size() int {
	return len(h)
}

func (h *MessageHash) String() string {
	return base58.Encode(h[:])
}

func (h *MessageHash) IsZero() bool {
	for _, b := range h {
		if b != 0 {
			return false
		}
	}
	return true
}

// UnmarshalType decodes payload type from given binary.
func UnmarshalType(data []byte) (Type, error) {
	buf := proto.NewBuffer(data)
	_, err := buf.DecodeVarint()
	if err != nil {
		return TypeUnknown, errors.Wrap(err, "failed to decode polymorph")
	}
	morph, err := buf.DecodeVarint()
	if err != nil {
		return TypeUnknown, errors.Wrap(err, "failed to decode polymorph")
	}
	return Type(morph), nil
}

func Marshal(payload Payload) ([]byte, error) {
	switch pl := payload.(type) {
	case *Meta:
		pl.Polymorph = uint32(TypeMeta)
		return pl.Marshal()
	case *Error:
		pl.Polymorph = uint32(TypeError)
		return pl.Marshal()
	case *ID:
		pl.Polymorph = uint32(TypeID)
		return pl.Marshal()
	case *State:
		pl.Polymorph = uint32(TypeState)
		return pl.Marshal()
	case *GetObject:
		pl.Polymorph = uint32(TypeGetObject)
		return pl.Marshal()
	case *PassState:
		pl.Polymorph = uint32(TypePassState)
		return pl.Marshal()
	case *Index:
		pl.Polymorph = uint32(TypeIndex)
		return pl.Marshal()
	case *Pass:
		pl.Polymorph = uint32(TypePass)
		return pl.Marshal()
	case *GetCode:
		pl.Polymorph = uint32(TypeGetCode)
		return pl.Marshal()
	case *Code:
		pl.Polymorph = uint32(TypeCode)
		return pl.Marshal()
	case *SetCode:
		pl.Polymorph = uint32(TypeSetCode)
		return pl.Marshal()
	case *GetFilament:
		pl.Polymorph = uint32(TypeGetFilament)
		return pl.Marshal()
	case *FilamentSegment:
		pl.Polymorph = uint32(TypeFilamentSegment)
		return pl.Marshal()
	case *SetIncomingRequest:
		pl.Polymorph = uint32(TypeSetIncomingRequest)
		return pl.Marshal()
	case *SetOutgoingRequest:
		pl.Polymorph = uint32(TypeSetOutgoingRequest)
		return pl.Marshal()
	case *SetResult:
		pl.Polymorph = uint32(TypeSetResult)
		return pl.Marshal()
	case *Activate:
		pl.Polymorph = uint32(TypeActivate)
		return pl.Marshal()
	case *RequestInfo:
		pl.Polymorph = uint32(TypeRequestInfo)
		return pl.Marshal()
	case *GotHotConfirmation:
		pl.Polymorph = uint32(TypeGotHotConfirmation)
		return pl.Marshal()
	case *Deactivate:
		pl.Polymorph = uint32(TypeDeactivate)
		return pl.Marshal()
	case *Update:
		pl.Polymorph = uint32(TypeUpdate)
		return pl.Marshal()
	}

	return nil, errors.New("unknown payload type")
}

func Unmarshal(data []byte) (Payload, error) {
	tp, err := UnmarshalType(data)
	if err != nil {
		return nil, err
	}
	switch tp {
	case TypeMeta:
		pl := Meta{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeError:
		pl := Error{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeID:
		pl := ID{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeState:
		pl := State{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeGetObject:
		pl := GetObject{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypePassState:
		pl := PassState{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeIndex:
		pl := Index{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypePass:
		pl := Pass{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeGetCode:
		pl := GetCode{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeCode:
		pl := Code{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeSetCode:
		pl := SetCode{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeGetFilament:
		pl := GetFilament{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeFilamentSegment:
		pl := FilamentSegment{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeSetIncomingRequest:
		pl := SetIncomingRequest{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeSetOutgoingRequest:
		pl := SetOutgoingRequest{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeSetResult:
		pl := SetResult{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeActivate:
		pl := Activate{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeRequestInfo:
		pl := RequestInfo{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeGotHotConfirmation:
		pl := GotHotConfirmation{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeDeactivate:
		pl := Deactivate{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeUpdate:
		pl := Update{}
		err := pl.Unmarshal(data)
		return &pl, err
	}

	return nil, errors.New("unknown payload type")
}

// UnmarshalFromMeta reads only payload skipping meta decoding. Use this instead of regular Unmarshal if you don't need
// Meta data.
func UnmarshalFromMeta(meta []byte) (Payload, error) {
	m := Meta{}
	// Can be optimized by using proto.NewBuffer.
	err := m.Unmarshal(meta)
	if err != nil {
		return nil, err
	}
	pl, err := Unmarshal(m.Payload)
	if err != nil {
		return nil, err
	}

	return pl, nil
}
