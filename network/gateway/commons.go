/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package gateway

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"

	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/network/state.messageBusLocker -o ./ -s _mock.go
type messageBusLocker interface {
	Lock(ctx context.Context)
	Unlock(ctx context.Context)
}

type commons struct {
	counter   uint64
	state     insolar.NetworkState
	stateLock sync.RWMutex
	span      *trace.Span

	Network  network.Gatewayer
	MBLocker messageBusLocker
}

// Acquire increases lock counter and locks message bus if it wasn't lock before
func Acquire(ctx context.Context, c commons) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Acquire")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Acquire in NetworkSwitcher: ", c.counter)
	c.counter = c.counter + 1
	if c.counter-1 == 0 {
		inslogger.FromContext(ctx).Info("Lock MB")
		ctx, c.span = instracer.StartSpan(context.Background(), "GIL Lock (Lock MB)")
		c.MBLocker.Lock(ctx)
	}
}

// Release decreases lock counter and unlocks message bus if it wasn't lock by someone else
func Release(ctx context.Context, c commons) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Release")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Release in NetworkSwitcher: ", c.counter)
	if c.counter == 0 {
		panic("Trying to unlock without locking")
	}
	c.counter = c.counter - 1
	if c.counter == 0 {
		inslogger.FromContext(ctx).Info("Unlock MB")
		c.MBLocker.Unlock(ctx)
		c.span.End()
	}
}
