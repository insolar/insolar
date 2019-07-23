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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

type Notification struct {
	Pulse insolar.PulseNumber
}

func NewRemoteTarget(transport Transport, receiver string) Target {
	return &remoteTarget{transport: transport, receiver: receiver}
}

type remoteTarget struct {
	transport Transport
	receiver  string
}

func (r *remoteTarget) Notify(ctx context.Context, pn insolar.PulseNumber) error {
	notification := Notification{
		Pulse: pn,
	}
	data, err := insolar.Serialize(&notification)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize Subscription request")
	}
	rawReply, err := r.transport.Send(ctx, r.receiver, "replica.Notify", data)
	if err != nil {
		return errors.Wrapf(err, "failed to send replica.Notify")
	}
	reply := GenericReply{}
	err = insolar.Deserialize(rawReply, &reply)
	if err != nil {
		return errors.Wrapf(err, "failed to deserialize reply on replica.Notify")
	}
	return reply.Error
}
