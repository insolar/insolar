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
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/messagebus"
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

// ChildIterator is used to iterate over objects children. During iteration children refs will be fetched from remote
// source (parent object).
//
// Data can be fetched only from Active Executor (AE), although children references can be stored on other nodes.
// To cope with this, we have a token system. Every time AE doesn't have data and asked for it, it will issue a token
// that will allow requester to fetch data from a different node. This node will return all children references it has,
// after which the requester has to go to AE again to fetch a new token. It will then be redirected to another node.
// E.i. children fetching happens like this:
// [R = requester, AE = active executor, LE = any light executor that has data, H = heavy executor]
// 1. R (get children 0 ... ) -> AE
// 2. AE (children 0 ... 3) -> R
// 3. R (get children 4 ...) -> AE
// 4. AE (redirect to LE) -> R
// 5. R (get children 4 ...) -> LE
// 6. LE (children 4 ... 5) -> R
// 7. R (get children 6 ...) -> AE
// 8. AE (redirect to H) -> R
// 9. R (get children 6 ...) -> H
// 10. H (children 6 ... 15 EOF) -> R
type ChildIterator struct {
	ctx         context.Context
	senderChain messagebus.Sender
	parent      insolar.Reference
	chunkSize   int
	fromPulse   *insolar.PulseNumber
	fromChild   *insolar.ID
	buff        []insolar.Reference
	buffIndex   int
	canFetch    bool
}

// NewChildIterator creates new child iterator.
func NewChildIterator(
	ctx context.Context,
	senderChain messagebus.Sender,
	parent insolar.Reference,
	fromPulse *insolar.PulseNumber,
	chunkSize int,
) (*ChildIterator, error) {
	iter := ChildIterator{
		ctx:         ctx,
		senderChain: senderChain,
		parent:      parent,
		fromPulse:   fromPulse,
		chunkSize:   chunkSize,
		canFetch:    true,
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
func (i *ChildIterator) Next() (*insolar.Reference, error) {
	// Get element from buffer.
	if !i.hasInBuffer() && i.canFetch {
		err := i.fetch()
		if err != nil {
			return nil, err
		}
	}

	ref := i.nextFromBuffer()
	if ref == nil {
		return nil, errors.New("failed to retrieve a child from buffer")
	}

	return ref, nil
}

func (i *ChildIterator) nextFromBuffer() *insolar.Reference {
	if !i.hasInBuffer() {
		return nil
	}
	ref := i.buff[i.buffIndex]
	i.buffIndex++
	return &ref
}

func (i *ChildIterator) fetch() error {
	if !i.canFetch {
		return errors.New("failed to fetch a children chunk")
	}

	genericReply, err := i.senderChain(i.ctx, &message.GetChildren{
		Parent:    i.parent,
		FromPulse: i.fromPulse,
		FromChild: i.fromChild,
		Amount:    i.chunkSize,
	}, nil)
	if err != nil {
		return err
	}
	rep, ok := genericReply.(*reply.Children)
	if !ok {
		return fmt.Errorf("unexpected reply: %#v", genericReply)
	}

	if rep.NextFrom == nil || rep.NextFrom.IsEmpty() {
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
