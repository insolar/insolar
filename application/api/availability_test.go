// Copyright 2020 Insolar Network Ltd.
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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
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
	defer testutils.LeakTester(t)

	cfg := configuration.NewLog()
	cfg.Level = "Debug"
	ctx, _ := inslogger.InitNodeLogger(context.Background(), cfg, "", "")

	counter := 0

	keeper := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch counter {
			case 0:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"available\": false}")
				require.NoError(t, err)
			case 1:
				w.WriteHeader(http.StatusBadRequest)
			case 2:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"test\": }")
				require.NoError(t, err)
			default:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"available\": true}")
				require.NoError(t, err)
			}
			counter += 1
		}))

	var checkPeriod uint = 1
	config := configuration.AvailabilityChecker{
		Enabled:        true,
		KeeperURL:      keeper.URL,
		RequestTimeout: 2,
		CheckPeriod:    checkPeriod,
	}

	nc := api.NewNetworkChecker(config)
	require.False(t, nc.IsAvailable(ctx))

	defer nc.Stop()
	err := nc.Start(ctx)
	require.NoError(t, err)

	// counter = 0
	require.False(t, nc.IsAvailable(ctx))
	time.Sleep(time.Duration(checkPeriod))

	// counter = 1
	require.False(t, nc.IsAvailable(ctx))
	time.Sleep(time.Duration(checkPeriod))

	// counter = 2, bad response body
	require.False(t, nc.IsAvailable(ctx))

	// counter default
	waitForStatus(t, nc, true)

	keeper.Close()
	waitForStatus(t, nc, false)
}
