package blockexplorer

import (
	"context"
	"errors"
	"github.com/insolar/insolar/core"
)

// ExplorerDescriptor represents meta info required to fetch all object data.
type ExplorerDescriptor struct {
	ctx context.Context
	be  *BlockExplorerManager

	head        core.RecordRef
	state       core.RecordID
	prototype   *core.RecordRef
	isPrototype bool
	nextPointer *core.RecordID
	memory      []byte
	parent      core.RecordRef
}

// IsPrototype determines if the object is a prototype.
func (d *ExplorerDescriptor) IsPrototype() bool {
	return d.isPrototype
}

// Code returns code reference.
func (d *ExplorerDescriptor) Code() (*core.RecordRef, error) {
	if !d.IsPrototype() {
		return nil, errors.New("object is not a prototype")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no code")
	}
	return d.prototype, nil
}

// Prototype returns prototype reference.
func (d *ExplorerDescriptor) Prototype() (*core.RecordRef, error) {
	if d.IsPrototype() {
		return nil, errors.New("object is not an instance")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no prototype")
	}
	return d.prototype, nil
}

// HeadRef returns reference to represented object record.
func (d *ExplorerDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateID returns reference to object state record.
func (d *ExplorerDescriptor) StateID() *core.RecordID {
	return &d.state
}

// NextPointer returns the next obj state for this object.
func (d *ExplorerDescriptor) NextPointer() *core.RecordID {
	return d.nextPointer
}

// Memory fetches latest memory of the object known to storage.
func (d *ExplorerDescriptor) Memory() []byte {
	return d.memory
}

// History returns object's history references.
func (d *ExplorerDescriptor) History(pulse *core.PulseNumber) (core.RefIterator, error) {
	return d.be.GetHistory(d.ctx, d.head, pulse)
}

// Parent returns object's parent.
func (d *ExplorerDescriptor) Parent() *core.RecordRef {
	return &d.parent
}
