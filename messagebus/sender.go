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

package messagebus

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/messagebus.sender -o .

// Sender is an internal interface used by recorder and player. It should not be publicated.
//
// Sender provides access to private MessageBus methods.
type sender interface {
	insolar.MessageBus
	CreateParcel(ctx context.Context, msg insolar.Message, token insolar.DelegationToken, currentPulse insolar.Pulse) (insolar.Parcel, error)
	SendParcel(ctx context.Context, msg insolar.Parcel, currentPulse insolar.Pulse, ops *insolar.MessageSendOptions) (insolar.Reply, error)
}
