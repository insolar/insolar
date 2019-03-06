package conveyour_getobject

import (
	"github.com/insolar/insolar/core"
)

type JetPayload struct {
	JetID core.RecordID
	Err   error
}

type GetObjectPayload struct {
	Memory []byte
	Err    error
}
