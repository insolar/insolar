//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package refmap

import (
	"math"
	"runtime"
	"sort"

	"github.com/insolar/insolar/reference"
)

func writeLocatorBuckets(wb *WriteBucketer) error {
	for bucketNo, max := 0, wb.BucketCount(); bucketNo < max; bucketNo++ {
		bucketContent := wb.GetBucketed(bucketNo)
		entriesL0, entriesL1 := splitL0andL1(bucketContent, wb.keyMap.GetInterned)
		runtime.KeepAlive(entriesL0)
		runtime.KeepAlive(entriesL1)
	}
	return nil
}

func splitL0andL1(selectors []ValueSelectorLocator, resolveFn BucketResolveFunc) ([]resolvedEntry, []resolvedEntry) {
	n := len(selectors)
	switch n {
	case 0:
		// empty bucket
		return nil, nil
	case 1:
		entry := resolveEntry(selectors[0], resolveFn)
		entries := []resolvedEntry{entry}
		if entry.countL1 == 0 /*leaf*/ {
			return entries, nil
		}
		return entries, entries
	}

	leafCount := 0
	entries := make([]resolvedEntry, n)
	for i, selector := range selectors {
		entry := resolveEntry(selector, resolveFn)
		entries[i] = entry
		if entry.countL1 == 0 {
			leafCount++
		}
	}
	sort.Sort(resolvedEntrySorter(entries))

	if leafCount == n {
		return entries, nil
	}
	nextL0idx := 0
	entriesL1 := make([]resolvedEntry, 0, n-leafCount)

	for i := 0; i != n; /* force panic on mismatched counts */ i++ {
		countL1 := entries[i].countL1
		if nextL0idx < i {
			entries[nextL0idx] = entries[i]
		}
		if countL1 == 0 {
			entries[nextL0idx].locator = 0
			nextL0idx++
			continue
		}

		entries[nextL0idx].locator = ValueLocator(len(entriesL1))
		nextL0idx++

		if countL1 == 1 {
			entriesL1 = append(entriesL1, entries[i])
			continue
		}

		base := i
		i += int(countL1 - 1) // consider loop increment
		entriesL1 = append(entriesL1, entries[base:i+1]...)
	}

	return entries[:nextL0idx], entriesL1
}

func resolveEntry(s ValueSelectorLocator, resolveFn BucketResolveFunc) resolvedEntry {
	if s.locator < 0 || s.locator > math.MaxInt64 {
		panic("illegal value")
	}
	entry := resolvedEntry{locator: s.locator}

	state := BucketState(0)
	entry.localRef, state = resolveFn(s.selector.LocalId)
	if s.selector.LocalId == s.selector.BaseId {
		entry.baseRef = entry.localRef
		if state == 1 {
			entry.countL1 = 0
			return entry
		}
	} else {
		entry.baseRef, _ = resolveFn(s.selector.BaseId)
	}
	entry.countL1 = uint32(state)
	return entry
}

type resolvedEntry struct {
	localRef *reference.Local
	baseRef  *reference.Local
	countL1  uint32
	locator  ValueLocator
}

type resolvedEntrySorter []resolvedEntry

func (v resolvedEntrySorter) Less(i, j int) bool {
	if !reference.LessLocal(v[i].localRef, v[j].localRef) {
		return false
	}
	return reference.LessLocal(v[i].baseRef, v[j].baseRef)
}

func (v resolvedEntrySorter) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v resolvedEntrySorter) Len() int {
	return len(v)
}
