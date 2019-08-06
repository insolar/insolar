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
