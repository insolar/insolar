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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetRequestInfo struct {
	dep    *proc.Dependencies
	meta   payload.Meta
	passed bool
}

func NewGetRequestInfo(dep *proc.Dependencies, meta payload.Meta, passed bool) *GetRequestInfo {
	return &GetRequestInfo{
		dep:    dep,
		meta:   meta,
		passed: passed,
	}
}

func (s *GetRequestInfo) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.RequestInfo{}
	err := msg.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal RequestInfo message")
	}

	// todo check passIfNotFound := !s.passed
	request := proc.NewGetRequestInfo(s.meta, msg.ObjectID, msg.RequestID)
	s.dep.GetRequestInfo(request)
	return f.Procedure(ctx, request, false)
}
