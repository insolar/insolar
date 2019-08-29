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

package contractrequester

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar/flow"
)

type handleResults struct {
	cr *ContractRequester

	Message *message.Message
}

func (s *handleResults) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *handleResults) Present(ctx context.Context, f flow.Flow) error {
	return s.cr.ReceiveResult(s.Message)
}

func (s *handleResults) Past(ctx context.Context, f flow.Flow) error {
	return s.Present(ctx, f)
}
