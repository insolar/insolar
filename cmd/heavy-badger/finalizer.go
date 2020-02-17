// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
