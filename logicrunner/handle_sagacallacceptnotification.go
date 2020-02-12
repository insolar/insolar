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

package logicrunner

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
)

type HandleSagaCallAcceptNotification struct {
	dep  *Dependencies
	meta payload.Meta
}

func (h *HandleSagaCallAcceptNotification) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.SagaCallAcceptNotification{}
	err := msg.Unmarshal(h.meta.Payload)
	if err != nil {
		return err
	}

	virtual := record.Virtual{}
	err = virtual.Unmarshal(msg.Request)
	if err != nil {
		return err
	}
	rec := record.Unwrap(&virtual)
	outgoing, ok := rec.(*record.OutgoingRequest)
	if !ok {
		return fmt.Errorf("unexpected request received %T", rec)
	}

	if err := checkOutgoingRequest(ctx, outgoing); err != nil {
		return err
	}

	outgoingReqRef := insolar.NewRecordReference(msg.DetachedRequestID)
	_, _, err = h.dep.OutgoingSender.SendOutgoingRequest(ctx, *outgoingReqRef, outgoing)
	return err
}
