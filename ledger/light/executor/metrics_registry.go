// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
