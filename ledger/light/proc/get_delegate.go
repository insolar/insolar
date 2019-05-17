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

package proc

import (
	"context"
	"errors"

	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
)

type GetDelegate struct {
	replyTo chan<- bus.Reply
	idx     *object.Lifeline
	msg     *message.GetDelegate
}

func NewGetDelegate(msg *message.GetDelegate, idx *object.Lifeline, rep chan<- bus.Reply) *GetDelegate {
	return &GetDelegate{
		msg:     msg,
		replyTo: rep,
		idx:     idx,
	}
}

func (s *GetDelegate) Proceed(ctx context.Context) error {
	delegateRef, ok := s.idx.DelegateByKey(s.msg.AsType)
	if !ok {
		err := errors.New("the object has no delegate for this type")
		s.replyTo <- bus.Reply{Reply: nil, Err: err}
		return err
	}

	s.replyTo <- bus.Reply{Reply: &reply.Delegate{Head: delegateRef}, Err: nil}
	return nil
}
