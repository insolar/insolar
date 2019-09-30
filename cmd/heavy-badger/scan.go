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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar"
	pulsedb "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
)

func (app *appCtx) scanCommand() *cobra.Command {
	var (
		showStats bool
		// from pulse
	)
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scans database and return various stats.",
		Run: func(_ *cobra.Command, _ []string) {
			if err := checkDirectory(app.dataDir); err != nil {
				fatalf("Database directory '%v' open failed. Error: \"%v\"", app.dataDir, err)
			}

			ops := badger.DefaultOptions(app.dataDir)
			dbWrapped, err := store.NewBadgerDB(ops)
			if err != nil {
				fatalf("failed open database directory %v: %v", app.dataDir, err)
			}
			defer func() {
				err := dbWrapped.Backend().Close()
				if err != nil {
					fatalf("failed close database directory %v: %v", app.dataDir, err)
				}
			}()
			showDBStat(dbWrapped)
		},
	}
	scanCmd.Flags().BoolVar(&showStats, "show-stats", true, "show badger stats")
	// _ = showStats
	return scanCmd
}

func getAllPulses(ctx context.Context, db *store.BadgerDB) (pulses []insolar.Pulse) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERROR: getAllPulses", err)
		}
	}()

	var pulseStore = pulsedb.NewDB(db)
	p, err := pulseStore.Latest(ctx)
	pulses = append(pulses, p)
	if err == pulsedb.ErrNotFound {
		return
	}

	if err != nil {
		fatalf("failed to get latest pulse: %v\n", err)
	}
	pulses = append(pulses, p)

	for ; err == nil; p, err = pulseStore.Backwards(ctx, p.PulseNumber, 1) {
		// fmt.Println("found:", p.PulseNumber)
		pulses = append(pulses, p)
	}
	if err != nil && err != pulsedb.ErrNotFound {
		fatalf("failed to get pulse on step %v: %v\n", len(pulses), err)
	}
	return
}

func printMinusMinusLine() {
	fmt.Println(strings.Repeat("-", 50))
}

type pulseView struct {
	insolar.Pulse
	// hide unrelevant insolar.Pulse fields from json serializer
	OriginID *bool `json:",omitempty"`
	Entropy  *bool `json:",omitempty"`
	Signs    *bool `json:",omitempty"`
}

type jsonCfg struct {
	pretty bool
	indent string
	prefix string
}

type jsonOpt func(*jsonCfg)

func jsonPretty(pretty bool) jsonOpt {
	return func(cfg *jsonCfg) { cfg.pretty = pretty }
}

func jsonPrefix(prefix string) jsonOpt {
	return func(cfg *jsonCfg) { cfg.prefix = prefix }
}

func printJSON(v interface{}, opts ...jsonOpt) {
	cfg := &jsonCfg{
		pretty: true,
		indent: "  ",
		prefix: "  ",
	}
	for _, o := range opts {
		o(cfg)
	}

	var (
		b   []byte
		err error
	)
	if cfg.pretty {
		b, err = json.MarshalIndent(v, cfg.prefix, cfg.indent)
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		panic(err)
	}

	if cfg.pretty {
		fmt.Printf(cfg.indent)
	}
	fmt.Printf("%s\n", b)
}

func showDBStat(db *store.BadgerDB) {
	var ctx = context.Background()

	// show pulses info
	allPulses := getAllPulses(ctx, db)
	fmt.Println("Pulses info:")
	printMinusMinusLine()
	fmt.Printf("found %v pulses\n", len(allPulses))
	if len(allPulses) > 0 {
		first, last := allPulses[len(allPulses)-1], allPulses[0]
		printPulse := func(p insolar.Pulse) {
			pv := pulseView{Pulse: p}
			printJSON(pv, jsonPrefix(""))
			fmt.Printf("pulse timestamp -> %v\n", time.Unix(pv.PulseTimestamp/1000000000, 0))
		}
		if first.PulseNumber != insolar.GenesisPulse.PulseNumber {
			panic("first pulse is a not genesis pulse:")
		}
		fmt.Print("Genesis pulse:")
		printPulse(first)
		fmt.Println()

		if len(allPulses) > 1 {
			first = allPulses[len(allPulses)-2]
		}
		d := int64(last.PulseNumber - first.PulseNumber)
		fmt.Printf("Pulses [%v:%v] (Δ=%v ≈%s)\n",
			first.PulseNumber, last.PulseNumber,
			d, time.Duration(int64(time.Second)*d))
		fmt.Print("first pulse:")
		printPulse(first)
		fmt.Print("last pulse:")
		printPulse(last)
	}
	fmt.Println()

	// allPulses = append(allPulses, *insolar.GenesisPulse)
	for i := len(allPulses) - 1; i >= 0; i-- {
		pn := allPulses[i].PulseNumber
		fmt.Printf("Pulse %v stats.\n", pn)

		lastKnownStat, _ := getKVStatForPrefix(ctx, db, pulseKey{store.ScopeLastKnownIndexPN, pn})
		fmt.Print("lastKnownIndexPNKey:")
		printJSON(lastKnownStat, jsonPretty(false))
		// prefix := bytes.Join([][]byte{pn.Bytes()}, nil)

		indexStat, _ := getKVStatForPrefix(ctx, db, pulseKey{store.ScopeIndex, pn})
		fmt.Print("indexes:            ")
		printJSON(indexStat, jsonPretty(false))

		recordsStat, recordsHist := getKVStatForPrefix(ctx, db, pulseKey{store.ScopeRecord, pn})
		fmt.Print("records:            ")
		printJSON(recordsStat, jsonPretty(false))
		fmt.Printf("Values size: %s\n", humanize.Bytes(uint64(recordsStat.ValuesTotalBytes)))
		fmt.Printf("Histogram of value sizes (in bytes)\n")
		recordsHist.valueSizeHistogram.printHistogram()

		fmt.Println()
	}
	fmt.Println()

}

type KVStat struct {
	Count            int64
	KeysTotalBytes   int64
	ValuesTotalBytes int64
}

type pulseKey struct {
	scope store.Scope
	pn    insolar.PulseNumber
}

func (k pulseKey) Scope() store.Scope {
	return k.scope
}

func (k pulseKey) ID() []byte {
	return k.pn.Bytes()
}

// func (k pulseKey) prefix() []byte {
// 	return append(k.Scope().Bytes(), k.ID()...)
// }
//

func getKVStatForPrefix(
	ctx context.Context,
	db *store.BadgerDB,
	start pulseKey,
) (KVStat, *sizeHistogram) {
	badgerHistogram := newSizeHistogram()

	// Collect key and value sizes.
	// for itr.Seek(keyPrefix); itr.ValidForPrefix(keyPrefix); itr.Next() {

	it := db.NewIterator(start, false)
	defer it.Close()

	var result KVStat
	for it.Next() {
		k := it.Key()
		if !bytes.HasPrefix(k, start.ID()) {
			break
		}
		v, err := it.Value()
		if err != nil {
			panic(err)
		}
		result.Count++
		keySize, valueSize := int64(len(k)), int64(len(v))
		result.KeysTotalBytes += keySize
		result.ValuesTotalBytes += valueSize

		badgerHistogram.keySizeHistogram.Update(keySize)
		badgerHistogram.valueSizeHistogram.Update(valueSize)

		// TODO: move to separate command like 'huge-records'
		if valueSize > 16777216 {
			fmt.Printf(">> Found Huge Record value. Size= %v (%s): ", valueSize, humanize.Bytes(uint64(valueSize)))
			fmt.Printf("ID=%s\n", insolar.NewIDFromBytes(k).String())

			var matRec record.Material
			err := matRec.Unmarshal(v)
			if err != nil {
				panic(err)
			}
			virtual := record.Unwrap(&matRec.Virtual)
			fmt.Printf("Material>Material: type=%T\n", virtual)
			// matRec.
		}
	}
	// fmt.Printf("last key after iter %x\n", k)
	return result, badgerHistogram
}
