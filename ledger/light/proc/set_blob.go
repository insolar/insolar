///
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
///

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
)

type SetBlob struct {
	replyTo chan<- bus.Reply
	code    insolar.Reference

	Dep struct {
		Bus            insolar.MessageBus
		RecordAccessor object.RecordAccessor
		Coordinator    jet.Coordinator
		BlobAccessor   blob.Accessor
	}
}

func NewSetBlob(code insolar.Reference, replyTo chan<- bus.Reply) *SetBlob {
	return &SetBlob{
		code:    code,
		replyTo: replyTo,
	}
}

func (p *SetBlob) Proceed(ctx context.Context) error {
	p.replyTo <- p.reply(ctx)
	return nil
}
