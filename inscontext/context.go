/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package inscontext

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

// Ctx is extension of common Context with added Logger.
type Ctx struct {
	ctx context.Context
	log core.Logger
}

// Deadline implements Deadline() from Context interface.
func (ctx *Ctx) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

// Done implements Done() from Context interface.
func (ctx *Ctx) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

// Err implements Err() from Context interface.
func (ctx *Ctx) Err() error {
	return ctx.ctx.Err()
}

// Value implements Value() from Context interface.
func (ctx *Ctx) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

// NewCtxFromContext creates new Ctx from context.Context.
func NewCtxFromContext(parent context.Context) *Ctx {
	ctx, ok := parent.(*Ctx)
	if !ok {
		ctx = &Ctx{
			ctx: parent,
			log: log.GlobalLogger,
		}
	}
	return ctx
}

// WithCancel is the same as context.WithCancel() but always returns Ctx.
func WithCancel(parent context.Context) (ctx *Ctx, cancel context.CancelFunc) {
	ctx = NewCtxFromContext(parent)
	ctx.ctx, cancel = context.WithCancel(ctx.ctx)
	return ctx, cancel
}

// WithDeadline is the same as context.WithDeadline() but always returns Ctx.
func WithDeadline(parent context.Context, d time.Time) (ctx *Ctx, cancel context.CancelFunc) {
	ctx = NewCtxFromContext(parent)
	ctx.ctx, cancel = context.WithCancel(ctx.ctx)
	return ctx, cancel
}

// WithTimeout is the same as context.WithTimeout() but always returns Ctx.
func WithTimeout(parent context.Context, timeout time.Duration) (ctx *Ctx, cancel context.CancelFunc) {
	ctx = NewCtxFromContext(parent)
	ctx.ctx, cancel = context.WithTimeout(ctx.ctx, timeout)
	return ctx, cancel
}

// Background is the same as context.Background() but always returns Ctx.
func Background() *Ctx {
	return &Ctx{
		ctx: context.Background(),
		log: log.GlobalLogger,
	}
}

// TODO is the same as context.TODO() but always returns Ctx.
func TODO() *Ctx {
	return &Ctx{
		ctx: context.TODO(),
		log: log.GlobalLogger,
	}
}

// WithValue is the same as context.WithValue() but always returns Ctx.
func WithValue(parent context.Context, key, val interface{}) *Ctx {
	ctx := NewCtxFromContext(parent)
	ctx.ctx = context.WithValue(ctx.ctx, key, val)
	return ctx
}

type key int

const (
	traceKey = key(0)
)

// WithTrace returns a copy of parent with added trace mark.
func WithTrace(parent context.Context, trace string) *Ctx {
	ctx := NewCtxFromContext(parent)
	return WithValue(ctx, traceKey, trace)
}

func (ctx *Ctx) TraceID() string {
	return ctx.Value(traceKey).(string)
}

// WithLog returns *Ctx with provided core.Logger,
// if parent is not *Ctx instance returns a new Ctx instance
// overwise just set logger for provided Ctx.
func WithLog(parent context.Context, clog core.Logger) *Ctx {
	ctx := NewCtxFromContext(parent)
	ctx.log = clog
	return ctx
}

// Log returns core.Logger provided by *Ctx.
func (ctx *Ctx) Log() core.Logger {
	return ctx.log
}
