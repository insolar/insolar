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

package metrics_test

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/metrics"
)

var (
	metricCount      = stats.Int64("some_count", "number of processed videos", stats.UnitDimensionless)
	metricCountValue = int64(0)

	metricDist      = stats.Int64("some_distribution", "size of processed video", stats.UnitBytes)
	metricDistValue = int64(0)

	someTag = insmetrics.MustTagKey("xyz")
)

func newTestMetrics(ctx context.Context, config configuration.Metrics) *metrics.Metrics {
	roleName := "testRole"
	m := metrics.NewMetrics(
		config,
		metrics.GetInsolarRegistry(roleName),
		roleName,
	)
	if err := m.Init(ctx); err != nil {
		panic(err)
	}
	return m
}

func TestMain(m *testing.M) {
	err := view.Register(
		&view.View{
			Name:        "some_metric_count",
			Measure:     metricCount,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{someTag},
		},
		&view.View{
			Name:        "some_metric_distribution",
			Measure:     metricDist,
			Aggregation: view.Distribution(0, 1<<16, 1<<32),
			TagKeys:     []tag.Key{someTag},
		},
	)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestMetrics_NewMetrics(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testm := newTestMetrics(ctx, configuration.Metrics{
		Namespace: "insolar",
	})

	// checks is metrics server properly exports metrics added with opencensus on prometheus http endpoint
	var (
		countRe = regexp.MustCompile(`insolar_some_metric_count{[^}]*xyz="11\.12\.13"[^}]*}`)
		distRe  = regexp.MustCompile(`insolar_some_metric_distribution_count{[^}]*xyz="11\.12\.13"[^}]*}`)
	)

	metricsCtx := insmetrics.ChangeTags(context.Background(), tag.Insert(someTag, "11.12.13"))

	cntAdd, distAdd := int64(rand.Intn(100)), int64(rand.Intn(1<<32))
	metricCountValue += cntAdd
	metricDistValue++
	stats.Record(metricsCtx,
		metricCount.M(cntAdd),
		metricDist.M(distAdd))

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	var (
		respCode    int
		lastCounter string
		lastDist    string
		found       int
	)

	// we need a loop here because counters are updated asynchronously
	for i := 0; i < 5; i++ {
		fmt.Println("fetch loop:", i)
		fetch := func() bool {
			rr := httptest.NewRecorder()
			testm.Handler().ServeHTTP(rr, req)

			respCode = rr.Code
			if http.StatusOK != respCode {
				return false
			}

			scanner := bufio.NewScanner(rr.Body)
			found = 0
			for scanner.Scan() {
				s := scanner.Text()
				if strings.HasPrefix(s, "insolar_some_metric_count") {
					lastCounter = s
					if !countRe.MatchString(s) {
						return false
					}
					if fmt.Sprintf("%v", metricCountValue) != metricValue(s) {
						return false
					}
					found++
				}

				if strings.HasPrefix(s, "insolar_some_metric_distribution_count") {
					lastDist = s
					if !distRe.MatchString(s) {
						return false
					}
					if fmt.Sprintf("%v", metricDistValue) != metricValue(s) {
						return false
					}
					found++
				}
			}
			return found == 2
		}
		if ok := fetch(); ok {
			break
		}
		// sleep for giving a chance to metrics framework update metrics
		time.Sleep(time.Millisecond * 500)
	}

	require.Equal(t, http.StatusOK, respCode, "fetched ok")
	assert.Equal(t, 2, found, "all metrics found")
	assert.Regexp(t, countRe, lastCounter, "counter value matches")
	assert.Equalf(t, fmt.Sprintf("%v", metricDistValue), metricValue(lastDist),
		"check value of %v", lastCounter)
	assert.Regexp(t, distRe, lastDist, "distribution counter value matches")
	assert.Equalf(t, fmt.Sprintf("%v", metricDistValue), metricValue(lastDist),
		"check value of %v", lastDist)
}

func metricValue(s string) string {
	return s[strings.LastIndex(s, " ")+1:]
}

func TestMetrics_ZPages(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testm := newTestMetrics(ctx, configuration.Metrics{
		ZpagesEnabled: true,
	})

	// One more thing... from https://github.com/rakyll/opencensus-grpc-demo
	// also check /debug/rpcz
	req, err := http.NewRequest("GET", "/debug/tracez", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	testm.Handler().ServeHTTP(rr, req)

	require.NoError(t, err, "fetch tracez page error check")
	require.Equal(t, http.StatusOK, rr.Code)
}

//
func TestMetrics_Status(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testm := newTestMetrics(ctx, configuration.Metrics{})

	req, err := http.NewRequest("GET", "/_status", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	testm.Handler().ServeHTTP(rr, req)

	require.NoError(t, err, "fetch status page error check")
	require.Equal(t, http.StatusOK, rr.Code)
}
