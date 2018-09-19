package event

import (
	"io"

	"github.com/insolar/insolar/core"
)

// ChildMessage is a event for saving contract's body as a child
type ChildMessage struct {
	baseEvent
	Into  core.RecordRef
	Class core.RecordRef
	Body  []byte
}

// GetOperatingRole returns operating jet role for given event type.
func (m *ChildMessage) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

// GetReference returns referenced object.
func (m *ChildMessage) GetReference() core.RecordRef {
	return m.Into
}

// Serialize serializes event.
func (m *ChildMessage) Serialize() (io.Reader, error) {
	return serialize(m, ChildEventType)
}
