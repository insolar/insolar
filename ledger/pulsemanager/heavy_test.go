package pulsemanager_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

func TestPulseManager_SendToHeavy(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	// Mock N1: LR mock do nothing
	lrMock := testutils.NewLogicRunnerMock(t)
	lrMock.OnPulseMock.Return(nil)

	// Mock N2: we are light material
	nodeMock := network.NewNodeMock(t)
	nodeMock.RoleMock.Return(core.RoleLightMaterial)

	// Mock N3: nodenet returns mocked node (above)
	// and add stub for GetActiveNodes
	nodenetMock := network.NewNodeNetworkMock(t)
	nodenetMock.GetActiveNodesMock.Return(nil)
	nodenetMock.GetOriginMock.Return(nodeMock)

	// Mock N4: message bus for Send method
	busMock := testutils.NewMessageBusMock(t)

	// mock bus.Mock method, store synced records, and calls count with HeavyRecord
	var synckeys []key
	var syncsended int
	type messageStat struct {
		size int
		keys []key
	}
	syncmessagesPerMessage := map[int]*messageStat{}
	busMock.SendFunc = func(ctx context.Context, msg core.Message) (core.Reply, error) {
		heavymsg, ok := msg.(*message.HeavyRecords)
		if ok {
			syncsended++
			var size int
			var keys []key

			// fmt.Printf("[%v] prepared message with keys:\n", syncsended)
			for _, rec := range heavymsg.Records {
				keys = append(keys, rec.K)
				size += len(rec.K) + len(rec.V)

				// k := key(rec.K)
				// fmt.Printf("  [%v] %v (pulse=%v)\n", syncsended, k, k.pulse())
			}
			synckeys = append(synckeys, keys...)
			syncmessagesPerMessage[syncsended] = &messageStat{
				size: size,
				keys: keys,
			}
		}
		return nil, nil
	}

	// build PulseManager
	pm := pulsemanager.NewPulseManager(db)
	pm.LR = lrMock
	pm.NodeNet = nodenetMock
	pm.Bus = busMock

	// start PulseManager
	err := pm.Start(ctx)
	assert.NoError(t, err)

	// store last pulse as light material and set next one
	lastpulse := core.FirstPulseNumber + 1
	err = setpulse(ctx, pm, lastpulse)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		// fmt.Printf("%v: call addRecords for pulse %v\n", t.Name(), lastpulse)
		for j := 0; j < 3; j++ {
			addRecords(ctx, t, db, core.PulseNumber(lastpulse))
		}
		lastpulse++
	}

	// fmt.Println(strings.Repeat("-", 55))
	// fmt.Println("Case1: sync after db fill and with new received pulses")
	err = setpulse(ctx, pm, lastpulse)
	require.NoError(t, err)

	// fmt.Println(strings.Repeat("-", 55))
	// fmt.Println("Case2: sync during db fill")
	for i := 0; i < 2; i++ {
		// fmt.Printf("%v: call addRecords for pulse %v\n", t.Name(), lastpulse)
		// fill DB with records, indexes (TODO: add blobs)
		for j := 0; j < 3; j++ {
			addRecords(ctx, t, db, core.PulseNumber(lastpulse))
		}

		lastpulse++
		err = setpulse(ctx, pm, lastpulse)
		require.NoError(t, err)
	}
	// time.Sleep(time.Second * 2)

	err = pm.Stop(ctx)
	assert.NoError(t, err)

	synckeys = sortkeys(synckeys)
	// fmt.Println("sort(syncmessages):", len(synckeys))
	// printkeys(synckeys, "  ")

	synckeys = uniqkeys(synckeys)
	// fmt.Println("uniq(syncmessages):", len(synckeys))
	// printkeys(synckeys, "  ")

	// fmt.Println("Per message:", syncsended)
	// for n := 1; n <= syncsended; n++ {
	// 	stat := syncmessagesPerMessage[n]
	// 	fmt.Printf("Message %v. Size=%v\n", n, stat.size)
	// 	printkeys(sortkeys(uniqkeys(stat.keys)), "  ")
	// }
	// fmt.Println("")

	recs, idxs := getallkeys(db.GetBadgerDB())
	// fmt.Println("getallkeys:", len(idxs)+len(recs))
	// fmt.Println(" records", len(recs))
	// printkeys(recs, "  ")
	// fmt.Println(" indexes", len(idxs))
	// printkeys(idxs, "  ")

	assert.Equal(t, len(recs)+len(idxs), len(synckeys), "uniq synced keys count are the same as records count in database")
}

func setpulse(ctx context.Context, pm core.PulseManager, pulsenum int) error {
	fmt.Println("CALL PulseManager.Set with pulse", pulsenum)
	return pm.Set(ctx, core.Pulse{PulseNumber: core.PulseNumber(pulsenum)})
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	db *storage.DB,
	pulsenum core.PulseNumber,
) {
	// set record
	parentID, err := db.SetRecord(
		ctx,
		pulsenum,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: testutils.RandomRef(),
			},
		},
	)
	require.NoError(t, err)

	// set index of record
	err = db.SetObjectIndex(ctx, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(t, err)
	return
}

var (
	scopeIDLifeline = byte(1)
	scopeIDRecord   = byte(2)
)

type key []byte

func getallkeys(db *badger.DB) (records []key, indexes []key) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		switch k[0] {
		case scopeIDRecord:
			records = append(records, k)
		case scopeIDLifeline:
			indexes = append(indexes, k)
		}
	}
	return
}

func (b key) pulse() core.PulseNumber {
	return core.NewPulseNumber(b[1 : 1+core.PulseNumberSize])
}

func (b key) String() string {
	return hex.EncodeToString(b)
}

func printkeys(keys []key, prefix string) {
	for _, k := range keys {
		fmt.Printf("%v%v (%v)\n", prefix, k, k.pulse())
	}
}

func uniqkeys(keys []key) (keyout []key) {
	uniq := map[string]bool{}
	for _, k := range keys {
		if uniq[string(k)] {
			continue
		}
		uniq[string(k)] = true
		keyout = append(keyout, k)
	}
	return
}

func sortkeys(keys []key) []key {
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	return keys
}
