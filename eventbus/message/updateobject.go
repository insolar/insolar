package message

import (
	"io"

	"github.com/insolar/insolar/core"
)

// UpdateObjectMessage is a message for calling constructor and obtain its response
type UpdateObjectMessage struct {
	baseEvent
	Object core.RecordRef
	Body   []byte
}

// GetOperatingRole returns operating jet role for given message type.
func (m *UpdateObjectMessage) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

// Get reference returns referenced object.
func (m *UpdateObjectMessage) GetReference() core.RecordRef {
	return m.Object
}

// Serialize serializes message.
func (m *UpdateObjectMessage) Serialize() (io.Reader, error) {
	return serialize(m, UpdateObjectEventType)
}
