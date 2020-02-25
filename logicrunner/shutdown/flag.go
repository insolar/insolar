// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
