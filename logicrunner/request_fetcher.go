//
// Copyright 2019 Insolar Technologies GmbH
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
//

package logicrunner

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/metrics"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestFetcher -o ./ -s _mock.go -g

type RequestFetcher interface {
	FetchPendings(ctx context.Context, trs chan<- *common.Transcript)
	Abort(ctx context.Context)
}

type requestFetcher struct {
	object insolar.Reference

	aborted      chan struct{}
	stopFetching func()

	broker           ExecutionBrokerI
	artifactsManager artifacts.Client
	outgoingsSender  OutgoingRequestSender
}

func NewRequestsFetcher(
	obj insolar.Reference, am artifacts.Client, br ExecutionBrokerI, os OutgoingRequestSender,
) RequestFetcher {
	aborted := make(chan struct{})
	once := sync.Once{}
	return &requestFetcher{
		object:           obj,
		broker:           br,
		artifactsManager: am,
		outgoingsSender:  os,
		aborted:          aborted,
		stopFetching:     func() { once.Do(func() { close(aborted) }) },
	}
}

func (rf *requestFetcher) Abort(ctx context.Context) {
	rf.stopFetching()
}

func (rf *requestFetcher) isAborted() bool {
	select {
	case <-rf.aborted:
		return true
	default:
		return false
	}
}

func (rf *requestFetcher) FetchPendings(ctx context.Context, trs chan<- *common.Transcript) {
	defer close(trs)

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"object": rf.object.String(),
	})

	logger.Debug("request fetcher starting")

	err := rf.fetch(ctx, trs)
	if err != nil {
		logger.Error("couldn't make fetch round: ", err.Error())
	}
}

// XXX: merge with value in GetPendings, duplicate
const limit = 100

func (rf *requestFetcher) fetch(ctx context.Context, trs chan<- *common.Transcript) error {
	logger := inslogger.FromContext(ctx)

	for {
		stats.Record(ctx, metrics.RequestFetcherFetchCall.M(1))
		reqRefs, err := rf.artifactsManager.GetPendings(ctx, rf.object)
		if err != nil {
			if err == insolar.ErrNoPendingRequest {
				logger.Debug("no more pendings on ledger")
				rf.broker.NoMoreRequestsOnLedger(ctx)
				return nil
			}
			return err
		}

		addedCount := 0
		for _, reqRef := range reqRefs {
			if !reqRef.IsRecordScope() {
				logger.Errorf("skipping request with bad reference, ref=%s", reqRef.String())
				continue
			} else if rf.broker.IsKnownRequest(ctx, reqRef) {
				logger.Debug("skipping known request ", reqRef.String())
				stats.Record(ctx, metrics.RequestFetcherFetchKnown.M(1))
				continue
			}
			addedCount++

			logger.Debug("getting request from ledger")
			stats.Record(ctx, metrics.RequestFetcherFetchUnique.M(1))
			request, err := rf.artifactsManager.GetAbandonedRequest(ctx, rf.object, reqRef)
			if err != nil {
				return errors.Wrap(err, "couldn't get request")
			}

			if rf.isAborted() {
				logger.Debug("request fetcher was aborted, not returning request")
				return nil
			}

			switch v := request.(type) {
			case *record.IncomingRequest:
				if err := checkIncomingRequest(ctx, v); err != nil {
					err = errors.Wrap(err, "failed to check incoming request")
					logger.Error(err.Error())

					continue
				}
				tr := common.NewTranscriptCloneContext(ctx, reqRef, *v)
				trs <- tr
			case *record.OutgoingRequest:
				if err := checkOutgoingRequest(ctx, v); err != nil {
					err = errors.Wrap(err, "failed to check outgoing request")
					logger.Error(err.Error())

					continue
				}
				// FIXME: limit there may slow down things, placing "go" here is not good too
				rf.outgoingsSender.SendAbandonedOutgoingRequest(ctx, reqRef, v)
			default:
				logger.Error("requestFetcher fetched unknown request")
			}
		}
		if addedCount == 0 {
			if len(reqRefs) < limit {
				logger.Debug("we guess that ledger has no more requests")
				rf.broker.NoMoreRequestsOnLedger(ctx)
				return nil
			} else {
				logger.Warn("we can not get more requests")
				return nil
			}
		}
	}
}
