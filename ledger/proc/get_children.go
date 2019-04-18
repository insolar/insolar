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

package proc

import (
	"context"

	"github.com/insolar/insolar/ledger/artifactmanager"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
)

// GetChildren Procedure
type GetChildren struct {
	Jet     insolar.ID
	Message bus.Message

	Dep struct {
		RecentStorageProvider  recentstorage.Provider
		IDLocker               storage.IDLocker
		IndexAccessor          object.IndexAccessor
		JetCoordinator         insolar.JetCoordinator
		JetStorage             jet.Storage
		DelegationTokenFactory insolar.DelegationTokenFactory
		RecordAccessor         object.RecordAccessor
		TreeUpdater            jet.TreeUpdater
		IndexSaver             artifactmanager.IndexSaver
	}

	Result struct {
		Reply insolar.Reply
	}
}

func (p *GetChildren) Proceed(ctx context.Context) error {
	jetID := p.Jet
	parcel := p.Message.Parcel
	msg := parcel.Message().(*message.GetChildren)

	p.Dep.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Parent.Record())

	p.Dep.IDLocker.Lock(msg.Parent.Record())
	defer p.Dep.IDLocker.Unlock(msg.Parent.Record())

	idx, err := p.Dep.IndexAccessor.ForID(ctx, *msg.Parent.Record())
	if err == object.ErrIndexNotFound {
		heavy, err := p.Dep.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return err
		}
		idx, err = p.Dep.IndexSaver.SaveIndexFromHeavy(ctx, jetID, msg.Parent, heavy)
		if err != nil {
			return errors.Wrap(err, "failed to fetch index from heavy")
		}
		if idx.ChildPointer == nil {
			p.Result.Reply = &reply.Children{Refs: nil, NextFrom: nil}
			return nil
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to fetch object index")
	}

	var (
		refs         []insolar.Reference
		currentChild *insolar.ID
	)

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		currentChild = msg.FromChild
	} else {
		currentChild = idx.ChildPointer
	}

	// The object has no children.
	if currentChild == nil {
		p.Result.Reply = &reply.Children{Refs: nil, NextFrom: nil}
		return nil
	}

	var childJet *insolar.ID
	onHeavy, err := p.Dep.JetCoordinator.IsBeyondLimit(ctx, parcel.Pulse(), currentChild.Pulse())
	if err != nil && err != pulse.ErrNotFound {
		return err
	}
	if onHeavy {
		node, err := p.Dep.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return err
		}
		p.Result.Reply, err = reply.NewGetChildrenRedirect(p.Dep.DelegationTokenFactory, parcel, node, *currentChild)
		return err
	}

	childJetID, actual := p.Dep.JetStorage.ForID(ctx, currentChild.Pulse(), *msg.Parent.Record())
	childJet = (*insolar.ID)(&childJetID)

	if !actual {
		actualJet, err := p.Dep.TreeUpdater.FetchJet(ctx, *msg.Parent.Record(), currentChild.Pulse())
		if err != nil {
			return err
		}
		childJet = actualJet
	}

	// Try to fetch the first child.
	_, err = p.Dep.RecordAccessor.ForID(ctx, *currentChild)

	if err == object.ErrNotFound {
		node, err := p.Dep.JetCoordinator.NodeForJet(ctx, *childJet, parcel.Pulse(), currentChild.Pulse())
		if err != nil {
			return err
		}
		p.Result.Reply, err = reply.NewGetChildrenRedirect(p.Dep.DelegationTokenFactory, parcel, node, *currentChild)
		return err
	}

	if err != nil {
		return errors.Wrap(err, "failed to fetch child")
	}

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			p.Result.Reply = &reply.Children{Refs: refs, NextFrom: currentChild}
			return nil
		}
		counter++

		rec, err := p.Dep.RecordAccessor.ForID(ctx, *currentChild)

		// We don't have this child reference. Return what was collected.
		if err == object.ErrNotFound {
			p.Result.Reply = &reply.Children{Refs: refs, NextFrom: currentChild}
			return nil
		}
		if err != nil {
			return errors.New("failed to retrieve children")
		}

		virtRec := rec.Record
		childRec, ok := virtRec.(*object.ChildRecord)
		if !ok {
			return errors.New("failed to retrieve children")
		}
		currentChild = childRec.PrevChild

		// Skip records later than specified pulse.
		recPulse := childRec.Ref.Record().Pulse()
		if msg.FromPulse != nil && recPulse > *msg.FromPulse {
			continue
		}
		refs = append(refs, childRec.Ref)
	}

	p.Result.Reply = &reply.Children{Refs: refs, NextFrom: nil}
	return nil
}
