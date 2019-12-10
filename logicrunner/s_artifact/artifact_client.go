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

package s_artifact

import (
	"context"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type ArtifactClientService interface {
	artifacts.Client
}

type ArtifactClientServiceAdapter struct {
	svc  ArtifactClientService
	exec smachine.ExecutionAdapter
}

func (a *ArtifactClientServiceAdapter) PrepareSync(
	ctx smachine.ExecutionContext,
	fn func(svc ArtifactClientService),
) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func(interface{}) smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ArtifactClientServiceAdapter) PrepareAsync(
	ctx smachine.ExecutionContext,
	fn func(svc ArtifactClientService) smachine.AsyncResultFunc,
) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func(interface{}) smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

func (a *ArtifactClientServiceAdapter) PrepareNotify(
	ctx smachine.ExecutionContext,
	fn func(svc ArtifactClientService),
) smachine.NotifyRequester {
	return a.exec.PrepareNotify(ctx, func(interface{}) {
		fn(a.svc)
	})
}

type artifactClientService struct {
	artifacts.Client
}

func CreateArtifactClientService(client artifacts.Client) *ArtifactClientServiceAdapter {
	ctx := context.Background()
	ae, ch := smachine.NewCallChannelExecutor(ctx, -1, false, 16)
	smachine.StartDynamicChannelWorker(ctx, ch, nil)

	return &ArtifactClientServiceAdapter{
		svc: artifactClientService{
			Client: client,
		},
		exec: smachine.NewExecutionAdapter("ArtifactClientService", ae),
	}
}
