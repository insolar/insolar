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

package s_jet_storage

import (
	"context"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
)

// TODO[bigbes]: port it to state machine
type JetStorageService interface {
	jet.Storage
}

type JetStorageServiceAdapter struct {
	svc  JetStorageService
	exec smachine.ExecutionAdapter
}

func (a *JetStorageServiceAdapter) PrepareSync(ctx smachine.ExecutionContext, fn func(svc JetStorageService)) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func(interface{}) smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *JetStorageServiceAdapter) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc JetStorageService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func(interface{}) smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

func (a *JetStorageServiceAdapter) PrepareNotify(ctx smachine.ExecutionContext, fn func(svc JetStorageService)) smachine.NotifyRequester {
	return a.exec.PrepareNotify(ctx, func(interface{}) { fn(a.svc) })
}

type jetStorageService struct {
	jet.Storage
	Accessor pulse.Accessor
}

func CreateJetStorageService(JetStorage jet.Storage) *JetStorageServiceAdapter {
	ctx := context.Background()
	ae, ch := smachine.NewCallChannelExecutor(ctx, 0, false, 5)
	smachine.StartChannelWorker(ctx, ch, nil)

	return &JetStorageServiceAdapter{
		svc: jetStorageService{
			Storage: JetStorage,
		},
		exec: smachine.NewExecutionAdapter("JetStorage", ae),
	}
}
