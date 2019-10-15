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

	smachine "github.com/insolar/insolar/conveyor/smachinev2"
)

type ArtifactClientService interface {
	GetLatestValidatedStateAndCode() (state, code ArtifactBinary)
}

type ArtifactBinary interface {
	GetReference() //reference
	GetCacheId() ArtifactCacheId
}

type ArtifactClientServiceAdapter struct {
	svc  ArtifactClientService
	exec smachine.ExecutionAdapter
}

func (a *ArtifactClientServiceAdapter) PrepareSync(ctx smachine.ExecutionContext, fn func(svc ArtifactClientService)) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func() smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ArtifactClientServiceAdapter) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc ArtifactClientService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

func CreateArtifactClientService() *ArtifactClientServiceAdapter {
	ach := NewChannelAdapter(context.Background(), 0, -1)
	ea := smachine.NewExecutionAdapter("ServiceA", &ach)

	go func() {
		for {
			select {
			case <-ach.Context().Done():
				return
			case t := <-ach.Channel():
				t.RunAndSendResult()
			}
		}
	}()

	return &ArtifactClientServiceAdapter{artifactClientService{}, ea}
}

type artifactClientService struct {
}

func (artifactClientService) GetLatestValidatedStateAndCode() (state, code ArtifactBinary) {
	return nil, nil
}
