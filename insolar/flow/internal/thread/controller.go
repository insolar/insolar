// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package thread

import (
	"sync"
)

type Controller struct {
	lock sync.Mutex
	// cancel will be closed on ClosePulse()
	cancel chan struct{}
	// canBegin will be closed on BeginPulse()
	canBegin chan struct{}
	// canProcess will be closed on BeginPulse() and new instance will be opened again on ClosePulse()
	canProcess chan struct{}
}

func NewController() *Controller {
	process := make(chan struct{})
	close(process)
	return &Controller{cancel: make(chan struct{}), canBegin: make(chan struct{}), canProcess: process}
}

func (c *Controller) Cancel() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.cancel
}

func (c *Controller) CanBegin() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.canBegin
}

func (c *Controller) CanProcess() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.canProcess
}

func (c *Controller) BeginPulse() {
	c.lock.Lock()
	defer c.lock.Unlock()

	toBegin := c.canBegin
	c.canBegin = make(chan struct{})
	close(toBegin)

	c.cancel = make(chan struct{})
	close(c.canProcess)
}

func (c *Controller) ClosePulse() {
	c.lock.Lock()
	defer c.lock.Unlock()

	close(c.cancel)
	c.canProcess = make(chan struct{})
}
