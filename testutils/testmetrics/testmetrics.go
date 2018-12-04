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

package testmetrics

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
)

// TestMetrics provides testing helpers for metrics.
type TestMetrics struct {
	ctx     context.Context
	Metrics *metrics.Metrics
}

var (
	oncecfg sync.Once
	cfg     configuration.Metrics
)

// Start configures, creates and starts metrics server,
// returns initialized TestMetrics object.
func Start(ctx context.Context) TestMetrics {
	inslog := inslogger.FromContext(ctx)
	oncecfg.Do(func() {
		cfg = configuration.NewMetrics()
	})
	host, _ := parseAddr(cfg.ListenAddress)

	// just use any available port
	cfg.ListenAddress = host + ":0"
	// don't wait too long in tests
	cfg.ReportingPeriod = time.Millisecond

	m, err := metrics.NewMetrics(ctx, cfg, metrics.GetInsolarRegistry())
	if err != nil {
		panic(err)
	}

	err = m.Start(ctx)
	if err != nil {
		inslog.Fatal("metrics server failed to start:", err)
	}

	return TestMetrics{
		ctx:     ctx,
		Metrics: m,
	}
}

func parseAddr(address string) (string, int32) {
	pair := strings.SplitN(address, ":", 2)
	currentPort, err := strconv.Atoi(pair[1])
	if err != nil {
		panic(err)
	}
	return pair[0], int32(currentPort)
}

// FetchContent fetches content from /metrics.
func (tm TestMetrics) FetchContent() (string, error) {
	code, content, err := tm.FetchURL("/metrics")
	if err != nil && code != http.StatusOK {
		return "", errors.New("got non 200 code")
	}
	return content, err
}

// FetchURL fetches content from provided relative url.
func (tm TestMetrics) FetchURL(relurl string) (int, string, error) {
	// to be sure metrics are available
	time.Sleep(time.Millisecond * 5)

	fetchurl := "http://" + tm.Metrics.AddrString() + relurl
	response, err := http.Get(fetchurl)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	return response.StatusCode, string(content), err
}

// Stop wraps metrics Stop method.
func (tm TestMetrics) Stop() error {
	return tm.Metrics.Stop(tm.ctx)
}
