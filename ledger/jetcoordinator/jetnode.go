package jetcoordinator

import (
	"github.com/insolar/insolar/core"
)

type JetNode struct {
	threshold uint64
	ref       *core.RecordRef
	left      *JetNode
	right     *JetNode
}

func (jn *JetNode) GetContaining(objRef uint64) *core.RecordRef {
	if jn.ref != nil {
		return jn.ref
	}

	if objRef < jn.threshold {
		return jn.left.GetContaining(objRef)
	} else {
		return jn.right.GetContaining(objRef)
	}
	return nil
}
