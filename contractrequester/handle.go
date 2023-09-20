package contractrequester

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/contractrequester/metrics"
	"github.com/insolar/insolar/insolar/flow"
)

type handleResults struct {
	cr *ContractRequester

	Message *message.Message
}

func (s *handleResults) Future(ctx context.Context, f flow.Flow) error {
	stats.Record(ctx, metrics.HandleFuture.M(1))
	return f.Migrate(ctx, s.Present)
}

func (s *handleResults) Present(ctx context.Context, f flow.Flow) error {
	handleStart := time.Now()
	stats.Record(ctx, metrics.HandleStarted.M(1))
	defer func() {
		stats.Record(ctx,
			metrics.HandleTiming.M(float64(time.Since(handleStart).Nanoseconds())/1e6))
	}()
	return s.cr.ReceiveResult(ctx, s.Message)
}

func (s *handleResults) Past(ctx context.Context, f flow.Flow) error {
	stats.Record(ctx, metrics.HandlePast.M(1))
	return s.Present(ctx, f)
}
