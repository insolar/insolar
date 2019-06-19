package executor

import (
	"github.com/insolar/insolar/insolar"
)

// JetInfo holds info about jet.
type JetInfo struct {
	ID insolar.JetID
	// SplitIntent indicates what jet has intention to do drop.
	SplitIntent bool
	// SplitPerformed indicates what jet was slitted,
	SplitPerformed bool

	// MineNext if not set pendings would be removed for this jet from recent storage. (legacy)
	MineNext bool
}
