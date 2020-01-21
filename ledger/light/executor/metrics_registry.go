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

package executor

import (
	"context"
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statAbandonedRequestAge = stats.Int64(
		"oldest_abandoned_request_age",
		"How many pulses passed from last abandoned request creation",
		stats.UnitDimensionless,
	)
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.MetricsRegistry -o ./ -s _mock.go -g
type MetricsRegistry interface {
	SetOldestAbandonedRequestAge(age int)
	UpdateMetrics(ctx context.Context)
}

type metricsRegistry struct {
	lock                      sync.Mutex
	oldestAbandonedRequestAge int64
}

func NewMetricsRegistry() MetricsRegistry {
	return &metricsRegistry{}
}

func (mr *metricsRegistry) SetOldestAbandonedRequestAge(age int) {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	if mr.oldestAbandonedRequestAge < int64(age) {
		mr.oldestAbandonedRequestAge = int64(age)
	}
}

func (mr *metricsRegistry) UpdateMetrics(ctx context.Context) {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	stats.Record(ctx, statAbandonedRequestAge.M(mr.oldestAbandonedRequestAge))
	mr.oldestAbandonedRequestAge = 0
}

func init() {
	err := view.Register(
		&view.View{
			Name:        statAbandonedRequestAge.Name(),
			Description: statAbandonedRequestAge.Description(),
			Measure:     statAbandonedRequestAge,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
