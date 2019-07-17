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
	"github.com/insolar/insolar/insolar/node"
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
	Bus                    insolar.MessageBus                 `inject:""`
	PCS                    insolar.PlatformCryptographyScheme `inject:""`
	JetCoordinator         jet.Coordinator                    `inject:""`
	CryptographyService    insolar.CryptographyService        `inject:""`
	DelegationTokenFactory insolar.DelegationTokenFactory     `inject:""`
	JetStorage             jet.Storage                        `inject:""`

	DropModifier drop.Modifier `inject:""`

	IndexLocker object.IndexLocker `inject:""`

	Records object.RecordStorage `inject:""`
	Nodes   node.Accessor        `inject:""`

	HotDataWaiter hot.JetWaiter   `inject:""`
	JetReleaser   hot.JetReleaser `inject:""`

	WriteAccessor hot.WriteAccessor

	IndexStorage object.IndexStorage

	PulseCalculator storage.PulseCalculator

	conf           *configuration.Ledger
	JetTreeUpdater jet.Fetcher

	Sender         bus.Sender
	FlowDispatcher *dispatcher.Dispatcher
	handlers       map[insolar.MessageType]insolar.MessageHandler

	filamentModifier   *executor.FilamentModifierDefault
	FilamentCalculator *executor.FilamentCalculatorDefault
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
			p.Dep.JetAccessor = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetUpdater = h.JetTreeUpdater
			p.Dep.JetFetcher = h.JetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		WaitHot: func(p *proc.WaitHot) {
			p.Dep.Waiter = h.HotDataWaiter
			p.Dep.Sender = h.Sender
		},
		GetIndex: func(p *proc.EnsureIndex) {
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
			p.Dep.Sender = h.Sender
		},
		CheckJet: func(p *proc.CheckJet) {
			p.Dep.JetAccessor = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.JetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		WaitHotWM: func(p *proc.WaitHotWM) {
			p.Dep.Waiter = h.HotDataWaiter
			p.Dep.Sender = h.Sender
		},
		EnsureIndex: func(p *proc.EnsureIndexWM) {
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
			p.Dep.Sender = h.Sender
		},
		SetRequest: func(p *proc.SetRequest) {
			p.Dep(
				h.WriteAccessor,
				h.Records,
				h.filamentModifier,
				h.Sender,
				h.IndexLocker,
			)
		},
		SetResult: func(p *proc.SetResult) {
			p.Dep(
				h.WriteAccessor,
				h.Sender,
				h.IndexLocker,
				h.filamentModifier,
			)
		},
		ActivateObject: func(p *proc.ActivateObject) {
			p.Dep(
				h.WriteAccessor,
				h.IndexLocker,
				h.Records,
				h.IndexStorage,
				h.filamentModifier,
				h.Sender,
			)
		},
		DeactivateObject: func(p *proc.DeactivateObject) {
			p.Dep(
				h.WriteAccessor,
				h.IndexLocker,
				h.Records,
				h.IndexStorage,
				h.filamentModifier,
				h.Sender,
			)
		},
		UpdateObject: func(p *proc.UpdateObject) {
			p.Dep(
				h.WriteAccessor,
				h.IndexLocker,
				h.Records,
				h.IndexStorage,
				h.filamentModifier,
				h.Sender,
			)
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.JetTreeUpdater
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.Records
			p.Dep.Sender = h.Sender
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.RecordAccessor = h.Records
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.JetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep.RecordAccessor = h.Records
			p.Dep.Sender = h.Sender
		},
		GetChildren: func(p *proc.GetChildren) {
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.DelegationTokenFactory = h.DelegationTokenFactory
			p.Dep.RecordAccessor = h.Records
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetTreeUpdater = h.JetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		RegisterChild: func(p *proc.RegisterChild) {
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.JetCoordinator = h.JetCoordinator
			p.Dep.RecordModifier = h.Records
			p.Dep.PCS = h.PCS
			p.Dep.Sender = h.Sender
		},
		GetPendingRequests: func(p *proc.GetPendingRequests) {
			p.Dep(h.IndexStorage, h.Sender)
		},
		GetPendingRequestID: func(p *proc.GetPendingRequestID) {
			p.Dep(h.FilamentCalculator, h.Sender)
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Sender = h.Sender
		},
		HotObjects: func(p *proc.HotObjects) {
			p.Dep.DropModifier = h.DropModifier
			p.Dep.MessageBus = h.Bus
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetFetcher = h.JetTreeUpdater
			p.Dep.JetReleaser = h.JetReleaser
			p.Dep.Sender = h.Sender
			p.Dep.Calculator = h.PulseCalculator
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(h.Sender, h.FilamentCalculator)
		},
		PassState: func(p *proc.PassState) {
			p.Dep.Sender = h.Sender
			p.Dep.Records = h.Records
		},
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(h.PCS)
		},
		SetCode: func(p *proc.SetCode) {
			p.Dep(h.WriteAccessor, h.Records, h.PCS, h.Sender)
		},
		GetDelegate: func(p *proc.GetDelegate) {
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.Sender = h.Sender
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

// Init initializes handlers and middleware.
func (h *MessageHandler) Init(ctx context.Context) error {
	h.filamentModifier = executor.NewFilamentModifier(
		h.IndexStorage,
		h.Records,
		h.PCS,
		h.FilamentCalculator,
		h.PulseCalculator,
	)

	return nil
}

func (h *MessageHandler) OnPulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.ChangePulse(ctx, pn)
}
