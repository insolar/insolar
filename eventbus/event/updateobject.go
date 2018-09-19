package event

import (
	"io"

	"github.com/insolar/insolar/core"
)

// UpdateObjectMessage is a event for calling constructor and obtain its response
type UpdateObjectMessage struct {
	baseEvent
	Object core.RecordRef
	Body   []byte
}

// GetOperatingRole returns operating jet role for given event type.
func (m *UpdateObjectMessage) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

// Get reference returns referenced object.
func (m *UpdateObjectMessage) GetReference() core.RecordRef {
	return m.Object
}

// Serialize serializes event.
func (m *UpdateObjectMessage) Serialize() (io.Reader, error) {
	return serialize(m, UpdateObjectEventType)
}
