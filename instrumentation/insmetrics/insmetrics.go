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

package insmetrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	censusprom "go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/inslogger"
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

// ChangeTags wrapper around opencensus tag.New for tags modifiacations.
//
// Panics on errors.
func ChangeTags(ctx context.Context, mutator ...tag.Mutator) context.Context {
	ctx, err := tag.New(ctx, mutator...)
	if err != nil {
		panic(err)
	}
	return ctx
}

// RegisterPrometheus creates prometheus exporter and registers it in opencensus view lib.
func RegisterPrometheus(
	ctx context.Context,
	namespace string,
	registry *prometheus.Registry,
	reportperiod time.Duration,
) (*censusprom.Exporter, error) {
	inslog := inslogger.FromContext(ctx)
	exporter, err := censusprom.NewExporter(censusprom.Options{
		Namespace: namespace,
		Registry:  registry,
		OnError: func(err error) {
			inslog.Error("Failed to export to Prometheus: ", err)
		},
	})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)
	// TODO: make reporting period configurable
	if reportperiod == 0 {
		reportperiod = time.Second
	}
	view.SetReportingPeriod(reportperiod)
	return exporter, nil
}
