package reply

import (
	"github.com/insolar/insolar/insolar"
)

// CallMethod - the most common reply
type CallMethod struct {
	Object *insolar.Reference
	Result []byte
}

// Type returns type of the reply
func (r *CallMethod) Type() insolar.ReplyType {
	return TypeCallMethod
}

type RegisterRequest struct {
	Request insolar.Reference
}

// Type returns type of the reply
func (r *RegisterRequest) Type() insolar.ReplyType {
	return TypeRegisterRequest
}
