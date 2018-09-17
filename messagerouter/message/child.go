package message

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
)

// ChildMessage is a message for saving contract's body as a child
type ChildMessage struct {
	baseMessage
	Into  core.RecordRef
	Class core.RecordRef
	Body  []byte
}

// GetOperatingRole returns operating jet role for given message type.
func (m *ChildMessage) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

// Get reference returns referenced object.
func (m *ChildMessage) GetReference() core.RecordRef {
	return m.Into
}

// Serialize serializes message.
func (m *ChildMessage) Serialize() (io.Reader, error) {
	buff := &bytes.Buffer{}
	buff.Write([]byte{byte(CallConstructorMessageType)})
	enc := gob.NewEncoder(buff)
	err := enc.Encode(m)
	return buff, err
}
