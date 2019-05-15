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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
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
	RecentStorageProvider      recentstorage.Provider             `inject:""`
	Bus                        insolar.MessageBus                 `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	CryptographyService        insolar.CryptographyService        `inject:""`
	DelegationTokenFactory     insolar.DelegationTokenFactory     `inject:""`
	JetStorage                 jet.Storage                        `inject:""`

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

	IndexStorage       object.IndexStorage
	IndexStateModifier object.ExtendedIndexModifier

	conf           *configuration.Ledger
	middleware     *middleware
	jetTreeUpdater jet.Fetcher

	FlowDispatcher *dispatcher.Dispatcher
	handlers       map[insolar.MessageType]insolar.MessageHandler
}

// NewMessageHandler creates new handler.
func NewMessageHandler(
	indexStorage object.IndexStorage,
	indexStateModifier object.ExtendedIndexModifier,
	conf *configuration.Ledger,
) *MessageHandler {

	h := &MessageHandler{
		handlers:           map[insolar.MessageType]insolar.MessageHandler{},
		conf:               conf,
		IndexStorage:       indexStorage,
		IndexStateModifier: indexStateModifier,
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
			p.Dep.IndexState = h.IndexStateModifier
			p.Dep.Locker = h.IDLocker
			p.Dep.Storage = h.IndexStorage
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.Bus = h.Bus
		},
		SetRecord: func(p *proc.SetRecord) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.PlatformCryptographyScheme = h.PlatformCryptographyScheme
			p.Dep.PendingRequestsLimit = h.conf.PendingRequestsLimit
		},
		SetBlob: func(p *proc.SetBlob) {
			p.Dep.BlobAccessor = h.BlobAccessor
			p.Dep.BlobModifier = h.BlobModifier
			p.Dep.PlatformCryptographyScheme = h.PlatformCryptographyScheme
		},
		SendObject: func(p *proc.SendObject) {
			p.Dep.Jets = h.JetStorage
			p.Dep.Blobs = h.Blobs
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.JetUpdater = h.jetTreeUpdater
			p.Dep.Bus = h.Bus
			p.Dep.RecordAccessor = h.RecordAccessor
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
			p.Dep.PlatformCryptographyScheme = h.PlatformCryptographyScheme
			p.Dep.IDLocker = h.IDLocker
			p.Dep.IndexStateModifier = h.IndexStateModifier
			p.Dep.IndexStorage = h.IndexStorage
		},
		RegisterChild: func(p *proc.RegisterChild) {
			p.Dep.IDLocker = h.IDLocker
			p.Dep.IndexStorage = h.IndexStorage
			p.Dep.JetCoordinator = h.JetCoordinator
			p.Dep.RecordModifier = h.RecordModifier
			p.Dep.IndexStateModifier = h.IndexStateModifier
			p.Dep.PlatformCryptographyScheme = h.PlatformCryptographyScheme
		},
		GetPendingRequests: func(p *proc.GetPendingRequests) {
			p.Dep.RecentStorageProvider = h.RecentStorageProvider
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep.Jets = h.JetStorage
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

	h.Bus.MustRegister(insolar.TypeGetChildren,
		BuildMiddleware(h.handleGetChildren,
			instrumentHandler("handleGetChildren"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(insolar.TypeSetRecord, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeRegisterChild, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeSetBlob, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetPendingRequests, h.FlowDispatcher.WrapBusHandle)
	h.Bus.MustRegister(insolar.TypeGetJet, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(insolar.TypeHotRecords,
		BuildMiddleware(h.handleHotRecords,
			instrumentHandler("handleHotRecords"),
			m.releaseHotDataWaiters))

	h.Bus.MustRegister(insolar.TypeGetRequest, h.FlowDispatcher.WrapBusHandle)

	h.Bus.MustRegister(
		insolar.TypeGetPendingRequestID,
		BuildMiddleware(
			h.handleGetPendingRequestID,
			instrumentHandler("handleGetPendingRequestID"),
			m.checkJet,
		),
	)

	h.Bus.MustRegister(insolar.TypeValidateRecord, h.handleValidateRecord)
}

func (h *MessageHandler) handleGetDelegate(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)
	jetID := jetFromContext(ctx)

	h.IndexStateModifier.SetUsageForPulse(ctx, *msg.Head.Record(), parcel.Pulse())

	h.IDLocker.Lock(msg.Head.Record())
	defer h.IDLocker.Unlock(msg.Head.Record())

	idx, err := h.IndexStorage.ForID(ctx, *msg.Head.Record())
	if err == object.ErrIndexNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Head, heavy)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch index from heavy")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	h.IndexStateModifier.SetUsageForPulse(ctx, *msg.Head.Record(), parcel.Pulse())

	delegateRef, ok := idx.Delegates[msg.AsType]
	if !ok {
		return nil, errors.New("the object has no delegate for this type")
	}

	rep := reply.Delegate{
		Head: delegateRef,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetChildren(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetChildren)
	jetID := jetFromContext(ctx)

	h.IDLocker.Lock(msg.Parent.Record())
	defer h.IDLocker.Unlock(msg.Parent.Record())

	idx, err := h.IndexStorage.ForID(ctx, *msg.Parent.Record())
	if err == object.ErrIndexNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Parent, heavy)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch index from heavy")
		}
		if idx.ChildPointer == nil {
			return &reply.Children{Refs: nil, NextFrom: nil}, nil
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	h.IndexStateModifier.SetUsageForPulse(ctx, *msg.Parent.Record(), parcel.Pulse())

	var (
		refs         []insolar.Reference
		currentChild *insolar.ID
	)

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		currentChild = msg.FromChild
	} else {
		currentChild = idx.ChildPointer
	}

	// The object has no children.
	if currentChild == nil {
		return &reply.Children{Refs: nil, NextFrom: nil}, nil
	}

	var childJet *insolar.ID
	onHeavy, err := h.JetCoordinator.IsBeyondLimit(ctx, parcel.Pulse(), currentChild.Pulse())
	if err != nil && err != pulse.ErrNotFound {
		return nil, err
	}
	if onHeavy {
		node, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetChildrenRedirect(h.DelegationTokenFactory, parcel, node, *currentChild)
	}

	childJetID, actual := h.JetStorage.ForID(ctx, currentChild.Pulse(), *msg.Parent.Record())
	childJet = (*insolar.ID)(&childJetID)

	if !actual {
		actualJet, err := h.jetTreeUpdater.Fetch(ctx, *msg.Parent.Record(), currentChild.Pulse())
		if err != nil {
			return nil, err
		}
		childJet = actualJet
	}

	// Try to fetch the first child.
	_, err = h.RecordAccessor.ForID(ctx, *currentChild)

	if err == object.ErrNotFound {
		node, err := h.JetCoordinator.NodeForJet(ctx, *childJet, parcel.Pulse(), currentChild.Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetChildrenRedirect(h.DelegationTokenFactory, parcel, node, *currentChild)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch child")
	}

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		counter++

		rec, err := h.RecordAccessor.ForID(ctx, *currentChild)

		// We don't have this child reference. Return what was collected.
		if err == object.ErrNotFound {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		if err != nil {
			return nil, errors.New("failed to retrieve children")
		}

		virtRec := rec.Record
		childRec, ok := virtRec.(*object.ChildRecord)
		if !ok {
			return nil, errors.New("failed to retrieve children")
		}
		currentChild = childRec.PrevChild

		// Skip records later than specified pulse.
		recPulse := childRec.Ref.Record().Pulse()
		if msg.FromPulse != nil && recPulse > *msg.FromPulse {
			continue
		}
		refs = append(refs, childRec.Ref)
	}

	return &reply.Children{Refs: refs, NextFrom: nil}, nil
}

func (h *MessageHandler) handleGetPendingRequestID(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	jetID := jetFromContext(ctx)
	msg := parcel.Message().(*message.GetPendingRequestID)

	requests := h.RecentStorageProvider.GetPendingStorage(ctx, jetID).GetRequestsForObject(msg.ObjectID)
	if len(requests) == 0 {
		return &reply.Error{ErrType: reply.ErrNoPendingRequests}, nil
	}

	rep := reply.ID{
		ID: requests[0],
	}

	return &rep, nil
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	return &reply.OK{}, nil
}

func (h *MessageHandler) saveIndexFromHeavy(
	ctx context.Context, jetID insolar.ID, obj insolar.Reference, heavy *insolar.Reference,
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
	err = h.IndexStorage.Set(ctx, *obj.Record(), idx)
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to save")
	}
	return idx, nil
}

func (h *MessageHandler) handleHotRecords(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	logger := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.HotData)
	jetID := insolar.JetID(*msg.Jet.Record())

	logger.WithFields(map[string]interface{}{
		"jet": jetID.DebugString(),
	}).Info("received hot data")

	err := h.DropModifier.Set(ctx, msg.Drop)
	if err == drop.ErrOverride {
		err = nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[jet]: drop error (pulse: %v)", msg.Drop.Pulse)
	}

	pendingStorage := h.RecentStorageProvider.GetPendingStorage(ctx, insolar.ID(jetID))
	logger.Debugf("received %d pending requests", len(msg.PendingRequests))

	var notificationList []insolar.ID
	for objID, objContext := range msg.PendingRequests {
		if !objContext.Active {
			notificationList = append(notificationList, objID)
		}

		objContext.Active = false
		pendingStorage.SetContextToObject(ctx, objID, objContext)
	}

	go func() {
		for _, objID := range notificationList {
			go func(objID insolar.ID) {
				rep, err := h.Bus.Send(ctx, &message.AbandonedRequestsNotification{
					Object: objID,
				}, nil)

				if err != nil {
					logger.Error("failed to notify about pending requests")
					return
				}
				if _, ok := rep.(*reply.OK); !ok {
					logger.Error("received unexpected reply on pending notification")
				}
			}(objID)
		}
	}()

	for id, meta := range msg.HotIndexes {
		decodedIndex, err := object.DecodeIndex(meta.Index)
		if err != nil {
			logger.Error(err)
			continue
		}

		err = h.IndexStateModifier.SetWithMeta(ctx, id, meta.LastUsed, decodedIndex)
		if err != nil {
			logger.Error(err)
			continue
		}
	}

	h.JetStorage.Update(
		ctx, msg.PulseNumber, true, insolar.JetID(jetID),
	)

	h.jetTreeUpdater.Release(ctx, jetID, msg.PulseNumber)

	return &reply.OK{}, nil
}
