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

	"github.com/insolar/insolar/ledger/index"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type internalHandler func(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	db              *storage.DB
	jetDropHandlers map[core.MessageType]internalHandler
}

// NewMessageHandler creates new handler.
func NewMessageHandler(db *storage.DB) (*MessageHandler, error) {
	return &MessageHandler{
		db:              db,
		jetDropHandlers: map[core.MessageType]internalHandler{},
	}, nil
}

// Link links external components.
func (h *MessageHandler) Link(components core.Components) error {
	bus := components.MessageBus

	bus.MustRegister(core.TypeGetCode, h.messagePersistingWrapper(h.handleGetCode))
	bus.MustRegister(core.TypeGetClass, h.messagePersistingWrapper(h.handleGetClass))
	bus.MustRegister(core.TypeGetObject, h.messagePersistingWrapper(h.handleGetObject))
	bus.MustRegister(core.TypeGetDelegate, h.messagePersistingWrapper(h.handleGetDelegate))
	bus.MustRegister(core.TypeGetChildren, h.messagePersistingWrapper(h.handleGetChildren))
	bus.MustRegister(core.TypeUpdateObject, h.messagePersistingWrapper(h.handleUpdateObject))
	bus.MustRegister(core.TypeRegisterChild, h.messagePersistingWrapper(h.handleRegisterChild))
	bus.MustRegister(core.TypeJetDrop, h.handleJetDrop)
	bus.MustRegister(core.TypeSetRecord, h.messagePersistingWrapper(h.handleSetRecord))
	bus.MustRegister(core.TypeSetBlob, h.messagePersistingWrapper(h.handleSetBlob))
	bus.MustRegister(core.TypeUpdateClass, h.messagePersistingWrapper(h.handleUpdateClass))
	bus.MustRegister(core.TypeValidateRecord, h.messagePersistingWrapper(h.handleValidateRecord))

	h.jetDropHandlers[core.TypeGetCode] = h.handleGetCode
	h.jetDropHandlers[core.TypeGetClass] = h.handleGetClass
	h.jetDropHandlers[core.TypeGetObject] = h.handleGetObject
	h.jetDropHandlers[core.TypeGetDelegate] = h.handleGetDelegate
	h.jetDropHandlers[core.TypeGetChildren] = h.handleGetChildren
	h.jetDropHandlers[core.TypeUpdateObject] = h.handleUpdateObject
	h.jetDropHandlers[core.TypeRegisterChild] = h.handleRegisterChild
	h.jetDropHandlers[core.TypeSetRecord] = h.handleSetRecord
	h.jetDropHandlers[core.TypeUpdateClass] = h.handleUpdateClass
	h.jetDropHandlers[core.TypeValidateRecord] = h.handleValidateRecord

	return nil
}

func (h *MessageHandler) messagePersistingWrapper(handler internalHandler) core.MessageHandler {
	return func(context context.Context, genericMsg core.SignedMessage) (core.Reply, error) {
		err := persistMessageToDb(h.db, genericMsg.Message())
		if err != nil {
			return nil, err
		}

		lastPulseNumber, err := h.db.GetLatestPulseNumber()
		if err != nil {
			return nil, err
		}

		return handler(context, lastPulseNumber, genericMsg)
	}
}

func (h *MessageHandler) handleSetRecord(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.SetRecord)
	id, err := h.db.SetRecord(pulseNumber, record.DeserializeRecord(msg.Record))
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleSetBlob(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.SetBlob)

	calculatedID, err := record.CalculateIDForBlob(pulseNumber, msg.Memory)
	if err != nil {
		return nil, err
	}
	_, err = h.db.GetBlob(calculatedID)
	if err == nil {
		return &reply.ID{ID: *calculatedID}, nil
	}
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	id, err := h.db.SetBlob(pulseNumber, msg.Memory)
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleGetCode(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetCode)

	codeRec, err := getCode(h.db, msg.Code.Record())
	if err != nil {
		return nil, err
	}
	code, err := h.db.GetBlob(codeRec.Code)
	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        code,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetClass(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetClass)

	_, stateID, state, err := getClass(h.db, msg.Head.Record(), msg.State)
	if err != nil {
		return nil, err
	}

	var code *core.RecordRef
	if state.GetCode() == nil {
		code = nil
	} else {
		code = state.GetCode()
	}

	rep := reply.Class{
		Head:        msg.Head,
		State:       *stateID,
		Code:        code,
		MachineType: state.GetMachineType(),
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetObject(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetObject)

	idx, stateID, state, err := getObject(h.db, msg.Head.Record(), msg.State, msg.Approved)
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
		Class:        idx.ClassRef,
		ChildPointer: childPointer,
	}

	if state.GetMemory() != nil {
		rep.Memory, err = h.db.GetBlob(state.GetMemory())
		if err != nil {
			return nil, err
		}
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetDelegate(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetDelegate)

	idx, _, _, err := getObject(h.db, msg.Head.Record(), nil, false)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := idx.Delegates[msg.AsClass]
	if !ok {
		return nil, ErrNotFound
	}

	rep := reply.Delegate{
		Head: delegateRef,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetChildren(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetChildren)

	idx, _, _, err := getObject(h.db, msg.Parent.Record(), nil, false)
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

		rec, err := h.db.GetRecord(currentChild)
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

func (h *MessageHandler) handleUpdateClass(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.UpdateClass)

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ClassState)
	if !ok {
		return nil, errors.New("wrong class state record")
	}

	var id *core.RecordID
	err := h.db.Update(func(tx *storage.TransactionManager) error {
		idx, err := getClassIndex(tx, msg.Class.Record(), true)
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

		id, err = tx.SetRecord(pulseNumber, rec)
		if err != nil {
			return err
		}
		idx.LatestState = id
		idx.State = state.State()
		return tx.SetClassIndex(msg.Class.Record(), idx)
	})
	if err != nil {
		if err == ErrClassDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}
	return &reply.ID{ID: *id}, nil
}

func (h *MessageHandler) handleUpdateObject(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.UpdateObject)

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ObjectState)
	if !ok {
		return nil, errors.New("wrong class state record")
	}

	var idx *index.ObjectLifeline
	err := h.db.Update(func(tx *storage.TransactionManager) error {
		var err error
		idx, err = getObjectIndex(tx, msg.Object.Record(), true)
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

		id, err := tx.SetRecord(pulseNumber, rec)
		if err != nil {
			return err
		}
		idx.LatestState = id
		idx.State = state.State()
		if state.State() == record.StateActivation {
			if msg.Class == nil {
				return errors.New("not enough data for activation provided")
			}
			idx.ClassRef = *msg.Class
		}
		return tx.SetObjectIndex(msg.Object.Record(), idx)
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
		Class:        idx.ClassRef,
		ChildPointer: idx.ChildPointer,
	}
	return &rep, nil
}

func (h *MessageHandler) handleRegisterChild(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.RegisterChild)

	rec := record.DeserializeRecord(msg.Record)
	childRec, ok := rec.(*record.ChildRecord)
	if !ok {
		return nil, errors.New("wrong child record")
	}

	var child *core.RecordID
	err := h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, msg.Parent.Record(), nil, false)
		if err != nil {
			return err
		}
		// Children exist and pointer does not match (preserving chain consistency).
		if idx.ChildPointer != nil && !childRec.PrevChild.Equal(idx.ChildPointer) {
			return errors.New("invalid child record")
		}

		child, err = tx.SetRecord(pulseNumber, childRec)
		if err != nil {
			return err
		}
		idx.ChildPointer = child
		if msg.AsClass != nil {
			idx.Delegates[*msg.AsClass] = msg.Child
		}
		err = tx.SetObjectIndex(msg.Parent.Record(), idx)
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

func (h *MessageHandler) handleJetDrop(ctx context.Context, genericMsg core.SignedMessage) (core.Reply, error) {
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

func (h *MessageHandler) handleValidateRecord(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.ValidateRecord)

	err := h.db.Update(func(tx *storage.TransactionManager) error {
		idx, err := tx.GetObjectIndex(msg.Object.Record(), true)
		if err != nil {
			return errors.Wrap(err, "inconsistent object index")
		}

		// Rewinding to validated record.
		currentID := idx.LatestState
		for currentID != nil {
			// We have passed an approved record.
			if currentID.Equal(idx.LatestStateApproved) {
				return errors.New("changing approved records is not allowed")
			}

			// Fetching actual record.
			rec, err := tx.GetRecord(currentID)
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
				err := tx.SetObjectIndex(msg.Object.Record(), idx)
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

func persistMessageToDb(db *storage.DB, genericMsg core.Message) error {
	lastPulse, err := db.GetLatestPulseNumber()
	if err != nil {
		return err
	}
	err = db.SetMessage(lastPulse, genericMsg)
	if err != nil {
		return err
	}

	return nil
}

func getCode(s storage.Store, id *core.RecordID) (*record.CodeRecord, error) {
	rec, err := s.GetRecord(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}

	return codeRec, nil
}

func getClass(
	s storage.Store, head *core.RecordID, state *core.RecordID,
) (*index.ClassLifeline, *core.RecordID, record.ClassState, error) {
	idx, err := s.GetClassIndex(head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}

	var stateID *core.RecordID
	if state != nil {
		stateID = state
	} else {
		stateID = idx.LatestState
	}

	rec, err := s.GetRecord(stateID)
	if err != nil {
		return nil, nil, nil, err
	}
	stateRec, ok := rec.(record.ClassState)
	if !ok {
		return nil, nil, nil, errors.New("invalid class record")
	}
	if stateRec.State() == record.StateDeactivation {
		return nil, nil, nil, ErrClassDeactivated
	}

	return idx, stateID, stateRec, nil
}

func getObject(
	s storage.Store, head *core.RecordID, state *core.RecordID, approved bool,
) (*index.ObjectLifeline, *core.RecordID, record.ObjectState, error) {
	idx, err := s.GetObjectIndex(head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
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

	rec, err := s.GetRecord(stateID)
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

func getClassIndex(s storage.Store, head *core.RecordID, forupdate bool) (*index.ClassLifeline, error) {
	idx, err := s.GetClassIndex(head, forupdate)
	if err == storage.ErrNotFound {
		return &index.ClassLifeline{State: record.StateUndefined}, nil
	}
	return idx, err
}

func getObjectIndex(s storage.Store, head *core.RecordID, forupdate bool) (*index.ObjectLifeline, error) {
	idx, err := s.GetObjectIndex(head, forupdate)
	if err == storage.ErrNotFound {
		return &index.ObjectLifeline{State: record.StateUndefined}, nil
	}
	return idx, err
}

func validateState(old record.State, new record.State) error {
	if old == record.StateDeactivation {
		return ErrClassDeactivated
	}
	if old == record.StateUndefined && new != record.StateActivation {
		return errors.New("object is not activated")
	}
	if old != record.StateUndefined && new == record.StateActivation {
		return errors.New("object is already activated")
	}
	return nil
}
