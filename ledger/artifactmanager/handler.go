/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/core/delegationtoken"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/jet"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/hack"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/record"
)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	RecentStorageProvider      recentstorage.Provider          `inject:""`
	Bus                        core.MessageBus                 `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`
	CryptographyService        core.CryptographyService        `inject:""`
	DelegationTokenFactory     core.DelegationTokenFactory     `inject:""`
	HeavySync                  core.HeavySync                  `inject:""`
	PulseStorage               core.PulseStorage               `inject:""`
	JetStorage                 storage.JetStorage              `inject:""`
	DropStorage                storage.DropStorage             `inject:""`
	ObjectStorage              storage.ObjectStorage           `inject:""`
	Nodes                      nodes.Accessor                  `inject:""`
	PulseTracker               storage.PulseTracker            `inject:""`
	DBContext                  storage.DBContext               `inject:""`
	HotDataWaiter              HotDataWaiter                   `inject:""`

	certificate    core.Certificate
	replayHandlers map[core.MessageType]core.MessageHandler
	conf           *configuration.Ledger
	middleware     *middleware
	jetTreeUpdater *jetTreeUpdater
	isHeavy        bool
}

// NewMessageHandler creates new handler.
func NewMessageHandler(conf *configuration.Ledger, certificate core.Certificate) *MessageHandler {
	return &MessageHandler{
		certificate:    certificate,
		replayHandlers: map[core.MessageType]core.MessageHandler{},
		conf:           conf,
	}
}

func instrumentHandler(name string) Handler {
	return func(handler core.MessageHandler) core.MessageHandler {
		return func(ctx context.Context, p core.Parcel) (core.Reply, error) {
			// TODO: add tags to log
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

	h.jetTreeUpdater = newJetTreeUpdater(h.Nodes, h.JetStorage, h.Bus, h.JetCoordinator)

	h.isHeavy = h.certificate.GetRole() == core.StaticRoleHeavyMaterial

	// core.StaticRoleUnknown - genesis
	if h.certificate.GetRole() == core.StaticRoleLightMaterial || h.certificate.GetRole() == core.StaticRoleUnknown {
		h.setHandlersForLight(m)
		h.setReplayHandlers(m)
	}

	if h.isHeavy {
		h.setHandlersForHeavy(m)
	}

	return nil
}

func (h *MessageHandler) setHandlersForLight(m *middleware) {
	// Generic.
	h.Bus.MustRegister(core.TypeGetCode, BuildMiddleware(h.handleGetCode))

	h.Bus.MustRegister(core.TypeGetObject,
		BuildMiddleware(h.handleGetObject,
			instrumentHandler("handleGetObject"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeGetDelegate,
		BuildMiddleware(h.handleGetDelegate,
			instrumentHandler("handleGetDelegate"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeGetChildren,
		BuildMiddleware(h.handleGetChildren,
			instrumentHandler("handleGetChildren"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeSetRecord,
		BuildMiddleware(h.handleSetRecord,
			instrumentHandler("handleSetRecord"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeUpdateObject,
		BuildMiddleware(h.handleUpdateObject,
			instrumentHandler("handleUpdateObject"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeRegisterChild,
		BuildMiddleware(h.handleRegisterChild,
			instrumentHandler("handleRegisterChild"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeSetBlob,
		BuildMiddleware(h.handleSetBlob,
			instrumentHandler("handleSetBlob"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeGetObjectIndex,
		BuildMiddleware(h.handleGetObjectIndex,
			instrumentHandler("handleGetObjectIndex"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeGetPendingRequests,
		BuildMiddleware(h.handleHasPendingRequests,
			instrumentHandler("handleHasPendingRequests"),
			m.addFieldsToLogger,
			m.checkJet,
			m.waitForHotData))

	h.Bus.MustRegister(core.TypeGetJet,
		BuildMiddleware(h.handleGetJet,
			instrumentHandler("handleGetJet")))

	h.Bus.MustRegister(core.TypeHotRecords,
		BuildMiddleware(h.handleHotRecords,
			instrumentHandler("handleHotRecords"),
			m.releaseHotDataWaiters))

	h.Bus.MustRegister(
		core.TypeGetRequest,
		BuildMiddleware(
			h.handleGetRequest,
			instrumentHandler("handleGetRequest"),
			m.checkJet,
		),
	)

	h.Bus.MustRegister(
		core.TypeGetPendingRequestID,
		BuildMiddleware(
			h.handleGetPendingRequestID,
			instrumentHandler("handleGetPendingRequestID"),
			m.checkJet,
		),
	)

	// Validation.
	h.Bus.MustRegister(core.TypeValidateRecord,
		BuildMiddleware(h.handleValidateRecord,
			m.addFieldsToLogger,
			m.checkJet))

	h.Bus.MustRegister(core.TypeValidationCheck,
		BuildMiddleware(h.handleValidationCheck,
			m.addFieldsToLogger,
			m.checkJet))

	h.Bus.MustRegister(core.TypeJetDrop,
		BuildMiddleware(h.handleJetDrop,
			m.addFieldsToLogger,
			m.checkJet))
}
func (h *MessageHandler) setReplayHandlers(m *middleware) {
	// Generic.
	h.replayHandlers[core.TypeGetCode] = BuildMiddleware(h.handleGetCode, m.addFieldsToLogger)
	h.replayHandlers[core.TypeGetObject] = BuildMiddleware(h.handleGetObject, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeGetDelegate] = BuildMiddleware(h.handleGetDelegate, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeGetChildren] = BuildMiddleware(h.handleGetChildren, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeSetRecord] = BuildMiddleware(h.handleSetRecord, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeUpdateObject] = BuildMiddleware(h.handleUpdateObject, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeRegisterChild] = BuildMiddleware(h.handleRegisterChild, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeSetBlob] = BuildMiddleware(h.handleSetBlob, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeGetObjectIndex] = BuildMiddleware(h.handleGetObjectIndex, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeGetPendingRequests] = BuildMiddleware(h.handleHasPendingRequests, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeGetJet] = BuildMiddleware(h.handleGetJet)

	// Validation.
	h.replayHandlers[core.TypeValidateRecord] = BuildMiddleware(h.handleValidateRecord, m.addFieldsToLogger, m.checkJet)
	h.replayHandlers[core.TypeValidationCheck] = BuildMiddleware(h.handleValidationCheck, m.addFieldsToLogger, m.checkJet)
}
func (h *MessageHandler) setHandlersForHeavy(m *middleware) {
	// Heavy.
	h.Bus.MustRegister(core.TypeHeavyStartStop,
		BuildMiddleware(h.handleHeavyStartStop,
			instrumentHandler("handleHeavyStartStop")))

	h.Bus.MustRegister(core.TypeHeavyReset,
		BuildMiddleware(h.handleHeavyReset,
			instrumentHandler("handleHeavyReset")))

	h.Bus.MustRegister(core.TypeHeavyPayload,
		BuildMiddleware(h.handleHeavyPayload,
			instrumentHandler("handleHeavyPayload")))

	// Generic.
	h.Bus.MustRegister(core.TypeGetCode,
		BuildMiddleware(h.handleGetCode))

	h.Bus.MustRegister(core.TypeGetObject,
		BuildMiddleware(h.handleGetObject,
			instrumentHandler("handleGetObject"),
			m.zeroJetForHeavy))

	h.Bus.MustRegister(core.TypeGetDelegate,
		BuildMiddleware(h.handleGetDelegate,
			instrumentHandler("handleGetDelegate"),
			m.zeroJetForHeavy))

	h.Bus.MustRegister(core.TypeGetChildren,
		BuildMiddleware(h.handleGetChildren,
			instrumentHandler("handleGetChildren"),
			m.zeroJetForHeavy))

	h.Bus.MustRegister(core.TypeGetObjectIndex,
		BuildMiddleware(h.handleGetObjectIndex,
			instrumentHandler("handleGetObjectIndex"),
			m.zeroJetForHeavy))

	h.Bus.MustRegister(
		core.TypeGetRequest,
		BuildMiddleware(h.handleGetRequest,
			instrumentHandler("handleGetRequest"),
			m.zeroJetForHeavy))
}

func (h *MessageHandler) handleSetRecord(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	msg := parcel.Message().(*message.SetRecord)
	rec := record.DeserializeRecord(msg.Record)
	jetID := jetFromContext(ctx)

	calculatedID := record.NewRecordIDFromRecord(h.PlatformCryptographyScheme, parcel.Pulse(), rec)

	switch r := rec.(type) {
	case record.Request:
		if h.RecentStorageProvider.Count() > h.conf.PendingRequestsLimit {
			return &reply.Error{ErrType: reply.ErrTooManyPendingRequests}, nil
		}
		recentStorage := h.RecentStorageProvider.GetPendingStorage(ctx, jetID)
		recentStorage.AddPendingRequest(ctx, r.GetObject(), *calculatedID)
	case *record.ResultRecord:
		recentStorage := h.RecentStorageProvider.GetPendingStorage(ctx, jetID)
		recentStorage.RemovePendingRequest(ctx, r.Object, *r.Request.Record())
	}

	id, err := h.ObjectStorage.SetRecord(ctx, jetID, parcel.Pulse(), rec)
	if err == storage.ErrOverride {
		inslogger.FromContext(ctx).WithField("type", fmt.Sprintf("%T", rec)).Warn("set record override")
		id = calculatedID
	} else if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleSetBlob(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	msg := parcel.Message().(*message.SetBlob)
	jetID := jetFromContext(ctx)
	calculatedID := record.CalculateIDForBlob(h.PlatformCryptographyScheme, parcel.Pulse(), msg.Memory)

	_, err := h.ObjectStorage.GetBlob(ctx, jetID, calculatedID)
	if err == nil {
		return &reply.ID{ID: *calculatedID}, nil
	}
	if err != nil && err != core.ErrNotFound {
		return nil, err
	}

	id, err := h.ObjectStorage.SetBlob(ctx, jetID, parcel.Pulse(), msg.Memory)
	if err == nil {
		return &reply.ID{ID: *id}, nil
	}
	if err == storage.ErrOverride {
		return &reply.ID{ID: *calculatedID}, nil
	}
	return nil, err
}

func (h *MessageHandler) handleGetCode(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetCode)
	jetID := *jet.NewID(0, nil)

	codeRec, err := h.getCode(ctx, msg.Code.Record())
	if err == core.ErrNotFound {
		// We don't have code record. Must be on another node.
		node, err := h.JetCoordinator.NodeForJet(ctx, jetID, parcel.Pulse(), msg.Code.Record().Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetCodeRedirect(h.DelegationTokenFactory, parcel, node)
	}
	if err != nil {
		return nil, err
	}
	code, err := h.ObjectStorage.GetBlob(ctx, jetID, codeRec.Code)
	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        code,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetObject(
	ctx context.Context, parcel core.Parcel,
) (core.Reply, error) {
	msg := parcel.Message().(*message.GetObject)
	jetID := jetFromContext(ctx)
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object": msg.Head.Record().DebugString(),
		"pulse":  parcel.Pulse(),
	})

	if !h.isHeavy {
		h.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Head.Record())
	}

	// Fetch object index. If not found redirect.
	idx, err := h.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
	if err == core.ErrNotFound {
		if h.isHeavy {
			return nil, fmt.Errorf("failed to fetch index for %s", msg.Head.Record().String())
		}

		logger.Debug("failed to fetch index (fetching from heavy)")
		node, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Head, node)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch index from heavy")
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index %s", msg.Head.Record().String())
	}

	// Determine object state id.
	var stateID *core.RecordID
	if msg.State != nil {
		stateID = msg.State
	} else {
		if msg.Approved {
			stateID = idx.LatestStateApproved
		} else {
			stateID = idx.LatestState
		}
	}
	if stateID == nil {
		return &reply.Error{ErrType: reply.ErrStateNotAvailable}, nil
	}

	var (
		stateJet *core.RecordID
	)
	if h.isHeavy {
		stateJet = &jetID
	} else {
		var actual bool
		onHeavy, err := h.JetCoordinator.IsBeyondLimit(ctx, parcel.Pulse(), stateID.Pulse())
		if err != nil {
			return nil, err
		}
		if onHeavy {
			node, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return nil, err
			}
			logger.WithFields(map[string]interface{}{
				"state":    stateID.DebugString(),
				"going_to": node.String(),
			}).Debug("fetching object (on heavy)")

			obj, err := h.fetchObject(ctx, msg.Head, *node, stateID, parcel.Pulse())
			if err != nil {
				if err == core.ErrDeactivated {
					return &reply.Error{ErrType: reply.ErrDeactivated}, nil
				}
				return nil, err
			}

			return &reply.Object{
				Head:         msg.Head,
				State:        *stateID,
				Prototype:    obj.Prototype,
				IsPrototype:  obj.IsPrototype,
				ChildPointer: idx.ChildPointer,
				Parent:       idx.Parent,
				Memory:       obj.Memory,
			}, nil
		}

		stateJet, actual = h.JetStorage.FindJet(ctx, stateID.Pulse(), *msg.Head.Record())
		if !actual {
			actualJet, err := h.jetTreeUpdater.fetchJet(ctx, *msg.Head.Record(), stateID.Pulse())
			if err != nil {
				return nil, err
			}
			stateJet = actualJet
		}
	}

	// Fetch state record.
	rec, err := h.ObjectStorage.GetRecord(ctx, *stateJet, stateID)
	if err == core.ErrNotFound {
		if h.isHeavy {
			return nil, fmt.Errorf("failed to fetch state for %v. jet: %v, state: %v", msg.Head.Record(), stateJet.DebugString(), stateID.DebugString())
		}
		// The record wasn't found on the current node. Return redirect to the node that contains it.
		// We get Jet tree for pulse when given state was added.
		node, err := h.JetCoordinator.NodeForJet(ctx, *stateJet, parcel.Pulse(), stateID.Pulse())
		if err != nil {
			return nil, err
		}
		logger.WithFields(map[string]interface{}{
			"state":    stateID.DebugString(),
			"going_to": node.String(),
		}).Debug("fetching object (record not found)")

		obj, err := h.fetchObject(ctx, msg.Head, *node, stateID, parcel.Pulse())
		if err != nil {
			if err == core.ErrDeactivated {
				return &reply.Error{ErrType: reply.ErrDeactivated}, nil
			}
			return nil, err
		}

		return &reply.Object{
			Head:         msg.Head,
			State:        *stateID,
			Prototype:    obj.Prototype,
			IsPrototype:  obj.IsPrototype,
			ChildPointer: idx.ChildPointer,
			Parent:       idx.Parent,
			Memory:       obj.Memory,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("invalid object record")
	}
	if state.State() == record.StateDeactivation {
		return &reply.Error{ErrType: reply.ErrDeactivated}, nil
	}

	var childPointer *core.RecordID
	if idx.ChildPointer != nil {
		childPointer = idx.ChildPointer
	}
	rep := reply.Object{
		Head:         msg.Head,
		State:        *stateID,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: childPointer,
		Parent:       idx.Parent,
	}

	if state.GetMemory() != nil {
		rep.Memory, err = h.ObjectStorage.GetBlob(ctx, *stateJet, state.GetMemory())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch blob")
		}
	}

	return &rep, nil
}

func (h *MessageHandler) handleHasPendingRequests(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetPendingRequests)
	jetID := jetFromContext(ctx)

	for _, reqID := range h.RecentStorageProvider.GetPendingStorage(ctx, jetID).GetRequestsForObject(*msg.Object.Record()) {
		if reqID.Pulse() < parcel.Pulse() {
			return &reply.HasPendingRequests{Has: true}, nil
		}
	}

	return &reply.HasPendingRequests{Has: false}, nil
}

func (h *MessageHandler) handleGetJet(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetJet)

	jetID, actual := h.JetStorage.FindJet(ctx, msg.Pulse, msg.Object)

	return &reply.Jet{ID: *jetID, Actual: actual}, nil
}

func (h *MessageHandler) handleGetDelegate(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)
	jetID := jetFromContext(ctx)

	if !h.isHeavy {
		h.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Head.Record())
	}

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
	if err == core.ErrNotFound {
		if h.isHeavy {
			return nil, fmt.Errorf("failed to fetch index for %v", msg.Head.Record())
		}

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
	ctx context.Context, parcel core.Parcel,
) (core.Reply, error) {
	msg := parcel.Message().(*message.GetChildren)
	jetID := jetFromContext(ctx)

	if !h.isHeavy {
		h.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Parent.Record())
	}

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
	if err == core.ErrNotFound {
		if h.isHeavy {
			return nil, fmt.Errorf("failed to fetch index for %v", msg.Parent.Record())
		}

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

	var (
		refs         []core.RecordRef
		currentChild *core.RecordID
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

	var childJet *core.RecordID
	if h.isHeavy {
		childJet = &jetID
	} else {
		var actual bool
		onHeavy, err := h.JetCoordinator.IsBeyondLimit(ctx, parcel.Pulse(), currentChild.Pulse())
		if err != nil {
			return nil, err
		}
		if onHeavy {
			node, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return nil, err
			}
			return reply.NewGetChildrenRedirect(h.DelegationTokenFactory, parcel, node, *currentChild)
		}

		childJet, actual = h.JetStorage.FindJet(ctx, currentChild.Pulse(), *msg.Parent.Record())
		if !actual {
			actualJet, err := h.jetTreeUpdater.fetchJet(ctx, *msg.Parent.Record(), currentChild.Pulse())
			if err != nil {
				return nil, err
			}
			childJet = actualJet
		}
	}

	// Try to fetch the first child.
	_, err = h.ObjectStorage.GetRecord(ctx, *childJet, currentChild)
	if err == core.ErrNotFound {
		if h.isHeavy {
			return nil, fmt.Errorf("failed to fetch child for %v. jet: %v, state: %v", msg.Parent.Record(), childJet.DebugString(), currentChild.DebugString())
		}
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

		rec, err := h.ObjectStorage.GetRecord(ctx, *childJet, currentChild)
		// We don't have this child reference. Return what was collected.
		if err == core.ErrNotFound {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		if err != nil {
			return nil, errors.New("failed to retrieve children")
		}

		childRec, ok := rec.(*record.ChildRecord)
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

func (h *MessageHandler) handleGetRequest(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	jetID := jetFromContext(ctx)
	msg := parcel.Message().(*message.GetRequest)

	rec, err := h.ObjectStorage.GetRecord(ctx, jetID, &msg.Request)
	if err != nil {
		return nil, errors.New("failed to fetch request")
	}

	req, ok := rec.(*record.RequestRecord)
	if !ok {
		return nil, errors.New("failed to decode request")
	}

	rep := reply.Request{
		ID:     msg.Request,
		Record: record.SerializeRecord(req),
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetPendingRequestID(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
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

func (h *MessageHandler) handleUpdateObject(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	msg := parcel.Message().(*message.UpdateObject)
	jetID := jetFromContext(ctx)
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object": msg.Object.Record().DebugString(),
		"pulse":  parcel.Pulse(),
	})

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("wrong object state record")
	}

	h.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Object.Record())

	// FIXME: temporary fix. If we calculate blob id on the client, pulse can change before message sending and this
	//  id will not match the one calculated on the server.
	blobID, err := h.ObjectStorage.SetBlob(ctx, jetID, parcel.Pulse(), msg.Memory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set blob")
	}

	switch s := state.(type) {
	case *record.ObjectActivateRecord:
		s.Memory = blobID
	case *record.ObjectAmendRecord:
		s.Memory = blobID
	}

	var idx *index.ObjectLifeline
	err = h.DBContext.Update(ctx, func(tx *storage.TransactionManager) error {
		var err error
		idx, err = tx.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
		// No index on our node.
		if err == core.ErrNotFound {
			if state.State() == record.StateActivation {
				// We are activating the object. There is no index for it anywhere.
				idx = &index.ObjectLifeline{State: record.StateUndefined}
			} else {
				logger.Debug("failed to fetch index (fetching from heavy)")
				// We are updating object. Index should be on the heavy executor.
				heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
				if err != nil {
					return err
				}
				idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Object, heavy)
				if err != nil {
					return errors.Wrap(err, "failed to fetch index from heavy")
				}
			}
		} else if err != nil {
			return err
		}

		if err = validateState(idx.State, state.State()); err != nil {
			return err
		}

		recID := record.NewRecordIDFromRecord(h.PlatformCryptographyScheme, parcel.Pulse(), rec)

		// Index exists and latest record id does not match (preserving chain consistency).
		// For the case when vm can't save or send result to another vm and it tries to update the same record again
		if idx.LatestState != nil && !state.PrevStateID().Equal(idx.LatestState) && idx.LatestState != recID {
			return errors.New("invalid state record")
		}

		id, err := tx.SetRecord(ctx, jetID, parcel.Pulse(), rec)
		if err == storage.ErrOverride {
			logger.WithField("type", fmt.Sprintf("%T", rec)).Warn("set record override (#1)")
			id = recID
		} else if err != nil {
			return err
		}
		idx.LatestState = id
		idx.State = state.State()
		if state.State() == record.StateActivation {
			idx.Parent = state.(*record.ObjectActivateRecord).Parent
		}

		idx.LatestUpdate = parcel.Pulse()
		return tx.SetObjectIndex(ctx, jetID, msg.Object.Record(), idx)
	})
	if err != nil {
		if err == ErrObjectDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}

	logger.WithField("state", idx.LatestState.DebugString()).Debug("saved object")

	rep := reply.Object{
		Head:         msg.Object,
		State:        *idx.LatestState,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: idx.ChildPointer,
		Parent:       idx.Parent,
	}
	return &rep, nil
}

func (h *MessageHandler) handleRegisterChild(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	logger := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.RegisterChild)
	jetID := jetFromContext(ctx)
	rec := record.DeserializeRecord(msg.Record)
	childRec, ok := rec.(*record.ChildRecord)
	if !ok {
		return nil, errors.New("wrong child record")
	}

	h.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Parent.Record())

	var child *core.RecordID
	err := h.DBContext.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, err := h.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
		if err == core.ErrNotFound {
			heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return err
			}
			idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Parent, heavy)
			if err != nil {
				return errors.Wrap(err, "failed to fetch index from heavy")
			}
		} else if err != nil {
			return err
		}

		recID := record.NewRecordIDFromRecord(h.PlatformCryptographyScheme, parcel.Pulse(), childRec)

		// Children exist and pointer does not match (preserving chain consistency).
		// For the case when vm can't save or send result to another vm and it tries to update the same record again
		if idx.ChildPointer != nil && !childRec.PrevChild.Equal(idx.ChildPointer) && idx.ChildPointer != recID {
			return errors.New("invalid child record")
		}

		child, err = tx.SetRecord(ctx, jetID, parcel.Pulse(), childRec)
		if err == storage.ErrOverride {
			logger.WithField("type", fmt.Sprintf("%T", rec)).Warn("set record override (#2)")
			child = recID
		} else if err != nil {
			return err
		}

		idx.ChildPointer = child
		if msg.AsType != nil {
			idx.Delegates[*msg.AsType] = msg.Child
		}
		idx.LatestUpdate = parcel.Pulse()
		err = tx.SetObjectIndex(ctx, jetID, msg.Parent.Record(), idx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *child}, nil
}

func (h *MessageHandler) handleJetDrop(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.JetDrop)

	if !hack.SkipValidation(ctx) {
		for _, parcelBuff := range msg.Messages {
			jetDropMsg, err := message.Deserialize(bytes.NewBuffer(parcelBuff))
			if err != nil {
				return nil, err
			}
			handler, ok := h.replayHandlers[jetDropMsg.Type()]
			if !ok {
				return nil, errors.New("unknown message type")
			}

			_, err = handler(ctx, &message.Parcel{Msg: jetDropMsg})
			if err != nil {
				return nil, err
			}
		}
	}

	err := h.JetStorage.AddJets(ctx, msg.JetID)
	if err != nil {
		return nil, err
	}

	h.JetStorage.UpdateJetTree(
		ctx, parcel.Pulse(), true, msg.JetID,
	)

	return &reply.OK{}, nil
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	msg := parcel.Message().(*message.ValidateRecord)
	jetID := jetFromContext(ctx)

	err := h.DBContext.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, err := tx.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
		if err == core.ErrNotFound {
			heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return err
			}
			idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Object, heavy)
			if err != nil {
				return errors.Wrap(err, "failed to fetch index from heavy")
			}
		} else if err != nil {
			return err
		}

		// Find node that has this state.
		node, err := h.JetCoordinator.NodeForJet(ctx, jetID, parcel.Pulse(), msg.Object.Record().Pulse())
		if err != nil {
			return err
		}

		// Send checking message.
		genericReply, err := h.Bus.Send(ctx, &message.ValidationCheck{
			Object:              msg.Object,
			ValidatedState:      msg.State,
			LatestStateApproved: idx.LatestStateApproved,
		}, &core.MessageSendOptions{
			Receiver: node,
		})
		if err != nil {
			return err
		}
		switch genericReply.(type) {
		case *reply.OK:
			if msg.IsValid {
				idx.LatestStateApproved = &msg.State
			} else {
				idx.LatestState = idx.LatestStateApproved
			}
			idx.LatestUpdate = parcel.Pulse()
			err = tx.SetObjectIndex(ctx, jetID, msg.Object.Record(), idx)
			if err != nil {
				return errors.Wrap(err, "failed to save object index")
			}
		case *reply.NotOK:
			return errors.New("validation sequence integrity failure")
		default:
			return errors.New("handleValidateRecord: unexpected reply")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) handleGetObjectIndex(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetObjectIndex)
	jetID := jetFromContext(ctx)

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}

	buf, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize index")
	}

	return &reply.ObjectIndex{Index: buf}, nil
}

func (h *MessageHandler) handleValidationCheck(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.ValidationCheck)
	jetID := jetFromContext(ctx)

	rec, err := h.ObjectStorage.GetRecord(ctx, jetID, &msg.ValidatedState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch state record")
	}
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("failed to fetch state record")
	}
	approved := msg.LatestStateApproved
	validated := state.PrevStateID()
	if !approved.Equal(validated) && approved != nil && validated != nil {
		return &reply.NotOK{}, nil
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) getCode(ctx context.Context, id *core.RecordID) (*record.CodeRecord, error) {
	jetID := *jet.NewID(0, nil)

	rec, err := h.ObjectStorage.GetRecord(ctx, jetID, id)
	if err != nil {
		return nil, err
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}

	return codeRec, nil
}

func validateState(old record.State, new record.State) error {
	if old == record.StateDeactivation {
		return ErrObjectDeactivated
	}
	if old == record.StateUndefined && new != record.StateActivation {
		return errors.New("object is not activated")
	}
	if old != record.StateUndefined && new == record.StateActivation {
		return errors.New("object is already activated")
	}
	return nil
}

func (h *MessageHandler) saveIndexFromHeavy(
	ctx context.Context, jetID core.RecordID, obj core.RecordRef, heavy *core.RecordRef,
) (*index.ObjectLifeline, error) {
	genericReply, err := h.Bus.Send(ctx, &message.GetObjectIndex{
		Object: obj,
	}, &core.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return nil, fmt.Errorf("failed to fetch object index: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}
	idx, err := index.DecodeObjectLifeline(rep.Index)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	err = h.ObjectStorage.SetObjectIndex(ctx, jetID, obj.Record(), idx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save")
	}
	return idx, nil
}

func (h *MessageHandler) fetchObject(
	ctx context.Context, obj core.RecordRef, node core.RecordRef, stateID *core.RecordID, pulse core.PulseNumber,
) (*reply.Object, error) {
	sender := BuildSender(
		h.Bus.Send,
		followRedirectSender(h.Bus),
		retryJetSender(pulse, h.JetStorage),
	)
	genericReply, err := sender(
		ctx,
		&message.GetObject{
			Head:     obj,
			Approved: false,
			State:    stateID,
		},
		&core.MessageSendOptions{
			Receiver: &node,
			Token:    &delegationtoken.GetObjectRedirectToken{},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object state")
	}
	if rep, ok := genericReply.(*reply.Error); ok {
		return nil, rep.Error()
	}

	rep, ok := genericReply.(*reply.Object)
	if !ok {
		return nil, fmt.Errorf("failed to fetch object state: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}
	return rep, nil
}

func (h *MessageHandler) handleHotRecords(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	if h.isHeavy {
		return nil, errors.New("heavy updates are forbidden")
	}

	logger := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.HotData)
	// FIXME: check split signatures.
	jetID := *msg.Jet.Record()

	logger.WithFields(map[string]interface{}{
		"jet": jetID.DebugString(),
	}).Info("received hot data")

	err := h.DropStorage.SetDrop(ctx, msg.DropJet, &msg.Drop)
	if err == storage.ErrOverride {
		err = nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[jet]: drop error (pulse: %v)", msg.Drop.Pulse)
	}

	err = h.DropStorage.SetDropSizeHistory(ctx, msg.DropJet, msg.JetDropSizeHistory)
	if err != nil {
		return nil, errors.Wrap(err, "[ handleHotRecords ] Can't SetDropSizeHistory")
	}

	pendingStorage := h.RecentStorageProvider.GetPendingStorage(ctx, jetID)
	logger.Debugf("received %d pending requests", len(msg.PendingRequests))

	var notificationList []core.RecordID
	for objID, objContext := range msg.PendingRequests {
		if !objContext.Active {
			notificationList = append(notificationList, objID)
		}

		objContext.Active = false
		pendingStorage.SetContextToObject(ctx, objID, objContext)
	}

	go func() {
		for _, objID := range notificationList {
			go func(objID core.RecordID) {
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

	indexStorage := h.RecentStorageProvider.GetIndexStorage(ctx, jetID)
	for id, meta := range msg.RecentObjects {
		decodedIndex, err := index.DecodeObjectLifeline(meta.Index)
		if err != nil {
			logger.Error(err)
			continue
		}

		err = h.ObjectStorage.SetObjectIndex(ctx, jetID, &id, decodedIndex)
		if err != nil {
			logger.Error(err)
			continue
		}

		indexStorage.AddObjectWithTLL(ctx, id, meta.TTL)
	}

	h.JetStorage.UpdateJetTree(
		ctx, msg.PulseNumber, true, jetID,
	)

	h.jetTreeUpdater.releaseJet(ctx, jetID, msg.PulseNumber)

	err = h.JetStorage.AddJets(ctx, jetID)
	if err != nil {
		logger.Error(errors.Wrap(err, "couldn't add jet"))
		return nil, err
	}

	return &reply.OK{}, nil
}
