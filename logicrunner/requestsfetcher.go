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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestsFetcher -o ./ -s _mock.go -g

type RequestsFetcher interface {
	FetchPendings(ctx context.Context)
	Abort(ctx context.Context)
}

type requestsFetcher struct {
	object insolar.Reference

	mu      sync.Mutex
	active  bool
	stopper func()

	broker ExecutionBrokerI
	am     artifacts.Client
	os     OutgoingRequestSender
}

func NewRequestsFetcher(
	obj insolar.Reference, am artifacts.Client, br ExecutionBrokerI, os OutgoingRequestSender,
) RequestsFetcher {
	return &requestsFetcher{
		object: obj,
		broker: br,
		am:     am,
		os:     os,
	}
}

func (rf *requestsFetcher) FetchPendings(ctx context.Context) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.active {
		inslogger.FromContext(ctx).Debug("requests fetcher is active, not starting")
		return
	}

	ctx, cancel := context.WithCancel(ctx)

	rf.active = true
	rf.stopper = cancel

	go func() {
		defer func() {
			rf.mu.Lock()
			defer rf.mu.Unlock()
			rf.active = false
			rf.stopper = nil
		}()

		err := rf.fetch(ctx)
		if err != nil {
			inslogger.FromContext(ctx).Error("couldn't make fetch round: ", err.Error())
		}
	}()
}

func (rf *requestsFetcher) Abort(ctx context.Context) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.active {
		rf.stopper()
	}
}

func (rf *requestsFetcher) fetch(ctx context.Context) error {
	reqRefs, err := rf.am.GetPendings(ctx, rf.object)
	if err != nil {
		if err == insolar.ErrNoPendingRequest {
			rf.broker.NoMoreRequestsOnLedger(ctx)
			return nil
		}
		return err
	}

	logger := inslogger.FromContext(ctx)
	for _, reqRef := range reqRefs {
		if rf.broker.IsKnownRequest(ctx, reqRef) {
			logger.Debug("skipping known request ", reqRef.String())
			continue
		}

		request, err := rf.am.GetAbandonedRequest(ctx, rf.object, reqRef)
		if err != nil {
			logger.Error("couldn't get request: ", err.Error())
			continue
		}

		select {
		case <-ctx.Done():
			logger.Debug("quiting fetching requests, was stopped")
			return nil
		default:
		}

		switch v := request.(type) {
		case *record.IncomingRequest:
			tr := common.NewTranscriptCloneContext(ctx, reqRef, *v)
			rf.broker.AddRequestsFromLedger(ctx, tr)
		case *record.OutgoingRequest:
			rf.os.SendAbandonedOutgoingRequest(ctx, reqRef, v)
		default:
			logger.Error("requestsFetcher.fetch: request is nor IncomingRequest or OutgoingRequest")
		}
	}

	return nil
}
