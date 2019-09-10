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

package metrics_test

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/testutils/testmetrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func testMetricsServerOutput(t *testing.T) {
	// checks is metrics server properly exports metrics added with opencensus on prometheus http endpoint
	ctx := inslogger.TestContext(t)
	testm, err := testmetrics.Start(ctx, t)
	require.NoError(t, err, "metrics server start")

	var (
		metricCount = stats.Int64("some_count", "number of processed videos", stats.UnitDimensionless)
		metricDist  = stats.Int64("some_distribution", "size of processed video", stats.UnitBytes)
	)
	someTag := insmetrics.MustTagKey("xyz")

	err = view.Register(
		&view.View{
			Name:        "some_metric_count",
			Measure:     metricCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{someTag},
		},
		&view.View{
			Name:        "some_metric_distribution",
			Measure:     metricDist,
			Aggregation: view.Distribution(0, 1<<16, 1<<32),
			TagKeys:     []tag.Key{someTag},
		},
	)
	require.NoError(t, err)

	var (
		countRe = regexp.MustCompile(`insolar_some_metric_count{[^}]*xyz="11\.12\.13"[^}]*} 1`)
		distRe  = regexp.MustCompile(`insolar_some_metric_distribution_count{[^}]*xyz="11\.12\.13"[^}]*} 1`)
	)

	metricsCtx := insmetrics.ChangeTags(context.Background(), tag.Insert(someTag, "11.12.13"))
	stats.Record(metricsCtx, metricCount.M(1), metricDist.M(rand.Int63()))

	var (
		content  string
		fetchErr error
	)
	// loop because at some strange circumstances at CI one fetch is not enough
	for i := 0; i < 1000; i++ {
		time.Sleep(500 * time.Millisecond)
		content, fetchErr = testm.FetchContent()
		if fetchErr != nil {
			continue
		}
		if strings.Contains(content, "insolar_some_metric") {
			break
		}
	}

	require.NoError(t, fetchErr, "fetch content failed")
	assert.Regexp(t,
		countRe,
		content,
		"counter value is equal to 1")
	assert.Regexp(t,
		distRe,
		content,
		"distribution counter value is equal to 1")

	assert.NoError(t, testm.Stop(), "metrics server is stopped")
}

func TestMetrics_NewMetrics(t *testing.T) {
	if os.Getenv("ISOLATE_METRICS_STATE") == "1" {
		testMetricsServerOutput(t)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMetrics_NewMetrics")
	cmd.Env = append(os.Environ(), "ISOLATE_METRICS_STATE=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Fatalf("Process failed with error '%v', expects os.Exit(0)", err)
	}
}

func TestMetrics_ZPages(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testm, err := testmetrics.Start(ctx, t)
	require.NoError(t, err, "metrics server start")

	// One more thing... from https://github.com/rakyll/opencensus-grpc-demo
	// also check /debug/rpcz
	code, content, err := testm.FetchURL("/debug/tracez")
	_ = content
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)
	// fmt.Println("/debug/tracez => ", content)

	assert.NoError(t, testm.Stop())
}

func TestMetrics_Status(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testm, err := testmetrics.Start(ctx, t)
	require.NoError(t, err, "metrics server start")

	code, _, err := testm.FetchURL("/_status")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, code)

	assert.NoError(t, testm.Stop())
}
