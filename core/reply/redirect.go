package reply

import (
	"github.com/insolar/insolar/core"
)

type Redirect struct {
	To   core.RecordRef
	Sign core.Signature
}

func (r *Redirect) Type() core.ReplyType {
	return TypeRedirect
}


type DefinedStateRedirect struct {
	Redirect
	StateID core.RecordID
}

func (r *DefinedStateRedirect) Type() core.ReplyType {
	return TypeDefinedStateRedirect
}
