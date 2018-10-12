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

package context

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

// Ctx is extension of common Context with added Logger.
type Ctx struct {
	goCtx context.Context
	log   core.Logger
}

// Deadline implements Deadline() from Context interface.
func (ctx *Ctx) Deadline() (deadline time.Time, ok bool) {
	return ctx.goCtx.Deadline()
}

// Done implements Done() from Context interface.
func (ctx *Ctx) Done() <-chan struct{} {
	return ctx.goCtx.Done()
}

// Err implements Err() from Context interface.
func (ctx *Ctx) Err() error {
	return ctx.goCtx.Err()
}

// Value implements Value() from Context interface.
func (ctx *Ctx) Value(key interface{}) interface{} {
	return ctx.goCtx.Value(key)
}

// NewCtxFromContext creates new Ctx from context.Context
func NewCtxFromContext(parent context.Context) *Ctx {
	ctx, ok := parent.(*Ctx)
	if !ok {
		ctx = &Ctx{
			goCtx: parent,
			log:   log.GlobalLogger,
		}
	}
	return ctx
}

// WithCancel is the same as context.WithCancel() but always returns Ctx
func WithCancel(parent context.Context) (ctx *Ctx, cancel context.CancelFunc) {
	ctx = NewCtxFromContext(parent)
	ctx.goCtx, cancel = context.WithCancel(ctx.goCtx)
	return ctx, cancel
}

// WithDeadline is the same as context.WithDeadline() but always returns Ctx
func WithDeadline(parent context.Context, d time.Time) (ctx *Ctx, cancel context.CancelFunc) {
	ctx = NewCtxFromContext(parent)
	ctx.goCtx, cancel = context.WithCancel(ctx.goCtx)
	return ctx, cancel
}
