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

package handle

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type Error struct {
	message *message.Message
}

func NewError(msg *message.Message) *Error {
	return &Error{
		message: msg,
	}
}

func (s *Error) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.UnmarshalFromMeta(s.message.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to unmarshal error"))
		return nil
	}
	p, ok := pl.(*payload.Error)
	if !ok {
		inslogger.FromContext(ctx).Errorf("unexpected error type %T", pl)
		return nil
	}

	inslogger.FromContext(ctx).WithField(
		"correlation_id",
		middleware.MessageCorrelationID(s.message),
	).Error("received error: ", p.Text)
	return nil
}
