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

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/testutils/testmetrics"
)

func TestLedgerArtifactManager_Metrics(t *testing.T) {
	// BEWARE: this test should not be run in parallel!
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.SendMock.Return(&reply.ID{}, nil)
	cs := testutils.NewPlatformCryptographyScheme()
	am := NewArtifactManger(db)
	am.PlatformCryptographyScheme = cs
	am.DefaultBus = mb

	tmetrics := testmetrics.Start(ctx)
	defer tmetrics.Stop()

	msg := message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa"}
	_, err := am.RegisterRequest(ctx, *am.GenesisRef(), &message.Parcel{Msg: &msg})
	require.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	_, _ = db, am
	content, err := tmetrics.FetchContent()
	require.NoError(t, err)

	assert.Contains(t, content, `insolar_artifactmanager_latency_count{method="RegisterRequest",result="2xx"} 1`)
	assert.Contains(t, content, `insolar_artifactmanager_calls{method="RegisterRequest",result="2xx"} 1`)
}
