package message

import (
	"fmt"

	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// BaseMessage base of message class family, do not use it standalone
type BaseMessage struct {
	Reference core.RecordRef
	Request   core.RecordRef
	Domain    core.RecordRef
}

func (m *BaseMessage) GetReference() core.RecordRef {
	return m.Reference
}

func (m *BaseMessage) Serialize() (io.Reader, error) {
	panic("BaseMessage is not usable object")
	buff := &bytes.Buffer{}
	buff.Write([]byte{byte(BaseMessageType)})
	enc := gob.NewEncoder(buff)
	err := enc.Encode(m)
	return buff, err
}

// MessageType is a enum type of message
type MessageType byte

const (
	BaseMessageType MessageType = iota
	CallMethodMessageType
	CallConstructorMessageType
	MessageTypesCount
)

// GetEmptyMessage constructs specified message
func GetEmptyMessage(mt MessageType) core.Message {
	switch mt {
	case 0:
		panic("working with message type == 0 is prohibited")
	case CallMethodMessageType:
		return &CallMethodMessage{}
	case CallConstructorMessageType:
		return &CallConstructorMessage{}
	default:
		panic(fmt.Sprintf("unimplemented messagetype %d", mt))
	}
}

// Deserialize returns a message
func Deserialize(buff io.Reader) (core.Message, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short slice for deserialize message")
	}

	m := GetEmptyMessage(MessageType(b[0]))
	enc := gob.NewDecoder(buff)
	err = enc.Decode(m)
	return m, err
}
