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
	"fmt"
	"math"
	"sort"

	"github.com/insolar/insolar/reference"
)

const headerBinarySize = 3 * 4
const entryL0BinarySize = 32 + 8
const entryL1BinarySize = 32 + 8

func writeLocatorBuckets(wb *WriteBucketer, maxChapterBinarySize int) error {
	for bucketNo, max := 0, wb.BucketCount(); bucketNo < max; bucketNo++ {
		bucketContent := wb.GetBucketed(bucketNo)

		fmt.Println("=== Bucket", bucketNo, " ============================================= ")

		if len(bucketContent) == 0 {
			fmt.Println("Len	0")
			// TODO writeHeader(bucketNo, 0, 0)
			continue
		}

		entriesL0, entriesL1 := splitL0andL1(bucketContent, wb.keyMap.GetInterned)
		// TODO writeHeader(bucketNo, len(entriesL0), len(entriesL1))
		fmt.Println("Len	", len(bucketContent), "L0	", len(entriesL0), "L1	", len(entriesL1))

		chapterNo, indexL0, indexL1 := 0, 0, 0
		_ = writeChapters(entriesL0, entriesL1, maxChapterBinarySize-headerBinarySize,
			func(entriesL0, entriesL1 []resolvedEntry) error {
				chapterNo++
				fmt.Println("=== Bucket", bucketNo, " / Chapter", chapterNo, "L0", len(entriesL0), "L1", len(entriesL1), "============================= ")
				for _, entry := range entriesL0 {
					entry.localRef.GetPulseNumber()
					fmt.Println("	L0[", indexL0, "]	", entry.localRef.GetPulseNumber(), "[", entry.countL1, "] =>", entry.locator)
					indexL0++
				}

				for _, entry := range entriesL1 {
					entry.localRef.GetPulseNumber()
					fmt.Println("	L1[", indexL1, "]	", entry.baseRef.GetPulseNumber(), "=>", entry.locator)
					indexL1++
				}
				return nil
			},
		)
		fmt.Println()
	}
	return nil
}

func writeChapters(entriesL0, entriesL1 []resolvedEntry, maxChapterBinarySize int, writeFn func(entriesL0, entriesL1 []resolvedEntry) error) error {

	batchL0 := maxChapterBinarySize / entryL0BinarySize
	batchL1 := maxChapterBinarySize / entryL1BinarySize

	if len(entriesL0) > 0 {
		for {
			n := len(entriesL0)
			if n > batchL0 {
				if err := writeFn(entriesL0[:batchL0], nil); err != nil {
					return err
				}
				entriesL0 = entriesL0[batchL0:]
				continue
			}

			entriesL1portion := entriesL1
			switch remainingForL1 := (maxChapterBinarySize - n*entryL0BinarySize) / entryL1BinarySize; {
			case remainingForL1 >= len(entriesL1):
				entriesL1 = nil
			case remainingForL1 < MinBucketPageSize || remainingForL1 < batchL1>>2:
				entriesL1portion = nil
			default:
				subBatch := batchL1 >> 2

				batchL1 -= batchL1 % subBatch
				remainingForL1 -= remainingForL1 % subBatch

				entriesL1portion = entriesL1[:remainingForL1]
				entriesL1 = entriesL1[remainingForL1:]
			}

			if err := writeFn(entriesL0, entriesL1portion); err != nil {
				return err
			}
			break
		}
	}

	if len(entriesL1) > 0 {
		for {
			n := len(entriesL1)
			if n > batchL1 {
				if err := writeFn(nil, entriesL1[:batchL1]); err != nil {
					return err
				}
				entriesL1 = entriesL1[batchL1:]
				continue
			}

			if err := writeFn(nil, entriesL1); err != nil {
				return err
			}
			break
		}
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
	switch cmp := v[i].localRef.Compare(*v[j].localRef); {
	case cmp < 0:
		return true
	case cmp > 0:
		return false
	}
	return v[i].baseRef.Compare(*v[j].baseRef) < 0
}

func (v resolvedEntrySorter) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v resolvedEntrySorter) Len() int {
	return len(v)
}
