package message

import (
	"github.com/insolar/insolar/core"
)

// BootstrapRequest is used for bootstrap records generation.
type BootstrapRequest struct {
	// Name should be unique for each bootstrap record.
	Name string
}

// Type implementation for bootstrap request.
func (*BootstrapRequest) Type() core.MessageType {
	return core.TypeBootstrapRequest
}

// Target implementation for bootstrap request.
func (m *BootstrapRequest) Target() *core.RecordRef {
	ref := core.NewRefFromBase58(m.Name)
	return &ref
}

// TargetRole implementation for bootstrap request.
func (*BootstrapRequest) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

// GetCaller implementation for bootstrap request.
func (*BootstrapRequest) GetCaller() *core.RecordRef {
	return nil
}
