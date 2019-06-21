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

package object

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/insolar/insolar/insolar"
	buswm "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network/storage"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.PendingModifier -o ./ -s _mock.go

// PendingModifier provides methods for modifying pending requests
type PendingModifier interface {
	// SetRequest sets a request for a specific object
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, reqID insolar.ID) error
	// SetResult sets a result for a specific object. Also, if there is a not closed request for a provided result,
	// the request will be closed
	SetResult(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, resID insolar.ID, res record.Result) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.PendingAccessor -o ./ -s _mock.go

// PendingAccessor provides methods for fetching pending requests.
type PendingAccessor interface {
	// OpenRequestsForObjID returns a specific number of open requests for a specific object
	OpenRequestsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]record.Request, error)
	// OpenRequestsIDsForObjID returns a specific number of ids of open requests for a specific object
	OpenRequestsIDsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]insolar.ID, error)
	// Records returns all the records for a provided object
	Records(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error)
}

type HeavyPendingAccessor interface {
	Records(ctx context.Context, readFrom insolar.PulseNumber, readUntil insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.FilamentCacheManager -o ./ -s _mock.go

type FilamentCacheManager interface {
	Gather(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
	SendAbandonedNotification(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.FilamentCacheCleaner -o ./ -s _mock.go

type FilamentCacheCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}

// FilamentCacheStorage is a in-memory storage, that stores a collection of IndexBuckets
type FilamentCacheStorage struct {
	idxAccessor     IndexAccessor
	idxModifier     IndexModifier
	idLocker        IDLocker
	recordStorage   RecordStorage
	coordinator     jet.Coordinator
	pcs             insolar.PlatformCryptographyScheme
	pulseCalculator storage.PulseCalculator
	bus             insolar.MessageBus
	busWM           buswm.Sender

	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*pendingMeta
}

func NewFilamentCacheStorage(
	idxAccessor IndexAccessor,
	idxModifier IndexModifier,
	idLocker IDLocker,
	recordStorage RecordStorage,
	coordinator jet.Coordinator,
	pcs insolar.PlatformCryptographyScheme,
	pulseCalculator storage.PulseCalculator,
	bus insolar.MessageBus,
	busWM buswm.Sender,
) *FilamentCacheStorage {
	return &FilamentCacheStorage{
		idxAccessor:     idxAccessor,
		idxModifier:     idxModifier,
		idLocker:        idLocker,
		recordStorage:   recordStorage,
		coordinator:     coordinator,
		pcs:             pcs,
		pulseCalculator: pulseCalculator,
		bus:             bus,
		busWM:           busWM,
		buckets:         map[insolar.PulseNumber]map[insolar.ID]*pendingMeta{},
	}
}

// pendingMeta contains info for calculating pending requests states
// The structure contains full chain of pendings (they are grouped by pulse).
// Groups are sorted by pulse, from a lowest to a highest
// There are a few maps inside, that help not to search through full fillament every SetRequest/SetResult
type pendingMeta struct {
	sync.RWMutex

	isStateCalculated bool

	fullFilament []chainLink

	notClosedRequestsIds      []insolar.ID
	notClosedRequestsIdsIndex map[insolar.PulseNumber]map[insolar.ID]struct{}
}

type chainLink struct {
	PN             insolar.PulseNumber
	MetaRecordsIDs []insolar.ID
}

func (i *FilamentCacheStorage) createPendingBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *pendingMeta {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &pendingMeta{
		fullFilament:              []chainLink{},
		notClosedRequestsIds:      []insolar.ID{},
		notClosedRequestsIdsIndex: map[insolar.PulseNumber]map[insolar.ID]struct{}{},
		isStateCalculated:         false,
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*pendingMeta{}
		i.buckets[pn] = objsByPn
	}

	_, ok = objsByPn[objID]
	if !ok {
		objsByPn[objID] = bucket
	}

	inslogger.FromContext(ctx).Debugf("[createPendingBucket] create bucket for obj - %v was created successfully", objID.DebugString())
	return bucket
}

func (i *FilamentCacheStorage) pendingBucket(pn insolar.PulseNumber, objID insolar.ID) *pendingMeta {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	return objsByPn[objID]
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *FilamentCacheStorage) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	delete(i.buckets, pn)
}

// SetRequest sets a request for a specific object
func (i *FilamentCacheStorage) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, reqID insolar.ID) error {
	idx := i.idxAccessor.Index(pn, objID)
	if idx == nil {
		return ErrLifelineNotFound
	}
	i.idLocker.Lock(&objID)
	defer i.idLocker.Unlock(&objID)

	pb := i.pendingBucket(pn, objID)
	if pb == nil {
		pb = i.createPendingBucket(ctx, pn, objID)
	}

	inslogger.FromContext(ctx).Debugf("SetRequest before %v pn : %v", objID.DebugString(), pn)
	pb.Lock()
	defer pb.Unlock()
	inslogger.FromContext(ctx).Debugf("SetRequest after %v pn : %v", objID.DebugString(), pn)

	lfl := idx.Lifeline

	if lfl.PendingPointer != nil && reqID.Pulse() < lfl.PendingPointer.Pulse() {
		return errors.New("request from the past")
	}

	pf := record.PendingFilament{
		RecordID:       reqID,
		PreviousRecord: idx.Lifeline.PendingPointer,
	}

	if lfl.EarliestOpenRequest == nil {
		lfl.EarliestOpenRequest = &pn
	}

	pfv := record.Wrap(pf)
	hash := record.HashVirtual(i.pcs.ReferenceHasher(), pfv)
	metaID := *insolar.NewID(pn, hash)

	err := i.recordStorage.Set(ctx, metaID, record.Material{Virtual: &pfv, JetID: jetID})
	if err != nil {
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	idx.PendingRecords = append(idx.PendingRecords, metaID)
	lfl.PendingPointer = &metaID
	idx.Lifeline = lfl

	err = i.idxModifier.SetIndex(ctx, pn, *idx)
	if err != nil {
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	pb.addMetaIDToFilament(pn, metaID)

	_, ok := pb.notClosedRequestsIdsIndex[pn]
	if !ok {
		pb.notClosedRequestsIdsIndex[pn] = map[insolar.ID]struct{}{}
	}
	pb.notClosedRequestsIdsIndex[pn][reqID] = struct{}{}
	pb.notClosedRequestsIds = append(pb.notClosedRequestsIds, reqID)

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

	inslogger.FromContext(ctx).Debugf("open requests - %v for - %v", len(pb.notClosedRequestsIds), objID.DebugString())

	return nil

}

func (b *pendingMeta) addMetaIDToFilament(pn insolar.PulseNumber, metaID insolar.ID) {
	isInserted := false
	for i, chainPart := range b.fullFilament {
		if chainPart.PN == pn {
			b.fullFilament[i].MetaRecordsIDs = append(b.fullFilament[i].MetaRecordsIDs, metaID)
			isInserted = true
		}
	}

	if !isInserted {
		b.fullFilament = append(b.fullFilament, chainLink{MetaRecordsIDs: []insolar.ID{metaID}, PN: pn})
	}
}

// SetResult sets a result for a specific object. Also, if there is a not closed request for a provided result,
// the request will be closed
func (i *FilamentCacheStorage) SetResult(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, resID insolar.ID, res record.Result) error {
	logger := inslogger.FromContext(ctx)

	idx := i.idxAccessor.Index(pn, objID)
	if idx == nil {
		return ErrLifelineNotFound
	}

	pb := i.pendingBucket(pn, objID)
	if pb == nil {
		pb = i.createPendingBucket(ctx, pn, objID)
	}

	inslogger.FromContext(ctx).Debugf("SetResult before %v pn : %v", objID.DebugString(), pn)
	pb.Lock()
	defer pb.Unlock()
	inslogger.FromContext(ctx).Debugf("SetResult after %v pn : %v", objID.DebugString(), pn)

	_, ok := pb.notClosedRequestsIdsIndex[res.Request.Record().Pulse()]
	if !ok {
		// TODO: https://insolar.atlassian.net/browse/INS-2705 @egorikas
		logger.Error(errors.Wrapf(ErrResultWithoutRequest, "no requests for %v", resID.DebugString()))
		// return ErrResultWithoutRequest
	}

	lfl := idx.Lifeline

	pf := record.PendingFilament{
		RecordID:       resID,
		PreviousRecord: lfl.PendingPointer,
	}

	pfv := record.Wrap(pf)
	hash := record.HashVirtual(i.pcs.ReferenceHasher(), pfv)
	metaID := *insolar.NewID(pn, hash)

	err := i.recordStorage.Set(ctx, metaID, record.Material{Virtual: &pfv, JetID: jetID})
	if err != nil {
		if err == ErrOverride {
			logger.Error(errors.Wrap(err, "failed to add a result to filament"))
		}
		// return errors.Wrap(err, "failed to add a result to filament")
	}

	pb.addMetaIDToFilament(pn, metaID)

	reqsIDs, ok := pb.notClosedRequestsIdsIndex[res.Request.Record().Pulse()]
	if ok {
		delete(reqsIDs, *res.Request.Record())
		for i := 0; i < len(pb.notClosedRequestsIds); i++ {
			if pb.notClosedRequestsIds[i] == *res.Request.Record() {
				pb.notClosedRequestsIds = append(pb.notClosedRequestsIds[:i], pb.notClosedRequestsIds[i+1:]...)
				break
			}
		}
	}

	// If no open requests and we know about a full filament
	if len(pb.notClosedRequestsIds) == 0 && pb.isStateCalculated {
		logger.Debugf("no open requests for - %v, pn: %v,", objID.DebugString(), pn)
		logger.Debugf("RefreshPendingFilament set EarliestOpenRequest - %v, val - %v", objID.DebugString(), lfl.EarliestOpenRequest)
		lfl.EarliestOpenRequest = nil
	} else {
		logger.Debugf("not closed reqs: %v, pn: %v, objID : %v", len(pb.notClosedRequestsIds), pn, objID.DebugString())
	}

	idx.Lifeline = lfl

	err = i.idxModifier.SetIndex(ctx, pn, *idx)
	if err != nil {
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	stats.Record(ctx,
		statObjectPendingResultsInMemoryAddedCount.M(int64(1)),
	)

	return nil
}

// SetFilament adds a slice of records to an object with provided id and pulse. It's assumed, that the method is
// called for setting records from another light, during the process of filling full chaing of pendings
func (i *FilamentCacheStorage) setFilament(ctx context.Context, pm *pendingMeta, filPN insolar.PulseNumber, recs []record.CompositeFilamentRecord) error {
	recsIds := make([]insolar.ID, len(recs))
	for idx, rec := range recs {
		recsIds[idx] = rec.MetaID

		err := i.recordStorage.Set(ctx, rec.MetaID, rec.Meta)
		if err != nil && err != ErrOverride {
			return errors.Wrap(err, "filament update failed")
		}
		err = i.recordStorage.Set(ctx, rec.RecordID, rec.Record)
		if err != nil && err != ErrOverride {
			return errors.Wrap(err, "filament update failed")
		}
	}

	pm.fullFilament = append(pm.fullFilament, chainLink{MetaRecordsIDs: recsIds, PN: filPN})
	sort.Slice(pm.fullFilament, func(i, j int) bool {
		return pm.fullFilament[i].PN < pm.fullFilament[j].PN
	})

	return nil
}

func (i *FilamentCacheStorage) Gather(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	i.idLocker.Lock(&objID)
	defer i.idLocker.Unlock(&objID)

	idx := i.idxAccessor.Index(pn, objID)
	if idx == nil {
		return ErrLifelineNotFound
	}

	pb := i.pendingBucket(pn, objID)
	if pb == nil {
		pb = i.createPendingBucket(ctx, pn, objID)
	}

	inslogger.FromContext(ctx).Debugf("Gather before %v pn : %v", objID.DebugString(), pn)
	pb.Lock()
	defer pb.Unlock()
	inslogger.FromContext(ctx).Debugf("Gather after %v pn : %v", objID.DebugString(), pn)

	logger := inslogger.FromContext(ctx)
	lfl := idx.Lifeline

	// state already calculated
	if pb.isStateCalculated {
		logger.Debugf("Gather filament. objID - %v, pn - %v. State is already calculated", objID, pn)
		return nil
	}

	// No pendings
	if lfl.PendingPointer == nil {
		logger.Debugf("Gather filament. objID - %v, pn - %v. No pendings", objID, pn)
		return nil
	}
	// No open pendings
	if lfl.EarliestOpenRequest == nil {
		logger.Debugf("Gather filament. objID - %v, pn - %v. No open pendings", objID, pn)
		return nil
	}
	// If an earliest pending created during a current pulse
	if lfl.EarliestOpenRequest != nil && *lfl.EarliestOpenRequest == pn {
		logger.Debugf("Gather filament. objID - %v, pn - %v. If an earliest pending created during a current pulse", objID, pn)
		return nil
	}

	fp, err := i.firstPending(ctx, pb)
	if err != nil {
		return err
	}

	if fp == nil || fp.PreviousRecord == nil {
		err = i.fillPendingFilament(ctx, pn, objID, lfl.PendingPointer.Pulse(), *lfl.EarliestOpenRequest, pb)
		if err != nil {
			return err
		}
	} else {
		err = i.fillPendingFilament(ctx, pn, objID, fp.PreviousRecord.Pulse(), *lfl.EarliestOpenRequest, pb)
		if err != nil {
			return err
		}
	}

	err = i.refresh(ctx, idx, pb)
	if err != nil {
		return err
	}

	err = i.idxModifier.SetIndex(ctx, pn, *idx)
	if err != nil {
		return err
	}

	return nil
}

func (i *FilamentCacheStorage) SendAbandonedNotification(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) error {
	logger := inslogger.FromContext(ctx)
	idx := i.idxAccessor.Index(currentPN, objID)
	if idx == nil {
		return ErrLifelineNotFound
	}

	if idx.Lifeline.EarliestOpenRequest == nil {
		return nil
	}

	notifyPoint, err := i.pulseCalculator.Backwards(ctx, currentPN, 2)
	if err == pulse.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if notifyPoint.PulseNumber < *idx.Lifeline.EarliestOpenRequest {
		return nil
	}

	rep, err := i.bus.Send(ctx, &message.AbandonedRequestsNotification{
		Object: objID,
	}, nil)
	if err != nil {
		logger.Error("failed to notify about pending requests")
		return err
	}
	if _, ok := rep.(*reply.OK); !ok {
		logger.Error("received unexpected reply on pending notification")
		return errors.New("received unexpected reply on pending notification")
	}
	return nil
}

func (i *FilamentCacheStorage) firstPending(ctx context.Context, pb *pendingMeta) (*record.PendingFilament, error) {
	if len(pb.fullFilament) == 0 {
		return nil, nil
	}
	if len(pb.fullFilament[0].MetaRecordsIDs) == 0 {
		return nil, nil
	}

	metaID := pb.fullFilament[0].MetaRecordsIDs[0]
	rec, err := i.recordStorage.ForID(ctx, metaID)
	if err != nil {
		return nil, err
	}

	return record.Unwrap(rec.Virtual).(*record.PendingFilament), nil
}

func (i *FilamentCacheStorage) fillPendingFilament(
	ctx context.Context,
	currentPN insolar.PulseNumber,
	objID insolar.ID,
	destPN insolar.PulseNumber,
	earlistOpenRequest insolar.PulseNumber,
	pm *pendingMeta,
) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("RefreshPendingFilament.fillPendingFilament"))
	defer span.End()

	logger := inslogger.FromContext(ctx)
	logger.Debugf("objID:%v fillPendingFilament start, pn:%v, destPN: %v, earliestOpenRequest: %v", objID.DebugString(), currentPN, destPN, earlistOpenRequest)
	continueFilling := true

	for continueFilling {
		node, err := i.coordinator.NodeForObject(ctx, objID, currentPN, destPN)
		if err != nil {
			return err
		}

		logger.Debugf("objID:%v fillPendingFilament goes to destPN:%v, pn:%v node: %v", objID.DebugString(), destPN, currentPN, node)
		var pl payload.Payload
		// TODO: temp hack waiting for INS-2597 INS-2598 @egorikas
		// Because a current node can be a previous LME for a object
		if *node == i.coordinator.Me() {
			res, err := i.Records(ctx, destPN, objID)
			if err != nil {
				return errors.Wrap(err, "can't get a previous segment of the filament")
			}

			pl = &payload.PendingFilament{
				ObjectID: objID,
				Records:  res,
			}
		} else {
			logger.Debugf("objID:%v fillPendingFilament &payload.GetPendingFilament, pn:%v", objID.DebugString(), currentPN)
			msg, err := payload.NewMessage(&payload.GetPendingFilament{
				ObjectID:  objID,
				StartFrom: destPN,
				ReadUntil: earlistOpenRequest,
			})
			if err != nil {
				return errors.Wrap(err, "failed to create a GetPendingFilament message")
			}

			logger.Debugf("objID:%v fillPendingFilament i.busWM.SendTarget, pn:%v", objID.DebugString(), currentPN)
			rep, done := i.busWM.SendTarget(ctx, msg, *node)
			defer done()
			var ok bool
			res, ok := <-rep
			logger.Debugf("objID:%v fillPendingFilament after res, ok := <-rep, pn:%v", objID.DebugString(), currentPN)
			if !ok {
				return errors.New("failed to get a pending filament. no reply")
			}
			logger.Debugf("objID:%v fillPendingFilament payload.UnmarshalFromMeta, pn:%v", objID.DebugString(), currentPN)
			pl, err = payload.UnmarshalFromMeta(res.Payload)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal reply")
			}

		}
		switch r := pl.(type) {
		case *payload.PendingFilament:
			if len(r.Records) == 0 {
				panic(fmt.Sprintf("unexpected behaviour, objID - %v, pn - %v, node - %v, startFrom %v  readUntil %v", objID.DebugString(), currentPN, node, destPN, earlistOpenRequest))
			}
			logger.Debugf("objID:%v fillPendingFilament setFilament, pn:%v", objID.DebugString(), currentPN)
			err := i.setFilament(ctx, pm, destPN, r.Records)
			if err != nil {
				return err
			}

			logger.Debugf("objID:%v fillPendingFilament record.Unwrap(r.Records[0].Meta.Virtual), pn:%v", objID.DebugString(), currentPN)
			firstRec := record.Unwrap(r.Records[0].Meta.Virtual).(*record.PendingFilament)
			logger.Debugf("objID:%v fillPendingFilament firstRec.PreviousRecord==%v, pn:%v", objID.DebugString(), firstRec.PreviousRecord.DebugString(), currentPN)
			if firstRec.PreviousRecord == nil {
				logger.Debugf("objID:%v fillPendingFilament finish, pn:%v", objID.DebugString(), currentPN)
				return nil
			}

			// If know border read to the start of the chain
			// In other words, we read until limit
			if firstRec.PreviousRecord.Pulse() >= earlistOpenRequest {
				destPN = firstRec.PreviousRecord.Pulse()
			} else {
				continueFilling = false
			}
			logger.Debugf("objID:%v fillPendingFilament continueFilling destPN:%v, continueFilling:%v, pn:%v", objID.DebugString(), destPN, continueFilling, currentPN)
		case *payload.Error:
			logger.Debugf("objID:%v fillPendingFilament *payload.Error, pn:%v", objID.DebugString(), r.Text, currentPN)
			return errors.New(r.Text)
		default:
			return fmt.Errorf("fillPendingFilament: unexpected reply: %#v", r)
		}
	}

	logger.Debugf("objID:%v fillPendingFilament finish, pn:%v", objID.DebugString(), currentPN)
	return nil
}

// RefreshState recalculates state of the chain, marks requests as closed and opened.
func (i *FilamentCacheStorage) refresh(ctx context.Context, idx *FilamentIndex, pb *pendingMeta) error {
	if pb.isStateCalculated {
		return nil
	}

	for _, chainLink := range pb.fullFilament {
		for _, metaID := range chainLink.MetaRecordsIDs {
			metaRec, err := i.recordStorage.ForID(ctx, metaID)
			if err != nil {
				return errors.Wrap(err, "failed to refresh an index state")
			}

			concreteMeta := record.Unwrap(metaRec.Virtual).(*record.PendingFilament)
			rec, err := i.recordStorage.ForID(ctx, concreteMeta.RecordID)
			if err != nil {
				return errors.Wrap(err, "failed to refresh an index state")
			}

			switch r := record.Unwrap(rec.Virtual).(type) {
			case *record.Request:
				_, ok := pb.notClosedRequestsIdsIndex[chainLink.PN]
				if !ok {
					pb.notClosedRequestsIdsIndex[chainLink.PN] = map[insolar.ID]struct{}{}
				}
				pb.notClosedRequestsIdsIndex[chainLink.PN][concreteMeta.RecordID] = struct{}{}
			case *record.Result:
				openReqs, ok := pb.notClosedRequestsIdsIndex[r.Request.Record().Pulse()]
				if ok {
					delete(openReqs, *r.Request.Record())
				}
			default:
				panic(fmt.Sprintf("unknow type - %v", r))
			}
		}
	}

	isEarliestFound := false

	for i, chainLink := range pb.fullFilament {
		if len(pb.notClosedRequestsIdsIndex[chainLink.PN]) != 0 {
			if !isEarliestFound {
				idx.Lifeline.EarliestOpenRequest = &pb.fullFilament[i].PN
				isEarliestFound = true
			}

			for openReqID := range pb.notClosedRequestsIdsIndex[chainLink.PN] {
				pb.notClosedRequestsIds = append(pb.notClosedRequestsIds, openReqID)
			}
		}
	}

	pb.isStateCalculated = true

	if len(pb.notClosedRequestsIds) == 0 {
		idx.Lifeline.EarliestOpenRequest = nil
	}

	return nil
}

func (i *FilamentCacheStorage) OpenRequestsIDsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]insolar.ID, error) {
	pb := i.pendingBucket(currentPN, objID)
	if pb == nil {
		return nil, ErrLifelineNotFound
	}

	inslogger.FromContext(ctx).Debugf("OpenRequestsForObjID before %v pn : %v", objID.DebugString(), currentPN)
	pb.RLock()
	defer pb.RUnlock()

	if !pb.isStateCalculated {
		return []insolar.ID{}, nil
	}

	inslogger.FromContext(ctx).Debugf("OpenRequestsForObjID after %v pn : %v", objID.DebugString(), currentPN)

	if len(pb.notClosedRequestsIds) < count {
		count = len(pb.notClosedRequestsIds)
	}

	return pb.notClosedRequestsIds[:count], nil
}

// OpenRequestsForObjID returns open requests for a specific object
func (i *FilamentCacheStorage) OpenRequestsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]record.Request, error) {
	pb := i.pendingBucket(currentPN, objID)
	if pb == nil {
		return nil, ErrLifelineNotFound
	}

	inslogger.FromContext(ctx).Debugf("OpenRequestsForObjID before %v pn : %v", objID.DebugString(), currentPN)
	pb.RLock()
	defer pb.RUnlock()

	if !pb.isStateCalculated {
		return []record.Request{}, nil
	}

	inslogger.FromContext(ctx).Debugf("OpenRequestsForObjID after %v pn : %v", objID.DebugString(), currentPN)

	if len(pb.notClosedRequestsIds) < count {
		count = len(pb.notClosedRequestsIds)
	}

	res := make([]record.Request, count)

	for idx := 0; idx < count; idx++ {
		rec, err := i.recordStorage.ForID(ctx, pb.notClosedRequestsIds[idx])
		if err != nil {
			return nil, err
		}

		switch r := record.Unwrap(rec.Virtual).(type) {
		case *record.Request:
			if r.Object == nil {
				panic(fmt.Sprintf("nothing wrong here %v", pb.notClosedRequestsIds[idx].DebugString()))
			}
			res[idx] = *r
		default:
			panic("filament is totally broken")
		}
	}

	return res, nil
}

// Records returns all the records for a provided object
func (i *FilamentCacheStorage) Records(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error) {
	b := i.pendingBucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	inslogger.FromContext(ctx).Debugf("Records before %v pn : %v", objID.DebugString(), currentPN)
	b.RLock()
	defer b.RUnlock()
	inslogger.FromContext(ctx).Debugf("Records before %v pn : %v", objID.DebugString(), currentPN)

	var res []record.CompositeFilamentRecord

	for _, pm := range b.fullFilament {
		var tempRes []record.CompositeFilamentRecord

		for _, recID := range pm.MetaRecordsIDs {
			metaRec, err := i.recordStorage.ForID(ctx, recID)
			if err != nil {
				return nil, err
			}

			concreteMeta := record.Unwrap(metaRec.Virtual).(*record.PendingFilament)
			rec, err := i.recordStorage.ForID(ctx, concreteMeta.RecordID)
			if err != nil {
				return nil, err
			}

			tempRes = append(tempRes, record.CompositeFilamentRecord{
				Record:   rec,
				RecordID: concreteMeta.RecordID,
				Meta:     metaRec,
				MetaID:   recID,
			})
		}

		res = append(tempRes, res...)
	}

	return res, nil
}
