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

package artifacts

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	opentracing "github.com/opentracing/opentracing-go"
)

type methodInstrumenter struct {
	ctx     context.Context
	span    opentracing.Span
	name    string
	start   time.Time
	errlink *error
}

func instrument(ctx context.Context, name string, err *error) (context.Context, *methodInstrumenter) {
	name = "artifacts." + name
	ctx, span := instracer.StartSpan(ctx, name)
	return ctx, &methodInstrumenter{
		ctx:     ctx,
		errlink: err,
		span:    span,
		start:   time.Now(),
		name:    name,
	}
}

func (mi *methodInstrumenter) end() {
	latency := time.Since(mi.start)
	inslog := inslogger.FromContext(mi.ctx)

	code := "2xx"
	if mi.errlink != nil && *mi.errlink != nil && *mi.errlink != flow.ErrCancelled {
		code = "5xx"
		inslog.Debug(mi.name, " method returned error: ", *mi.errlink)
		instracer.AddError(mi.span, *mi.errlink)
	}

	ctx := insmetrics.InsertTag(mi.ctx, tagMethod, mi.name)
	ctx = insmetrics.ChangeTags(
		ctx,
		tag.Insert(tagMethod, mi.name),
		tag.Insert(tagResult, code),
	)
	stats.Record(ctx, statCalls.M(1), statLatency.M(latency.Nanoseconds()/1e6))
	mi.span.Finish()
}
