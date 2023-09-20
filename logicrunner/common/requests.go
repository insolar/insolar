package common

import (
	"github.com/insolar/insolar/insolar/record"
)

type OutgoingRequest struct {
	Request  record.IncomingRequest
	Response []byte
	Error    error
}
