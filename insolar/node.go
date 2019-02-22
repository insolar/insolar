package insolar

import (
	"github.com/insolar/insolar/core"
)

// Node represents insolar node.
type Node struct {
	ID   core.RecordRef
	Role core.StaticRole
}
