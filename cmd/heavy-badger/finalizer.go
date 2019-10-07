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

package main

import (
	"fmt"
	"os"
	"os/signal"
)

type finalizersHolder struct {
	finalizers []func() error
}

func (f *finalizersHolder) add(fn func() error) {
	f.finalizers = append(f.finalizers, fn)
}

func (f *finalizersHolder) run() {
	for _, fn := range f.finalizers {
		err := fn()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		}
	}
}

func (f *finalizersHolder) onSignals(signals ...os.Signal) <-chan struct{} {
	ret := make(chan struct{})
	if len(signals) < 1 {
		close(ret)
		return ret
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)

	go func() {
		<-sigs
		f.run()
		close(ret)
	}()
	return ret
}
