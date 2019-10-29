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

// +build slowtest

package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/configuration"
)

func waitForStatus(t *testing.T, nc *api.NetworkChecker, expected bool) {
	ctx := context.Background()
	var available bool
	for i := 0; i < 10; i++ {
		available = nc.IsAvailable(ctx)
		if available == expected {
			return
		}
		time.Sleep(time.Second)
	}
	require.Fail(t, "Status not passed, expected: ", expected)
}

func TestAvailabilityChecker_UpdateStatus(t *testing.T) {
	ctx := context.Background()

	keeper := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprintln(w, "{\"available\": true}")
			require.NoError(t, err)
		}))

	config := configuration.AvailabilityChecker{
		Enabled:        true,
		KeeperURL:      keeper.Config.Addr + keeper.URL,
		RequestTimeout: 2,
		CheckPeriod:    1,
	}

	nc := api.NewNetworkChecker(config)
	defer nc.Stop()

	require.False(t, nc.IsAvailable(ctx))

	err := nc.Start(ctx)
	require.NoError(t, err)
	waitForStatus(t, nc, true)

	keeper.Close()
	waitForStatus(t, nc, false)
}
