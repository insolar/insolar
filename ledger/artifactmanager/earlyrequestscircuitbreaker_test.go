/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestEarlyRequestCircuitBreakerProvider_GetBreaker_CreateNew(t *testing.T) {
	provider := &earlyRequestCircuitBreakerProvider{
		breakers: map[core.RecordID]*requestCircuitBreakerProvider{},
	}
	expectedJet := testutils.RandomJet()

	breaker := provider.getBreaker(context.TODO(), expectedJet)

	require.Equal(t, 1, len(provider.breakers))
	require.Equal(t, breaker, provider.breakers[expectedJet])
	require.NotNil(t, provider.breakers[expectedJet].timeoutChannel)
	require.NotNil(t, provider.breakers[expectedJet].hotDataChannel)
}
