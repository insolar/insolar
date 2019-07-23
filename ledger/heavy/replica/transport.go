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

package replica

import (
	"context"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/replica.Transport -o ./ -s _mock.go

// Transport provides methods for sending a message,
// registering the recipient of the message, and obtaining an identity.
type Transport interface {
	// Send performs message sending.
	Send(ctx context.Context, receiver, method string, data []byte) ([]byte, error)
	// Register performs message recipient registering.
	Register(method string, handle Handle)
	// Me obtains transport identity.
	Me() string
}

type Handle func(data []byte) ([]byte, error)
