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

	"github.com/insolar/insolar/insolar/message"

	"github.com/insolar/insolar/insolar/flow/bus"
)

/*
func (h *MessageHandler) handleRegisterChild(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	logger := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.RegisterChild)
	jetID := jetFromContext(ctx)
	r, err := object.DecodeVirtual(msg.Record)
	if err != nil {
		return nil, errors.Wrap(err, "can't deserialize record")
	}
	childRec, ok := r.(*object.ChildRecord)
	if !ok {
		return nil, errors.New("wrong child record")
	}

	h.IDLocker.Lock(msg.Parent.Record())
	defer h.IDLocker.Unlock(msg.Parent.Record())

	var child *insolar.ID
	idx, err := h.IndexStorage.ForID(ctx, *msg.Parent.Record())
	if err == object.ErrIndexNotFound {
		heavy, err := h.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = h.saveIndexFromHeavy(ctx, jetID, msg.Parent, heavy)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch index from heavy")
		}
	} else if err != nil {
		return nil, err
	}
	h.IndexStateModifier.SetUsageForPulse(ctx, *msg.Parent.Record(), parcel.Pulse())

	recID := object.NewRecordIDFromRecord(h.PlatformCryptographyScheme, parcel.Pulse(), childRec)

	// Children exist and pointer does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if idx.ChildPointer != nil && !childRec.PrevChild.Equal(*idx.ChildPointer) && idx.ChildPointer != recID {
		return nil, errors.New("invalid child record")
	}

	child = object.NewRecordIDFromRecord(h.PlatformCryptographyScheme, parcel.Pulse(), childRec)
	rec := record.MaterialRecord{
		Record: childRec,
		JetID:  insolar.JetID(jetID),
	}

	err = h.RecordModifier.Set(ctx, *child, rec)

	if err == object.ErrOverride {
		logger.WithField("type", fmt.Sprintf("%T", r)).Warn("set record override (#2)")
		child = recID
	} else if err != nil {
		return nil, errors.Wrap(err, "can't save record into storage")
	}

	idx.ChildPointer = child
	if msg.AsType != nil {
		idx.Delegates[*msg.AsType] = msg.Child
	}
	idx.LatestUpdate = parcel.Pulse()
	idx.JetID = insolar.JetID(jetID)
	err = h.IndexStorage.Set(ctx, *msg.Parent.Record(), idx)
	if err != nil {
		return nil, err
	}

	return &reply.ID{ID: *child}, nil
}
*/

type RegisterChild struct {
	replyTo chan<- bus.Reply
	message *message.RegisterChild

	Dep struct {
		// RecordAccessor object.RecordAccessor
	}
}

func NewRegisterChild(msg *message.RegisterChild, replyTo chan<- bus.Reply) *RegisterChild {
	return &RegisterChild{
		message: msg,
		replyTo: replyTo,
	}
}

func (p *RegisterChild) Proceed(ctx context.Context) error {
	return nil
	// TODO TODO actually implement
	/*	jetID := p.Jet
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
			actualJet, err := p.Dep.TreeUpdater.Fetch(ctx, *msg.Parent.Record(), currentChild.Pulse())
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
		return nil */
}
