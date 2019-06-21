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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/internal/ledger/artifact.Manager -o ./ -s _gen_mock.go

// Manager implements methods required for direct ledger access.
type Manager interface {
	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (w/o deactivation check).
	GetObject(ctx context.Context, head insolar.Reference) (ObjectDescriptor, error)

	// RegisterRequest creates request record in storage.
	RegisterRequest(ctx context.Context, req record.Request) (*insolar.ID, error)

	// RegisterResult saves payload result in storage (emulates of save call result by VM).
	RegisterResult(ctx context.Context, obj, request insolar.Reference, payload []byte) (*insolar.ID, error)

	// ActivateObject creates activate object record in storage.
	// If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivateObject(
		ctx context.Context,
		domain, obj, parent, prototype insolar.Reference,
		asDelegate bool,
		memory []byte,
	) (ObjectDescriptor, error)

	// ActivatePrototype creates activate object record in storage.
	// Provided prototype reference will be used as objects prototype memory as memory of created object.
	ActivatePrototype(
		ctx context.Context,
		domain, obj, parent, code insolar.Reference,
		memory []byte,
	) (ObjectDescriptor, error)

	// UpdateObject creates amend object record in storage.
	// Provided reference should be a reference to the head of the object.
	// Provided memory well be the new object memory.
	//
	// Returned descriptor is the latest object state (exact) reference.
	UpdateObject(ctx context.Context, domain, request insolar.Reference, obj ObjectDescriptor, memory []byte) (ObjectDescriptor, error)

	// DeployCode creates new code record in storage (code records are used to activate prototypes).
	DeployCode(
		ctx context.Context,
		domain insolar.Reference,
		request insolar.Reference,
		code []byte,
		machineType insolar.MachineType,
	) (*insolar.ID, error)
}

// Scope implements Manager interface.
type Scope struct {
	PulseNumber insolar.PulseNumber

	PCS insolar.PlatformCryptographyScheme

	BlobStorage blob.Storage

	RecordModifier object.RecordModifier
	RecordAccessor object.RecordAccessor

	IndexAccessor object.IndexAccessor
	IndexModifier object.IndexModifier

	LifelineModifier object.LifelineModifier
	LifelineAccessor object.LifelineAccessor
}

// GetObject returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (w/o deactivation check).
func (m *Scope) GetObject(
	ctx context.Context,
	head insolar.Reference,
) (ObjectDescriptor, error) {

	idx, err := m.IndexAccessor.ForID(nil, m.PulseNumber, *head.Record())
	if err != nil {
		return nil, err
	}

	rec, err := m.RecordAccessor.ForID(ctx, *idx.Lifeline.LatestState)
	if err != nil {
		return nil, err
	}

	concrete := record.Unwrap(rec.Virtual)
	state, ok := concrete.(record.State)
	if !ok {
		return nil, errors.New("invalid object record")
	}

	desc := &objectDescriptor{
		head:         head,
		state:        *idx.Lifeline.LatestState,
		prototype:    state.GetImage(),
		isPrototype:  state.GetIsPrototype(),
		childPointer: idx.Lifeline.ChildPointer,
		parent:       idx.Lifeline.Parent,
	}
	if state.GetMemory() != nil {
		b, err := m.BlobStorage.ForID(ctx, *state.GetMemory())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch blob")
		}
		desc.memory = b.Value
	}
	return desc, nil
}

// RegisterRequest creates request record in storage.
func (m *Scope) RegisterRequest(ctx context.Context, req record.Request) (*insolar.ID, error) {
	virtRec := record.Wrap(req)
	return m.setRecord(ctx, virtRec)
}

// RegisterResult saves payload result in storage (emulates of save call result by VM).
func (m *Scope) RegisterResult(
	ctx context.Context, obj, request insolar.Reference, payload []byte,
) (*insolar.ID, error) {
	res := record.Result{
		Object:  *obj.Record(),
		Request: request,
		Payload: payload,
	}
	virtRec := record.Wrap(res)

	return m.setRecord(ctx, virtRec)
}

// ActivateObject creates activate object record in storage.
// If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *Scope) ActivateObject(
	ctx context.Context,
	domain, obj, parent, prototype insolar.Reference,
	asDelegate bool,
	memory []byte,
) (ObjectDescriptor, error) {
	return m.activateObject(ctx, domain, obj, prototype, false, parent, asDelegate, memory)
}

// ActivatePrototype creates activate object record in storage.
// Provided prototype reference will be used as objects prototype memory as memory of created object.
func (m *Scope) ActivatePrototype(
	ctx context.Context,
	domain, obj, parent, code insolar.Reference,
	memory []byte,
) (ObjectDescriptor, error) {
	return m.activateObject(ctx, domain, obj, code, true, parent, false, memory)
}

func (m *Scope) activateObject(
	ctx context.Context,
	domain insolar.Reference,
	obj insolar.Reference,
	prototype insolar.Reference,
	isPrototype bool,
	parent insolar.Reference,
	asDelegate bool,
	memory []byte,
) (ObjectDescriptor, error) {
	parentIdx, err := m.IndexAccessor.ForID(ctx, m.PulseNumber, *parent.Record())
	if err != nil {
		return nil, errors.Wrapf(err, "not found parent index for activated object: %v", parent.String())
	}

	stateRecord := record.Activate{
		Domain:      domain,
		Request:     obj,
		Image:       prototype,
		IsPrototype: isPrototype,
		Parent:      parent,
		IsDelegate:  asDelegate,
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
		parentIdx.Lifeline.ChildPointer,
		asType,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate")
	}

	return stateObj, nil
}

// UpdateObject creates amend object record in storage.
// Provided reference should be a reference to the head of the object.
// Provided memory well be the new object memory.
//
// Returned descriptor is the latest object state (exact) reference.
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

	amendRecord := record.Amend{
		Domain:      domain,
		Request:     request,
		Image:       *image,
		IsPrototype: objDesc.IsPrototype(),
		PrevState:   *objDesc.StateID(),
	}

	return m.updateStateObject(ctx, *objDesc.HeadRef(), amendRecord, memory)
}

// DeployCode creates new code record in storage (code records are used to activate prototypes).
func (m *Scope) DeployCode(
	ctx context.Context,
	domain insolar.Reference,
	request insolar.Reference,
	code []byte,
	machineType insolar.MachineType,
) (*insolar.ID, error) {
	codeRec := record.Code{
		Domain:      domain,
		Request:     request,
		Code:        code,
		MachineType: machineType,
	}

	return m.setRecord(
		ctx,
		record.Wrap(codeRec),
	)
}

func (m *Scope) setRecord(ctx context.Context, rec record.Virtual) (*insolar.ID, error) {
	hash := record.HashVirtual(m.PCS.ReferenceHasher(), rec)
	id := insolar.NewID(m.PulseNumber, hash)

	matRec := record.Material{
		Virtual: &rec,
		JetID:   insolar.ZeroJetID,
	}
	return id, m.RecordModifier.Set(ctx, *id, matRec)
}

func (m *Scope) setBlob(ctx context.Context, memory []byte) (*insolar.ID, error) {
	blobID := object.CalculateIDForBlob(m.PCS, m.PulseNumber, memory)
	err := m.BlobStorage.Set(
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
	idx, err := m.IndexAccessor.ForID(ctx, m.PulseNumber, *parent.Record())
	if err != nil {
		return err
	}

	m.IndexModifier.UpdateIndex(ctx, m.PulseNumber, *parent.Record(), func(updIndex object.FilamentIndex) (index object.FilamentIndex, e error) {

	})

	childRec := record.Child{Ref: obj}
	if prevChild != nil && prevChild.NotEmpty() {
		childRec.PrevChild = *prevChild
	}

	hash := record.HashVirtual(m.PCS.ReferenceHasher(), record.Wrap(childRec))
	recID := insolar.NewID(m.PulseNumber, hash)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if idx.Lifeline.ChildPointer != nil && !childRec.PrevChild.Equal(*idx.Lifeline.ChildPointer) && idx.Lifeline.ChildPointer != recID {
		return errors.New("invalid child record")
	}

	child, err := m.setRecord(ctx, record.Wrap(childRec))
	if err != nil {
		return err
	}

	idx.Lifeline.ChildPointer = child
	if asType != nil {
		idx.Lifeline.SetDelegate(*asType, obj)
	}
	idx.Lifeline.LatestUpdate = m.PulseNumber
	idx.JetID = insolar.JetID(jetID)
	return m.LifelineModifier.Set(ctx, m.PulseNumber, *parent.Record(), idx)
}

func (m *Scope) updateStateObject(
	ctx context.Context,
	objRef insolar.Reference,
	stateObject record.State,
	memory []byte,
) (ObjectDescriptor, error) {
	var jetID = insolar.ID(insolar.ZeroJetID)
	blobID, err := m.setBlob(ctx, memory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update blob")
	}

	var virtRecord record.Virtual
	switch so := stateObject.(type) {
	case record.Activate:
		so.Memory = *blobID
		virtRecord = record.Wrap(so)
	case record.Amend:
		so.Memory = *blobID
		virtRecord = record.Wrap(so)
	default:
		panic("unknown state object type")
	}

	idx, err := m.LifelineAccessor.ForID(ctx, m.PulseNumber, *objRef.Record())
	// No index on our node.
	if err != nil {
		if err != object.ErrLifelineNotFound {
			return nil, errors.Wrap(err, "failed get index for updating state object")
		}
		if stateObject.ID() != record.StateActivation {
			return nil, errors.Wrap(err, "index not found for updating non Activation state object")
		}
		// We are activating the object. There is no index for it yet.
		idx = object.Lifeline{StateID: record.StateUndefined}
	}

	id, err := m.setRecord(ctx, virtRecord)
	if err != nil {
		return nil, errors.Wrap(err, "fail set record for state object")
	}

	// update index
	idx.StateID = stateObject.ID()
	idx.LatestState = id
	idx.LatestUpdate = m.PulseNumber
	if stateObject.ID() == record.StateActivation {
		idx.Parent = stateObject.(record.Activate).Parent
	}
	idx.JetID = insolar.JetID(jetID)
	err = m.LifelineModifier.Set(ctx, m.PulseNumber, *objRef.Record(), idx)
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
