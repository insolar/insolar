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
	wmbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetRequestWM struct {
	message   payload.Meta
	requestID insolar.ID

	Dep struct {
		RecordAccessor object.RecordAccessor
		Sender         wmbus.Sender
	}
}

func NewGetRequestWM(msg payload.Meta, requestID insolar.ID) *GetRequestWM {
	return &GetRequestWM{
		requestID: requestID,
		message:   msg,
	}
}

func (p *GetRequestWM) Proceed(ctx context.Context) error {
	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.requestID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch request")
	}

	concrete := record.Unwrap(rec.Virtual)
	_, isIncoming := concrete.(*record.IncomingRequest)
	_, isOutgoing := concrete.(*record.IncomingRequest)
	if !isIncoming && !isOutgoing {
		return errors.New("failed to decode request")
	}

	msg, err := payload.NewMessage(&payload.Request{
		RequestID: p.requestID,
		Request:   *rec.Virtual,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.Dep.Sender.Reply(ctx, p.message, msg)
	inslogger.FromContext(ctx).Info("sending request")

	return nil
}
