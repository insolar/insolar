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
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
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

	LifelineIndex         object.LifelineIndex
	IndexBucketModifier   object.IndexBucketModifier
	LifelineStateModifier object.LifelineStateModifier

	conf           *configuration.Ledger
	middleware     *middleware
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
	conf *configuration.Ledger,
) *MessageHandler {

	h := &MessageHandler{
		handlers:              map[insolar.MessageType]insolar.MessageHandler{},
		conf:                  conf,
		LifelineIndex:         index,
		IndexBucketModifier:   indexBucketModifier,
		LifelineStateModifier: indexStateModifier,
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
		GetIndex: func(p *proc.GetIndex) {
			p.Dep.IndexState = h.LifelineStateModifier
			p.Dep.Locker = h.IDLocker
			p.Dep.Index = h.LifelineIndex
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
		},
		FetchJetWM: func(p *proc.FetchJetWM) {
			p.Dep.JetAccessor = h.JetStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetUpdater = h.jetTreeUpdater
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
		SetRecord: func(p *proc.SetRecord) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.PCS = h.PCS
			p.Dep.PendingRequestsLimit = h.conf.PendingRequestsLimit
		},
		SetBlob: func(p *proc.SetBlob) {
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.PlatformCryptographyScheme = h.PCS
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Blobs = h.Blobs
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetUpdater = h.jetTreeUpdater
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.Sender = h.Sender
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.BlobAccessor = h.BlobAccessor
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep.RecordAccessor = h.RecordAccessor
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
		},
		GetChildren: func(p *proc.GetChildren) {
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.DelegationTokenFactory = h.DelegationTokenFactory
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetTreeUpdater = h.jetTreeUpdater
		},
		RegisterChild: func(p *proc.RegisterChild) {
			p.Dep.IDLocker = h.IDLocker
			p.Dep.LifelineIndex = h.LifelineIndex
			p.Dep.JetCoordinator = h.JetCoordinator
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.LifelineStateModifier = h.LifelineStateModifier
			p.Dep.PCS = h.PCS
		},
		GetPendingRequests: func(p *proc.GetPendingRequests) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
		},
		GetPendingRequestID: func(p *proc.GetPendingRequestID) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep.Jets = h.JetStorage
		},
		HotData: func(p *proc.HotData) {
			p.Dep.DropModifier = h.DropModifier
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.MessageBus = h.Bus
			p.Dep.IndexBucketModifier = h.IndexBucketModifier
			p.Dep.JetStorage = h.JetStorage
			p.Dep.JetFetcher = h.jetTreeUpdater
			p.Dep.JetReleaser = h.JetReleaser
		},
	}

	initHandle := func(msg bus.Message) *handle.Init {
		return &handle.Init{
			Dep:     dep,
			Message: msg,
		}
	}

	h.FlowDispatcher = dispatcher.NewDispatcher(func(msg bus.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg bus.Message) flow.Handle {
		return initHandle(msg).Future
	})
	return h
}

func instrumentHandler(name string) Handler {
	return func(handler insolar.MessageHandler) insolar.MessageHandler {
		return func(ctx context.Context, p insolar.Parcel) (insolar.Reply, error) {
			inslog := inslogger.FromContext(ctx)
			start := time.Now()
			code := "2xx"
			ctx = insmetrics.InsertTag(ctx, tagMethod, name)

			repl, err := handler(ctx, p)

			latency := time.Since(start)
			if err != nil {
				code = "5xx"
				inslog.Errorf("AM's handler %v returns error: %v", name, err)
			}
			inslog.Debugf("measured time of AM method %v is %v", name, latency)

			ctx = insmetrics.ChangeTags(
				ctx,
				tag.Insert(tagMethod, name),
				tag.Insert(tagResult, code),
			)
			stats.Record(ctx, statCalls.M(1), statLatency.M(latency.Nanoseconds()/1e6))

			return repl, err
		}
	}
}

// Init initializes handlers and middleware.
func (h *MessageHandler) Init(ctx context.Context) error {
	m := newMiddleware(h)
	h.middleware = m

	h.jetTreeUpdater = jet.NewFetcher(h.Nodes, h.JetStorage, h.Bus, h.JetCoordinator)

	h.setHandlersForLight(m)

	return nil
}

func (h *MessageHandler) OnPulse(ctx context.Context, pn insolar.Pulse) {
	h.FlowDispatcher.ChangePulse(ctx, pn)
}

func (h *MessageHandler) setHandlersForLight(m *middleware) {
	// Generic.

	h.Bus.MustRegister(insolar.TypeGetCode, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetObject, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeUpdateObject, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(insolar.TypeGetDelegate,
		BuildMiddleware(h.handleGetDelegate,
			instrumentHandler("handleGetDelegate"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

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

func (h *MessageHandler) handleGetDelegate(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)
	jetID := jetFromContext(ctx)

	h.IDLocker.Lock(msg.Head.Record())
	defer h.IDLocker.Unlock(msg.Head.Record())

	idx, err := h.LifelineIndex.ForID(ctx, parcel.Pulse(), *msg.Head.Record())
	if err == object.ErrLifelineNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, parcel.Pulse(), jetID, msg.Head, heavy)
		if err != nil {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"jet": jetID.DebugString(),
				"pn":  parcel.Pulse(),
			}).Error(errors.Wrapf(err, "failed to fetch index from heavy - %v", *msg.Head.Record()))
			return nil, errors.Wrapf(err, "failed to fetch index from heavy")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	err = h.LifelineStateModifier.SetLifelineUsage(ctx, parcel.Pulse(), *msg.Head.Record())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	err = h.LifelineStateModifier.SetLifelineUsage(ctx, parcel.Pulse(), *msg.Head.Record())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}

	delegateRef, ok := idx.DelegateByKey(msg.AsType)
	if !ok {
		return nil, errors.New("the object has no delegate for this type")
	}

	rep := reply.Delegate{
		Head: delegateRef,
	}

	return &rep, nil
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	return &reply.OK{}, nil
}

func (h *MessageHandler) saveIndexFromHeavy(
	ctx context.Context, parcelPN insolar.PulseNumber, jetID insolar.ID, obj insolar.Reference, heavy *insolar.Reference,
) (object.Lifeline, error) {
	genericReply, err := h.Bus.Send(ctx, &message.GetObjectIndex{
		Object: obj,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to send")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return object.Lifeline{}, fmt.Errorf("failed to fetch object index: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}
	idx, err := object.DecodeIndex(rep.Index)
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to decode")
	}

	idx.JetID = insolar.JetID(jetID)
	err = h.LifelineIndex.Set(ctx, parcelPN, *obj.Record(), idx)
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to save")
	}
	return idx, nil
}
