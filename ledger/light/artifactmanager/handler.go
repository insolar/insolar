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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/blob"
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

	BlobModifier blob.Modifier `inject:""`
	BlobAccessor blob.Accessor `inject:""`
	Blobs        blob.Storage  `inject:""`

	IndexLocker object.IndexLocker `inject:""`

	Records object.RecordStorage `inject:""`
	Nodes   node.Accessor        `inject:""`

	HotDataWaiter hot.JetWaiter   `inject:""`
	JetReleaser   hot.JetReleaser `inject:""`

	WriteAccessor hot.WriteAccessor

	IndexStorage object.IndexStorage

	PulseCalculator storage.PulseCalculator

	conf           *configuration.Ledger
	jetTreeUpdater jet.Fetcher

	Sender         wbus.Sender
	FlowDispatcher *dispatcher.Dispatcher
	handlers       map[insolar.MessageType]insolar.MessageHandler

	filamentModifier   *executor.FilamentModifierDefault
	filamentCalculator *executor.FilamentCalculatorDefault
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
			p.Dep.JetUpdater = h.jetTreeUpdater
			p.Dep.JetFetcher = h.jetTreeUpdater
		},
		WaitHot: func(p *proc.WaitHot) {
			p.Dep.Waiter = h.HotDataWaiter
		},
		GetIndex: func(p *proc.EnsureIndex) {
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
		},
		CheckJet: func(p *proc.CheckJet) {
			p.Dep.JetAccessor = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		WaitHotWM: func(p *proc.WaitHotWM) {
			p.Dep.Waiter = h.HotDataWaiter
			p.Dep.Sender = h.Sender
		},
		GetIndexWM: func(p *proc.EnsureIndexWM) {
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
		SetActivationRequest: func(p *proc.SetActivationRequest) {
			p.Dep(
				h.PCS,
				h.WriteAccessor,
				h.Records,
				h.Sender,
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
		SetBlob: func(p *proc.SetBlob) {
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.PCS = h.PCS
			p.Dep.WriteAccessor = h.WriteAccessor
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Blobs = h.Blobs
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.Records
			p.Dep.Sender = h.Sender
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.RecordAccessor = h.Records
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep.RecordAccessor = h.Records
		},
		UpdateObject: func(p *proc.UpdateObject) {
			p.Dep.RecordModifier = h.Records
			p.Dep.Bus = h.Bus
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.PCS = h.PCS
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.WriteAccessor = h.WriteAccessor
			p.Dep.Filaments = h.filamentModifier
		},
		GetChildren: func(p *proc.GetChildren) {
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.DelegationTokenFactory = h.DelegationTokenFactory
			p.Dep.RecordAccessor = h.Records
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetTreeUpdater = h.jetTreeUpdater
		},
		RegisterChild: func(p *proc.RegisterChild) {
			p.Dep.IndexLocker = h.IndexLocker
			p.Dep.IndexAccessor = h.IndexStorage
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.JetCoordinator = h.JetCoordinator
			p.Dep.RecordModifier = h.Records
			p.Dep.PCS = h.PCS
		},
		GetPendingRequests: func(p *proc.GetPendingRequests) {
			p.Dep(h.IndexStorage)
		},
		GetPendingRequestID: func(p *proc.GetPendingRequestID) {
			p.Dep(h.filamentCalculator)
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep.Jets = h.JetStorage
		},
		HotData: func(p *proc.HotData) {
			p.Dep.DropModifier = h.DropModifier
			p.Dep.MessageBus = h.Bus
			p.Dep.IndexModifier = h.IndexStorage
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.JetReleaser = h.JetReleaser
			p.Dep.Calculator = h.PulseCalculator
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(h.Sender, h.filamentCalculator)
		},
		PassState: func(p *proc.PassState) {
			p.Dep.Blobs = h.BlobAccessor
			p.Dep.Sender = h.Sender
			p.Dep.Records = h.Records
		},
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(h.PCS)
		},
		SetCode: func(p *proc.SetCode) {
			p.Dep(h.WriteAccessor, h.Records, h.BlobModifier, h.PCS, h.Sender)
		},
		GetDelegate: func(p *proc.GetDelegate) {
			p.Dep.IndexAccessor = h.IndexStorage
		},
	}

	initHandle := func(msg bus.Message) *handle.Init {
		return handle.NewInit(dep, h.Sender, msg)
	}

	h.FlowDispatcher = dispatcher.NewDispatcher(func(msg bus.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg bus.Message) flow.Handle {
		return initHandle(msg).Future
	})
	return h
}

// Init initializes handlers and middleware.
func (h *MessageHandler) Init(ctx context.Context) error {
	h.jetTreeUpdater = jet.NewFetcher(h.Nodes, h.JetStorage, h.Bus, h.JetCoordinator)
	h.filamentCalculator = executor.NewFilamentCalculator(
		h.IndexStorage,
		h.Records,
		h.JetCoordinator,
		h.jetTreeUpdater,
		h.Sender,
	)
	h.filamentModifier = executor.NewFilamentModifier(
		h.IndexStorage,
		h.Records,
		h.PCS,
		h.filamentCalculator,
	)
	h.setHandlersForLight()

	return nil
}

func (h *MessageHandler) OnPulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.ChangePulse(ctx, pn)
}

func (h *MessageHandler) setHandlersForLight() {
	// Generic.

	h.Bus.MustRegister(insolar.TypeGetCode, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetObject, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeUpdateObject, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetDelegate, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(insolar.TypeGetChildren, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(insolar.TypeSetRecord, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeRegisterChild, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeSetBlob, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetPendingRequests, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetJet, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeHotRecords, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetRequest, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetPendingRequestID, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(insolar.TypeValidateRecord, h.handleValidateRecord)
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	return &reply.OK{}, nil
}
