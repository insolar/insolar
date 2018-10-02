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
	var err error
	bus := components.MessageBus

	if err = bus.Register(message.TypeGetCode, h.handleGetCode); err != nil {
		return err
	}
	if err = bus.Register(message.TypeGetClass, h.handleGetClass); err != nil {
		return err
	}
	if err = bus.Register(message.TypeGetObject, h.handleGetObject); err != nil {
		return err
	}
	if err = bus.Register(message.TypeGetDelegate, h.handleGetDelegate); err != nil {
		return err
	}
	if err = bus.Register(message.TypeDeclareType, h.handleDeclareType); err != nil {
		return err
	}
	if err = bus.Register(message.TypeDeployCode, h.handleDeployCode); err != nil {
		return err
	}
	if err = bus.Register(message.TypeActivateClass, h.handleActivateClass); err != nil {
		return err
	}
	if err = bus.Register(message.TypeDeactivateClass, h.handleDeactivateClass); err != nil {
		return err
	}
	if err = bus.Register(message.TypeUpdateClass, h.handleUpdateClass); err != nil {
		return err
	}
	if err = bus.Register(message.TypeActivateObject, h.handleActivateObject); err != nil {
		return err
	}
	if err = bus.Register(message.TypeActivateObjectDelegate, h.handleActivateObjectDelegate); err != nil {
		return err
	}
	if err = bus.Register(message.TypeDeactivateObject, h.handleDeactivateObject); err != nil {
		return err
	}
	if err = bus.Register(message.TypeUpdateObject, h.handleUpdateObject); err != nil {
		return err
	}

	return err
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

	react := reply.Code{
		Code:        code,
		MachineType: mt,
	}

	return &react, nil
}

func (h *MessageHandler) handleGetClass(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetClass)

	_, stateRef, state, err := getClass(h.db, &msg.Head, msg.State)
	if err != nil {
		return nil, err
	}

	var code *core.RecordRef
	if state.GetCode() == nil {
		code = nil
	} else {
		code = state.GetCode().CoreRef()
	}

	react := reply.Class{
		Head:  msg.Head,
		State: *stateRef,
		Code:  code,
	}

	return &react, nil
}

func (h *MessageHandler) handleGetObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetObject)

	idx, stateRef, state, err := getObject(h.db, &msg.Head, msg.State)
	if err != nil {
		return nil, err
	}

	children := make([]core.RecordRef, 0, len(idx.Children))
	for _, c := range idx.Children {
		children = append(children, *c.CoreRef())
	}

	react := reply.Object{
		Head:     msg.Head,
		State:    *stateRef,
		Class:    *idx.ClassRef.CoreRef(),
		Memory:   state.GetMemory(),
		Children: children,
	}

	return &react, nil
}

func (h *MessageHandler) handleGetDelegate(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.GetDelegate)

	idx, _, _, err := getObject(h.db, &msg.Head, nil)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := idx.Delegates[msg.AsClass]
	if !ok {
		return nil, ErrNotFound
	}

	react := reply.Delegate{
		Head: *delegateRef.CoreRef(),
	}

	return &react, nil
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
	codeID, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}

	return &reply.Reference{Ref: *getReference(&msg.Request, codeID)}, nil
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
	codeRef, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return &reply.Reference{Ref: *getReference(&msg.Request, codeRef)}, nil
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
		idx, _, _, err := getClass(tx, &msg.Class, nil)
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

	return &reply.Reference{Ref: *getReference(&msg.Request, deactivationID)}, nil
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
		idx, _, _, err := getClass(tx, &msg.Class, nil)
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

	return &reply.Reference{Ref: *getReference(&msg.Request, amendID)}, nil
}

func (h *MessageHandler) handleActivateObject(genericMsg core.Message) (core.Reply, error) {
	msg := genericMsg.(*message.ActivateObject)

	domainRef := record.Core2Reference(msg.Domain)
	requestRef := record.Core2Reference(msg.Request)
	classRef := record.Core2Reference(msg.Class)
	parentRef := record.Core2Reference(msg.Parent)

	var err error
	_, _, _, err = getClass(h.db, &msg.Class, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &msg.Parent, nil)
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
		parentIdx.Children = append(parentIdx.Children, record.Reference{
			Record: *objID,
			Domain: record.Core2Reference(msg.Request).Domain,
		})
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
	_, _, _, err = getClass(h.db, &msg.Class, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &msg.Parent, nil)
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
		idx, _, _, err := getObject(tx, &msg.Object, nil)
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

	return &reply.Reference{Ref: *getReference(&msg.Request, deactivationID)}, nil
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
		idx, _, _, err := getObject(tx, &msg.Object, nil)
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

	return &reply.Reference{Ref: *getReference(&msg.Request, amendID)}, nil
}

func getReference(request *core.RecordRef, id *record.ID) *core.RecordRef {
	ref := record.Reference{
		Record: *id,
		Domain: record.Core2Reference(*request).Domain,
	}
	return ref.CoreRef()
}

func getClass(
	s storage.Store, head *core.RecordRef, state *core.RecordRef,
) (*index.ClassLifeline, *core.RecordRef, record.ClassState, error) {
	headRef := record.Core2Reference(*head)
	idx, err := s.GetClassIndex(&headRef.Record)
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

	return idx, getReference(head, &stateID), stateRec, nil
}

func getObject(
	s storage.Store, head *core.RecordRef, state *core.RecordRef,
) (*index.ObjectLifeline, *core.RecordRef, record.ObjectState, error) {
	headRef := record.Core2Reference(*head)
	idx, err := s.GetObjectIndex(&headRef.Record)
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

	return idx, getReference(head, &stateID), stateRec, nil
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
