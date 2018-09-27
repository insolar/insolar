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
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// EventHandler processes events for local storage interaction.
type EventHandler struct {
	db *storage.DB
}

// NewEventHandler creates new handler.
func NewEventHandler(db *storage.DB) (*EventHandler, error) {
	return &EventHandler{db: db}, nil
}

// Handle performs event processing.
func (h *EventHandler) Handle(genericEvent core.Event) (core.Reaction, error) {
	switch e := genericEvent.(type) {
	case *event.GetCode:
		return h.handleGetCode(e)
	case *event.GetClass:
		return h.handleGetClass(e)
	case *event.GetObject:
		return h.handleGetObject(e)
	case *event.GetDelegate:
		return h.handleGetDelegate(e)
	case *event.DeclareType:
		return h.handleDeclareType(e)
	case *event.DeployCode:
		return h.handleDeployCode(e)
	case *event.ActivateClass:
		return h.handleActivateClass(e)
	case *event.DeactivateClass:
		return h.handleDeactivateClass(e)
	case *event.UpdateClass:
		return h.handleUpdateClass(e)
	case *event.ActivateObject:
		return h.handleActivateObject(e)
	case *event.ActivateObjectDelegate:
		return h.handleActivateObjectDelegate(e)
	case *event.DeactivateObject:
		return h.handleDeactivateObject(e)
	case *event.UpdateObject:
		return h.handleUpdateObject(e)
	}

	return nil, errors.New("no handler for this event")
}

func (h *EventHandler) handleGetCode(ev *event.GetCode) (core.Reaction, error) {
	codeRef := record.Core2Reference(ev.Code)
	rec, err := h.db.GetRecord(&codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}
	code, mt, err := codeRec.GetCode(ev.MachinePref)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code from record")
	}

	react := reaction.Code{
		Code:        code,
		MachineType: mt,
	}

	return &react, nil
}

func (h *EventHandler) handleGetClass(ev *event.GetClass) (core.Reaction, error) {
	_, stateRef, state, err := getClass(h.db, &ev.Head, ev.State)
	if err != nil {
		return nil, err
	}

	var code *core.RecordRef
	if state.GetCode() == nil {
		code = nil
	} else {
		code = state.GetCode().CoreRef()
	}

	react := reaction.Class{
		Head:  ev.Head,
		State: *stateRef,
		Code:  code,
	}

	return &react, nil
}

func (h *EventHandler) handleGetObject(ev *event.GetObject) (core.Reaction, error) {
	idx, stateRef, state, err := getObject(h.db, &ev.Head, ev.State)
	if err != nil {
		return nil, err
	}

	children := make([]core.RecordRef, 0, len(idx.Children))
	for _, c := range idx.Children {
		children = append(children, *c.CoreRef())
	}

	react := reaction.Object{
		Head:     ev.Head,
		State:    *stateRef,
		Class:    *idx.ClassRef.CoreRef(),
		Memory:   state.GetMemory(),
		Children: children,
	}

	return &react, nil
}

func (h *EventHandler) handleGetDelegate(ev *event.GetDelegate) (core.Reaction, error) {
	idx, _, _, err := getObject(h.db, &ev.Head, nil)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := idx.Delegates[ev.AsClass]
	if !ok {
		return nil, ErrNotFound
	}

	react := reaction.Delegate{
		Head: *delegateRef.CoreRef(),
	}

	return &react, nil
}

func (h *EventHandler) handleDeclareType(ev *event.DeclareType) (core.Reaction, error) {
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)

	rec := record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TypeDeclaration: ev.TypeDec,
	}
	codeRef, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}

	return &reaction.Reference{Ref: *codeRef.CoreRef()}, nil
}

func (h *EventHandler) handleDeployCode(ev *event.DeployCode) (core.Reaction, error) {
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)

	rec := record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TargetedCode: ev.CodeMap,
	}
	codeRef, err := h.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return &reaction.Reference{Ref: *codeRef.CoreRef()}, nil
}

func (h *EventHandler) handleActivateClass(ev *event.ActivateClass) (core.Reaction, error) {
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)

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
	var classRef *record.Reference
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		classRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetClassIndex(classRef, &index.ClassLifeline{
			LatestStateRef: *classRef,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *classRef.CoreRef()}, nil
}

func (h *EventHandler) handleDeactivateClass(ev *event.DeactivateClass) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	classRef := record.Core2Reference(ev.Class)

	var deactivationRef *record.Reference
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getClass(tx, &ev.Class, nil)
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
				HeadRecord:    classRef,
				AmendedRecord: idx.LatestStateRef,
			},
		}

		deactivationRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestStateRef = *deactivationRef
		err = tx.SetClassIndex(&classRef, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *deactivationRef.CoreRef()}, nil
}

func (h *EventHandler) handleUpdateClass(ev *event.UpdateClass) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	classRef := record.Core2Reference(ev.Class)
	migrationRefs := make([]record.Reference, 0, len(ev.Class))
	for _, migration := range ev.Migrations {
		migrationRefs = append(migrationRefs, record.Core2Reference(migration))
	}

	err = validateCode(h.db, &ev.Code)
	if err != nil {
		return nil, err
	}
	for _, migration := range ev.Migrations {
		err = validateCode(h.db, &migration)
		if err != nil {
			return nil, err
		}
	}

	var amendRef *record.Reference
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getClass(tx, &ev.Class, nil)
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
				HeadRecord:    classRef,
				AmendedRecord: idx.LatestStateRef,
			},
			NewCode:    record.Core2Reference(ev.Code),
			Migrations: migrationRefs,
		}

		amendRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestStateRef = *amendRef
		idx.AmendRefs = append(idx.AmendRefs, *amendRef)
		err = tx.SetClassIndex(&classRef, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *amendRef.CoreRef()}, nil
}

func (h *EventHandler) handleActivateObject(ev *event.ActivateObject) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	classRef := record.Core2Reference(ev.Class)
	parentRef := record.Core2Reference(ev.Parent)

	_, _, _, err = getClass(h.db, &ev.Class, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &ev.Parent, nil)
	if err != nil {
		return nil, err
	}

	var objRef *record.Reference
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
			Memory:              ev.Memory,
			Parent:              parentRef,
			Delegate:            false,
		}

		// save new record and it's index
		objRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objRef, &index.ObjectLifeline{
			ClassRef:       classRef,
			LatestStateRef: *objRef,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's children
		parentIdx, err := tx.GetObjectIndex(&parentRef)
		if err != nil {
			if err == ErrNotFound {
				parentIdx = &index.ObjectLifeline{}
			} else {
				return errors.Wrap(err, "inconsistent index")
			}
		}
		parentIdx.Children = append(parentIdx.Children, *objRef)
		err = tx.SetObjectIndex(&parentRef, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &reaction.Reference{Ref: *objRef.CoreRef()}, nil
}

func (h *EventHandler) handleActivateObjectDelegate(ev *event.ActivateObjectDelegate) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	classRef := record.Core2Reference(ev.Class)
	parentRef := record.Core2Reference(ev.Parent)

	_, _, _, err = getClass(h.db, &ev.Class, nil)
	if err != nil {
		return nil, err
	}
	_, _, _, err = getObject(h.db, &ev.Parent, nil)
	if err != nil {
		return nil, err
	}

	var objRef *record.Reference
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
			Memory:              ev.Memory,
			Parent:              parentRef,
			Delegate:            true,
		}

		// save new record and it's index
		objRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objRef, &index.ObjectLifeline{
			ClassRef:       classRef,
			LatestStateRef: *objRef,
			Delegates:      map[core.RecordRef]record.Reference{},
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's delegates
		parentIdx, err := tx.GetObjectIndex(&parentRef)
		if err != nil {
			return errors.Wrap(err, "inconsistent index")
		}
		if _, ok := parentIdx.Delegates[ev.Class]; ok {
			return ErrClassDelegateAlreadyExists
		}
		parentIdx.Delegates[ev.Class] = *objRef
		err = tx.SetObjectIndex(&parentRef, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *objRef.CoreRef()}, nil
}

func (h *EventHandler) handleDeactivateObject(ev *event.DeactivateObject) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	objRef := record.Core2Reference(ev.Object)

	var deactivationRef *record.Reference
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &ev.Object, nil)
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
				HeadRecord:    objRef,
				AmendedRecord: idx.LatestStateRef,
			},
		}
		deactivationRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestStateRef = *deactivationRef
		err = tx.SetObjectIndex(&objRef, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *deactivationRef.CoreRef()}, nil
}

func (h *EventHandler) handleUpdateObject(ev *event.UpdateObject) (core.Reaction, error) {
	var err error
	domainRef := record.Core2Reference(ev.Domain)
	requestRef := record.Core2Reference(ev.Request)
	objRef := record.Core2Reference(ev.Object)

	var amendRef *record.Reference
	err = h.db.Update(func(tx *storage.TransactionManager) error {
		idx, _, _, err := getObject(tx, &ev.Object, nil)
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
				HeadRecord:    objRef,
				AmendedRecord: idx.LatestStateRef,
			},
			NewMemory: ev.Memory,
		}

		amendRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		idx.LatestStateRef = *amendRef
		err = tx.SetObjectIndex(&objRef, idx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reaction.Reference{Ref: *amendRef.CoreRef()}, nil
}

func getClass(
	s storage.Store, head *core.RecordRef, state *core.RecordRef,
) (*index.ClassLifeline, *core.RecordRef, record.ClassState, error) {
	headRef := record.Core2Reference(*head)
	idx, err := s.GetClassIndex(&headRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}

	var stateRef record.Reference
	if state != nil {
		stateRef = record.Core2Reference(*state)
	} else {
		stateRef = idx.LatestStateRef
	}

	rec, err := s.GetRecord(&stateRef)
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

	return idx, stateRef.CoreRef(), stateRec, nil
}

func getObject(
	s storage.Store, head *core.RecordRef, state *core.RecordRef,
) (*index.ObjectLifeline, *core.RecordRef, record.ObjectState, error) {
	headRef := record.Core2Reference(*head)
	idx, err := s.GetObjectIndex(&headRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}

	var stateRef record.Reference
	if state != nil {
		stateRef = record.Core2Reference(*state)
	} else {
		stateRef = idx.LatestStateRef
	}

	rec, err := s.GetRecord(&stateRef)
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

	return idx, stateRef.CoreRef(), stateRec, nil
}

func validateCode(s storage.Store, ref *core.RecordRef) error {
	codeRef := record.Core2Reference(*ref)
	rec, err := s.GetRecord(&codeRef)
	if err != nil {
		return err
	}

	if _, ok := rec.(*record.CodeRecord); !ok {
		return errors.New("invalid code reference")
	}

	return nil
}
