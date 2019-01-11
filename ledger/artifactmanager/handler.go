/*
 *    Copyright 2018 Insolar
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

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/heavy"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/pkg/errors"

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
	HeavyJetTreeSync           heavy.JetTreeSync               `inject:""`

	db             *storage.DB
	replayHandlers map[core.MessageType]core.MessageHandler
	conf           *configuration.Ledger
}

// NewMessageHandler creates new handler.
func NewMessageHandler(
	db *storage.DB, conf *configuration.Ledger,
) *MessageHandler {
	return &MessageHandler{
		db:             db,
		replayHandlers: map[core.MessageType]core.MessageHandler{},
		conf:           conf,
	}
}

// Init initializes handlers and middleware.
func (h *MessageHandler) Init(ctx context.Context) error {
	m := newMiddleware(h.conf, h.db, h.JetCoordinator, h.Bus)

	// Generic.
	h.replayHandlers[core.TypeGetCode] = m.checkJet(h.handleGetCode)
	h.replayHandlers[core.TypeGetObject] = m.checkJet(h.handleGetObject)
	h.replayHandlers[core.TypeGetDelegate] = m.checkJet(h.handleGetDelegate)
	h.replayHandlers[core.TypeGetChildren] = m.checkJet(h.handleGetChildren)
	h.replayHandlers[core.TypeSetRecord] = m.checkJet(h.handleSetRecord)
	h.replayHandlers[core.TypeUpdateObject] = m.checkJet(h.handleUpdateObject)
	h.replayHandlers[core.TypeRegisterChild] = m.checkJet(h.handleRegisterChild)
	h.replayHandlers[core.TypeSetBlob] = m.checkJet(h.handleSetBlob)
	h.replayHandlers[core.TypeGetObjectIndex] = m.checkJet(h.handleGetObjectIndex)
	h.replayHandlers[core.TypeGetPendingRequests] = m.checkJet(h.handleHasPendingRequests)
	h.replayHandlers[core.TypeGetJet] = h.handleGetJet

	// Validation.
	h.replayHandlers[core.TypeValidateRecord] = m.checkJet(h.handleValidateRecord)
	h.replayHandlers[core.TypeValidationCheck] = m.checkJet(h.handleValidationCheck)
	h.replayHandlers[core.TypeHotRecords] = h.handleHotRecords

	// Generic.
	h.Bus.MustRegister(core.TypeGetCode, m.checkJet(m.waitForDrop(m.saveParcel(h.handleGetCode))))
	h.Bus.MustRegister(core.TypeGetObject, m.checkJet(m.waitForDrop(m.saveParcel(h.handleGetObject))))
	h.Bus.MustRegister(core.TypeGetDelegate, m.checkJet(m.waitForDrop(m.saveParcel(h.handleGetDelegate))))
	h.Bus.MustRegister(core.TypeGetChildren, m.checkJet(m.waitForDrop(m.saveParcel(h.handleGetChildren))))
	h.Bus.MustRegister(core.TypeSetRecord, m.checkJet(m.checkHeavySync(m.waitForDrop(m.saveParcel(h.handleSetRecord)))))
	h.Bus.MustRegister(core.TypeUpdateObject, m.checkJet(m.checkHeavySync(m.waitForDrop(m.saveParcel(h.handleUpdateObject)))))
	h.Bus.MustRegister(core.TypeRegisterChild, m.checkJet(m.checkHeavySync(m.waitForDrop(m.saveParcel(h.handleRegisterChild)))))
	h.Bus.MustRegister(core.TypeSetBlob, m.checkJet(m.checkHeavySync(m.waitForDrop(m.saveParcel(h.handleSetBlob)))))
	h.Bus.MustRegister(core.TypeGetObjectIndex, m.checkJet(m.waitForDrop(m.saveParcel(h.handleGetObjectIndex))))
	h.Bus.MustRegister(core.TypeGetPendingRequests, m.checkJet(m.waitForDrop(m.saveParcel(h.handleHasPendingRequests))))
	h.Bus.MustRegister(core.TypeGetJet, h.handleGetJet)

	// Validation.
	h.Bus.MustRegister(core.TypeValidateRecord, m.checkJet(m.waitForDrop(m.saveParcel(h.handleValidateRecord))))
	h.Bus.MustRegister(core.TypeValidationCheck, m.checkJet(m.waitForDrop(m.saveParcel(h.handleValidationCheck))))
	h.Bus.MustRegister(core.TypeHotRecords, m.checkJet(m.unlockDropWaiters(m.saveParcel(h.handleHotRecords))))
	h.Bus.MustRegister(core.TypeJetDrop, m.checkJet(h.handleJetDrop))

	// Heavy.
	h.Bus.MustRegister(core.TypeHeavyStartStop, h.handleHeavyStartStop)
	h.Bus.MustRegister(core.TypeHeavyReset, h.handleHeavyReset)
	h.Bus.MustRegister(core.TypeHeavyPayload, h.handleHeavyPayload)
	h.Bus.MustRegister(core.TypeHeavyJetTree, h.handleHeavyJetTree)

	return nil
}

func (h *MessageHandler) handleSetRecord(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.SetRecord)
	rec := record.DeserializeRecord(msg.Record)
	jetID := jetFromContext(ctx)

	id, err := h.db.SetRecord(ctx, jetID, parcel.Pulse(), rec)
	if err != nil {
		return nil, err
	}

	recentStorage := h.RecentStorageProvider.GetStorage(jetID)
	if request, ok := rec.(record.Request); ok {
		recentStorage.AddPendingRequest(request.GetObject(), *id)
	}
	if result, ok := rec.(*record.ResultRecord); ok {
		recentStorage.RemovePendingRequest(result.Object, *result.Request.Record())
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleSetBlob(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.SetBlob)
	jetID := jetFromContext(ctx)
	calculatedID := record.CalculateIDForBlob(h.PlatformCryptographyScheme, parcel.Pulse(), msg.Memory)

	_, err := h.db.GetBlob(ctx, jetID, calculatedID)
	if err == nil {
		return &reply.ID{ID: *calculatedID}, nil
	}
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}

	id, err := h.db.SetBlob(ctx, jetID, parcel.Pulse(), msg.Memory)
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
	jetID := jetFromContext(ctx)

	codeRec, err := getCode(ctx, h.db, jetID, msg.Code.Record())
	if err == storage.ErrNotFound {
		// The record wasn't found on the current node. Return redirect to the node that contains it.
		var node *core.RecordRef
		node, err := h.nodeForJet(ctx, jetID, parcel.Pulse(), msg.Code.Record().Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetCodeRedirect(h.DelegationTokenFactory, parcel, node)
	}
	if err != nil {
		return nil, err
	}
	code, err := h.db.GetBlob(ctx, jetID, codeRec.Code)
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

	var (
		idx *index.ObjectLifeline
		err error
	)

	// Fetch object index. If not found redirect.
	idx, err = h.db.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
	if err == storage.ErrNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		_, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Head, heavy)
		if err != nil {
			return nil, err
		}
		// Add requested object to recent.
		h.RecentStorageProvider.GetStorage(jetID).AddObject(*msg.Head.Record())
		return reply.NewGetObjectRedirectReply(h.DelegationTokenFactory, parcel, heavy, msg.State)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	// Add requested object to recent.
	h.RecentStorageProvider.GetStorage(jetID).AddObject(*msg.Head.Record())

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

	// Fetch state record.
	rec, err := h.db.GetRecord(ctx, jetID, stateID)
	if err == storage.ErrNotFound {
		// The record wasn't found on the current node. Return redirect to the node that contains it.
		node, err := h.nodeForJet(ctx, jetID, parcel.Pulse(), stateID.Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetObjectRedirectReply(h.DelegationTokenFactory, parcel, node, stateID)
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
		rep.Memory, err = h.db.GetBlob(ctx, jetID, state.GetMemory())
		if err != nil {
			return nil, err
		}
	}

	return &rep, nil
}

func (h *MessageHandler) handleHasPendingRequests(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetPendingRequests)
	jetID := jetFromContext(ctx)

	for _, reqID := range h.RecentStorageProvider.GetStorage(jetID).GetRequestsForObject(*msg.Object.Record()) {
		if reqID.Pulse() < parcel.Pulse() {
			return &reply.HasPendingRequests{Has: true}, nil
		}
	}

	return &reply.HasPendingRequests{Has: false}, nil
}

func (h *MessageHandler) handleGetJet(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetJet)
	tree, err := h.db.GetJetTree(ctx, msg.Object.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jet tree")
	}
	jetID, actual := tree.Find(msg.Object)
	if err != nil {
		return nil, err
	}

	return &reply.Jet{ID: *jetID, Actual: actual}, nil
}

func (h *MessageHandler) handleGetDelegate(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)
	jetID := jetFromContext(ctx)

	var (
		idx *index.ObjectLifeline
		err error
	)

	idx, err = h.db.GetObjectIndex(ctx, jetID, msg.Head.Record(), false)
	if err == storage.ErrNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Head, heavy)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}

	h.RecentStorageProvider.GetStorage(jetID).AddObject(*msg.Head.Record())

	delegateRef, ok := idx.Delegates[msg.AsType]
	if !ok {
		return nil, ErrNotFound
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

	idx, err := h.db.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
	if err == storage.ErrNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		_, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Parent, heavy)
		if err != nil {
			return nil, err
		}
		return reply.NewGetChildrenRedirect(h.DelegationTokenFactory, parcel, heavy)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	h.RecentStorageProvider.GetStorage(jetID).AddObject(*msg.Parent.Record())

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

	// We don't have this child reference.
	if currentChild != nil && currentChild.Pulse() != parcel.Pulse() {
		node, err := h.nodeForJet(ctx, jetID, parcel.Pulse(), currentChild.Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetChildrenRedirect(h.DelegationTokenFactory, parcel, node)
	}

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		counter++

		rec, err := h.db.GetRecord(ctx, jetID, currentChild)
		// We don't have this child reference. Return what was collected.
		if err == storage.ErrNotFound {
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

func (h *MessageHandler) handleUpdateObject(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.UpdateObject)
	jetID := jetFromContext(ctx)

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("wrong object state record")
	}

	recentStorage := h.RecentStorageProvider.GetStorage(jetID)
	var idx *index.ObjectLifeline
	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		var err error
		inslog.Debugf("Get index for: %v, jet: %v", msg.Object.Record(), jetID)
		idx, err = tx.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
		// No index on our node.
		if err == storage.ErrNotFound {
			if state.State() == record.StateActivation {
				// We are activating the object. There is no index for it anywhere.
				idx = &index.ObjectLifeline{State: record.StateUndefined}
			} else {
				inslog.Debugf("Not found index for: %v, jet: %v", msg.Object.Record(), jetID)
				// We are updating object. Index should be on the heavy executor.
				heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
				if err != nil {
					return err
				}
				idx, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Object, heavy)
				if err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		}
		if err = validateState(idx.State, state.State()); err != nil {
			return err
		}
		// Index exists and latest record id does not match (preserving chain consistency).
		if idx.LatestState != nil && !state.PrevStateID().Equal(idx.LatestState) {
			return errors.New("invalid state record")
		}

		recentStorage.AddObject(*msg.Object.Record())

		id, err := tx.SetRecord(ctx, jetID, parcel.Pulse(), rec)
		if err != nil {
			return err
		}
		idx.LatestState = id
		idx.State = state.State()
		if state.State() == record.StateActivation {
			idx.Parent = state.(*record.ObjectActivateRecord).Parent
		}
		inslog.Debugf("Save index for: %v, jet: %v", msg.Object.Record(), jetID)
		return tx.SetObjectIndex(ctx, jetID, msg.Object.Record(), idx)
	})
	if err != nil {
		if err == ErrObjectDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}

	recentStorage.AddObject(*msg.Object.Record())

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
	msg := parcel.Message().(*message.RegisterChild)

	jetID := jetFromContext(ctx)

	rec := record.DeserializeRecord(msg.Record)
	childRec, ok := rec.(*record.ChildRecord)
	if !ok {
		return nil, errors.New("wrong child record")
	}

	recentStorage := h.RecentStorageProvider.GetStorage(jetID)
	var child *core.RecordID
	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, err := h.db.GetObjectIndex(ctx, jetID, msg.Parent.Record(), false)
		if err == storage.ErrNotFound {
			heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return err
			}
			idx, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Parent, heavy)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		recentStorage.AddObject(*msg.Parent.Record())

		// Children exist and pointer does not match (preserving chain consistency).
		if idx.ChildPointer != nil && !childRec.PrevChild.Equal(idx.ChildPointer) {
			return errors.New("invalid child record")
		}

		child, err = tx.SetRecord(ctx, jetID, parcel.Pulse(), childRec)
		if err != nil {
			return err
		}
		idx.ChildPointer = child
		if msg.AsType != nil {
			idx.Delegates[*msg.AsType] = msg.Child
		}
		err = tx.SetObjectIndex(ctx, jetID, msg.Parent.Record(), idx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	recentStorage.AddObject(*msg.Parent.Record())

	return &reply.ID{ID: *child}, nil
}

func (h *MessageHandler) handleJetDrop(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.JetDrop)

	if !hack.SkipValidation(ctx) {
		for _, parcelBuff := range msg.Messages {
			parcel, err := message.Deserialize(bytes.NewBuffer(parcelBuff))
			if err != nil {
				return nil, err
			}

			handler, ok := h.replayHandlers[parcel.Message().Type()]
			if !ok {
				return nil, errors.New("unknown message type")
			}

			_, err = handler(ctx, parcel)
			if err != nil {
				return nil, err
			}
		}
	}

	err := h.db.AddJets(ctx, msg.JetID)
	if err != nil {
		return nil, err
	}

	err = h.db.UpdateJetTree(
		ctx,
		parcel.Pulse(),
		true,
		msg.JetID,
	)
	if err != nil {
		return nil, err
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.ValidateRecord)
	jetID := jetFromContext(ctx)

	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, err := tx.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
		if err == storage.ErrNotFound {
			heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return err
			}
			idx, err = h.saveIndexFromHeavy(ctx, h.db, jetID, msg.Object, heavy)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Find node that has this state.
		node, err := h.nodeForJet(ctx, jetID, parcel.Pulse(), msg.Object.Record().Pulse())
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
			err = tx.SetObjectIndex(ctx, jetID, msg.Object.Record(), idx)
			if err != nil {
				return errors.Wrap(err, "failed to save object index")
			}
		case *reply.NotOK:
			return errors.New("validation sequence integrity failure")
		default:
			return errors.New("unexpected reply")
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

	idx, err := h.db.GetObjectIndex(ctx, jetID, msg.Object.Record(), true)
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

	rec, err := h.db.GetRecord(ctx, jetID, &msg.ValidatedState)
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

func getCode(ctx context.Context, s storage.Store, jetID core.RecordID, id *core.RecordID) (*record.CodeRecord, error) {

	rec, err := s.GetRecord(ctx, jetID, id)
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
	ctx context.Context, s storage.Store, jetID core.RecordID, obj core.RecordRef, heavy *core.RecordRef,
) (*index.ObjectLifeline, error) {

	genericReply, err := h.Bus.Send(ctx, &message.GetObjectIndex{
		Object: obj,
	}, &core.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return nil, errors.New("failed to fetch object index: unexpected reply")
	}
	idx, err := index.DecodeObjectLifeline(rep.Index)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	err = s.SetObjectIndex(ctx, jetID, obj.Record(), idx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}
	return idx, nil
}

func (h *MessageHandler) handleHotRecords(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	inslog := inslogger.FromContext(ctx)
	if hack.SkipValidation(ctx) {
		return &reply.OK{}, nil
	}

	msg := parcel.Message().(*message.HotData)
	// FIXME: check split signatures.
	jetID := *msg.Jet.Record()

	err := h.db.SetDrop(ctx, jetID, &msg.Drop)
	if err != nil {
		return nil, errors.Wrap(err, "[ handleHotRecords ] Can't SetDrop")
	}
	err = h.db.UpdateJetTree(
		ctx,
		parcel.Pulse(),
		true,
		jetID,
	)
	if err != nil {
		return nil, err
	}

	// TODO: @andreyromancev. 09.01.2019. Remove after multijet works properly.
	err = h.db.UpdateJetTree(
		ctx,
		parcel.Pulse(),
		true,
		*jet.NewID(2, []byte{1 << 7}), // 10
		*jet.NewID(2, []byte{1 << 6}), // 01
	)
	if err != nil {
		return nil, err
	}

	recentStorage := h.RecentStorageProvider.GetStorage(jetID)
	for objID, requests := range msg.PendingRequests {
		for reqID, request := range requests {
			newID, err := h.db.SetRecord(ctx, jetID, reqID.Pulse(), record.DeserializeRecord(request))
			if err != nil {
				inslog.Error(err)
				continue
			}
			if !bytes.Equal(reqID.Bytes(), newID.Bytes()) {
				inslog.Errorf(
					"Problems with saving the pending request, ids don't match - %v  %v",
					reqID.Bytes(),
					newID.Bytes(),
				)
				continue
			}
			recentStorage.AddPendingRequest(objID, reqID)
		}
	}

	for id, meta := range msg.RecentObjects {
		decodedIndex, err := index.DecodeObjectLifeline(meta.Index)
		if err != nil {
			inslog.Error(err)
			continue
		}

		err = h.db.SetObjectIndex(ctx, jetID, &id, decodedIndex)
		if err != nil {
			inslog.Error(err)
			continue
		}

		meta.TTL--
		recentStorage.AddObjectWithTLL(id, meta.TTL)
	}

	err = h.db.AddJets(ctx, jetID)
	if err != nil {
		return nil, err
	}

	err = h.db.SetDropSizeHistory(ctx, jetID, msg.JetDropSizeHistory)
	if err != nil {
		return nil, errors.Wrap(err, "[ handleHotRecords ] Can't SetDropSizeHistory")
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) nodeForJet(
	ctx context.Context, jetID core.RecordID, parcelPulse, targetPulse core.PulseNumber,
) (*core.RecordRef, error) {
	if targetPulse == core.PulseNumberCurrent {
		targetPulse = parcelPulse
	}
	if parcelPulse-targetPulse < h.conf.LightChainLimit {
		return h.JetCoordinator.LightExecutorForJet(ctx, jetID, targetPulse)
	}
	return h.JetCoordinator.Heavy(ctx, parcelPulse)
}
