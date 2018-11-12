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

	"github.com/insolar/insolar/instrumentation/hack"
	"github.com/insolar/insolar/ledger/index"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type internalHandler func(ctx context.Context, pulseNumber core.PulseNumber, parcel core.Parcel) (core.Reply, error)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	db              *storage.DB
	jetDropHandlers map[core.MessageType]internalHandler
	Bus             core.MessageBus `inject:""`
}

// NewMessageHandler creates new handler.
func NewMessageHandler(db *storage.DB) *MessageHandler {
	return &MessageHandler{
		db:              db,
		jetDropHandlers: map[core.MessageType]internalHandler{},
	}
}

// Start initializes handlers.
func (h *MessageHandler) Start(ctx context.Context) error {
	h.Bus.MustRegister(core.TypeGetCode, h.messagePersistingWrapper(h.handleGetCode))
	h.Bus.MustRegister(core.TypeGetObject, h.messagePersistingWrapper(h.handleGetObject))
	h.Bus.MustRegister(core.TypeGetDelegate, h.messagePersistingWrapper(h.handleGetDelegate))
	h.Bus.MustRegister(core.TypeGetChildren, h.messagePersistingWrapper(h.handleGetChildren))
	h.Bus.MustRegister(core.TypeUpdateObject, h.messagePersistingWrapper(h.handleUpdateObject))
	h.Bus.MustRegister(core.TypeRegisterChild, h.messagePersistingWrapper(h.handleRegisterChild))
	h.Bus.MustRegister(core.TypeJetDrop, h.handleJetDrop)
	h.Bus.MustRegister(core.TypeSetRecord, h.messagePersistingWrapper(h.handleSetRecord))
	h.Bus.MustRegister(core.TypeSetBlob, h.messagePersistingWrapper(h.handleSetBlob))
	h.Bus.MustRegister(core.TypeValidateRecord, h.messagePersistingWrapper(h.handleValidateRecord))

	h.jetDropHandlers[core.TypeGetCode] = h.handleGetCode
	h.jetDropHandlers[core.TypeGetObject] = h.handleGetObject
	h.jetDropHandlers[core.TypeGetDelegate] = h.handleGetDelegate
	h.jetDropHandlers[core.TypeGetChildren] = h.handleGetChildren
	h.jetDropHandlers[core.TypeUpdateObject] = h.handleUpdateObject
	h.jetDropHandlers[core.TypeRegisterChild] = h.handleRegisterChild
	h.jetDropHandlers[core.TypeSetRecord] = h.handleSetRecord
	h.jetDropHandlers[core.TypeValidateRecord] = h.handleValidateRecord

	return nil
}

func (h *MessageHandler) messagePersistingWrapper(handler internalHandler) core.MessageHandler {
	return func(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
		err := persistMessageToDb(ctx, h.db, genericMsg.Message())
		if err != nil {
			return nil, err
		}

		lastPulseNumber, err := h.db.GetLatestPulseNumber(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, lastPulseNumber, genericMsg)
	}
}

func (h *MessageHandler) handleSetRecord(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.SetRecord)
	id, err := h.db.SetRecord(ctx, pulseNumber, record.DeserializeRecord(msg.Record))
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleSetBlob(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.SetBlob)

	calculatedID := record.CalculateIDForBlob(pulseNumber, msg.Memory)
	_, err := h.db.GetBlob(ctx, calculatedID)
	if err == nil {
		return &reply.ID{ID: *calculatedID}, nil
	}
	if err != nil && err != storage.ErrNotFound {
		return nil, err
	}

	id, err := h.db.SetBlob(ctx, pulseNumber, msg.Memory)
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleGetCode(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetCode)

	codeRec, err := getCode(ctx, h.db, msg.Code.Record())
	if err != nil {
		return nil, err
	}
	code, err := h.db.GetBlob(ctx, codeRec.Code)
	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        code,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetObject(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetObject)

	idx, stateID, state, err := getObject(ctx, h.db, msg.Head.Record(), msg.State, msg.Approved)
	if err != nil {
		switch err {
		case ErrObjectDeactivated:
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		case ErrStateNotAvailable:
			return &reply.Error{ErrType: reply.ErrStateNotAvailable}, nil
		default:
			return nil, err
		}
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
		rep.Memory, err = h.db.GetBlob(ctx, state.GetMemory())
		if err != nil {
			return nil, err
		}
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetDelegate(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetDelegate)

	idx, _, _, err := getObject(ctx, h.db, msg.Head.Record(), nil, false)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := idx.Delegates[msg.AsType]
	if !ok {
		return nil, ErrNotFound
	}

	rep := reply.Delegate{
		Head: delegateRef,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetChildren(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetChildren)

	idx, _, _, err := getObject(ctx, h.db, msg.Parent.Record(), nil, false)
	if err != nil {
		return nil, err
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

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		counter++

		rec, err := h.db.GetRecord(ctx, currentChild)
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

func (h *MessageHandler) handleUpdateObject(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.UpdateObject)

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("wrong object state record")
	}

	var idx *index.ObjectLifeline
	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		var err error
		idx, err = getObjectIndex(ctx, tx, msg.Object.Record(), true)
		if err != nil {
			return err
		}
		if err = validateState(idx.State, state.State()); err != nil {
			return err
		}
		// Index exists and latest record id does not match (preserving chain consistency).
		if idx.LatestState != nil && !state.PrevStateID().Equal(idx.LatestState) {
			return errors.New("invalid state record")
		}

		id, err := tx.SetRecord(ctx, pulseNumber, rec)
		if err != nil {
			return err
		}
		idx.LatestState = id
		idx.State = state.State()
		if state.State() == record.StateActivation {
			idx.Parent = state.(*record.ObjectActivateRecord).Parent
		}
		return tx.SetObjectIndex(ctx, msg.Object.Record(), idx)
	})
	if err != nil {
		if err == ErrObjectDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}

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

func (h *MessageHandler) handleRegisterChild(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.RegisterChild)

	rec := record.DeserializeRecord(msg.Record)
	childRec, ok := rec.(*record.ChildRecord)
	if !ok {
		return nil, errors.New("wrong child record")
	}

	var child *core.RecordID
	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(ctx, tx, msg.Parent.Record(), nil, false)
		if err != nil {
			return err
		}
		// Children exist and pointer does not match (preserving chain consistency).
		if idx.ChildPointer != nil && !childRec.PrevChild.Equal(idx.ChildPointer) {
			return errors.New("invalid child record")
		}

		child, err = tx.SetRecord(ctx, pulseNumber, childRec)
		if err != nil {
			return err
		}
		idx.ChildPointer = child
		if msg.AsType != nil {
			idx.Delegates[*msg.AsType] = msg.Child
		}
		err = tx.SetObjectIndex(ctx, msg.Parent.Record(), idx)
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

func (h *MessageHandler) handleJetDrop(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
	if hack.SkipValidation(ctx) {
		return &reply.OK{}, nil
	}
	msg := genericMsg.Message().(*message.JetDrop)

	for _, rawMessage := range msg.Messages {
		parsedMessage, err := message.Deserialize(bytes.NewBuffer(rawMessage))
		if err != nil {
			return nil, err
		}

		handler, ok := h.jetDropHandlers[parsedMessage.Message().Type()]
		if !ok {
			return nil, errors.New("unknown message type")
		}

		_, err = handler(ctx, msg.PulseNumber, parsedMessage)
		if err != nil {
			return nil, err
		}
	}

	return &reply.OK{}, nil
}

func (h *MessageHandler) handleValidateRecord(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.Parcel) (core.Reply, error) {
	msg := genericMsg.Message().(*message.ValidateRecord)

	err := h.db.Update(ctx, func(tx *storage.TransactionManager) error {
		idx, err := tx.GetObjectIndex(ctx, msg.Object.Record(), true)
		if err != nil {
			return errors.Wrap(err, "failed to fetch object index")
		}

		// Rewinding to validated record.
		currentID := idx.LatestState
		for currentID != nil {
			// We have passed an approved record.
			if currentID.Equal(idx.LatestStateApproved) {
				return errors.New("changing approved records is not allowed")
			}

			// Fetching actual record.
			rec, err := tx.GetRecord(ctx, currentID)
			if err != nil {
				return nil
			}
			currentState, ok := rec.(record.ObjectState)
			if !ok {
				return errors.New("invalid object record")
			}

			// Validated record found.
			if currentID.Equal(&msg.State) {
				if msg.IsValid {
					idx.LatestStateApproved = currentID
				} else {
					idx.LatestState = currentState.PrevStateID()
				}
				err := tx.SetObjectIndex(ctx, msg.Object.Record(), idx)
				if err != nil {
					return err
				}
				break
			}

			currentID = currentState.PrevStateID()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reply.OK{}, nil
}

func persistMessageToDb(ctx context.Context, db *storage.DB, genericMsg core.Message) error {
	lastPulse, err := db.GetLatestPulseNumber(ctx)
	if err != nil {
		return err
	}
	err = db.SetMessage(ctx, lastPulse, genericMsg)
	if err != nil {
		return err
	}

	return nil
}

func getCode(ctx context.Context, s storage.Store, id *core.RecordID) (*record.CodeRecord, error) {
	rec, err := s.GetRecord(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}

	return codeRec, nil
}

func getObject(
	ctx context.Context,
	s storage.Store,
	head *core.RecordID,
	state *core.RecordID,
	approved bool,
) (*index.ObjectLifeline, *core.RecordID, record.ObjectState, error) {
	idx, err := s.GetObjectIndex(ctx, head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to fetch object index")
	}

	var stateID *core.RecordID
	if state != nil {
		stateID = state
	} else {
		if approved {
			stateID = idx.LatestStateApproved
		} else {
			stateID = idx.LatestState
		}
	}

	if stateID == nil {
		return nil, nil, nil, ErrStateNotAvailable
	}

	rec, err := s.GetRecord(ctx, stateID)
	if err != nil {
		return nil, nil, nil, err
	}
	stateRec, ok := rec.(record.ObjectState)
	if !ok {
		return nil, nil, nil, errors.New("invalid object record")
	}
	if stateRec.State() == record.StateDeactivation {
		return nil, nil, nil, ErrObjectDeactivated
	}

	return idx, stateID, stateRec, nil
}

func getObjectIndex(ctx context.Context, s storage.Store, head *core.RecordID, forupdate bool) (*index.ObjectLifeline, error) {
	idx, err := s.GetObjectIndex(ctx, head, forupdate)
	if err == storage.ErrNotFound {
		return &index.ObjectLifeline{State: record.StateUndefined}, nil
	}
	return idx, err
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
