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

package shutdown

import (
	"context"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/shutdown.Flag -o ./ -s _mock.go -g
type Flag interface {
	Stop(ctx context.Context) func()
	Done(ctx context.Context, isDone func() bool)

	IsStopped() bool
}

type flag struct {
	stopLock  sync.Mutex
	isStopped bool

	stopChannel chan struct{}
}

func NewFlag() Flag {
	return &flag{
		stopLock:    sync.Mutex{},
		isStopped:   false,
		stopChannel: make(chan struct{}),
	}
}

func (g *flag) IsStopped() bool {
	g.stopLock.Lock()
	defer g.stopLock.Unlock()

	return g.isStopped
}

func (g *flag) Stop(ctx context.Context) func() {
	logger := inslogger.FromContext(ctx)
	logger.Debug("shutdown initiated")

	g.stopLock.Lock()
	defer g.stopLock.Unlock()

	g.isStopped = true

	return func() {
		logger.Debug("waiting for successful shutdown")
		<-g.stopChannel
		logger.Debug("waited for shutdown to be finished")
	}
}

func (g *flag) Done(ctx context.Context, isDone func() bool) {
	logger := inslogger.FromContext(ctx)

	g.stopLock.Lock()
	defer g.stopLock.Unlock()

	if g.isStopped && isDone() {
		logger.Debug("ready to shut down")

		select {
		case _, ok := <-g.stopChannel:
			if ok {
				panic("unexpected message was written to channel")
			}
		default:
			close(g.stopChannel)
		}
	}
}
