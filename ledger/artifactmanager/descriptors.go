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
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor struct {
	machineType core.MachineType
	code        []byte
	ref         core.RecordRef
	machinePref []core.MachineType
}

// Ref returns reference to represented code record.
func (d *CodeDescriptor) Ref() *core.RecordRef {
	return &d.ref
}

// MachineType returns first available machine type for provided machine preference.
func (d *CodeDescriptor) MachineType() core.MachineType {
	return d.machineType
}

// Code returns code for first available machine type for provided machine preference.
func (d *CodeDescriptor) Code() []byte {
	return d.code
}

// ClassDescriptor represents meta info required to fetch all class data.
type ClassDescriptor struct {
	cache struct {
		codeDescriptor core.CodeDescriptor
	}

	am core.ArtifactManager

	head  core.RecordRef
	state core.RecordRef
	code  *core.RecordRef // Can be nil.
}

// HeadRef returns head reference to represented class record.
func (d *ClassDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateRef returns reference to represented class state record.
func (d *ClassDescriptor) StateRef() *core.RecordRef {
	return &d.state
}

// CodeDescriptor returns descriptor for fetching object's code data.
func (d *ClassDescriptor) CodeDescriptor(machinePref []core.MachineType) (core.CodeDescriptor, error) {
	if d.cache.codeDescriptor != nil {
		return d.cache.codeDescriptor, nil
	}

	if d.code == nil {
		return nil, errors.New("class has no code")
	}

	return d.am.GetCode(*d.code, machinePref)
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	cache struct {
		classDescriptor core.ClassDescriptor
	}
	am core.ArtifactManager

	head     core.RecordRef
	state    core.RecordRef
	class    core.RecordRef
	memory   []byte
	children []core.RecordRef
}

// HeadRef returns reference to represented object record.
func (d *ObjectDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateRef returns reference to object state record.
func (d *ObjectDescriptor) StateRef() *core.RecordRef {
	return &d.state
}

// Memory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) Memory() []byte {
	return d.memory
}

// Children returns object's children references.
func (d *ObjectDescriptor) Children() core.RefIterator {
	return &RefIterator{elements: d.children}
}

// ClassDescriptor returns descriptor for fetching object's class data.
func (d *ObjectDescriptor) ClassDescriptor(state *core.RecordRef) (core.ClassDescriptor, error) {
	if d.cache.classDescriptor != nil {
		return d.cache.classDescriptor, nil
	}

	return d.am.GetClass(d.class, state)
}

type RefIterator struct {
	elements     []core.RecordRef
	currentIndex int
}

func (i *RefIterator) HasNext() bool {
	return len(i.elements) > i.currentIndex
}

func (i *RefIterator) Next() (core.RecordRef, error) {
	el := i.elements[i.currentIndex]
	i.currentIndex++
	return el, nil
}
