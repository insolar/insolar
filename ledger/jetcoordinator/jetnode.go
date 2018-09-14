package jetcoordinator

import (
	"bytes"

	"github.com/insolar/insolar/core"
)

type JetNode struct {
	ref   core.RecordRef
	left  *JetNode
	right *JetNode
}

func (jn *JetNode) GetContaining(objRef *core.RecordRef) *core.RecordRef {
	if jn.left == nil || jn.right == nil {
		return &jn.ref
	}

	// Ignore pulse number when selecting jet affinity. Object reference can be generated without knowing its pulse.
	if bytes.Compare(objRef[core.PulseNumberSize:], jn.ref[core.PulseNumberSize:]) < 0 {
		return jn.left.GetContaining(objRef)
	}
	return jn.right.GetContaining(objRef)
}
