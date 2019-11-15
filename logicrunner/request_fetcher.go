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
	FetchPendings(ctx context.Context) <-chan *common.Transcript
	Abort(ctx context.Context)
}

type requestFetcher struct {
	mu sync.Mutex

	object insolar.Reference

	aborted      chan struct{}
	stopFetching func()
	skipSlice    []insolar.ID

	artifactsManager artifacts.Client
	outgoingsSender  OutgoingRequestSender
}

func NewRequestsFetcher(
	obj insolar.Reference, am artifacts.Client, os OutgoingRequestSender,
) RequestFetcher {
	aborted := make(chan struct{})
	once := sync.Once{}
	return &requestFetcher{
		object:           obj,
		artifactsManager: am,
		outgoingsSender:  os,
		aborted:          aborted,
		stopFetching:     func() { once.Do(func() { close(aborted) }) },
		skipSlice:        make([]insolar.ID, 0),
	}
}

func (rf *requestFetcher) Abort(ctx context.Context) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

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

func (rf *requestFetcher) FetchPendings(ctx context.Context) <-chan *common.Transcript {
	// TODO: move to const
	trs := make(chan *common.Transcript, 10)

	aborted := make(chan struct{})
	once := sync.Once{}

	rf.mu.Lock()
	rf.aborted = aborted
	rf.stopFetching = func() { once.Do(func() { close(aborted) }) }
	rf.mu.Unlock()

	logger := inslogger.FromContext(ctx)
	logger.Debug("request fetcher starting")

	go func() {
		err := rf.fetch(ctx, trs)
		if err != nil {
			logger.Error("couldn't make fetch round: ", err.Error())
		}
	}()

	return trs
}

const skipSliceSizeLimit = 300

func (rf *requestFetcher) fetch(ctx context.Context, trs chan<- *common.Transcript) error {
	defer close(trs)

	logger := inslogger.FromContext(ctx)

	for {
		stats.Record(ctx, metrics.RequestFetcherFetchCall.M(1))

		if len(rf.skipSlice) > skipSliceSizeLimit {
			rf.skipSlice = rf.skipSlice[len(rf.skipSlice):]
		}

		reqRefs, err := rf.artifactsManager.GetPendings(ctx, rf.object, rf.skipSlice)
		if err != nil {
			if err == insolar.ErrNoPendingRequest {
				logger.Debug("no more pendings on ledger")
				trs <- nil
				return nil
			}
			return err
		}

		for _, reqRef := range reqRefs {
			rf.skipSlice = append(rf.skipSlice, *reqRef.GetLocal())

			if !reqRef.IsRecordScope() {
				logger.Errorf("skipping request with bad reference, ref=%s", reqRef.String())
				continue
			}

			logger.Debug("getting request from ledger")
			stats.Record(ctx, metrics.RequestFetcherFetchUnique.M(1))
			request, err := rf.artifactsManager.GetRequest(ctx, rf.object, reqRef)
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
	}
}
