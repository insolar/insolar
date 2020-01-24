///
// Copyright 2020 Insolar Technologies GmbH
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
///

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/pkg/errors"
)

// data from task
//--------------------------
var numPulses = 6500000
var numRecords = 150000000
var numIndexes = 150000000
var numDrops = 208000000

//--------------------------

var recordRatio = numRecords / numPulses
var indexRatio = numIndexes / numPulses
var dropRatio = numDrops / numPulses

var idChan = make(chan []byte, 1000)

func genIDS() {
	for i := 0; ; i++ {
		hash := gen.ID().Hash()
		idChan <- hash
	}
}

var signature = gen.Signature(256)

func writeRecord(ctx context.Context, db *object.RecordDB, pulseNumber insolar.PulseNumber) {
	records := make([]record.Material, 0, recordRatio)
	for i := 0; i < recordRatio; i++ {
		hash := <-idChan
		id := insolar.NewID(pulseNumber, hash)
		records = append(records, record.Material{ID: *id, Signature: signature})
	}
	err := db.BatchSet(ctx, records)
	if err != nil {
		panic(err)
	}
}

func writeIndex(ctx context.Context, db *object.IndexDB, pulseNumber insolar.PulseNumber, wg *sync.WaitGroup) {
	for i := 0; i < indexRatio; i++ {
		hash := <-idChan
		err := db.SetIndex(ctx, pulseNumber, record.Index{ObjID: *insolar.NewID(pulseNumber, hash)})
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func writeDrop(ctx context.Context, db *drop.DB, pulseNumber insolar.PulseNumber, wg *sync.WaitGroup) {
	numJets := 32
	for i := 0; i < numJets; i++ {
		hash := <-idChan
		err := db.Set(ctx, drop.Drop{Pulse: pulseNumber, JetID: *insolar.NewJetID(uint8(len(hash)), hash)})
		if err != nil {
			panic(err)
		}
	}

	wg.Done()
}

var entropyGenerator = &entropygenerator.StandardEntropyGenerator{}
var pulseDelta = insolar.PulseNumber(10)
var originID = [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0}
var entropy = entropyGenerator.GenerateEntropy()

func makeNewPulse(newPulseNumber insolar.PulseNumber) *insolar.Pulse {
	return &insolar.Pulse{
		PulseNumber:      newPulseNumber,
		Entropy:          entropy,
		NextPulseNumber:  newPulseNumber + pulseDelta,
		EpochPulseNumber: newPulseNumber.AsEpoch(),
		OriginID:         originID,
		PulseTimestamp:   int64(newPulseNumber),
		Signs:            map[string]insolar.PulseSenderConfirmation{},
	}
}

func main() {

	globalStart := time.Now()
	go genIDS()

	ctx := context.Background()
	bdb, err := store.NewBadgerDB(badger.DefaultOptions("/tmp/LATENCY_BADGER"))
	if err != nil {
		panic(errors.Wrap(err, "failed to open badger"))
	}

	defer bdb.Stop(ctx)

	pulseDB := pulse.NewDB(bdb)
	recordDB := object.NewRecordDB(bdb)
	indexDB := object.NewIndexDB(bdb, recordDB)
	dropDB := drop.NewDB(bdb)

	numIterations := numPulses / 1000 // write only 1/1000 part of oll data ( for tests )

	wg := &sync.WaitGroup{}
	wg.Add(2 * numIterations)
	for i := 0; i < numIterations; i++ {
		if i%29 == 0 {
			fmt.Println("iter: ", i)
		}

		start := time.Now()
		nextPulseNumber := insolar.GenesisPulse.PulseNumber + (insolar.PulseNumber(i) * pulseDelta)
		nextPulse := makeNewPulse(nextPulseNumber)
		fmt.Printf("Make new pulse took %s\n", time.Since(start))

		start = time.Now()
		go writeDrop(ctx, dropDB, nextPulseNumber, wg)
		fmt.Printf("writeDrop took %s\n", time.Since(start))

		start = time.Now()
		go writeIndex(ctx, indexDB, nextPulseNumber, wg)
		fmt.Printf("writeIndex took %s\n", time.Since(start))

		start = time.Now()
		writeRecord(ctx, recordDB, nextPulseNumber)
		fmt.Printf("writeRecord took %s\n", time.Since(start))

		start = time.Now()
		err = pulseDB.Append(ctx, *nextPulse)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Append took %s\n", time.Since(start))
		fmt.Println(" ")
	}

	wg.Wait()

	fmt.Printf("\nTotal time %s, number of iterations: %d\n", time.Since(globalStart), numIterations)
}
