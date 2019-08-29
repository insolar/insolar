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
	"fmt"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetObject struct {
	dep *proc.Dependencies

	meta   payload.Meta
	passed bool
}

func NewGetObject(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetObject {
	return &GetObject{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (s *GetObject) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetObject message")
	}
	msg, ok := pl.(*payload.GetObject)
	if !ok {
		return fmt.Errorf("wrong request type: %T", pl)
	}

	ctx, _ = inslogger.WithField(ctx, "object", msg.ObjectID.DebugString())

	passIfNotExecutor := !s.passed
	jet := proc.NewFetchJet(msg.ObjectID, flow.Pulse(ctx), s.meta, passIfNotExecutor)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		if err == proc.ErrNotExecutor && passIfNotExecutor {
			return nil
		}
		return err
	}
	objJetID := jet.Result.Jet

	hot := proc.NewWaitHot(objJetID, flow.Pulse(ctx), s.meta)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	// To ensure, that we have the index. Because index can be on a heavy node.
	// If we don't have it and heavy does, GetObject fails because it should update light's index state.
	ensureIdx := proc.NewEnsureIndex(msg.ObjectID, objJetID, s.meta, flow.Pulse(ctx))
	s.dep.EnsureIndex(ensureIdx)
	if err := f.Procedure(ctx, ensureIdx, false); err != nil {
		return err
	}

	send := proc.NewSendObject(s.meta, msg.ObjectID)
	s.dep.SendObject(send)
	return f.Procedure(ctx, send, false)
}
