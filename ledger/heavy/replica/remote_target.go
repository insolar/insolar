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

func NewRemoteTarget(transport Transport, receiver string) Target {
	return &target{transport: transport, receiver: receiver}
}

type target struct {
	transport Transport
	receiver  string
}

func (r *target) Notify() error {
	rawReply, err := r.transport.Send(r.receiver, "replica.Notify", nil)
	if err != nil {
		return errors.Wrapf(err, "failed to send replica.Notify")
	}
	reply := Reply{}
	err = insolar.Deserialize(rawReply, &reply)
	if err != nil {
		return errors.Wrapf(err, "failed to deserialize reply on replica.Notify")
	}
	return reply.Error
}
