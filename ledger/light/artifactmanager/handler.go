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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/ledger/light/handle"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"

	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	RecentStorageProvider  recentstorage.Provider             `inject:""`
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

	IDLocker object.IDLocker `inject:""`

	RecordModifier object.RecordModifier `inject:""`
	RecordAccessor object.RecordAccessor `inject:""`
	Nodes          node.Accessor         `inject:""`

	HotDataWaiter hot.JetWaiter   `inject:""`
	JetReleaser   hot.JetReleaser `inject:""`

	WriteAccessor hot.WriteAccessor

	LifelineIndex         object.LifelineIndex
	IndexBucketModifier   object.IndexBucketModifier
	LifelineStateModifier object.LifelineStateModifier
	PendingModifier       object.PendingModifier
	PendingAccessor       object.PendingAccessor

	conf           *configuration.Ledger
	jetTreeUpdater jet.Fetcher

	Sender         wbus.Sender
	FlowDispatcher *dispatcher.Dispatcher
	handlers       map[insolar.MessageType]insolar.MessageHandler
}

// NewMessageHandler creates new handler.
func NewMessageHandler(
	index object.LifelineIndex,
	indexBucketModifier object.IndexBucketModifier,
	indexStateModifier object.LifelineStateModifier,
	pendingModifier object.PendingModifier,
	pendingAccessor object.PendingAccessor,
	conf *configuration.Ledger,
) *MessageHandler {

	h := &MessageHandler{
		handlers:              map[insolar.MessageType]insolar.MessageHandler{},
		conf:                  conf,
		LifelineIndex:         index,
		IndexBucketModifier:   indexBucketModifier,
		LifelineStateModifier: indexStateModifier,
		PendingModifier:       pendingModifier,
		PendingAccessor:       pendingAccessor,
	}

	dep := &proc.Dependencies{
		FetchJet: func(p *proc.FetchJet) {
			p.Dep.JetAccessor = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetUpdater = h.jetTreeUpdater
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		WaitHot: func(p *proc.WaitHot) {
			p.Dep.Waiter = h.HotDataWaiter
			p.Dep.Sender = h.Sender
		},
		GetIndex: func(p *proc.GetIndex) {
			p.Dep.IndexState = h.LifelineStateModifier
			p.Dep.Locker = h.IDLocker
			p.Dep.Index = h.LifelineIndex
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
			p.Dep.Sender = h.Sender
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
		GetIndexWM: func(p *proc.GetIndexWM) {
			p.Dep.IndexState = h.LifelineStateModifier
			p.Dep.Locker = h.IDLocker
			p.Dep.Index = h.LifelineIndex
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
			p.Dep.Sender = h.Sender
		},
		SetRequest: func(p *proc.SetRequest) {
			p.Dep(
				h.WriteAccessor,
				h.RecordModifier,
				h.RecentStorageProvider,
				h.PendingModifier,
				h.Sender,
			)
		},
		SetActivationRequest: func(p *proc.SetActivationRequest) {
			p.Dep(
				h.PCS,
				h.WriteAccessor,
				h.RecordModifier,
				h.Sender,
			)
		},
		SetResult: func(p *proc.SetResult) {
			p.Dep(
				h.WriteAccessor,
				h.RecordModifier,
				h.RecentStorageProvider,
				h.PendingModifier,
				h.Sender,
			)
		},
		SetBlob: func(p *proc.SetBlob) {
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.Sender = h.Sender
			p.Dep.PCS = h.PCS
			p.Dep.WriteAccessor = h.WriteAccessor
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Blobs = h.Blobs
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.Sender = h.Sender
			p.Dep.PendingAccessor = h.PendingAccessor
			p.Dep.PendingModifier = h.PendingModifier
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.Sender = h.Sender
		},
		UpdateObject: func(p *proc.UpdateObject) {
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.Bus = h.Bus
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.PCS = h.PCS
			p.Dep.IDLocker = h.IDLocker
			p.Dep.LifelineStateModifier = h.LifelineStateModifier
			p.Dep.LifelineIndex = h.LifelineIndex
			p.Dep.Sender = h.Sender
			p.Dep.WriteAccessor = h.WriteAccessor
			p.Dep.PendingModifier = h.PendingModifier
		},
		GetChildren: func(p *proc.GetChildren) {
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.DelegationTokenFactory = h.DelegationTokenFactory
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetTreeUpdater = h.jetTreeUpdater
			p.Dep.Sender = h.Sender
		},
		RegisterChild: func(p *proc.RegisterChild) {
			p.Dep.IDLocker = h.IDLocker
			p.Dep.LifelineIndex = h.LifelineIndex
			p.Dep.JetCoordinator = h.JetCoordinator
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.LifelineStateModifier = h.LifelineStateModifier
			p.Dep.PCS = h.PCS
			p.Dep.Sender = h.Sender
		},
		GetPendingRequests: func(p *proc.GetPendingRequests) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.Sender = h.Sender
		},
		GetPendingRequestID: func(p *proc.GetPendingRequestID) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.Sender = h.Sender
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Sender = h.Sender
		},
		HotData: func(p *proc.HotData) {
			p.Dep.DropModifier = h.DropModifier
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.MessageBus = h.Bus
			p.Dep.IndexBucketModifier = h.IndexBucketModifier
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.JetReleaser = h.JetReleaser
			p.Dep.Sender = h.Sender
		},
		GetPendingFilament: func(p *proc.GetPendingFilament) {
			p.Dep.PendingAccessor = h.PendingAccessor
			p.Dep.LifelineAccessor = h.LifelineIndex
		},
		RefreshPendingFilament: func(p *proc.RefreshPendingFilament) {
			p.Dep.LifelineAccessor = h.LifelineIndex
			p.Dep.PendingModifier = h.PendingModifier
			p.Dep.PendingAccessor = h.PendingAccessor
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
		},
		PassState: func(p *proc.PassState) {
			p.Dep.Blobs = h.BlobAccessor
			p.Dep.Sender = h.Sender
			p.Dep.Records = h.RecordAccessor
		},
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(h.PCS)
		},
		SetCode: func(p *proc.SetCode) {
			p.Dep(h.WriteAccessor, h.RecordModifier, h.BlobModifier, h.PCS, h.Sender)
		},
		GetDelegate: func(p *proc.GetDelegate) {
			p.Dep.Sender = h.Sender
		},
		Sender: h.Sender,
	}

	initHandle := func(msg *watermillMsg.Message) *handle.Init {
		return handle.NewInit(dep, h.Sender, msg)
	}

	h.FlowDispatcher = dispatcher.NewDispatcher(func(msg *watermillMsg.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg *watermillMsg.Message) flow.Handle {
		return initHandle(msg).Future
	})
	return h
}

// Init initializes handlers and middleware.
func (h *MessageHandler) Init(ctx context.Context) error {
	h.jetTreeUpdater = jet.NewFetcher(h.Nodes, h.JetStorage, h.Bus, h.JetCoordinator)
	h.setHandlersForLight()

	return nil
}

func (h *MessageHandler) OnPulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.ChangePulse(ctx, pn)
}

func (h *MessageHandler) setHandlersForLight() {
	// Generic.
	h.Bus.MustRegister(insolar.TypeValidateRecord, h.handleValidateRecord)
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	return &reply.OK{}, nil
}
