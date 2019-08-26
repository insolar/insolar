package common

import (
	"context"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func ConvertQueueToMessageQueue(ctx context.Context, queue []*Transcript) []*payload.ExecutionQueueElement {
	mq := make([]*payload.ExecutionQueueElement, 0)
	var traces string
	for _, elem := range queue {
		mq = append(mq, &payload.ExecutionQueueElement{
			RequestRef:  elem.RequestRef,
			Incoming:    elem.Request,
			ServiceData: ServiceDataFromContext(elem.Context),
		})

		traces += inslogger.TraceID(elem.Context) + ", "
	}

	inslogger.FromContext(ctx).Debug("ConvertQueueToMessageQueue: ", traces)

	return mq
}
