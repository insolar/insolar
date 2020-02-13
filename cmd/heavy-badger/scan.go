// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
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

	perPulseStat bool
	statGraph    bool

	disableProgressbar bool

	enableRecordsTypesStat bool

	limitPulses int
	skipPulses  int

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
	nonStrictFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.nonStrict, "non-strict", false,
			"non strict mode (skip fail on some controversial error)")
	}
	perPulseStatFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.perPulseStat, "per-pulse", false,
			"show stat for every pulse")
	}
	graphSizesHistoryFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.statGraph, "graph", false,
			"show per pulse db size graph (generates html in tmp dir and runs browser by open command)")
	}
	limitPulsesFlag := func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&scan.limitPulses, "limit", "l", 0,
			"scan only provided count of pulses")
	}
	skipPulsesFlag := func(cmd *cobra.Command) {
		cmd.Flags().IntVar(&scan.skipPulses, "skip", 0,
			"skip provided count of pulses")
	}
	disableProgressBarFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&scan.disableProgressbar, "no-progress-bar", false,
			"don't show pulses progress bar")
	}

	commonFlags := func(cmd *cobra.Command) {
		nonStrictFlag(cmd)
		perPulseStatFlag(cmd)
		graphSizesHistoryFlag(cmd)
		limitPulsesFlag(cmd)
		skipPulsesFlag(cmd)
		disableProgressBarFlag(cmd)
	}

	var scopeName string
	var scanScopeByNameCmd = &cobra.Command{
		Use:   "scope-pulses",
		Short: "scan in scope by pules and generate report",
		Run: func(_ *cobra.Command, _ []string) {
			if scan.perPulseStat || scan.searchValuesGreaterThan > 0 {
				scan.disableProgressbar = true
			}
			scan.openDB(app.dataDir)
			scan.scanScopePulesByName(scopeName)
		},
	}
	scanScopeByNameCmd.Flags().StringVarP(&scopeName, "scope-name", "s", "ScopeRecord",
		"scope name")
	scanScopeByNameCmd.Flags().Int64Var(&scan.searchValuesGreaterThan, "print-value-gt-size", 0,
		"search values greater than provided bytes")
	scanScopeByNameCmd.Flags().BoolVar(&scan.enableRecordsTypesStat, "record-types-stat", false,
		"parse values in ScopeRecord and accumulate stat by records stat")
	commonFlags(scanScopeByNameCmd)

	var fastScan bool
	var names []string
	var ids []int
	var scopesStatCmd = &cobra.Command{
		Use:   "scopes-stat",
		Short: "show statistic by scope (by default scans all scopes)",
		Run: func(_ *cobra.Command, _ []string) {
			scan.openDB(app.dataDir)
			scan.scopesReport(fastScan, names, ids)
		},
	}
	scopesStatCmd.Flags().BoolVar(&fastScan, "fast", false, "scan and stat only keys")
	scopesStatCmd.Flags().StringSliceVarP(&names, "scope-name", "s", nil,
		"scope name")
	scopesStatCmd.Flags().IntSliceVarP(&ids, "scope-id", "i", nil,
		"scope id (check command 'scopes' output)")
	graphSizesHistoryFlag(scopesStatCmd)

	var printAll bool
	var scanPulses = &cobra.Command{
		Use:   "pulses",
		Short: "report on pulses",
		Run: func(_ *cobra.Command, _ []string) {
			scan.openDB(app.dataDir)
			scan.pulsesReport(printAll)
		},
	}
	nonStrictFlag(scanPulses)
	scanPulses.Flags().BoolVar(&printAll, "all", false, "print all pulses")

	scanCmd.AddCommand(
		scanPulses,
		scanScopeByNameCmd,
		scopesStatCmd,
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

	var pulseStore = pulsedb.NewBadgerDB(dbs.db)
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

func (dbs *dbScanner) pulsesReport(all bool) {
	allPulses := dbs.getAllPulses()
	if all {
		dbs.printAllPulsesInfo(allPulses)
		return
	}
	dbs.printShortPulsesInfo(allPulses)
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
	dbs.printShortPulsesInfo(pulses)

	fmt.Printf("scan %s ...\n", scope.String())
	if dbs.skipPulses > 0 {
		if dbs.skipPulses >= len(pulses) {
			fmt.Printf("Limit %v exceeds pulses count %v. Nothing to scan.\n", dbs.skipPulses, len(pulses))
			return
		}
		// note: we expect pulses are in reverse order here!
		pulses = pulses[:len(pulses)-dbs.skipPulses]
	}

	limit := len(pulses)
	if dbs.limitPulses > 0 && dbs.limitPulses < limit {
		limit = dbs.limitPulses
		fmt.Printf("scan %v of %v pulses\n", limit, len(pulses))
	}

	var sumPrinters []printer

	baseIters := make([]iteration, 0, 8)
	if dbs.searchValuesGreaterThan > 0 {
		baseIters = append(baseIters, valuesExceedSize(scope, dbs.searchValuesGreaterThan))
	}
	if scope == store.ScopeRecord {
		if dbs.enableRecordsTypesStat {
			statByType := newRecordsStatByType()
			baseIters = append(baseIters, statByType.iter)
			sumPrinters = append(sumPrinters, statByType)
		}
	}

	sumValuesHistogram := newHistogram("Summary Sizes")
	sumPrinters = append(sumPrinters, sumValuesHistogram)
	baseIters = append(baseIters, sumValuesHistogram.iter)
	pulsePrinter := func() {}

	var graphImpl Grapher = StubDrawer{}
	statGraphAdd := func() {}
	if dbs.statGraph {
		graphImpl = &webGraph{
			Title:       "Heavy Storage Consumption",
			DataHeaders: []string{"pulse", "record's values Mb"},
		}
	}

	bar := createProgressBar(limit, dbs.disableProgressbar)

	var i = 0
	for ; i < limit; i++ {
		bar.Increment()

		pulse := pulses[len(pulses)-i-1]
		pn := pulse.PulseNumber

		iters := make([]iteration, 0, len(baseIters)+2)
		iters = append(iters, baseIters...)

		var valuesH *histogram
		if dbs.statGraph || dbs.perPulseStat {
			valuesH = newHistogram("Per Pulse Sizes")
			iters = append(iters, valuesH.iter)
			if dbs.statGraph {
				statGraphAdd = func() {
					graphImpl.Add(pulse, float64(valuesH.values.sum)/mb)
				}
			}
		}
		if dbs.perPulseStat {
			pulsePrinter = func() {
				fmt.Printf("Pulse %v stats.\n", pn)
				printLine("-")
				valuesH.Print()
				fmt.Println()
			}
		}

		// actual iteration happens here
		iterate(
			dbs.db,
			pulseKey{store.ScopeRecord, pn},
			nil,
			iters...,
		)
		statGraphAdd()
		pulsePrinter()
	}
	bar.Finish()

	if dbs.limitPulses > 0 && i == dbs.limitPulses {
		fmt.Printf("LIMIT %v reached. STOP.\n", dbs.limitPulses)
	}

	// reports
	for _, p := range sumPrinters {
		p.Print()
	}
	fmt.Printf("\n%v keys in %v pulses scanned.\n\n", sumValuesHistogram.keys.totalCount, i)
	graphImpl.Draw()
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
