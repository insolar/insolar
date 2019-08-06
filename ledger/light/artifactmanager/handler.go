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

package artifactmanager

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow/dispatcher"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/network/storage"
)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	PCS            insolar.PlatformCryptographyScheme `inject:""`
	JetCoordinator jet.Coordinator                    `inject:""`
	JetStorage     jet.Storage                        `inject:""`
	JetReleaser    hot.JetReleaser                    `inject:""`
	DropModifier   drop.Modifier                      `inject:""`
	IndexLocker    object.IndexLocker                 `inject:""`
	Records        object.AtomicRecordStorage         `inject:""`
	HotDataWaiter  hot.JetWaiter                      `inject:""`

	WriteAccessor      hot.WriteAccessor
	IndexStorage       object.MemoryIndexStorage
	PulseCalculator    storage.PulseCalculator
	JetTreeUpdater     executor.JetFetcher
	Sender             bus.Sender
	FlowDispatcher     *dispatcher.Dispatcher
	FilamentCalculator *executor.FilamentCalculatorDefault
	RequestChecker     *executor.RequestCheckerDefault

	conf     *configuration.Ledger
	handlers map[insolar.MessageType]insolar.MessageHandler
}

// NewMessageHandler creates new handler.
func NewMessageHandler(
	conf *configuration.Ledger,
) *MessageHandler {

	h := &MessageHandler{
		handlers: map[insolar.MessageType]insolar.MessageHandler{},
		conf:     conf,
	}

	dep := &proc.Dependencies{
		FetchJet: func(p *proc.FetchJet) {
			p.Dep(
				h.JetStorage,
				h.JetTreeUpdater,
				h.JetCoordinator,
				h.Sender,
			)
		},
		WaitHot: func(p *proc.WaitHot) {
			p.Dep(
				h.HotDataWaiter,
				h.Sender,
			)
		},
		EnsureIndex: func(p *proc.EnsureIndex) {
			p.Dep(
				h.IndexLocker,
				h.IndexStorage,
				h.JetCoordinator,
				h.Sender,
			)
		},
		SetRequest: func(p *proc.SetRequest) {
			p.Dep(
				h.WriteAccessor,
				h.FilamentCalculator,
				h.Sender,
				h.IndexLocker,
				h.IndexStorage,
				h.Records,
				h.PCS,
				h.RequestChecker,
				h.JetCoordinator,
			)
		},
		SetResult: func(p *proc.SetResult) {
			p.Dep(
				h.WriteAccessor,
				h.Sender,
				h.IndexLocker,
				h.FilamentCalculator,
				h.Records,
				h.IndexStorage,
				h.PCS,
			)
		},
		HasPendings: func(p *proc.HasPendings) {
			p.Dep(
				h.IndexStorage,
				h.Sender,
			)
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep(
				h.JetCoordinator,
				h.JetStorage,
				h.JetTreeUpdater,
				h.Records,
				h.IndexStorage,
				h.Sender,
			)
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep(
				h.Records,
				h.JetCoordinator,
				h.JetTreeUpdater,
				h.Sender,
			)
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep(
				h.Records,
				h.Sender,
				h.JetCoordinator,
				h.JetTreeUpdater,
			)
		},
		GetPendings: func(p *proc.GetPendings) {
			p.Dep(
				h.FilamentCalculator,
				h.Sender,
			)
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep(
				h.JetStorage,
				h.Sender,
			)
		},
		HotObjects: func(p *proc.HotObjects) {
			p.Dep(
				h.DropModifier,
				h.IndexStorage,
				h.JetStorage,
				h.JetTreeUpdater,
				h.JetReleaser,
				h.JetCoordinator,
				h.PulseCalculator,
				h.Sender,
			)
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(
				h.Sender,
				h.FilamentCalculator,
			)
		},
		PassState: func(p *proc.PassState) {
			p.Dep(
				h.Records,
				h.Sender,
			)
		},
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(h.PCS)
		},
		SetCode: func(p *proc.SetCode) {
			p.Dep(
				h.WriteAccessor,
				h.Records,
				h.PCS,
				h.Sender,
			)
		},
	}

	initHandle := func(msg *message.Message) *handle.Init {
		return handle.NewInit(dep, h.Sender, msg)
	}

	h.FlowDispatcher = dispatcher.NewDispatcher(func(msg *message.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg *message.Message) flow.Handle {
		return initHandle(msg).Future
	}, func(msg *message.Message) flow.Handle {
		return initHandle(msg).Past
	})
	return h
}

func (h *MessageHandler) BeginPulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.BeginPulse(ctx, pn)
}

func (h *MessageHandler) ClosePulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.ClosePulse(ctx, pn)
}
