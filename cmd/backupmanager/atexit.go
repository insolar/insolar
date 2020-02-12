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

package main

import (
	"os"
	"sync"

	"github.com/insolar/insolar/insolar"
)

type ExitContextCallback func() error

type ExitContext struct {
	logger    insolar.Logger
	callbacks map[string]ExitContextCallback
	once      sync.Once
}

var (
	exitContext     ExitContext
	exitContextOnce sync.Once
)

func InitExitContext(logger insolar.Logger) {
	initExitContext := func() {
		exitContext.callbacks = make(map[string]ExitContextCallback)
		exitContext.logger = logger
	}

	exitContextOnce.Do(initExitContext)
}

func AtExit(name string, cb ExitContextCallback) {
	exitContext.callbacks[name] = cb
}

func Exit(code int) {
	exit := func() {
		for name, cb := range exitContext.callbacks {
			err := cb()
			if err != nil {
				exitContext.logger.Errorf("Failed to call atexit %s: %s", name, err.Error())
			}
		}
		os.Exit(code)
	}
	exitContext.once.Do(exit)
}
