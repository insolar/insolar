// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insmetrics

import (
	"context"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	prometheusclient "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// MustTagKey creates new tag.Key, panics on error
func MustTagKey(key string) tag.Key {
	k, err := tag.NewKey(key)
	if err != nil {
		panic(err)
	}
	return k
}

// InsertTag inserts (add) tag with provided value into context.
//
// Panics on error.
func InsertTag(ctx context.Context, key tag.Key, value string) context.Context {
	return ChangeTags(ctx, tag.Insert(key, value))
}

// ChangeTags wrapper around opencensus tag.New which panics on any tag creation error.
//
// Panics on errors.
func ChangeTags(ctx context.Context, mutator ...tag.Mutator) context.Context {
	ctx, err := tag.New(ctx, mutator...)
	if err != nil {
		panic(err)
	}
	return ctx
}

// Errorer is a logger with error
type Errorer interface {
	Error(...interface{})
}

// RegisterPrometheus creates prometheus exporter and registers it in opencensus view lib.
func RegisterPrometheus(
	namespace string,
	registry *prometheusclient.Registry,
	reportperiod time.Duration,
	logger Errorer,
	nodeRole string,
) (*prometheus.Exporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: namespace,
		Registry:  registry,
		OnError: func(err error) {
			logger.Error("Failed to export to Prometheus: ", err)
		},
		ConstLabels: prometheusclient.Labels{"role": nodeRole},
	})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)
	if reportperiod == 0 {
		reportperiod = time.Second
	}
	view.SetReportingPeriod(reportperiod)
	return exporter, nil
}
