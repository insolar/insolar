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

package requestsqueue

import (
	"context"
	"sync"

	"github.com/insolar/insolar/logicrunner/common"
)

type RequestSource int

const (
	FromLedger RequestSource = iota
	FromPreviousExecutor
	FromThisPulse

	numberOfSources
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/requestsqueue.RequestsQueue -o ./ -s _mock.go -g
type RequestsQueue interface {
	// Append adds request(s) from provided source to the queue
	Append(ctx context.Context, from RequestSource, transcripts ...*common.Transcript)
	// TakeFirst extracts first for processing request from the queue, prioritizing
	// requests from ledger, previous executor and only after those fresh requests
	TakeFirst(ctx context.Context) *common.Transcript
	// TakeAllOriginatedFrom extracts all requests that originated from particular source
	TakeAllOriginatedFrom(ctx context.Context, from RequestSource) []*common.Transcript
	// NumberOfOld returns quantity of requests from ledger and previous executor
	// (e.g not fresh, aka old)
	NumberOfOld(ctx context.Context) int
	// Clean cleans queue
	Clean(ctx context.Context)
}

type queue struct {
	lock  sync.Mutex
	lists [numberOfSources][]*common.Transcript
}

func New() RequestsQueue {
	return &queue{}
}

func (q *queue) Append(
	_ context.Context, from RequestSource, trs ...*common.Transcript,
) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.lists[from] = append(q.lists[from], trs...)
}

func (q *queue) TakeFirst(_ context.Context) *common.Transcript {
	q.lock.Lock()
	defer q.lock.Unlock()

	for i := 0; i < int(numberOfSources); i++ {
		if len(q.lists[i]) > 0 {
			var res *common.Transcript
			res, q.lists[i] = q.lists[i][0], q.lists[i][1:]
			return res
		}
	}

	return nil
}

func (q *queue) TakeAllOriginatedFrom(_ context.Context, from RequestSource) []*common.Transcript {
	q.lock.Lock()
	defer q.lock.Unlock()

	res := q.lists[from]
	q.lists[from] = nil
	return res
}

func (q *queue) NumberOfOld(_ context.Context) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.lists[FromLedger]) + len(q.lists[FromPreviousExecutor])
}

func (q *queue) Clean(_ context.Context) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for i := 0; i < int(numberOfSources); i++ {
		q.lists[i] = nil
	}
}
