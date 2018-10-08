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
	"github.com/insolar/insolar/ledger/index"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

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

	bus.MustRegister(core.TypeGetCode, h.handleGetCode)
	bus.MustRegister(core.TypeGetClass, h.handleGetClass)
	bus.MustRegister(core.TypeGetObject, h.handleGetObject)
	bus.MustRegister(core.TypeGetDelegate, h.handleGetDelegate)
	bus.MustRegister(message.TypeGetChildren, h.handleGetChildren)
	bus.MustRegister(core.TypeDeclareType, h.handleDeclareType)
	bus.MustRegister(core.TypeDeployCode, h.handleDeployCode)
	bus.MustRegister(core.TypeActivateClass, h.handleActivateClass)
	bus.MustRegister(core.TypeDeactivateClass, h.handleDeactivateClass)
	bus.MustRegister(core.TypeUpdateClass, h.handleUpdateClass)
	bus.MustRegister(core.TypeActivateObject, h.handleActivateObject)
	bus.MustRegister(core.TypeActivateObjectDelegate, h.handleActivateObjectDelegate)
	bus.MustRegister(core.TypeDeactivateObject, h.handleDeactivateObject)
	bus.MustRegister(core.TypeUpdateObject, h.handleUpdateObject)
	bus.MustRegister(core.TypeRegisterChild, h.handleRegisterChild)
	bus.MustRegister(core.TypeRequestCall, h.handleRegisterRequest)

	return nil
}

func (h *MessageHandler) handleRegisterRequest(
	genericMsg core.Message,
) (core.Reply, error) {
	msg := genericMsg.(*message.RequestCall)
	requestRec := &record.CallRequest{
		Payload: message.MustSerializeBytes(msg.Message),
	}
	id, err := h.db.SetRequest(requestRec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set request record")
	}
	return &reply.ID{ID: *id.CoreID()}, nil
}

func (h *MessageHandler) handleGetCode(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetCode)
	codeRef := record.Core2Reference(msg.Code)
	rec, err := h.db.GetRecord(&codeRef.Record)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}
	code, mt, err := codeRec.GetCode(msg.MachinePref)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code from record")
	}

	rep := reply.Code{
		Code:        code,
		MachineType: mt,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetClass(genericMsg core.Message) (core.Reply, error) {
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
		Head:  msg.Head,
		State: *stateID,
		Code:  code,
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetObject)
	headRef := record.Core2Reference(msg.Head)

	idx, stateID, state, err := getObject(h.db, &headRef.Record, msg.State)
	if err != nil {
		return nil, err
	}

	rep := reply.Object{
		Head:   msg.Head,
		State:  *stateID,
		Class:  *idx.ClassRef.CoreRef(),
		Memory: state.GetMemory(),
	}

	return &rep, nil
}

func (h *MessageHandler) handleGetDelegate(genericMsg core.Message) (core.Reply, error) {
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

	return &rep, nil
}

func (h *MessageHandler) handleGetChildren(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetChildren)
	parentRef := record.Core2Reference(msg.Parent)

	idx, _, _, err := getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var (
		refs      []core.RecordRef
		fromChild *record.ID
	)

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		id := record.Bytes2ID(msg.FromChild[:])
		fromChild = &id
	} else {
		fromChild = idx.LatestChild
	}

	i := storage.NewChainIterator(h.db, fromChild)
	counter := 0
	for i.HasNext() {
		id, rec, err := i.Next()
		if err != nil {
			return nil, errors.New("failed to retrieve children")
		}

		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: id.CoreID()}, nil
		}
		counter++

		child, ok := rec.(*record.ChildRecord)
		if !ok {
			return nil, errors.New("failed to retrieve children")
		}

		// Skip records later than specified pulse.
		if msg.FromPulse != nil && child.Ref.Record.Pulse > *msg.FromPulse {
			continue
		}

		refs = append(refs, *child.Ref.CoreRef())
	}

	return &reply.Children{Refs: refs, NextFrom: nil}, nil
}

func (h *MessageHandler) handleDeclareType(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.DeclareType)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)

	rec := record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TypeDeclaration: msg.TypeDec,
	}
	typeID, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}

	return &reply.Reference{Ref: *getReference(&msg.Request, typeID)}, nil
}

func (h *MessageHandler) handleDeployCode(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.DeployCode)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)

	rec := record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TargetedCode: msg.CodeMap,
	}
	codeID, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return &reply.Reference{Ref: *getReference(&msg.Request, codeID)}, nil
}

func (h *MessageHandler) handleActivateClass(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.ActivateClass)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)

	rec := record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
	}

	var err error
	var classID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		classID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetClassIndex(classID, &index.ClassLifeline{
			LatestState: *classID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reply.Reference{Ref: *getReference(&msg.Request, classID)}, nil
}

func (h *MessageHandler) handleDeactivateClass(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.DeactivateClass)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)

	var (
		err            error
		deactivationID *record.ID
	)
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getClass(tx, &classRef.Record, nil)
		if err != nil {
			return err
		}
		rec := record.DeactivationRecord{
			AmendRecord: record.AmendRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
				AmendedRecord: idx.LatestState,
			},
		}

		deactivationID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestState = *deactivationID
		err = tx.SetClassIndex(&classRef.Record, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *deactivationID.CoreID()}, nil
}

func (h *MessageHandler) handleUpdateClass(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.UpdateClass)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	migrationRefs := make([]record.Reference, 0, len(msg.Class))
	for _, migration := range msg.Migrations {
		migrationRefs = append(migrationRefs, record.Core2Reference(migration))
	}

	var err error
	err = validateCode(h.db, &msg.Code)
	if err != nil {
		return nil, err
	}
	for _, migration := range msg.Migrations {
		err = validateCode(h.db, &migration)
		if err != nil {
			return nil, err
		}
	}

	var amendID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getClass(tx, &classRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.ClassAmendRecord{
			AmendRecord: record.AmendRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
				AmendedRecord: idx.LatestState,
			},
			NewCode:    record.Core2Reference(msg.Code),
			Migrations: migrationRefs,
		}

		amendID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestState = *amendID
		idx.AmendRefs = append(idx.AmendRefs, *amendID)
		err = tx.SetClassIndex(&classRef.Record, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *amendID.CoreID()}, nil
}

func (h *MessageHandler) handleActivateObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.ActivateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	parentRef := record.Core2Reference(msg.Parent)

	var err error
	_, _, _, err = getClass(h.db, &classRef.Record, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var objID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		rec := record.ObjectActivateRecord{
			ActivationRecord: record.ActivationRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
			},
			ClassActivateRecord: classRef,
			Memory:              msg.Memory,
			Parent:              parentRef,
			Delegate:            false,
		}

		// save new record and it's index
		objID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objID, &index.ObjectLifeline{
			ClassRef:    classRef,
			LatestState: *objID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's children
		parentIdx, err := tx.GetObjectIndex(&parentRef.Record)
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

	return &reply.Reference{Ref: *getReference(&msg.Request, objID)}, nil
}

func (h *MessageHandler) handleActivateObjectDelegate(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.ActivateObjectDelegate)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	parentRef := record.Core2Reference(msg.Parent)

	var err error
	_, _, _, err = getClass(h.db, &classRef.Record, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &parentRef.Record, nil)
	if err != nil {
		return nil, err
	}

	var objID *record.ID
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		rec := record.ObjectActivateRecord{
			ActivationRecord: record.ActivationRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
			},
			ClassActivateRecord: classRef,
			Memory:              msg.Memory,
			Parent:              parentRef,
			Delegate:            true,
		}

		// save new record and it's index
		objID, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objID, &index.ObjectLifeline{
			ClassRef:    classRef,
			LatestState: *objID,
			Delegates:   map[core.RecordRef]record.Reference{},
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's delegates
		parentIdx, err := tx.GetObjectIndex(&parentRef.Record)
		if err != nil {
			return errors.Wrap(err, "inconsistent index")
		}
		if _, ok := parentIdx.Delegates[msg.Class]; ok {
			return ErrClassDelegateAlreadyExists
		}
		parentIdx.Delegates[msg.Class] = record.Reference{
			Record: *objID,
			Domain: record.Core2Reference(msg.Request).Domain,
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

	return &reply.Reference{Ref: *getReference(&msg.Request, objID)}, nil
}

func (h *MessageHandler) handleDeactivateObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.DeactivateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	objRef := record.Core2Reference(msg.Object)

	var (
		err            error
		deactivationID *record.ID
	)
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &objRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.DeactivationRecord{
			AmendRecord: record.AmendRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
				AmendedRecord: idx.LatestState,
			},
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

	return &reply.ID{ID: *deactivationID.CoreID()}, nil
}

func (h *MessageHandler) handleUpdateObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.UpdateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	objRef := record.Core2Reference(msg.Object)

	var (
		err     error
		amendID *record.ID
	)
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &objRef.Record, nil)
		if err != nil {
			return err
		}

		rec := record.ObjectAmendRecord{
			AmendRecord: record.AmendRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						DomainRecord:  domainRef,
						RequestRecord: requestRef,
					},
				},
				AmendedRecord: idx.LatestState,
			},
			NewMemory: msg.Memory,
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

	return &reply.ID{ID: *amendID.CoreID()}, nil
}

func (h *MessageHandler) handleRegisterChild(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.RegisterChild)
	parentRef := record.Core2Reference(msg.Parent)

	var child *record.ID
	err := h.db.Update(func(tx *storage.TransactionManager) error {
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

	return &reply.ID{ID: *child.CoreID()}, nil
}

func getReference(request *core.RecordRef, id *record.ID) *core.RecordRef {
	ref := record.Reference{
		Record: *id,
		Domain: record.Core2Reference(*request).Domain,
	}
	return ref.CoreRef()
}

func getClass(
	s storage.Store, head *record.ID, state *core.RecordRef,
) (*index.ClassLifeline, *core.RecordID, record.ClassState, error) {
	idx, err := s.GetClassIndex(head)
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
	idx, err := s.GetObjectIndex(head)
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

func validateCode(s storage.Store, ref *core.RecordRef) error {
	codeRef := record.Core2Reference(*ref)
	rec, err := s.GetRecord(&codeRef.Record)
	if err != nil {
		return err
	}

	if _, ok := rec.(*record.CodeRecord); !ok {
		return errors.New("invalid code reference")
	}

	return nil
}
