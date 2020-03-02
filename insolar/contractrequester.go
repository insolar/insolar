// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"context"
)

//go:generate minimock -i github.com/insolar/insolar/insolar.ContractRequester -o ../testutils -s _mock.go -g

// Payload represents any kind of data that can be encoded in consistent manner.
type Payload interface {
	Marshal() ([]byte, error)
}

// ContractRequester is the global contract requester handler. Other system parts communicate with contract requester through it.
type ContractRequester interface {
	SendRequest(ctx context.Context, msg Payload) (Reply, *Reference, error)
	Call(ctx context.Context, ref *Reference, method string, argsIn []interface{}, pulse PulseNumber) (Reply, *Reference, error)
}
