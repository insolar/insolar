package api

import (
	"context"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/insolar"
)

type RunnerWrapper struct {
	API      insolar.APIRunner
	AdminAPI insolar.APIRunner
}

// NewWrapper is C-tor for wrapper of API Runner
func NewWrapper(publicAPI, adminAPI insolar.APIRunner) *RunnerWrapper {
	return &RunnerWrapper{
		API:      publicAPI,
		AdminAPI: adminAPI,
	}
}

// Start runs api servers
func (w *RunnerWrapper) Start(ctx context.Context) error {
	if starter, ok := w.API.(component.Starter); ok {
		err := starter.Start(ctx)
		if err != nil {
			return err
		}
	}
	if starter, ok := w.AdminAPI.(component.Starter); ok {
		err := starter.Start(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Start stops api servers
func (w *RunnerWrapper) Stop(ctx context.Context) error {
	var (
		first  error
		second error
	)
	if stopper, ok := w.API.(component.Stopper); ok {
		first = stopper.Stop(ctx)
	}
	if stopper, ok := w.AdminAPI.(component.Stopper); ok {
		second = stopper.Stop(ctx)
	}
	if first != nil {
		return first
	}
	if second != nil {
		return second
	}
	return nil
}
