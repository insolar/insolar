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

package handle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/payload"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/ledger/light/proc"
)

type PassState struct {
	dep  *proc.Dependencies
	meta payload.Meta
}

func NewPassState(dep *proc.Dependencies, meta payload.Meta) *PassState {
	return &PassState{
		dep:  dep,
		meta: meta,
	}
}

func (s *PassState) Present(ctx context.Context, f flow.Flow) error {
	// Pass state unmarshal pl
	pl, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	passState, ok := pl.(*payload.PassState)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	// Origin message unmarshal
	pl, err = payload.Unmarshal(passState.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal origin payload")
	}
	origin, ok := pl.(*payload.Meta)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	// Origin message unmarshal pl
	pl, err = payload.Unmarshal(origin.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	_, ok = pl.(*payload.GetObject)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	state := proc.NewPassState(s.meta, passState.StateID, *origin)
	s.dep.PassState(state)
	return f.Procedure(ctx, state, false)
}
