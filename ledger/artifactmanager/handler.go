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
	"sync"
	"time"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

var handlers = map[core.MessageType]core.MessageHandler{}
var handlersRWLock = sync.RWMutex{}

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	db *storage.DB
}

// NewMessageHandler creates new handler.
func NewMessageHandler(db *storage.DB) (*MessageHandler, error) {
	return &MessageHandler{db: db}, nil
}

// Link links external components.
func (h *MessageHandler) Link(components core.Components) error {
	bus := components.MessageBus

	handlersRWLock.Lock()
	handlers[core.TypeGetCode] = h.handleGetCode
	handlers[core.TypeGetClass] = h.handleGetClass
	handlers[core.TypeGetObject] = h.handleGetObject
	handlers[core.TypeGetDelegate] = h.handleGetDelegate
	handlers[core.TypeGetChildren] = h.handleGetChildren
	handlers[core.TypeActivateObject] = h.handleActivateObject
	handlers[core.TypeActivateObjectDelegate] = h.handleActivateObjectDelegate
	handlers[core.TypeDeactivateObject] = h.handleDeactivateObject
	handlers[core.TypeUpdateObject] = h.handleUpdateObject
	handlers[core.TypeRegisterChild] = h.handleRegisterChild
	handlers[core.TypeJetDrop] = h.handleJetDrop
	handlers[core.TypeSetRecord] = h.handleSetRecord
	handlers[core.TypeUpdateClass] = h.handleUpdateClass
	handlersRWLock.Unlock()

	handlersRWLock.RLock()
	for handlerType, handler := range handlers {
		bus.MustRegister(handlerType, handler)
	}
	handlersRWLock.RUnlock()

	return nil
}

func logTimeInside(start time.Time, funcName string) {
	if time.Since(start) > time.Second {
		log.Debugf("Handle takes too long: %s: time inside - %s", funcName, time.Since(start))
	}
}

func (h *MessageHandler) handleSetRecord(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.SetRecord)

	id, err := h.db.SetRecord(record.DeserializeRecord(msg.Record))
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *id.CoreID()}, nil
}

func (h *MessageHandler) handleGetCode(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.GetCode)
	codeRef := record.Core2Reference(msg.Code)

	codeRec, err := getCode(h.db, codeRef.Record)
	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        codeRec.Code,
		MachineType: codeRec.MachineType,
	}

	logTimeInside(start, "handleGetCode")

	return &rep, nil
}

func (h *MessageHandler) handleGetClass(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.GetClass)
	headRef := record.Core2Reference(msg.Head)

	_, stateID, state, err := getClass(h.db, &headRef.Record, msg.State)
	if err != nil {
		return nil, err
	}

	var code *core.RecordRef
	if state.GetCode() == nil {
		code = nil
	} else {
		code = state.GetCode().CoreRef()
	}

	rep := reply.Class{
		Head:        msg.Head,
		State:       *stateID,
		Code:        code,
		MachineType: state.GetMachineType(),
	}

	logTimeInside(start, "handleGetClass")

	return &rep, nil
}

func (h *MessageHandler) handleGetObject(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.GetObject)
	headRef := record.Core2Reference(msg.Head)

	idx, stateID, state, err := getObject(h.db, &headRef.Record, msg.State)
	if err != nil {
		if err == ErrObjectDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}

	rep := reply.Object{
		Head:   msg.Head,
		State:  *stateID,
		Class:  *idx.ClassRef.CoreRef(),
		Memory: state.GetMemory(),
	}

	logTimeInside(start, "handleGetObject")

	return &rep, nil
}

func (h *MessageHandler) handleGetDelegate(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.GetDelegate)
	headRef := record.Core2Reference(msg.Head)

	idx, _, _, err := getObject(h.db, &headRef.Record, nil)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := idx.Delegates[msg.AsClass]
	if !ok {
		return nil, ErrNotFound
	}

	rep := reply.Delegate{
		Head: *delegateRef.CoreRef(),
	}

	logTimeInside(start, "handleGetDelegate")

	return &rep, nil
}

func (h *MessageHandler) handleGetChildren(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.GetChildren)
	parentRef := record.Core2Reference(msg.Parent)

	idx, _, _, err := getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var (
		refs         []core.RecordRef
		currentChild *record.ID
	)

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		id := record.Bytes2ID(msg.FromChild[:])
		currentChild = &id
	} else {
		currentChild = idx.LatestChild
	}

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: currentChild.CoreID()}, nil
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
		if msg.FromPulse != nil && childRec.Ref.Record.Pulse > *msg.FromPulse {
			continue
		}
		refs = append(refs, *childRec.Ref.CoreRef())
	}

	logTimeInside(start, "handleGetChildren")

	return &reply.Children{Refs: refs, NextFrom: nil}, nil
}

func (h *MessageHandler) handleUpdateClass(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.UpdateClass)
	classCoreID := msg.Class.GetRecordID()
	classID := record.Bytes2ID(classCoreID[:])

	rec := record.DeserializeRecord(msg.Record)
	state, ok := rec.(record.ClassState)
	if !ok {
		return nil, errors.New("wrong class state record")
	}

	var id *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, err := getClassIndex(tx, &classID, true)
		if err != nil {
			return err
		}
		if idx.Deactivated {
			return ErrClassDeactivated
		}
		id, err = tx.SetRecord(rec)
		if err != nil {
			return err
		}
		idx.LatestState = *id
		idx.Deactivated = state.IsDeactivation()
		return tx.SetClassIndex(&classID, idx)
	})
	if err != nil {
		if err == ErrClassDeactivated {
			return &reply.Error{ErrType: reply.ErrDeactivated}, nil
		}
		return nil, err
	}
	return &reply.ID{ID: *id.CoreID()}, nil
}

func (h *MessageHandler) handleActivateObject(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.ActivateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	parentRef := record.Core2Reference(msg.Parent)

	_, _, _, err = getClass(h.db, &classRef.Record, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var activateID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		rec := record.ObjectActivateRecord{
			ResultRecord: record.ResultRecord{
				Domain:  domainRef,
				Request: requestRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: msg.Memory,
			},
			Class:    classRef,
			Parent:   parentRef,
			Delegate: false,
		}

		// save new record and it's index
		activateID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(&requestRef.Record, &index.ObjectLifeline{
			ClassRef:    classRef,
			LatestState: *activateID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's children
		parentIdx, err := tx.GetObjectIndex(&parentRef.Record, true)
		if err != nil {
			if err == ErrNotFound {
				parentIdx = &index.ObjectLifeline{}
			} else {
				return errors.Wrap(err, "inconsistent index")
			}
		}
		err = tx.SetObjectIndex(&parentRef.Record, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logTimeInside(start, "handleActivateObject")

	return &reply.ID{ID: *activateID.CoreID()}, nil
}

func (h *MessageHandler) handleActivateObjectDelegate(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.ActivateObjectDelegate)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	parentRef := record.Core2Reference(msg.Parent)

	_, _, _, err = getClass(h.db, &classRef.Record, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var activationID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		rec := record.ObjectActivateRecord{
			ResultRecord: record.ResultRecord{
				Domain:  domainRef,
				Request: requestRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: msg.Memory,
			},
			Class:    classRef,
			Parent:   parentRef,
			Delegate: true,
		}

		// save new record and it's index
		activationID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(&requestRef.Record, &index.ObjectLifeline{
			ClassRef:    classRef,
			LatestState: *activationID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's delegates
		parentIdx, err := tx.GetObjectIndex(&parentRef.Record, true)
		if err != nil {
			return errors.Wrap(err, "inconsistent index")
		}
		if _, ok := parentIdx.Delegates[msg.Class]; ok {
			return ErrClassDelegateAlreadyExists
		}
		parentIdx.Delegates[msg.Class] = requestRef
		err = tx.SetObjectIndex(&parentRef.Record, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logTimeInside(start, "handleActivateObjectDelegate")

	return &reply.ID{ID: *activationID.CoreID()}, nil
}

func (h *MessageHandler) handleDeactivateObject(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.DeactivateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	objRef := record.Core2Reference(msg.Object)

	var (
		deactivationID *record.ID
	)
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &objRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.DeactivationRecord{
			ResultRecord: record.ResultRecord{
				Domain:  domainRef,
				Request: requestRef,
			},
			PrevState: idx.LatestState,
		}
		deactivationID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestState = *deactivationID
		err = tx.SetObjectIndex(&objRef.Record, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logTimeInside(start, "handleDeactivateObject")

	return &reply.ID{ID: *deactivationID.CoreID()}, nil
}

func (h *MessageHandler) handleUpdateObject(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.UpdateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	objRef := record.Core2Reference(msg.Object)

	var (
		amendID *record.ID
	)
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &objRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.ObjectAmendRecord{
			ResultRecord: record.ResultRecord{
				Domain:  domainRef,
				Request: requestRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: msg.Memory,
			},
			PrevState: idx.LatestState,
		}

		amendID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestState = *amendID
		err = tx.SetObjectIndex(&objRef.Record, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logTimeInside(start, "handleUpdateObject")

	return &reply.ID{ID: *amendID.CoreID()}, nil
}

func (h *MessageHandler) handleRegisterChild(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	start := time.Now()

	err := persistMessageToDb(h.db, genericMsg)
	if err != nil {
		return nil, err
	}

	msg := genericMsg.(*message.RegisterChild)
	parentRef := record.Core2Reference(msg.Parent)

	var child *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &parentRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.ChildRecord{
			PrevChild: idx.LatestChild,
			Ref:       record.Core2Reference(msg.Child),
		}
		child, err = tx.SetRecord(&rec)
		if err != nil {
			return err
		}
		idx.LatestChild = child
		err = tx.SetObjectIndex(&parentRef.Record, idx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logTimeInside(start, "handleRegisterChild")

	return &reply.ID{ID: *child.CoreID()}, nil
}

func (h *MessageHandler) handleJetDrop(ctx core.Context, genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.JetDrop)

	for _, rawMessage := range msg.Messages {
		parsedMessage, err := message.Deserialize(bytes.NewBuffer(rawMessage))
		if err != nil {
			return nil, err
		}
		handlersRWLock.RLock()
		handler, ok := handlers[parsedMessage.Type()]
		handlersRWLock.RUnlock()
		if !ok {
			return nil, errors.New("unknown message type")
		}

		_, err = handler(ctx, parsedMessage)
		if err != nil {
			return nil, err
		}
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

func getReference(request *core.RecordRef, id *record.ID) *core.RecordRef {
	ref := record.Reference{
		Record: *id,
		Domain: record.Core2Reference(*request).Domain,
	}
	return ref.CoreRef()
}

func getCode(s storage.Store, id record.ID) (*record.CodeRecord, error) {
	rec, err := s.GetRecord(&id)
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
	s storage.Store, head *record.ID, state *core.RecordRef,
) (*index.ClassLifeline, *core.RecordID, record.ClassState, error) {
	idx, err := s.GetClassIndex(head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}

	var stateID record.ID
	if state != nil {
		stateID = record.Core2Reference(*state).Record
	} else {
		stateID = idx.LatestState
	}

	rec, err := s.GetRecord(&stateID)
	if err != nil {
		return nil, nil, nil, err
	}
	stateRec, ok := rec.(record.ClassState)
	if !ok {
		return nil, nil, nil, errors.New("invalid class record")
	}
	if stateRec.IsDeactivation() {
		return nil, nil, nil, ErrClassDeactivated
	}

	return idx, stateID.CoreID(), stateRec, nil
}

func getObject(
	s storage.Store, head *record.ID, state *core.RecordRef,
) (*index.ObjectLifeline, *core.RecordID, record.ObjectState, error) {
	idx, err := s.GetObjectIndex(head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}

	var stateID record.ID
	if state != nil {
		stateID = record.Core2Reference(*state).Record
	} else {
		stateID = idx.LatestState
	}

	rec, err := s.GetRecord(&stateID)
	if err != nil {
		return nil, nil, nil, err
	}
	stateRec, ok := rec.(record.ObjectState)
	if !ok {
		return nil, nil, nil, errors.New("invalid object record")
	}
	if stateRec.IsDeactivation() {
		return nil, nil, nil, ErrObjectDeactivated
	}

	return idx, stateID.CoreID(), stateRec, nil
}

func getClassIndex(s storage.Store, head *record.ID, forupdate bool) (*index.ClassLifeline, error) {
	idx, err := s.GetClassIndex(head, forupdate)
	if err == storage.ErrNotFound {
		return &index.ClassLifeline{}, nil
	}
	return idx, err
}
