package message

import (
	"github.com/insolar/insolar/core"
)

// SignedMessageHeader is a struct with meta for the signed message
type SignedMessageHeader struct {
	Target core.RecordRef
	Role   core.JetRole
}

// NewSignedMessageHeader creates header from the message-body
func NewSignedMessageHeader(msg core.Message) SignedMessageHeader {
	return SignedMessageHeader{
		Target: ExtractTarget(msg),
		Role:   ExtractRole(msg),
	}
}
