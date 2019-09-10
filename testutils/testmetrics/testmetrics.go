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

package testmetrics

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

// TestMetrics provides testing helpers for metrics.
type TestMetrics struct {
	ctx     context.Context
	Metrics *metrics.Metrics
}

// Start configures, creates and starts metrics server,
// returns initialized TestMetrics object.
func Start(ctx context.Context, t *testing.T) (*TestMetrics, error) {
	cfg := configuration.NewMetrics()
	cfg.ReportingPeriod = time.Millisecond * 100
	host, _ := parseAddr(cfg.ListenAddress)

	// just use any available port
	cfg.ListenAddress = host + ":0"
	// don't wait too long in tests
	cfg.ReportingPeriod = time.Millisecond

	m, err := metrics.NewMetrics(ctx, cfg, metrics.GetInsolarRegistry("test"), "test")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new metrics server")
	}

	cert := testutils.NewCertificateMock(t)
	cert.GetRoleMock.Return(insolar.StaticRoleVirtual)

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateMock.Return(cert)

	err = m.Start(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "metrics server failed to start on host %v", host)
	}

	return &TestMetrics{
		ctx:     ctx,
		Metrics: m,
	}, nil
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
func (tm *TestMetrics) FetchContent() (string, error) {
	code, content, err := tm.FetchURL("/metrics")
	if err != nil && code != http.StatusOK {
		return "", errors.New("got non 200 code")
	}
	return content, err
}

// FetchURL fetches content from provided relative url.
func (tm *TestMetrics) FetchURL(relurl string) (int, string, error) {
	// to be sure metrics are available
	time.Sleep(time.Millisecond * 5)

	fetchurl := "http://" + tm.Metrics.AddrString() + relurl
	response, err := http.Get(fetchurl) //nolint: gosec
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	return response.StatusCode, string(content), err
}

// Stop wraps metrics Stop method.
func (tm *TestMetrics) Stop() error {
	return tm.Metrics.Stop(tm.ctx)
}
