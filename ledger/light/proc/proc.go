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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
)

type Dependencies struct {
	FetchJet   func(*FetchJet) *FetchJet
	WaitHot    func(*WaitHot) *WaitHot
	GetIndex   func(*GetIndex) *GetIndex
	SendObject func(p *SendObject) *SendObject
	GetCode    func(*GetCode) *GetCode
}

type ReturnReply struct {
	ReplyTo chan<- bus.Reply
	Err     error
	Reply   insolar.Reply
}

func (p *ReturnReply) Proceed(ctx context.Context) error {
	select {
	case p.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}:
	case <-ctx.Done():
	}
	return nil
}
