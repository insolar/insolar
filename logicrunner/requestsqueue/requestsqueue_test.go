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

func TestQueue_New(t *testing.T) {
	q := New()
	require.NotNil(t, q)
}

func TestQueue_Append_TakeFirst(t *testing.T) {
	ctx := inslogger.TestContext(t)
	q := New()

	r1 := &common.Transcript{RequestRef: gen.Reference()}
	r2 := &common.Transcript{RequestRef: gen.Reference()}
	r3 := &common.Transcript{RequestRef: gen.Reference()}
	r4 := &common.Transcript{RequestRef: gen.Reference()}

	q.Append(ctx, FromLedger, r1)
	q.Append(ctx, FromPreviousExecutor, r2, r3)
	q.Append(ctx, FromThisPulse, r4)

	require.Equal(t, r1, q.TakeFirst(ctx))
	require.Equal(t, r2, q.TakeFirst(ctx))
	require.Equal(t, r3, q.TakeFirst(ctx))
	require.Equal(t, r4, q.TakeFirst(ctx))
	require.Nil(t, q.TakeFirst(ctx))
}

func TestQueue_TakeAllOriginatedFrom(t *testing.T) {
	ctx := inslogger.TestContext(t)
	q := New()

	r1 := &common.Transcript{RequestRef: gen.Reference()}
	r2 := &common.Transcript{RequestRef: gen.Reference()}
	r3 := &common.Transcript{RequestRef: gen.Reference()}
	r4 := &common.Transcript{RequestRef: gen.Reference()}

	q.Append(ctx, FromLedger, r1)
	q.Append(ctx, FromPreviousExecutor, r2, r3)
	q.Append(ctx, FromThisPulse, r4)

	res := q.TakeAllOriginatedFrom(ctx, FromPreviousExecutor)
	require.Equal(t, []*common.Transcript{r2, r3}, res)

	require.Equal(t, r1, q.TakeFirst(ctx))
	require.Equal(t, r4, q.TakeFirst(ctx))
	require.Nil(t, q.TakeFirst(ctx))

	// appending doesn't change what we got
	q.Append(ctx, FromPreviousExecutor, r1, r4)
	require.Equal(t, []*common.Transcript{r2, r3}, res)
}

func TestQueue_NumberOfOld(t *testing.T) {
	ctx := inslogger.TestContext(t)
	q := New()

	r1 := &common.Transcript{RequestRef: gen.Reference()}

	q.Append(ctx, FromLedger, r1)
	require.Equal(t, 1, q.NumberOfOld(ctx))

	q.Append(ctx, FromPreviousExecutor, r1, r1)
	require.Equal(t, 3, q.NumberOfOld(ctx))

	q.Append(ctx, FromThisPulse, r1)
	require.Equal(t, 3, q.NumberOfOld(ctx))

	q.TakeAllOriginatedFrom(ctx, FromLedger)
	require.Equal(t, 2, q.NumberOfOld(ctx))
}

func TestQueue_Clean(t *testing.T) {
	ctx := inslogger.TestContext(t)
	q := New()

	r1 := &common.Transcript{RequestRef: gen.Reference()}

	q.Append(ctx, FromLedger, r1)
	q.Append(ctx, FromPreviousExecutor, r1, r1)
	q.Append(ctx, FromThisPulse, r1)
	q.Clean(ctx)
	require.Nil(t, q.TakeFirst(ctx))
}
