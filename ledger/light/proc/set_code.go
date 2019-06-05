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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetCode struct {
	message  *message.Message
	record   record.Code
	code     []byte
	recordID insolar.ID
	jetID    insolar.JetID

	dep struct {
		writer  hot.WriteAccessor
		records object.RecordModifier
		blobs   blob.Modifier
		pcs     insolar.PlatformCryptographyScheme
		sender  bus.Sender
	}
}

func NewSetCode(msg *message.Message, rec record.Code, code []byte, recID insolar.ID, jetID insolar.JetID) *SetCode {
	return &SetCode{
		message:  msg,
		record:   rec,
		code:     code,
		recordID: recID,
		jetID:    jetID,
	}
}

func (p *SetCode) Dep(
	w hot.WriteAccessor,
	r object.RecordModifier,
	b blob.Modifier,
	pcs insolar.PlatformCryptographyScheme,
	s bus.Sender,
) {
	p.dep.writer = w
	p.dep.records = r
	p.dep.blobs = b
	p.dep.pcs = pcs
	p.dep.sender = s
}

func (p *SetCode) Proceed(ctx context.Context) error {
	done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == hot.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return err
	}
	defer done()

	h := p.dep.pcs.ReferenceHasher()
	_, err = h.Write(p.code)
	if err != nil {
		return errors.Wrap(err, "failed to calculate code id")
	}
	blobID := *insolar.NewID(flow.Pulse(ctx), h.Sum(nil))
	if blobID != p.record.Code {
		return fmt.Errorf(
			"received blob id %s does not match with %s",
			p.record.Code.DebugString(),
			blobID.DebugString(),
		)
	}
	err = p.dep.blobs.Set(ctx, blobID, blob.Blob{Value: p.code, JetID: p.jetID})
	if err != nil {
		return errors.Wrap(err, "failed to store blob")
	}

	virtual := record.Wrap(p.record)
	material := record.Material{
		Virtual: &virtual,
	}
	err = p.dep.records.Set(ctx, p.recordID, material)
	if err != nil {
		return errors.Wrap(err, "failed to store record")
	}

	msg, err := payload.NewMessage(&payload.ID{ID: p.recordID})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}
	go p.dep.sender.Reply(ctx, p.message, msg)

	return nil
}
