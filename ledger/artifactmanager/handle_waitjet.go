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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"
)

type WaitJet struct {
	dep *DepInjector

	Message bus.Message

	Res struct {
		Jet insolar.JetID
		Err error
	}
}

func (s *WaitJet) Present(ctx context.Context, f flow.Flow) error {
	jet := s.dep.FetchJet(&FetchJet{Parcel: s.Message.Parcel})
	if err := f.Procedure(ctx, jet); err != nil {
		if err == flow.ErrCancelled {
			return err
		}
		return err
	}

	if jet.Res.Miss {
		rep := &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet.Res.Jet)},
		}
		if err := f.Procedure(ctx, rep); err != nil {
			return err
		}
		return errors.New("jet miss")
	}

	hot := s.dep.WaitHot(&WaitHot{
		Parcel: s.Message.Parcel,
		JetID:  jet.Res.Jet,
	})
	if err := f.Procedure(ctx, hot); err != nil {
		return err
	}
	if hot.Res.Timeout {
		rep := &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.Error{ErrType: reply.ErrHotDataTimeout},
		}
		if err := f.Procedure(ctx, rep); err != nil {
			return err
		}
		return errors.New("hot waiter timeout")
	}

	s.Res.Jet = jet.Res.Jet

	return nil
}
