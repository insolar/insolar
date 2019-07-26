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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetIndex struct {
	objectID insolar.ID

	Result struct {
		Lifeline record.Lifeline
	}

	dep struct {
		indices object.IndexAccessor
	}
}

func NewGetIndex(objectID insolar.ID) *GetIndex {
	return &GetIndex{objectID: objectID}
}

func (p *GetIndex) Dep(indices object.IndexAccessor) {
	p.dep.indices = indices
}

func (p *GetIndex) Proceed(ctx context.Context) error {
	idx, err := p.dep.indices.ForID(ctx, flow.Pulse(ctx), p.objectID)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}
	p.Result.Lifeline = idx.Lifeline

	return nil
}
