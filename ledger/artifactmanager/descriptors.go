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
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor struct {
	code        []byte
	machineType core.MachineType
	ref         core.RecordRef

	ctx context.Context
	am  core.ArtifactManager
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
	return d.code, nil
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	ctx context.Context
	am  *LedgerArtifactManager

	head         core.RecordRef
	state        core.RecordID
	prototype    *core.RecordRef
	isPrototype  bool
	childPointer *core.RecordID // can be nil.
	memory       []byte
	parent       core.RecordRef
}

// IsPrototype determines if the object is a prototype.
func (d *ObjectDescriptor) IsPrototype() bool {
	return d.isPrototype
}

// Code returns code reference.
func (d *ObjectDescriptor) Code() (*core.RecordRef, error) {
	if !d.IsPrototype() {
		return nil, errors.New("object is not a prototype")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no code")
	}
	return d.prototype, nil
}

// Prototype returns prototype reference.
func (d *ObjectDescriptor) Prototype() (*core.RecordRef, error) {
	if d.IsPrototype() {
		return nil, errors.New("object is not an instance")
	}
	if d.prototype == nil {
		return nil, errors.New("object has no prototype")
	}
	return d.prototype, nil
}

// HeadRef returns reference to represented object record.
func (d *ObjectDescriptor) HeadRef() *core.RecordRef {
	return &d.head
}

// StateID returns reference to object state record.
func (d *ObjectDescriptor) StateID() *core.RecordID {
	return &d.state
}

// ChildPointer returns the latest child for this object.
func (d *ObjectDescriptor) ChildPointer() *core.RecordID {
	return d.childPointer
}

// Memory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) Memory() []byte {
	return d.memory
}

// Children returns object's children references.
func (d *ObjectDescriptor) Children(pulse *core.PulseNumber) (core.RefIterator, error) {
	return d.am.GetChildren(d.ctx, d.head, pulse)
}

// Parent returns object's parent.
func (d *ObjectDescriptor) Parent() *core.RecordRef {
	return &d.parent
}

// HasPendingRequests returns true if the object has unclosed requests.
func (d *ObjectDescriptor) HasPendingRequests() bool {
	return false
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
	ctx          context.Context
	messageBus   core.MessageBus
	currentPulse core.Pulse
	parent       core.RecordRef
	chunkSize    int
	fromPulse    *core.PulseNumber
	fromChild    *core.RecordID
	buff         []core.RecordRef
	buffIndex    int
	canFetch     bool
}

// NewChildIterator creates new child iterator.
func NewChildIterator(
	ctx context.Context,
	mb core.MessageBus,
	parent core.RecordRef,
	fromPulse *core.PulseNumber,
	chunkSize int,
	currentPulse core.Pulse,
) (*ChildIterator, error) {
	iter := ChildIterator{
		ctx:          ctx,
		messageBus:   mb,
		parent:       parent,
		fromPulse:    fromPulse,
		chunkSize:    chunkSize,
		canFetch:     true,
		currentPulse: currentPulse,
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
	genericReply, err := i.sendAndFollowRedirect(
		i.ctx,
		&message.GetChildren{
			Parent:    i.parent,
			FromPulse: i.fromPulse,
			FromChild: i.fromChild,
			Amount:    i.chunkSize,
		},
		i.currentPulse,
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

func (i *ChildIterator) sendAndFollowRedirect(ctx context.Context, msg core.Message, currentPulse core.Pulse) (core.Reply, error) {
	rep, err := i.messageBus.Send(ctx, msg, currentPulse, nil)
	if err != nil {
		return nil, err
	}

	if redirect, ok := rep.(core.RedirectReply); ok {
		redirected := redirect.Redirected(msg)
		rep, err = i.messageBus.Send(
			ctx,
			redirected,
			currentPulse,
			&core.MessageSendOptions{
				Token:    redirect.GetToken(),
				Receiver: redirect.GetReceiver(),
			},
		)
		if err != nil {
			return nil, err
		}
		if _, ok = rep.(core.RedirectReply); ok {
			return nil, errors.New("double redirects are forbidden")
		}
		return rep, nil
	}

	return rep, err
}
