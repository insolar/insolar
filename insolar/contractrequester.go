//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	Call(ctx context.Context, msg Payload) (Reply, *Reference, error)
	SendRequest(ctx context.Context, ref *Reference, method string, argsIn []interface{}) (Reply, error)
	SendRequestWithPulse(ctx context.Context, ref *Reference, method string, argsIn []interface{}, pulse PulseNumber) (Reply, *Reference, error)
}
