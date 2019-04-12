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

package pulse

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

const testPulse = insolar.PulseNumber(42)

func TestContextWith(t *testing.T) {
	t.Parallel()
	ctx := ContextWith(context.Background(), testPulse)
	require.Equal(t, testPulse, ctx.Value(contextKey{}))
}

func TestFromContext(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), contextKey{}, testPulse)
	result := FromContext(ctx)
	require.Equal(t, testPulse, result)
}
