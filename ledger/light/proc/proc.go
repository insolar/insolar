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

package proc

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

type Dependencies struct {
	FetchJet       func(*FetchJet)
	WaitHot        func(*WaitHot)
	EnsureIndex    func(*EnsureIndex)
	SendObject     func(*SendObject)
	GetCode        func(*GetCode)
	GetRequest     func(*GetRequest)
	GetRequestInfo func(*SendRequestInfo)
	SetRequest     func(*SetRequest)
	SetResult      func(*SetResult)
	GetPendings    func(*GetPendings)
	GetJet         func(*GetJet)
	HotObjects     func(*HotObjects)
	PassState      func(*PassState)
	CalculateID    func(*CalculateID)
	SetCode        func(*SetCode)
	SendRequests   func(*SendRequests)
	HasPendings    func(*HasPendings)
}

func NewDependencies(
	// Common components.
	pcs insolar.PlatformCryptographyScheme,
	jetCoordinator jet.Coordinator,
	jetStorage jet.Storage,
	pulseCalculator pulse.Calculator,
	sender bus.Sender,

	// Ledger components.
	dropModifier drop.Modifier,
	indexLocker object.IndexLocker,
	recordStorage object.AtomicRecordStorage,
	indexStorage object.MemoryIndexStorage,

	// Executor components.
	jetReleaser executor.JetReleaser,
	hotWaiter executor.JetWaiter,
	writeAccessor executor.WriteAccessor,
	jetFetcher executor.JetFetcher,
	filaments executor.FilamentCalculator,
	requestChecker executor.RequestChecker,
) *Dependencies {
	dep := &Dependencies{
		FetchJet: func(p *FetchJet) {
			p.Dep(
				jetStorage,
				jetFetcher,
				jetCoordinator,
				sender,
			)
		},
		WaitHot: func(p *WaitHot) {
			p.Dep(
				hotWaiter,
				sender,
			)
		},
		EnsureIndex: func(p *EnsureIndex) {
			p.Dep(
				indexLocker,
				indexStorage,
				jetCoordinator,
				sender,
			)
		},
		SetRequest: func(p *SetRequest) {
			p.Dep(
				writeAccessor,
				filaments,
				sender,
				indexLocker,
				indexStorage,
				recordStorage,
				pcs,
				requestChecker,
				jetCoordinator,
			)
		},
		SetResult: func(p *SetResult) {
			p.Dep(
				writeAccessor,
				sender,
				indexLocker,
				filaments,
				recordStorage,
				indexStorage,
				pcs,
			)
		},
		HasPendings: func(p *HasPendings) {
			p.Dep(
				indexStorage,
				sender,
			)
		},
		SendObject: func(p *SendObject) {
			p.Dep(
				jetCoordinator,
				jetStorage,
				jetFetcher,
				recordStorage,
				indexStorage,
				sender,
			)
		},
		GetCode: func(p *GetCode) {
			p.Dep(
				recordStorage,
				jetCoordinator,
				jetFetcher,
				sender,
			)
		},
		GetRequest: func(p *GetRequest) {
			p.Dep(
				recordStorage,
				sender,
				jetCoordinator,
				jetFetcher,
			)
		},
		GetRequestInfo: func(p *SendRequestInfo) {
			p.Dep(
				filaments,
				sender,
				indexLocker,
			)
		},
		GetPendings: func(p *GetPendings) {
			p.Dep(
				filaments,
				sender,
			)
		},
		GetJet: func(p *GetJet) {
			p.Dep(
				jetStorage,
				sender,
			)
		},
		HotObjects: func(p *HotObjects) {
			p.Dep(
				dropModifier,
				indexStorage,
				jetStorage,
				jetFetcher,
				jetReleaser,
				jetCoordinator,
				pulseCalculator,
				sender,
			)
		},
		SendRequests: func(p *SendRequests) {
			p.Dep(
				sender,
				filaments,
			)
		},
		PassState: func(p *PassState) {
			p.Dep(
				recordStorage,
				sender,
			)
		},
		CalculateID: func(p *CalculateID) {
			p.Dep(pcs)
		},
		SetCode: func(p *SetCode) {
			p.Dep(
				writeAccessor,
				recordStorage,
				pcs,
				sender,
			)
		},
	}
	return dep
}

// NewDependenciesMock returns all dependencies for handlers.
// It's all empty.
// Use it ONLY for tests.
func NewDependenciesMock() *Dependencies {
	return &Dependencies{
		FetchJet:       func(*FetchJet) {},
		WaitHot:        func(*WaitHot) {},
		EnsureIndex:    func(*EnsureIndex) {},
		SendObject:     func(*SendObject) {},
		GetCode:        func(*GetCode) {},
		SetRequest:     func(*SetRequest) {},
		SetResult:      func(*SetResult) {},
		GetPendings:    func(*GetPendings) {},
		GetJet:         func(*GetJet) {},
		HotObjects:     func(*HotObjects) {},
		PassState:      func(*PassState) {},
		CalculateID:    func(*CalculateID) {},
		SetCode:        func(*SetCode) {},
		SendRequests:   func(*SendRequests) {},
		HasPendings:    func(*HasPendings) {},
		GetRequestInfo: func(*SendRequestInfo) {},
	}
}
