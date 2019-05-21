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

type Payload interface{}

func Marshal(payload Payload) ([]byte, error) {
	switch pl := payload.(type) {
	case *Error:
		pl.Polymorph = uint32(TypeError)
		return pl.Marshal()
	case *ID:
		pl.Polymorph = uint32(TypeID)
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
	}

	return nil, errors.New("unknown payload type")
}
