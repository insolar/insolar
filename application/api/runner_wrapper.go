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

package api

import (
	"context"

	"github.com/insolar/component-manager"
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
