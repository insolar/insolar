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

package common

import (
	"context"
	"sync/atomic"
)

type RunnerRestrainer interface {
	TryTake(ctx context.Context) bool
	Release(ctx context.Context)
	Available(ctx context.Context) int
}

type atomicLimiter struct {
	count      int32
	countLimit int32
}

func (b *atomicLimiter) TryTake(ctx context.Context) bool {
	for {
		if count := atomic.LoadInt32(&b.count); count < b.countLimit {
			success := atomic.CompareAndSwapInt32(&b.count, count, count+1)
			if success {
				return true
			}
		} else {
			return false
		}
	}
}

func (b *atomicLimiter) Release(ctx context.Context) {
	atomic.AddInt32(&b.count, -1)
}

func (b *atomicLimiter) Available(ctx context.Context) int {
	return int(b.countLimit - atomic.LoadInt32(&b.count))
}

func NewRunnerRestrainer(limit int) RunnerRestrainer {
	return &atomicLimiter{countLimit: int32(limit)}
}
