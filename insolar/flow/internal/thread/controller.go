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
	lock    sync.Mutex
	cancel  chan struct{}
	begin   chan struct{}
	process chan struct{}
}

func NewController() *Controller {
	process := make(chan struct{})
	close(process)
	return &Controller{cancel: make(chan struct{}), begin: make(chan struct{}), process: process}
}

func (c *Controller) Cancel() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.cancel
}

func (c *Controller) Begin() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.begin
}

func (c *Controller) Process() <-chan struct{} {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.process
}

func (c *Controller) BeginPulse() {
	c.lock.Lock()
	defer c.lock.Unlock()

	toBegin := c.begin
	c.begin = make(chan struct{})
	close(toBegin)

	c.cancel = make(chan struct{})
	close(c.process)
}

func (c *Controller) ClosePulse() {
	c.lock.Lock()
	defer c.lock.Unlock()

	close(c.cancel)
	c.process = make(chan struct{})
}
