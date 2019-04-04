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
package artifact

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/pkg/errors"
)

type Contract struct {
	Name        string // ???
	Domain      insolar.Reference
	MachineType insolar.MachineType
	Binary      []byte
}

//go:generate minimock -i github.com/insolar/insolar/internal/ledger/artifact.Manager -o ./ -s _gen_mock.go

type Manager interface {
	RegisterRequest(ctx context.Context, objectRef insolar.Reference, parcel insolar.Parcel) (*insolar.ID, error)
	ActivateObject(ctx context.Context, domain, obj, parent, prototype insolar.Reference, asDelegate bool, memory []byte) (ObjectDescriptor, error)
	UpdateObject(ctx context.Context, domain, request insolar.Reference, obj ObjectDescriptor, memory []byte) (ObjectDescriptor, error)
}

type Scope struct {
	PulseNumber                insolar.PulseNumber
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`

	BlobModifier    blob.Modifier         `inject:""`
	RecordsModifier object.RecordModifier `inject:""`

	// ObjectStorage
	// Depricated, should be removed after indices storage would be done.
	ObjectStorage storage.ObjectStorage `inject:""`
}

// NewScope creates new scope instance.
func NewScope(pn insolar.PulseNumber) *Scope {
	return &Scope{
		PulseNumber: pn,
	}
}

func (m *Scope) RegisterRequest(ctx context.Context, objectRef insolar.Reference, parcel insolar.Parcel) (*insolar.ID, error) {
	rec := &object.RequestRecord{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: m.hashParcel(parcel),
		Object:      *objectRef.Record(),
	}
	return m.setRecord(ctx, rec)
}

func (m *Scope) ActivateObject(
	ctx context.Context,
	domain, obj, parent, prototype insolar.Reference,
	asDelegate bool,
	memory []byte,
) (ObjectDescriptor, error) {
	var jetID = insolar.ID(insolar.ZeroJetID)
	parentIdx, err := m.ObjectStorage.GetObjectIndex(ctx, jetID, parent.Record())
	if err != nil {
		return nil, errors.Wrap(err, "not found parent index for activated object")
	}

	stateRecord := &object.ActivateRecord{
		SideEffectRecord: object.SideEffectRecord{
			Domain:  domain,
			Request: obj,
		},
		StateRecord: object.StateRecord{
			Image:       prototype,
			IsPrototype: false,
		},
		Parent:     parent,
		IsDelegate: asDelegate,
	}
	stateObj, err := m.updateStateObject(ctx, obj, stateRecord, memory)
	if err != nil {
		return nil, errors.Wrap(err, "fail to store activation state")
	}

	asType := &prototype
	if !asDelegate {
		asType = nil
	}
	err = m.registerChild(
		ctx,
		obj,
		parent,
		parentIdx.ChildPointer,
		asType,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate")
	}

	return stateObj, nil
}

func (m *Scope) UpdateObject(
	ctx context.Context,
	domain, request insolar.Reference,
	objDesc ObjectDescriptor,
	memory []byte,
) (ObjectDescriptor, error) {
	if objDesc.IsPrototype() {
		return nil, errors.New("object is not an instance")
	}

	var (
		image *insolar.Reference
		err   error
	)
	if objDesc.IsPrototype() {
		image, err = objDesc.Code()
	} else {
		image, err = objDesc.Prototype()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to update object")
	}

	amendRecord := &object.AmendRecord{
		SideEffectRecord: object.SideEffectRecord{
			Domain:  domain,
			Request: request,
		},
		StateRecord: object.StateRecord{
			Image:       *image,
			IsPrototype: objDesc.IsPrototype(),
		},
		PrevState: *objDesc.StateID(),
	}
	return m.updateStateObject(ctx, *objDesc.HeadRef(), amendRecord, memory)
}

func (m *Scope) setRecord(ctx context.Context, rec record.VirtualRecord) (*insolar.ID, error) {
	id := object.NewRecordIDFromRecord(m.PlatformCryptographyScheme, m.PulseNumber, rec)
	matRec := record.MaterialRecord{
		Record: rec,
		JetID:  insolar.ZeroJetID,
	}
	return id, m.RecordsModifier.Set(ctx, *id, matRec)
}

func (m *Scope) updateBlob(ctx context.Context, memory []byte) (*insolar.ID, error) {
	blobID := object.CalculateIDForBlob(m.PlatformCryptographyScheme, m.PulseNumber, memory)
	err := m.BlobModifier.Set(
		ctx,
		*blobID,
		blob.Blob{
			JetID: insolar.ZeroJetID,
			Value: memory,
		},
	)
	if err != nil && err != blob.ErrOverride {
		return nil, err
	}
	return blobID, nil
}

func (m *Scope) registerChild(
	ctx context.Context,
	obj insolar.Reference,
	parent insolar.Reference,
	prevChild *insolar.ID,
	asType *insolar.Reference,
) error {
	var jetID = insolar.ID(insolar.ZeroJetID)
	idx, err := m.ObjectStorage.GetObjectIndex(ctx, jetID, parent.Record())
	if err != nil {
		return err
	}

	childRec := &object.ChildRecord{
		PrevChild: prevChild,
		Ref:       obj,
	}

	recID := object.NewRecordIDFromRecord(m.PlatformCryptographyScheme, m.PulseNumber, childRec)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if idx.ChildPointer != nil && !childRec.PrevChild.Equal(*idx.ChildPointer) && idx.ChildPointer != recID {
		return errors.New("invalid child record")
	}

	child, err := m.setRecord(ctx, childRec)
	if err != nil {
		return err
	}

	idx.ChildPointer = child
	if asType != nil {
		idx.Delegates[*asType] = obj
	}
	idx.LatestUpdate = m.PulseNumber
	return m.ObjectStorage.SetObjectIndex(ctx, jetID, parent.Record(), idx)
}

func (m *Scope) updateStateObject(
	ctx context.Context,
	objRef insolar.Reference,
	stateObject object.State,
	memory []byte,
) (ObjectDescriptor, error) {
	var jetID = insolar.ID(insolar.ZeroJetID)
	blobID, err := m.updateBlob(ctx, memory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update blob")
	}

	var virtRecord record.VirtualRecord
	switch so := stateObject.(type) {
	case *object.ActivateRecord:
		so.Memory = blobID
		virtRecord = so
	case *object.AmendRecord:
		virtRecord = so
		so.Memory = blobID
	default:
		panic("unknown state object type")
	}

	idx, err := m.ObjectStorage.GetObjectIndex(ctx, jetID, objRef.Record())
	// No index on our node.
	if err != nil {
		if err != insolar.ErrNotFound {
			return nil, errors.Wrap(err, "failed get index for updating state object")
		}
		if stateObject.ID() != object.StateActivation {
			return nil, errors.Wrap(err, "index not found for updating non Activation state object")
		}
		// We are activating the object. There is no index for it yet.
		idx = &object.Lifeline{State: object.StateUndefined}
	}
	// TODO: validateState

	// TODO: validate index consistency
	id, err := m.setRecord(ctx, virtRecord)
	if err != nil {
		return nil, errors.Wrap(err, "fail set record for state object")
	}

	// update index
	idx.State = stateObject.ID()
	idx.LatestState = id
	idx.LatestUpdate = m.PulseNumber
	if stateObject.ID() == object.StateActivation {
		idx.Parent = stateObject.(*object.ActivateRecord).Parent
	}

	err = m.ObjectStorage.SetObjectIndex(ctx, jetID, objRef.Record(), idx)
	if err != nil {
		return nil, errors.Wrap(err, "fail set index for state object")
	}

	return &objectDescriptor{
		head:         objRef,
		state:        *idx.LatestState,
		prototype:    stateObject.GetImage(),
		isPrototype:  stateObject.GetIsPrototype(),
		childPointer: idx.ChildPointer,
		parent:       idx.Parent,
	}, nil
}

func (m *Scope) hashParcel(parcel insolar.Parcel) []byte {
	return m.PlatformCryptographyScheme.IntegrityHasher().Hash(message.MustSerializeBytes(parcel.Message()))
}
