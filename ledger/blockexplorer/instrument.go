package blockexplorer

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
)

type methodInstrumenter struct {
	ctx     context.Context
	name    string
	start   time.Time
	errlink *error
}

func instrument(ctx context.Context, name string) *methodInstrumenter {
	return &methodInstrumenter{
		ctx:   ctx,
		start: time.Now(),
		name:  name,
	}
}

func (mi *methodInstrumenter) err(err *error) *methodInstrumenter {
	mi.errlink = err
	return mi
}

func (mi *methodInstrumenter) end() {
	latency := time.Since(mi.start)
	inslog := inslogger.FromContext(mi.ctx)

	code := "2xx"
	if mi.errlink != nil && *mi.errlink != nil {
		code = "5xx"
		inslog.Error(*mi.errlink)
	}

	inslog.Debugf("measured time of BE method %v is %v", mi.name, latency)

	ctx := insmetrics.InsertTag(mi.ctx, tagMethod, mi.name)
	ctx = insmetrics.ChangeTags(
		ctx,
		tag.Insert(tagMethod, mi.name),
		tag.Insert(tagResult, code),
	)
	stats.Record(ctx, statCalls.M(1), statLatency.M(latency.Nanoseconds()/1e6))
}
