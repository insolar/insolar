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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/testutils/testmetrics"
)

func TestLedgerArtifactManager_Metrics(t *testing.T) {
	// BEWARE: this test should not be the parallel!
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	tmetrics := testmetrics.Start(ctx)
	defer tmetrics.Stop()

	msg := message.GenesisRequest{Name: "my little message"}
	_, err := am.RegisterRequest(ctx, &message.Parcel{Msg: &msg})
	require.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	_, _ = db, am
	content, err := tmetrics.FetchContent()
	require.NoError(t, err)

	assert.Contains(t, content, `insolar_artifactmanager_latency_count{method="RegisterRequest",result="2xx"} 1`)
	assert.Contains(t, content, `insolar_artifactmanager_calls{method="RegisterRequest",result="2xx"} 1`)
}
