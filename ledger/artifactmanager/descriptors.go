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
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/inscontext"
	"github.com/pkg/errors"
)

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor struct {
	cache struct {
		code []byte
	}
	machineType core.MachineType
	ref         core.RecordRef

	am core.ArtifactManager
}

// Ref returns reference to represented code record.
func (d *CodeDescriptor) Ref() *core.RecordRef {
	return &d.ref
}

// MachineType returns code machine type for represented code.
func (d *CodeDescriptor) MachineType() core.MachineType {
	return d.machineType
}

// Code returns code data.
func (d *CodeDescriptor) Code() ([]byte, error) {
	ctx := inscontext.TODO()
	if d.cache.code == nil {
		desc, err := d.am.GetCode(ctx, d.ref)
		if err != nil {
			return nil, err
		}
		code, err := desc.Code()
		if err != nil {
			return nil, err
		}
		d.cache.code = code
	}

	return d.cache.code, nil
}

// ClassDescriptor represents meta info required to fetch all class data.
type ClassDescriptor struct {
	cache struct {
		codeDescriptor core.CodeDescriptor
	}

	am core.ArtifactManager

	head        core.RecordRef
	state       core.RecordID
	code        *core.RecordRef // Can be nil.
	machineType core.MachineType
}

// HeadRef returns head reference to represented class record.
func (d *ClassDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateID returns reference to represented class state record.
func (d *ClassDescriptor) StateID() *core.RecordID {
	return &d.state
}

// CodeDescriptor returns descriptor for fetching object's code data.
func (d *ClassDescriptor) CodeDescriptor() core.CodeDescriptor {
	if d.cache.codeDescriptor == nil {
		d.cache.codeDescriptor = &CodeDescriptor{
			ref:         *d.code,
			machineType: d.machineType,
		}
	}

	return d.cache.codeDescriptor
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	cache struct {
		classDescriptor core.ClassDescriptor
	}
	am *LedgerArtifactManager

	head     core.RecordRef
	state    core.RecordID
	class    core.RecordRef
	memory   []byte
	children []core.RecordRef
}

// HeadRef returns reference to represented object record.
func (d *ObjectDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateID returns reference to object state record.
func (d *ObjectDescriptor) StateID() *core.RecordID {
	return &d.state
}

// Memory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) Memory() []byte {
	return d.memory
}

// Children returns object's children references.
func (d *ObjectDescriptor) Children(pulse *core.PulseNumber) (core.RefIterator, error) {
	ctx := inscontext.TODO()
	return d.am.GetChildren(ctx, d.head, pulse)
}

// ClassDescriptor returns descriptor for fetching object's class data.
func (d *ObjectDescriptor) ClassDescriptor(state *core.RecordRef) (core.ClassDescriptor, error) {
	ctx := inscontext.TODO()
	if d.cache.classDescriptor != nil {
		return d.cache.classDescriptor, nil
	}

	return d.am.GetClass(ctx, d.class, state)
}

// ChildIterator is used to iterate over objects children.
//
// During iteration children refs will be fetched from remote source (parent object).
type ChildIterator struct {
	messageBus core.MessageBus
	parent     core.RecordRef
	chunkSize  int
	fromPulse  *core.PulseNumber
	fromChild  *core.RecordID
	buff       []core.RecordRef
	buffIndex  int
	canFetch   bool
}

// NewChildIterator creates new child iterator.
func NewChildIterator(
	mb core.MessageBus, parent core.RecordRef, fromPulse *core.PulseNumber, chunkSize int,
) (*ChildIterator, error) {
	iter := ChildIterator{
		messageBus: mb,
		parent:     parent,
		fromPulse:  fromPulse,
		chunkSize:  chunkSize,
		canFetch:   true,
	}
	err := iter.fetch()
	if err != nil {
		return nil, err
	}
	return &iter, nil
}

// HasNext checks if any elements left in iterator.
func (i *ChildIterator) HasNext() bool {
	return i.hasInBuffer() || i.canFetch
}

// Next returns next element.
func (i *ChildIterator) Next() (*core.RecordRef, error) {
	// Get element from buffer.
	if !i.hasInBuffer() && i.canFetch {
		err := i.fetch()
		if err != nil {
			return nil, err
		}
	}

	ref := i.nextFromBuffer()
	if ref == nil {
		return nil, errors.New("failed to fetch record")
	}

	return ref, nil
}

func (i *ChildIterator) nextFromBuffer() *core.RecordRef {
	if !i.hasInBuffer() {
		return nil
	}
	ref := i.buff[i.buffIndex]
	i.buffIndex++
	return &ref
}

func (i *ChildIterator) fetch() error {
	if !i.canFetch {
		return errors.New("failed to fetch record")
	}
	genericReply, err := i.messageBus.Send(
		inscontext.TODO(),
		&message.GetChildren{
			Parent:    i.parent,
			FromPulse: i.fromPulse,
			FromChild: i.fromChild,
			Amount:    i.chunkSize,
		},
	)
	if err != nil {
		return err
	}
	rep, ok := genericReply.(*reply.Children)
	if !ok {
		return errors.New("failed to fetch record")
	}

	if rep.NextFrom == nil {
		i.canFetch = false
	}
	i.buff = rep.Refs
	i.buffIndex = 0
	i.fromChild = rep.NextFrom

	return nil
}

func (i *ChildIterator) hasInBuffer() bool {
	return i.buffIndex < len(i.buff)
}
