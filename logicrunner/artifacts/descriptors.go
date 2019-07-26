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

package artifacts

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

func NewCodeDescriptor(code []byte, machineType insolar.MachineType, ref insolar.Reference) CodeDescriptor {
	return &codeDescriptor{
		code:        code,
		machineType: machineType,
		ref:         ref,
	}
}

// CodeDescriptor represents meta info required to fetch all code data.
type codeDescriptor struct {
	code        []byte
	machineType insolar.MachineType
	ref         insolar.Reference
}

// Ref returns reference to represented code record.
func (d *codeDescriptor) Ref() *insolar.Reference {
	return &d.ref
}

// MachineType returns code machine type for represented code.
func (d *codeDescriptor) MachineType() insolar.MachineType {
	return d.machineType
}

// Code returns code data.
func (d *codeDescriptor) Code() ([]byte, error) {
	return d.code, nil
}

func NewObjectDescriptor(head insolar.Reference, state insolar.ID, prototype *insolar.Reference, isPrototype bool,
	childPointer *insolar.ID, memory []byte, parent insolar.Reference) ObjectDescriptor {

	return &objectDescriptor{
		head:         head,
		state:        state,
		prototype:    prototype,
		isPrototype:  isPrototype,
		childPointer: childPointer,
		memory:       memory,
		parent:       parent,
	}
}

// ObjectDescriptor represents meta info required to fetch all object data.
type objectDescriptor struct {
	head         insolar.Reference
	state        insolar.ID
	prototype    *insolar.Reference
	isPrototype  bool
	childPointer *insolar.ID // can be nil.
	memory       []byte
	parent       insolar.Reference
}

// IsPrototype determines if the object is a prototype.
func (d *objectDescriptor) IsPrototype() bool {
	return d.isPrototype
}

// Code returns code reference.
func (d *objectDescriptor) Code() (*insolar.Reference, error) {
	if !d.IsPrototype() {
		return nil, errors.New("object is not a prototype")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no code")
	}
	return d.prototype, nil
}

// Prototype returns prototype reference.
func (d *objectDescriptor) Prototype() (*insolar.Reference, error) {
	if d.IsPrototype() {
		return nil, errors.New("object is not an instance")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no prototype")
	}
	return d.prototype, nil
}

// HeadRef returns reference to represented object record.
func (d *objectDescriptor) HeadRef() *insolar.Reference {
	return &d.head
}

// StateID returns reference to object state record.
func (d *objectDescriptor) StateID() *insolar.ID {
	return &d.state
}

// ChildPointer returns the latest child for this object.
func (d *objectDescriptor) ChildPointer() *insolar.ID {
	return d.childPointer
}

// Memory fetches latest memory of the object known to storage.
func (d *objectDescriptor) Memory() []byte {
	return d.memory
}

// Parent returns object's parent.
func (d *objectDescriptor) Parent() *insolar.Reference {
	return &d.parent
}
