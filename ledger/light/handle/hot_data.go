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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type HotData struct {
	dep     *proc.Dependencies
	replyTo chan<- bus.Reply
	message *message.HotData
}

func NewHotData(dep *proc.Dependencies, rep chan<- bus.Reply, msg *message.HotData) *HotData {
	return &HotData{
		dep:     dep,
		replyTo: rep,
		message: msg,
	}
}

func (s *HotData) Present(ctx context.Context, f flow.Flow) error {
	hdProc := proc.NewHotData(s.message, s.replyTo)
	s.dep.HotData(hdProc)
	if err := f.Procedure(ctx, hdProc, false); err != nil {
		panic(errors.Wrap(err, "something broken"))
		return err
	}

	for _, meta := range s.message.HotIndexes {
		go func(hi message.HotIndex) {
			refreshPendingsState := proc.NewRefreshPendingFilament(s.replyTo, flow.Pulse(ctx), meta.ObjID)
			s.dep.RefreshPendingFilament(refreshPendingsState)
			if err := f.Procedure(ctx, refreshPendingsState, false); err != nil {
				panic(errors.Wrap(err, "something broken"))
			}

			lfl := object.Lifeline{}
			err := lfl.Unmarshal(meta.Index)
			if err != nil {
				panic(errors.Wrap(err, "something broken"))
			}
			if lfl.EarliestOpenRequest != nil {
				expirePendings := proc.NewExpirePending(s.replyTo, *lfl.EarliestOpenRequest, meta.ObjID, insolar.JetID(*s.message.Jet.Record()))
				s.dep.ExpirePending(expirePendings)
				if err := f.Procedure(ctx, refreshPendingsState, false); err != nil {
					panic(errors.Wrap(err, "something broken"))
				}
			}

		}(meta)
	}

	return nil
}
