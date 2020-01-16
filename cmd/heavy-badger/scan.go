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

import "github.com/insolar/insolar/insolar/store"

// AALEKSEEV TODO get rid of heavy-badger

import (
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

var mb = float64(1 << 20)

type dbScanner struct {
	db *store.BadgerDB

	nonStrict bool

	perPulseStat bool
	statGraph    bool

	disableProgressbar bool

	enableRecordsTypesStat bool

	limitPulses int
	skipPulses  int

	searchValuesGreaterThan int64
}

func (dbs *dbScanner) openDB(dataDir string) {
	db, _ := openDB(dataDir)
	dbs.db = db
}

func (dbs *dbScanner) closeDB() {
	err := dbs.db.Backend().Close()
	if err != nil {
		fatalf("failed close database: %v", err)
	}
}

// pulseView hides irrelevant insolar.Pulse fields from json serializer
type pulseView struct {
	insolar.Pulse
	OriginID *bool `json:",omitempty"`
	Entropy  *bool `json:",omitempty"`
	Signs    *bool `json:",omitempty"`
}

func pulseTime(p insolar.Pulse) time.Time {
	return time.Unix(p.PulseTimestamp/1000000000, 0)
}

func (dbs *dbScanner) printAllPulsesInfo(pulses []insolar.Pulse) {
	for i := 0; i < len(pulses); i++ {
		p := pulses[len(pulses)-i-1]
		fmt.Printf("pulse=%v, timestamp=%v\n", p.PulseNumber, pulseTime(p))
	}
}

func (dbs *dbScanner) printShortPulsesInfo(pulses []insolar.Pulse) {
	fmt.Println("Pulses info:")
	printLine("=")
	if len(pulses) == 0 {
		fmt.Printf("found %v pulses\n", len(pulses))
		return
	}

	first, last := pulses[len(pulses)-1], pulses[0]
	printPulse := func(prefix string, p insolar.Pulse) {
		pv := pulseView{Pulse: p}
		fmt.Printf(prefix+" (timestamp -> %v): ", pulseTime(p))
		printJSON(pv, jsonPrefix(""), setPretty(true))
	}
	if first.PulseNumber != insolar.GenesisPulse.PulseNumber {
		dbs.failIfStrictf("first pulse %v is a not genesis pulse", first.PulseNumber)
	}
	printPulse("Genesis pulse:", first)
	fmt.Println()

	if len(pulses) > 1 {
		first = pulses[len(pulses)-2]
	}
	d := int64(last.PulseNumber - first.PulseNumber)

	fmt.Printf("Found pulses %v, [%v:%v] (Δ=%v ≈%s)\n",
		len(pulses), first.PulseNumber, last.PulseNumber, d, time.Duration(int64(time.Second)*d))
	printPulse("first pulse:", first)
	printPulse("last pulse:", last)
}

func (dbs *dbScanner) scopesReport(fast bool, names []string, ids []int) {
	seen := map[store.Scope]struct{}{}
	var toScan []store.Scope
	for _, name := range names {
		scope, err := scopeFromName(name)
		if err != nil {
			fatalf("scan failed: %v", err)
		}
		toScan = append(toScan, scope)
		seen[scope] = struct{}{}
	}
	for id := range ids {
		scope := store.Scope(id)
		if _, ok := seen[scope]; ok {
			continue
		}
		toScan = append(toScan, scope)
	}
	if len(toScan) == 0 {
		toScan = allScopes()
	}
	for _, scope := range toScan {
		dbs.scanWholeScope(scope, fast)
	}
}

func (dbs *dbScanner) scanWholeScope(scope store.Scope, fast bool) {
	printLine("=")
	fmt.Printf("Scan scope %s\n", scope)
	printLine("=")

	h := newHistogram("Summary Sizes")

	var opts = &iterOptions{
		counter:  true,
		keysOnly: fast,
	}
	iterate(dbs.db, scopeKey{scope}, opts, h.iter)

	h.PrintKeys()
	if !fast {
		printLine("-")
		h.PrintValues()
	}
	fmt.Println()
}

type printer interface {
	Print()
}

func valuesExceedSize(scope store.Scope, size int64) iteration {
	if size == 0 {
		return nil
	}
	return func(k, v []byte) error {
		valueSize := int64(len(v))
		if valueSize < size {
			return nil
		}

		info := []string{fmt.Sprintf("key=%x - FOUND VALUE: size=%v", k, humanize.Bytes(uint64(valueSize)))}
		if scope == store.ScopeRecord {
			var matRec record.Material
			err := matRec.Unmarshal(v)
			if err != nil {
				panic(err)
			}
			virtual := record.Unwrap(&matRec.Virtual)
			info = append(info, fmt.Sprintf("type=%T", virtual))
		}
		if scope == store.ScopeRecord {
			id := insolar.NewIDFromBytes(k)
			info = append(info, "id="+id.String(), "pulse="+id.Pulse().String())
		}

		fmt.Println(strings.Join(info, " "))
		return nil
	}
}

type recordsStatByType struct {
	stats map[string]*histogram
}

func newRecordsStatByType() *recordsStatByType {
	return &recordsStatByType{
		stats: make(map[string]*histogram),
	}
}

func (rs *recordsStatByType) iter(k, v []byte) error {
	var matRec record.Material
	err := matRec.Unmarshal(v)
	if err != nil {
		return err
	}
	virtual := record.Unwrap(&matRec.Virtual)
	recType := fmt.Sprintf("	%T", virtual)

	hist, ok := rs.stats[recType]
	if !ok {
		hist = newHistogram(recType)
		rs.stats[recType] = hist
	}
	hist.values.Update(int64(len(v)))
	return nil
}

func (rs *recordsStatByType) Print() {
	fmt.Println("Records Stat By Type")
	printLine("=")
	for _, hist := range rs.stats {
		hist.PrintValues()
		printLine("-")
	}
	fmt.Println()
}
