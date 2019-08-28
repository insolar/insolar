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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetRequestInfo struct {
	dep  *proc.Dependencies
	meta payload.Meta
}

func NewGetRequestInfo(dep *proc.Dependencies, meta payload.Meta) *GetRequestInfo {
	return &GetRequestInfo{
		dep:  dep,
		meta: meta,
	}
}

func (s *GetRequestInfo) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	msg, ok := pl.(*payload.GetRequestInfo)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	jet := proc.NewFetchJet(msg.ObjectID, msg.Pulse, s.meta, false)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}
	objJetID := jet.Result.Jet

	hot := proc.NewWaitHot(objJetID, msg.Pulse, s.meta)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	ensureIdx := proc.NewEnsureIndex(msg.ObjectID, objJetID, s.meta, msg.Pulse)
	s.dep.EnsureIndex(ensureIdx)
	if err := f.Procedure(ctx, ensureIdx, false); err != nil {
		return err
	}

	request := proc.NewSendRequestInfo(s.meta, msg.ObjectID, msg.RequestID, msg.Pulse)
	s.dep.GetRequestInfo(request)
	return f.Procedure(ctx, request, false)
}
