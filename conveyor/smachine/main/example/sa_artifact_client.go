//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package example

import (
	"context"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
)

type ArtifactClientService interface {
	GetLatestValidatedStateAndCode() (state, code ArtifactBinary)
}

type ArtifactBinary interface {
	GetReference() insolar.Reference
	GetCacheId() ArtifactCacheId
}

type ArtifactClientServiceAdapter struct {
	svc  ArtifactClientService
	exec smachine.ExecutionAdapter
}

func (a *ArtifactClientServiceAdapter) PrepareSync(ctx smachine.ExecutionContext, fn func(svc ArtifactClientService)) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func(interface{}) smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ArtifactClientServiceAdapter) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc ArtifactClientService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func(interface{}) smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

func CreateArtifactClientService() *ArtifactClientServiceAdapter {
	ctx := context.Background()
	ae, ch := smachine.NewCallChannelExecutor(ctx, 0, false, 5)
	ea := smachine.NewExecutionAdapter("ServiceA", ae)

	smachine.StartChannelWorker(ctx, ch, nil)
	return &ArtifactClientServiceAdapter{&artifactClientService{}, ea}
}

var _ ArtifactClientService = &artifactClientService{}

type artifactClientService struct {
}

func (*artifactClientService) GetLatestValidatedStateAndCode() (state, code ArtifactBinary) {
	panic("implement me")
}
