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

const MaxFetchCount = 20

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestFetcher -o ./ -s _mock.go -g

type RequestFetcher interface {
	FetchPendings(ctx context.Context)
	Abort(ctx context.Context)
}

type requestFetcher struct {
	object insolar.Reference

	isActiveLock sync.Mutex
	isActive     bool
	stopFetching func()

	broker ExecutionBrokerI
	am     artifacts.Client
	os     OutgoingRequestSender
}

func NewRequestsFetcher(
	obj insolar.Reference, am artifacts.Client, br ExecutionBrokerI, os OutgoingRequestSender,
) RequestFetcher {
	return &requestFetcher{
		object: obj,
		broker: br,
		am:     am,
		os:     os,
	}
}

func (rf *requestFetcher) tryTakeActive(ctx context.Context) (context.Context, bool) {
	rf.isActiveLock.Lock()
	defer rf.isActiveLock.Unlock()

	if rf.isActive {
		return ctx, false
	}

	ctx, cancelFunc := context.WithCancel(ctx)

	rf.isActive = true
	rf.stopFetching = cancelFunc

	return ctx, true
}

func (rf *requestFetcher) releaseActive(_ context.Context) {
	rf.isActiveLock.Lock()
	defer rf.isActiveLock.Unlock()

	rf.isActive = false
	rf.stopFetching = nil
}

func (rf *requestFetcher) Abort(ctx context.Context) {
	rf.isActiveLock.Lock()
	defer rf.isActiveLock.Unlock()

	if rf.isActive {
		rf.stopFetching()
	}
}

func (rf *requestFetcher) FetchPendings(ctx context.Context) {
	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"object": rf.object.String(),
	})

	ctx, success := rf.tryTakeActive(ctx)
	if !success {
		logger.Debug("requestFetcher already started")
		return
	}

	go rf.fetchWrapper(ctx)
}

func (rf *requestFetcher) fetchWrapper(ctx context.Context) {
	defer rf.releaseActive(ctx)

	logger := inslogger.FromContext(ctx)
	logger.Debug("requestFetcher starting")

	err := rf.fetch(ctx)
	if err != nil {
		logger.Error("couldn't make fetch round: ", err.Error())
	}
}

func (rf *requestFetcher) fetch(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	stats.Record(ctx, metrics.RequestFetcherFetchCall.M(1))
	reqRefs, err := rf.am.GetPendings(ctx, rf.object)
	if err != nil {
		if err == insolar.ErrNoPendingRequest {
			logger.Debug("no more pendings on ledger")
			rf.broker.NoMoreRequestsOnLedger(ctx)
			return nil
		}
		return err
	}

	var (
		uniqueTaken        = 0
		uniqueLimitReached = false
	)

	for _, reqRef := range reqRefs {
		// limit count of unique and unknown taken requests to MaxFetchCount
		if uniqueTaken >= MaxFetchCount {
			uniqueLimitReached = true
			break
		}

		if !reqRef.IsRecordScope() {
			logger.Errorf("skipping request with bad reference, ref=%s", reqRef.String())
		} else if rf.broker.IsKnownRequest(ctx, reqRef) {
			logger.Debug("skipping known request ", reqRef.String())
			stats.Record(ctx, metrics.RequestFetcherFetchKnown.M(1))
			continue
		}

		logger.Debug("getting abandoned request from ledger")
		stats.Record(ctx, metrics.RequestFetcherFetchUnique.M(1))
		request, err := rf.am.GetAbandonedRequest(ctx, rf.object, reqRef)
		if err != nil {
			return errors.Wrap(err, "couldn't get request")
		}

		select {
		case <-ctx.Done():
			logger.Debug("requestFetcher stopping")
			return nil
		default:
		}

		uniqueTaken++

		switch v := request.(type) {
		case *record.IncomingRequest:
			logger.Debug("get abandoned IncomingRequest from ledger: ", v.String())
			if err := checkIncomingRequest(ctx, v); err != nil {
				err = errors.Wrap(err, "failed to check incoming request")
				logger.Error(err.Error())

				continue
			}
			tr := common.NewTranscriptCloneContext(ctx, reqRef, *v)
			rf.broker.AddRequestsFromLedger(ctx, tr)
		case *record.OutgoingRequest:
			logger.Debug("get abandoned OutgoingRequest from ledger: ", v.String())
			if err := checkOutgoingRequest(ctx, v); err != nil {
				err = errors.Wrap(err, "failed to check outgoing request")
				logger.Error(err.Error())

				continue
			}
			rf.os.SendAbandonedOutgoingRequest(ctx, reqRef, v)
		default:
			logger.Error("requestFetcher fetched unknown request")
		}
	}

	if !uniqueLimitReached {
		logger.Debug("no more pendings on ledger, we've fetched everything")

		rf.broker.NoMoreRequestsOnLedger(ctx)
	}

	return nil
}
