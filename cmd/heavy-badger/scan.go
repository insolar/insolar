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
	"context"
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar"
	pulsedb "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
)

var mb = float64(1 << 20)

type dbScanner struct {
	db *store.BadgerDB

	nonStrict bool

	showValuesHistogram bool
	perPulseStat        bool
	statGraph           string

	disableProgressbar bool

	enableRecordsTypesStat bool

	// TODO: add fromPulse
	limitPulses int

	searchValuesPrint       bool
	searchValuesGreaterThan int64
}

func (app *appCtx) scanCommand() *cobra.Command {
	scan := &dbScanner{}
	defer func() {
		if scan.db != nil {
			scan.closeDB()
		}
	}()

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "scans database commands (check scan -h)",
	}
	showHistogramFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.showValuesHistogram, "histograms", false,
			"show key/values histograms")
	}
	nonStrictFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.nonStrict, "non-strict", false,
			"non strict mode (skip fail on some controversial error)")
	}
	perPulseStatFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.perPulseStat, "per-pulse", false,
			"show stat for every pulse")
	}
	graphSizesHistoryFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&scan.statGraph, "graph", "",
			"show pulse size graph over time")
	}
	progressBarFlag := func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&scan.limitPulses, "limit", "l", 0,
			"scan only provided count of pulses")
	}

	commonFlags := func(cmd *cobra.Command) {
		showHistogramFlag(cmd)
		nonStrictFlag(cmd)
		perPulseStatFlag(cmd)
		graphSizesHistoryFlag(cmd)
		progressBarFlag(cmd)
	}

	// TODO: add scan by scope id
	var scopeName string
	var scanScopeByNameCmd = &cobra.Command{
		Use:   "pulse-scan",
		Short: "scan in scope by pules and generate report",
		Run: func(_ *cobra.Command, _ []string) {
			scan.openDB(app.dataDir)
			scan.scanScopePulesByName(scopeName)
		},
	}
	scanScopeByNameCmd.Flags().StringVarP(&scopeName, "scope-name", "s", "ScopeRecord",
		"scan provided scope name")
	scanScopeByNameCmd.Flags().Int64Var(&scan.searchValuesGreaterThan, "value-greater", 0,
		"search values greater than provided bytes")
	scanScopeByNameCmd.Flags().BoolVar(&scan.enableRecordsTypesStat, "record-types-stat", false,
		"parse values in ScopeRecord and accumulate stat by records stat")
	commonFlags(scanScopeByNameCmd)

	var fastScan bool
	var names []string
	var scanAllScopes = &cobra.Command{
		Use:   "scope",
		Short: "show histograms by scope",
		Run: func(_ *cobra.Command, _ []string) {
			scan.openDB(app.dataDir)
			scan.scopesReport(fastScan, names)
		},
	}
	scanAllScopes.Flags().BoolVar(&fastScan, "fast", false, "scan and stat only keys")
	scanAllScopes.Flags().StringSliceVarP(&names, "scope-name", "s", nil, "scope name (by default scans all scopes)")

	var scanPulses = &cobra.Command{
		Use:   "pulses",
		Short: "collects stat per prefixes",
		Run: func(_ *cobra.Command, _ []string) {
			scan.openDB(app.dataDir)
			scan.pulsesReport()
		},
	}
	nonStrictFlag(scanPulses)

	scanCmd.AddCommand(
		scanPulses,
		scanScopeByNameCmd,
		scanAllScopes,
	)
	return scanCmd
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

func (dbs *dbScanner) getAllPulses() (pulses []insolar.Pulse) {
	if dbs.nonStrict {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("ERROR: getAllPulses", err)
			}
		}()
	}
	ctx := context.Background()

	var pulseStore = pulsedb.NewDB(dbs.db)
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

func (dbs *dbScanner) printPulsesInfo(pulses []insolar.Pulse) {
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
		printJSON(pv, jsonPrefix(""))
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

func (dbs *dbScanner) pulsesReport() {
	allPulses := dbs.getAllPulses()
	dbs.printPulsesInfo(allPulses)
}

func (dbs *dbScanner) scopesReport(fast bool, names []string) {
	var toScan []store.Scope
	for _, name := range names {
		scope, err := scopeFromName(name)
		if err != nil {
			fatalf("scan failed: %v", err)
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
	h := newHistogram("Summary Sizes")
	var opts = &iterOptions{
		counter:  true,
		keysOnly: fast,
	}

	printLine("=")
	fmt.Printf("Scan scope %s\n", scope)
	printLine("=")
	iterate(dbs.db, scopeKey{scope}, opts, h.iter)
	h.PrintKeys()
	if !fast {
		printLine("-")
		h.PrintValues()
	}
	fmt.Println()
}

func (dbs *dbScanner) scanScopePulesByName(scopeName string) {
	scope, err := scopeFromName(scopeName)
	if err != nil {
		fatalf("scan failed: %v", err)
	}
	fmt.Printf("Start scan in scope %s (%d)\n", scopeName, scope)
	dbs.scanScopePules(scope)
}

func (dbs *dbScanner) scanScopePules(scope store.Scope) {
	pulses := dbs.getAllPulses()
	dbs.printPulsesInfo(pulses)

	fmt.Printf("scan %s ...\n", scope.String())
	limit := len(pulses)
	if dbs.limitPulses > 0 && dbs.limitPulses < limit {
		limit = dbs.limitPulses
		fmt.Printf("scan %v of %v pulses\n", limit, len(pulses))
	}

	// var exceedIter func(k, v []byte) error
	var sumPrinters []printer

	baseIters := make([]iteration, 0, 8)
	if scope == store.ScopeRecord {
		if dbs.searchValuesGreaterThan > 0 {
			baseIters = append(baseIters, valuesExceedSize(scope, dbs.searchValuesGreaterThan))
		}
		if dbs.enableRecordsTypesStat {
			statByType := newRecordsStatByType()
			baseIters = append(baseIters, statByType.iter)
			sumPrinters = append(sumPrinters, statByType)
		}
	}

	if dbs.showValuesHistogram {
		sumValuesHistogram := newHistogram("Summary Sizes")
		sumPrinters = append(sumPrinters, sumValuesHistogram)
		baseIters = append(baseIters, sumValuesHistogram.iter)
	} else {
		sumKVS := &KVStat{Desc: "KV Summary Stat:"}
		sumPrinters = append(sumPrinters, sumKVS)
		baseIters = append(baseIters, sumKVS.iter)
	}

	var drawer Grapher
	adder := func(x insolar.Pulse, v int64) {
		if x.PulseNumber == insolar.GenesisPulse.PulseNumber {
			return
		}
		drawer.Add(x, float64(v)/mb)
	}

	switch dbs.statGraph {
	case "": // do nothing
		adder = func(x insolar.Pulse, v int64) {}
		drawer = StubDrawer{}
	case "console":
		drawer = &ConsoleGraph{}
	case "web":
		drawer = &webGraph{
			Title:       "Heavy Storage Consumption",
			DataHeaders: []string{"pulse", "record's values Mb"},
		}
	default:
		fatalf("unknown graph output type: %v", dbs.statGraph)
	}

	bar := pb.StartNew(limit)

	var i = 0
	for ; i < limit; i++ {
		pulse := pulses[len(pulses)-i-1]

		if !dbs.disableProgressbar {
			bar.Increment()
		}
		pn := pulse.PulseNumber

		var pulsePrinter printer
		iters := make([]iteration, 0, 8)
		for _, it := range baseIters {
			iters = append(iters, it)
		}

		// we need per pulse stat for graph
		add := func() {}
		if dbs.showValuesHistogram {
			valuesH := newHistogram("Per Pulse Sizes")
			add = func() { adder(pulse, valuesH.values.sum) }
			iters = append(iters, valuesH.iter)
			pulsePrinter = valuesH
		} else {
			kvs := &KVStat{}
			add = func() { adder(pulse, kvs.ValuesTotalBytes) }
			iters = append(iters, kvs.iter)
			pulsePrinter = kvs
		}

		// actual iteration happens here
		iterate(
			dbs.db,
			pulseKey{store.ScopeRecord, pn},
			nil,
			iters...,
		)
		add()

		if dbs.perPulseStat {
			fmt.Printf("Pulse %v stats.\n", pn)
			pulsePrinter.Print()
			fmt.Println()
		}
	}
	if dbs.limitPulses > 0 && i == dbs.limitPulses {
		fmt.Printf("LIMIT %v reached. STOP.\n", dbs.limitPulses)
	}

	if !dbs.disableProgressbar {
		bar.Finish()
	}

	for _, p := range sumPrinters {
		p.Print()
	}
	fmt.Println()

	drawer.Draw()
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

		extra := ""
		if scope == store.ScopeRecord {
			var matRec record.Material
			err := matRec.Unmarshal(v)
			if err != nil {
				panic(err)
			}
			virtual := record.Unwrap(&matRec.Virtual)
			extra = fmt.Sprintf(": Type=%T", virtual)
		}

		fmt.Printf(">> big value: size=%v (%s > %v) ID=%s: %s\n",
			valueSize, humanize.Bytes(uint64(valueSize)), size,
			insolar.NewIDFromBytes(k).String(), extra,
		)
		return nil
	}
}

type KVStat struct {
	Desc             string
	Count            int64
	KeysTotalBytes   int64
	ValuesTotalBytes int64
}

func (kvs *KVStat) iter(k, v []byte) error {
	keySize, valueSize := int64(len(k)), int64(len(v))
	kvs.KeysTotalBytes += keySize
	kvs.ValuesTotalBytes += valueSize
	kvs.Count++
	return nil
}

func (kvs *KVStat) Print() {
	if kvs.Desc != "" {
		fmt.Println(kvs.Desc)
	}
	fmt.Printf("found: %v, keys total size: %v, values total size: %v\n",
		kvs.Count,
		humanize.Bytes(uint64(kvs.KeysTotalBytes)),
		humanize.Bytes(uint64(kvs.ValuesTotalBytes)),
	)
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
