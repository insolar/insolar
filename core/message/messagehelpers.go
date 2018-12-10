package message

import (
	"github.com/insolar/insolar/core"
)

func ExtractTarget(msg core.Message) *core.RecordRef {
	return msg.ExtractTarget()
}

func ExtractRole(msg core.Message) core.DynamicRole {
	return msg.ExtractRole()
}

// ExtractAllowedSenderObjectAndRole extracts information from message
// verify sender required to 's "caller" for sender
// verification purpose. If nil then check of sender's role is not
// provided by the message bus
func ExtractAllowedSenderObjectAndRole(msg core.Message) (*core.RecordRef, core.DynamicRole) {
	return msg.ExtractAllowedSenderObjectAndRole()
}
