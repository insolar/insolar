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
