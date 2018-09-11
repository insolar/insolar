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

	switch bytes.Compare(objRef[:], jn.ref[:]) {
	case 0:
		return &jn.ref
	case -1:
		return jn.left.GetContaining(objRef)
	default:
		return jn.right.GetContaining(objRef)
	}
}
