package message

import (
	"github.com/insolar/insolar/core"
)

// Token is an auth token for coorditaning messages
type Token struct {
	To core.RecordRef
	From core.RecordRef
}
