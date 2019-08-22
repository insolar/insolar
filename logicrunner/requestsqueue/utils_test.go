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

package requestsqueue

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/common"
)

func TestFirstNFromMany(t *testing.T) {
	ctx := inslogger.TestContext(t)

	t.Run("empty", func(t *testing.T) {
		q1 := New()
		q2 := New()
		requests, hasMore := FirstNFromMany(ctx, 2, q1, q2)
		require.Empty(t, requests)
		require.False(t, hasMore)
	})

	t.Run("order", func(t *testing.T) {
		r1 := &common.Transcript{RequestRef: gen.Reference()}
		r2 := &common.Transcript{RequestRef: gen.Reference()}
		r3 := &common.Transcript{RequestRef: gen.Reference()}
		r4 := &common.Transcript{RequestRef: gen.Reference()}
		r5 := &common.Transcript{RequestRef: gen.Reference()}
		r6 := &common.Transcript{RequestRef: gen.Reference()}

		q1 := New()
		q1.Append(ctx, FromLedger, r1)
		q1.Append(ctx, FromPreviousExecutor, r2)
		q1.Append(ctx, FromThisPulse, r3)

		q2 := New()
		q2.Append(ctx, FromLedger, r4)
		q2.Append(ctx, FromPreviousExecutor, r5)
		q2.Append(ctx, FromThisPulse, r6)

		requests, hasMore := FirstNFromMany(ctx, 6, q1, q2)
		require.Equal(t, []*common.Transcript{r1, r4, r2, r5, r3, r6}, requests)
		require.False(t, hasMore)
	})

	t.Run("size limit", func(t *testing.T) {
		r1 := &common.Transcript{RequestRef: gen.Reference()}
		r2 := &common.Transcript{RequestRef: gen.Reference()}
		r3 := &common.Transcript{RequestRef: gen.Reference()}
		r4 := &common.Transcript{RequestRef: gen.Reference()}
		r5 := &common.Transcript{RequestRef: gen.Reference()}
		r6 := &common.Transcript{RequestRef: gen.Reference()}

		q1 := New()
		q1.Append(ctx, FromLedger, r1)
		q1.Append(ctx, FromPreviousExecutor, r2)
		q1.Append(ctx, FromThisPulse, r3)

		q2 := New()
		q2.Append(ctx, FromLedger, r4)
		q2.Append(ctx, FromPreviousExecutor, r5)
		q2.Append(ctx, FromThisPulse, r6)

		requests, hasMore := FirstNFromMany(ctx, 4, q1, q2)
		require.Equal(t, []*common.Transcript{r1, r4, r2, r5}, requests)
		require.True(t, hasMore)
	})
}
