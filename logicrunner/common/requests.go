package common

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

type OutgoingRequest struct {
	Request   record.IncomingRequest
	NewObject *insolar.Reference
	Response  []byte
	Error     error
}
