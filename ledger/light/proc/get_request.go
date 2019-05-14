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
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetRequest struct {
	replyTo chan<- bus.Reply
	request insolar.ID

	Dep struct {
		RecordAccessor object.RecordAccessor
	}
}

func NewGetRequest(request insolar.ID, replyTo chan<- bus.Reply) *GetRequest {
	return &GetRequest{
		request: request,
		replyTo: replyTo,
	}
}

func (p *GetRequest) Proceed(ctx context.Context) error {
	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.request)
	if err != nil {
		return errors.Wrap(err, "failed to fetch request")
	}

	virtRec := rec.Record
	req, ok := virtRec.(*object.RequestRecord)
	if !ok {
		return errors.New("failed to decode request")
	}

	rep := &reply.Request{
		ID:     p.request,
		Record: object.EncodeVirtual(req),
	}

	p.replyTo <- bus.Reply{Reply: rep}
	return nil
}
