package payload

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

type Type uint32

//go:generate stringer -type=Type

const (
	TypeUnknown   Type = 0
	TypeError     Type = 100
	TypeID        Type = 101
	TypeJet       Type = 102
	TypeGetObject Type = 103
	TypeObjIndex  Type = 104
	TypeObjState  Type = 105
)

// Payload represents any kind of data that can be encoded in consistent manner.
type Payload interface {
	Marshal() ([]byte, error)
}

func Marshal(payload Payload) ([]byte, error) {
	switch pl := payload.(type) {
	case *Error:
		pl.Polymorph = uint32(TypeError)
		return pl.Marshal()
	case *ID:
		pl.Polymorph = uint32(TypeID)
		return pl.Marshal()
	case *GetObject:
		pl.Polymorph = uint32(TypeGetObject)
		return pl.Marshal()
	}

	return nil, errors.New("unknown payload type")
}

func Unmarshal(data []byte) (Payload, error) {
	buf := proto.NewBuffer(data)
	_, err := buf.DecodeVarint()
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode polymorph")
	}
	morph, err := buf.DecodeVarint()
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode polymorph")
	}

	switch Type(morph) {
	case TypeError:
		pl := Error{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeID:
		pl := ID{}
		err := pl.Unmarshal(data)
		return &pl, err
	case TypeGetObject:
		pl := GetObject{}
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
