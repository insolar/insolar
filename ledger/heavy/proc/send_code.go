// Copyright 2020 Insolar Network Ltd.
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

package proc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/object"
)

type SendCode struct {
	message payload.Meta

	Dep struct {
		RecordAccessor object.RecordAccessor
		Sender         bus.Sender
	}
}

func NewSendCode(msg payload.Meta) *SendCode {
	return &SendCode{
		message: msg,
	}
}

func (p *SendCode) Proceed(ctx context.Context) error {
	getCode := payload.GetCode{}
	err := getCode.Unmarshal(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetCode message")
	}

	rec, err := p.Dep.RecordAccessor.ForID(ctx, getCode.CodeID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch record")
	}
	buf, err := rec.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal record")
	}
	msg, err := payload.NewMessage(&payload.Code{
		Record: buf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}

	p.Dep.Sender.Reply(ctx, p.message, msg)

	return nil
}
