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

type Subscription struct {
	Target string
	At     Page
}

type PullRequest struct {
	Page Page
}

type GenericReply struct {
	Data  []byte
	Error error
}

type PullReply struct {
	Data  []byte
	Total uint32
}

func NewRemoteParent(transport Transport, receiver string) Parent {
	return &remoteParent{transport: transport, receiver: receiver}
}

type remoteParent struct {
	transport Transport
	receiver  string
}

func (r *remoteParent) Subscribe(ctx context.Context, _ Target, at Page) error {
	sub := Subscription{
		Target: r.transport.Me(),
		At:     at,
	}
	data, err := insolar.Serialize(&sub)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize Subscription request")
	}
	rawReply, err := r.transport.Send(ctx, r.receiver, "replica.Subscribe", data)
	if err != nil {
		return errors.Wrapf(err, "failed to send replica.Subscribe request")
	}
	reply := GenericReply{}
	err = insolar.Deserialize(rawReply, &reply)
	if err != nil {
		return errors.Wrapf(err, "failed to deserialize Subscribe reply")
	}
	return reply.Error
}

func (r *remoteParent) Pull(ctx context.Context, from Page) ([]byte, uint32, error) {
	pr := PullRequest{
		Page: from,
	}
	data, err := insolar.Serialize(&pr)
	if err != nil {
		return []byte{}, 0, errors.Wrapf(err, "failed to serialize Pull request")
	}
	res, err := r.transport.Send(ctx, r.receiver, "replica.Pull", data)
	if err != nil {
		return []byte{}, 0, errors.Wrapf(err, "failed to send replica.Pull request")
	}
	reply := GenericReply{}
	err = insolar.Deserialize(res, &reply)
	if err != nil {
		return []byte{}, 0, errors.Wrapf(err, "failed to deserialize reply from Pull")
	}
	ext := PullReply{}
	err = insolar.Deserialize(reply.Data, &ext)
	if err != nil {
		return []byte{}, 0, errors.Wrapf(err, "failed to deserialize PullReply")
	}

	return ext.Data, ext.Total, reply.Error
}
