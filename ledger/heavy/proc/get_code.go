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
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetCode struct {
	message *message.Message

	Dep struct {
		RecordAccessor object.RecordAccessor
		BlobAccessor   blob.Accessor
		Sender         bus.Sender
	}
}

func NewGetCode(msg *message.Message) *GetCode {
	return &GetCode{
		message: msg,
	}
}

func (p *GetCode) Proceed(ctx context.Context) error {
	pl, err := payload.UnmarshalFromMeta(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	getCode, ok := pl.(*payload.GetCode)
	if !ok {
		return fmt.Errorf("unexpected payload type: %T", pl)
	}

	rec, err := p.Dep.RecordAccessor.ForID(ctx, getCode.CodeID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch record")
	}
	virtual := record.Unwrap(rec.Virtual)
	code, ok := virtual.(*record.Code)
	if !ok {
		return fmt.Errorf("invalid code record %#v", virtual)
	}
	b, err := p.Dep.BlobAccessor.ForID(ctx, code.Code)
	if err != nil {
		return errors.Wrap(err, "failed to fetch code blob")
	}
	buf, err := rec.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal record")
	}
	msg, err := payload.NewMessage(&payload.Code{
		Record: buf,
		Code:   b.Value,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}
	go p.Dep.Sender.Reply(ctx, p.message, msg)

	return nil
}
