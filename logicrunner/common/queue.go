package common

import (
	"context"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func ConvertQueueToMessageQueue(ctx context.Context, queue []*Transcript) []*payload.ExecutionQueueElement {
	logger := inslogger.FromContext(ctx)

	mq := make([]*payload.ExecutionQueueElement, len(queue))
	var traces string
	if len(queue) > 0 {
		for i, elem := range queue {
			mq[i] = &payload.ExecutionQueueElement{
				RequestRef:  elem.RequestRef,
				Incoming:    elem.Request,
				ServiceData: ServiceDataFromContext(elem.Context),
			}

			traces += inslogger.TraceID(elem.Context) + ", "
		}

		logger.Debug("ConvertQueueToMessageQueue: ", traces)
	} else {
		logger.Debug("ConvertQueueToMessageQueue: empty queue ")
	}

	return mq
}
