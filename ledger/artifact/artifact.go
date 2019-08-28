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
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/artifact.Manager -o ./ -s _gen_mock.go -g

// Manager implements methods required for direct ledger access.
type Manager interface {
	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (w/o deactivation check).
	GetObject(ctx context.Context, head insolar.Reference) (ObjectDescriptor, error)

	// RegisterRequest creates request record in storage.
	RegisterRequest(ctx context.Context, req record.IncomingRequest) (*insolar.ID, error)

	// RegisterResult saves payload result in storage (emulates of save call result by VM).
	RegisterResult(ctx context.Context, obj, request insolar.Reference, payload []byte) (*insolar.ID, error)

	// ActivateObject creates activate object record in storage.
	// If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivateObject(
		ctx context.Context,
		domain, obj, parent, prototype insolar.Reference,
		memory []byte,
	) error

	// UpdateObject creates amend object record in storage.
	// Provided reference should be a reference to the head of the object.
	// Provided memory well be the new object memory.
	//
	// Returned descriptor is the latest object state (exact) reference.
	UpdateObject(ctx context.Context, domain, request insolar.Reference, obj ObjectDescriptor, memory []byte) error

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

	RecordModifier object.RecordModifier
	RecordAccessor object.RecordAccessor

	IndexAccessor object.IndexAccessor
	IndexModifier object.IndexModifier
}

// GetObject returns descriptor for provided state.
//
// If provided state is nil, the latest state will be returned (w/o deactivation check).
func (m *Scope) GetObject(
	ctx context.Context,
	head insolar.Reference,
) (ObjectDescriptor, error) {

	idx, err := m.IndexAccessor.ForID(ctx, m.PulseNumber, *head.GetLocal())
	if err != nil {
		return nil, err
	}

	rec, err := m.RecordAccessor.ForID(ctx, *idx.Lifeline.LatestState)
	if err != nil {
		return nil, err
	}

	concrete := record.Unwrap(&rec.Virtual)
	state, ok := concrete.(record.State)
	if !ok {
		return nil, errors.New("invalid object record")
	}

	desc := &objectDescriptor{
		head:        head,
		state:       *idx.Lifeline.LatestState,
		prototype:   state.GetImage(),
		isPrototype: state.GetIsPrototype(),
		parent:      idx.Lifeline.Parent,
	}
	if state.GetMemory() != nil {
		desc.memory = state.GetMemory()
	}
	return desc, nil
}

// RegisterRequest creates request record in storage.
func (m *Scope) RegisterRequest(ctx context.Context, req record.IncomingRequest) (*insolar.ID, error) {
	virtRec := record.Wrap(&req)
	return m.setRecord(ctx, virtRec)
}

// RegisterResult saves payload result in storage (emulates of save call result by VM).
func (m *Scope) RegisterResult(
	ctx context.Context, obj, request insolar.Reference, payload []byte,
) (*insolar.ID, error) {
	res := record.Result{
		Object:  *obj.GetLocal(),
		Request: request,
		Payload: payload,
	}
	virtRec := record.Wrap(&res)

	return m.setRecord(ctx, virtRec)
}

// ActivateObject creates activate object record in storage.
// If memory is not provided, the prototype default memory will be used.
//
// Request reference will be this object's identifier and referred as "object head".
func (m *Scope) ActivateObject(
	ctx context.Context,
	domain, obj, parent, prototype insolar.Reference,
	memory []byte,
) error {
	return m.activateObject(ctx, domain, obj, prototype, false, parent, memory)
}

func (m *Scope) activateObject(
	ctx context.Context,
	domain insolar.Reference,
	obj insolar.Reference,
	prototype insolar.Reference,
	isPrototype bool,
	parent insolar.Reference,
	memory []byte,
) error {
	_, err := m.IndexAccessor.ForID(ctx, m.PulseNumber, *parent.GetLocal())
	if err != nil {
		return errors.Wrapf(err, "not found parent index for activated object: %v", parent.String())
	}

	stateRecord := record.Activate{
		Domain:      domain,
		Request:     obj,
		Image:       prototype,
		IsPrototype: isPrototype,
		Parent:      parent,
	}
	err = m.updateStateObject(ctx, obj, &stateRecord, memory)
	if err != nil {
		return errors.Wrap(err, "fail to store activation state")
	}
	return nil
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
) error {
	if objDesc.IsPrototype() {
		return errors.New("object is not an instance")
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
		return errors.Wrap(err, "failed to update object")
	}

	amendRecord := record.Amend{
		Domain:      domain,
		Request:     request,
		Image:       *image,
		IsPrototype: objDesc.IsPrototype(),
		PrevState:   *objDesc.StateID(),
	}

	return m.updateStateObject(ctx, *objDesc.HeadRef(), &amendRecord, memory)
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
		record.Wrap(&codeRec),
	)
}

func (m *Scope) setRecord(ctx context.Context, rec record.Virtual) (*insolar.ID, error) {
	hash := record.HashVirtual(m.PCS.ReferenceHasher(), rec)
	id := insolar.NewID(m.PulseNumber, hash)

	matRec := record.Material{
		Virtual: rec,
		JetID:   insolar.ZeroJetID,
		ID:      *id,
	}
	return id, m.RecordModifier.Set(ctx, matRec)
}

func (m *Scope) updateStateObject(
	ctx context.Context,
	objRef insolar.Reference,
	stateObject record.State,
	memory []byte,
) error {
	var virtRecord record.Virtual

	switch so := stateObject.(type) {
	case *record.Activate:
		so.Memory = memory
		virtRecord = record.Wrap(so)
	case *record.Amend:
		so.Memory = memory
		virtRecord = record.Wrap(so)
	default:
		panic("unknown state object type")
	}

	idx, err := m.IndexAccessor.ForID(ctx, m.PulseNumber, *objRef.GetLocal())
	// No index on our node.
	if err != nil {
		if err != object.ErrIndexNotFound {
			return errors.Wrap(err, "failed get index for updating state object")
		}
		if stateObject.ID() != record.StateActivation {
			return errors.Wrap(err, "index not found for updating non Activation state object")
		}
		// We are activating the object. There is no index for it yet.
		idx = record.Index{
			Lifeline:       record.Lifeline{StateID: record.StateUndefined},
			PendingRecords: []insolar.ID{},
			ObjID:          *objRef.GetLocal(),
		}
	}

	id, err := m.setRecord(ctx, virtRecord)
	if err != nil {
		return errors.Wrap(err, "fail set record for state object")
	}

	// update index
	idx.Lifeline.StateID = stateObject.ID()
	idx.Lifeline.LatestState = id
	if stateObject.ID() == record.StateActivation {
		idx.Lifeline.Parent = stateObject.(*record.Activate).Parent
	}
	err = m.IndexModifier.SetIndex(ctx, m.PulseNumber, idx)
	if err != nil {
		return errors.Wrap(err, "fail set index for state object")
	}

	return nil
}
