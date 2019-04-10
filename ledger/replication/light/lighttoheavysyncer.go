/*
 *    Copyright 2019 Insolar Technologies
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

package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/utils/backoff"
)

type ToHeavySyncer interface {
	SyncPulse(pn insolar.PulseNumber) error
}

type Cleaner interface {
	Clean(ctx context.Context)
}

type DataGatherer interface {
	ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) (*message.HeavyPayload, error)
}

type toHeavySyncer struct {
	waitingQueueChecker *time.Ticker
	problemQueueChecker *time.Ticker

	syncWaitingPulses    pulseLinkedList
	sendingProblemPulses pulseLinkedList

	jetAccessor     jet.Accessor
	dropAccessor    drop.Accessor
	blobsAccessor   blob.CollectionAccessor
	recsAccessor    object.RecordCollectionAccessor
	indexesAccessor object.IndexCollectionAccessor

	cleaner Cleaner

	msgBus insolar.MessageBus
}

type pulseLinkedList struct {
	pulsesMux sync.Mutex
	head      *listNode
	tail      *listNode
}

type listNode struct {
	current insolar.PulseNumber
	jetID   insolar.JetID

	next *listNode

	backoff *backoff.Backoff
}

func (p *pulseLinkedList) addPulse(pn insolar.PulseNumber) {
	p.pulsesMux.Lock()
	defer p.pulsesMux.Unlock()

	if p.head == nil {
		p.head = &listNode{
			current: pn,
		}
		p.tail = p.head
		return
	}

	p.tail.next = &listNode{
		current: pn,
	}
	p.tail = p.tail.next
}

func (p *pulseLinkedList) enqueuePulseWithJet(pn insolar.PulseNumber, jetID insolar.JetID) {
	p.pulsesMux.Lock()
	defer p.pulsesMux.Unlock()

	if p.head == nil {
		p.head = &listNode{
			current: pn,
			jetID:   jetID,
		}
		p.tail = p.head
		return
	}

	p.tail.next = &listNode{
		current: pn,
		jetID:   jetID,
	}
	p.tail = p.tail.next
}

func (p *pulseLinkedList) dequeuePulse() (insolar.PulseNumber, bool) {
	p.pulsesMux.Lock()
	defer p.pulsesMux.Unlock()

	if p.head == nil {
		return 0, false
	}

	h := p.head
	p.head = p.head.next

	return h.current, true
}

func (t *toHeavySyncer) SyncPulse(ctx context.Context, pn insolar.PulseNumber) {
	t.syncWaitingPulses.enqueuePulse(pn)

}

func (t *toHeavySyncer) lazyInit(ctx context.Context) {
	if t.waitingQueueChecker == nil {
		t.waitingQueueChecker = time.NewTicker(500 * time.Millisecond)
		go func() {
			for range t.waitingQueueChecker.C {
				go t.sync(ctx)
			}
		}()
	}
}

func (t *toHeavySyncer) sync(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	pnForSync, ok := t.syncWaitingPulses.dequeuePulse()
	if !ok {
		logger.Infof("Sync queue is empty")
		return
	}

	jets := t.jetAccessor.All(ctx, pnForSync)
	for _, jID := range jets {
		msg, err := t.gatherForPnAndJet(ctx, pnForSync, jID)
		if err != nil {
			logger.Error(fmt.Sprintf("Problems with gather data for a pulse - %v and jet - %v", pnForSync, jID.DebugString()))
			continue
		}
		err = t.sendToHeavy(ctx, msg)
		if err != nil {
			logger.Errorf("Problems with sending msg to a heavy node", err)
			t.sendingProblemPulses.enqueuePulseWithJet(pnForSync, jID)
		}
	}

}

func (t *toHeavySyncer) gatherForPnAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) (*message.HeavyPayload, error) {
	dr, err := t.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		inslogger.FromContext(ctx).Error("synchronize: can't fetch a drop")
		return nil, err
	}

	bls := t.blobsAccessor.ForPulse(ctx, jetID, pn)
	records := t.recsAccessor.ForPulse(ctx, jetID, pn)

	indexes := t.indexesAccessor.ForPulseAndJet(ctx, jetID, pn)
	resIdx := map[insolar.ID][]byte{}
	for id, idx := range indexes {
		resIdx[id] = object.EncodeIndex(idx)
	}

	return &message.HeavyPayload{
		JetID:    jetID,
		PulseNum: pn,
		Indexes:  resIdx,
		Drop:     drop.MustEncode(&dr),
		Blobs:    convertBlobs(bls),
		Records:  convertRecords(records),
	}, nil
}

func (t *toHeavySyncer) sendToHeavy(ctx context.Context, data *message.HeavyPayload) error {
	rep, err := t.msgBus.Send(ctx, data, nil)
	if err != nil {
		return err
	}
	if rep != nil {
		err, ok := rep.(*reply.HeavyError)
		if ok {
			return err
		}
	}
	return nil
}

func convertBlobs(blobs []blob.Blob) [][]byte {
	var res [][]byte
	for _, b := range blobs {
		res = append(res, blob.MustEncode(&b))
	}
	return res
}

func convertRecords(records []record.MaterialRecord) (result [][]byte) {
	for _, r := range records {
		result = append(result, object.EncodeMaterial(r))
	}
	return
}
