// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
