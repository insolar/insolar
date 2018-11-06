package insmetrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	censusprom "go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

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
func RegisterPrometheus(namespace string, registry *prometheus.Registry) (*censusprom.Exporter, error) {
	exporter, err := censusprom.NewExporter(censusprom.Options{
		Namespace: namespace,
		Registry:  registry,
	})
	if err != nil {
		return nil, err
	}
	view.RegisterExporter(exporter)
	// TODO: make reporting period configurable
	view.SetReportingPeriod(1 * time.Second)
	return exporter, nil
}
