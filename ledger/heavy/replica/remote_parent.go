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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

type Subscription struct {
	Target string
	At     Position
}

type PullRequest struct {
	Scope byte
	From  Position
	Limit uint32
}

type Reply struct {
	Data  []byte
	Error error
}

func NewRemoteParent(transport Transport, receiver string) Parent {
	return &parent{transport: transport, receiver: receiver}
}

type parent struct {
	transport Transport
	receiver  string
}

func (r *parent) Subscribe(child Target, at Position) error {
	sub := Subscription{
		Target: r.transport.Me(),
		At:     at,
	}
	data, err := insolar.Serialize(&sub)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize Subscription request")
	}
	rawReply, err := r.transport.Send(r.receiver, "replica.Subscribe", data)
	if err != nil {
		return errors.Wrapf(err, "failed to send replica.Subscribe request")
	}
	reply := Reply{}
	err = insolar.Deserialize(rawReply, &reply)
	if err != nil {
		return errors.Wrapf(err, "failed to deserialize Subscribe reply")
	}
	return reply.Error
}

func (r *parent) Pull(scope byte, from Position, limit uint32) ([]byte, error) {
	pr := PullRequest{
		Scope: scope,
		From:  from,
		Limit: limit,
	}
	data, err := insolar.Serialize(&pr)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to serialize Pull request")
	}
	res, err := r.transport.Send(r.receiver, "replica.Pull", data)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to send replica.Pull request")
	}
	reply := Reply{}
	err = insolar.Deserialize(res, &reply)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to deserialize Pull reply")
	}
	return reply.Data, reply.Error
}
