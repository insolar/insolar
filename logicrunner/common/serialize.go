package common

import (
	"github.com/insolar/insolar/insolar"
)

type Serializer interface {
	Serialize(interface{}, *[]byte) error
	Deserialize([]byte, interface{}) error
}

type CBORSerializer struct{}

func (s *CBORSerializer) Serialize(what interface{}, to *[]byte) (err error) {
	*to, err = insolar.Serialize(what)
	return err
}

func (s *CBORSerializer) Deserialize(from []byte, to interface{}) error {
	return insolar.Deserialize(from, to)
}

func NewCBORSerializer() *CBORSerializer {
	return &CBORSerializer{}
}
