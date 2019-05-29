package payload

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

type Type uint32

//go:generate stringer -type=Type

const (
	TypeUnknown   Type = 0
	TypeError     Type = 1
	TypeID        Type = 2
	TypeObject    Type = 3
	TypeState     Type = 4
	TypeGetObject Type = 5
	TypePassState Type = 6
	TypeObjIndex  Type = 7
	TypeObjState  Type = 8
	TypeIndex     Type = 9
	TypePass      Type = 10
)

// Payload represents any kind of data that can be encoded in consistent manner.
type Payload interface {
	Marshal() ([]byte, error)
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
	case *Error:
		pl.Polymorph = uint32(TypeError)
		return pl.Marshal()
	case *ID:
		pl.Polymorph = uint32(TypeID)
		return pl.Marshal()
	case *Object:
		pl.Polymorph = uint32(TypeObject)
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
	}

	return nil, errors.New("unknown payload type")
}

func Unmarshal(data []byte) (Payload, error) {
	tp, err := UnmarshalType(data)
	if err != nil {
		return nil, err
	}
	switch tp {
	case TypeError:
		pl := Error{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeID:
		pl := ID{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeObject:
		pl := Object{}
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

// UnmarshalTypeFromMeta decodes payload type from given meta binary.
func UnmarshalTypeFromMeta(data []byte) (Type, error) {
	m := Meta{}
	// Can be optimized by using proto.NewBuffer.
	err := m.Unmarshal(data)
	if err != nil {
		return TypeUnknown, err
	}

	return UnmarshalType(m.Payload)
}
