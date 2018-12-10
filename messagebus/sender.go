/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagebus

import (
	"context"

	"github.com/insolar/insolar/core"
)

// Sender is an internal interface used by recorder and player. It should not be publicated.
//
// Sender provides access to private MessageBus methods.
//go:generate minimock -i github.com/insolar/insolar/messagebus.sender -o .
type sender interface {
	core.MessageBus
	CreateParcel(ctx context.Context, msg core.Message, token core.DelegationToken, currentPulse core.Pulse) (core.Parcel, error)
	SendParcel(ctx context.Context, msg core.Parcel, currentPulse core.Pulse, ops *core.MessageSendOptions) (core.Reply, error)
}
