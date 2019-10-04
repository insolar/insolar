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
	"fmt"
	"math"

	"github.com/dustin/go-humanize"
)

// histogramData stores information about a histogram
type histogramData struct {
	bins        []int64
	countPerBin []int64
	sumPerBin   []int64
	totalCount  int64
	min         int64
	max         int64
	sum         int64
}

func newKeyHistogram() *histogramData {
	keyBins := createHistogramBins(1, 16)
	return &histogramData{
		bins:        keyBins,
		countPerBin: make([]int64, len(keyBins)+1),
		sumPerBin:   make([]int64, len(keyBins)+1),
		max:         math.MinInt64,
		min:         math.MaxInt64,
		sum:         0,
	}
}

func newValueHistogram() *histogramData {
	valueBins := createHistogramBins(1, 30)
	return &histogramData{
		bins:        valueBins,
		countPerBin: make([]int64, len(valueBins)+1),
		sumPerBin:   make([]int64, len(valueBins)+1),
		max:         math.MinInt64,
		min:         math.MaxInt64,
		sum:         0,
	}
}

// createHistogramBins creates bins for an histogram. The bin sizes are powers
// of two of the form [2^min_exponent, ..., 2^max_exponent].
func createHistogramBins(minExponent, maxExponent uint32) []int64 {
	var bins []int64
	for i := minExponent; i <= maxExponent; i++ {
		bins = append(bins, int64(1)<<i)
	}
	return bins
}

// Update the min and max fields if value is less than or greater than the
// current min/max value.
func (h *histogramData) Update(value int64) {
	if value > h.max {
		h.max = value
	}
	if value < h.min {
		h.min = value
	}

	h.sum += value
	h.totalCount++

	for index := 0; index <= len(h.bins); index++ {
		// Allocate value in the last buckets if we reached the end of the Bounds array.
		if index == len(h.bins) {
			h.countPerBin[index]++
			h.sumPerBin[index] += value
			break
		}

		// Check if the value should be added to the "index" bin
		if value < h.bins[index] {
			h.countPerBin[index]++
			h.sumPerBin[index] += value
			break
		}
	}
}

// printHistogram prints the histogram data in a human-readable format.
func (h histogramData) Print() {
	fmt.Printf("Total count: %d\n", h.totalCount)
	if h.totalCount == 0 {
		return
	}
	fmt.Printf("Min value: %d\n", h.min)
	fmt.Printf("Max value: %d\n", h.max)
	fmt.Printf("Mean: %.2f\n", float64(h.sum)/float64(h.totalCount))
	fmt.Printf("%24s %9s %12s\n", "Range", "Count", "Sum")

	numBins := len(h.bins)
	for index, count := range h.countPerBin {
		if count == 0 {
			continue
		}

		// The last bin represents the bin that contains the range from
		// the last bin up to infinity so it's processed differently than the
		// other bins.
		if index == len(h.countPerBin)-1 {
			lowerBound := int(h.bins[numBins-1])
			fmt.Printf("[%10d, %10s) %9d\n", lowerBound, "infinity", count)
			continue
		}

		upperBound := int(h.bins[index])
		lowerBound := 0
		if index > 0 {
			lowerBound = int(h.bins[index-1])
		}

		fmt.Printf("[%10d, %10d) %9d %12s\n",
			lowerBound, upperBound, count, humanize.Bytes(uint64(h.sumPerBin[index])))
	}
	fmt.Println()
}

func newHistogram(description string) *histogram {
	return &histogram{
		descr:  description,
		keys:   newKeyHistogram(),
		values: newValueHistogram(),
	}
}

type histogram struct {
	descr  string
	keys   *histogramData
	values *histogramData
}

func (h *histogram) iter(k, v []byte) error {
	keySize := int64(len(k))
	valueSize := int64(len(v))
	h.keys.Update(keySize)
	h.values.Update(valueSize)
	return nil
}

func (h *histogram) Print() {
	h.PrintKeys()
	h.PrintValues()
}

func (h *histogram) PrintKeys() {
	fmt.Println("KEYS", h.descr)
	h.keys.Print()
	fmt.Printf("total size: %v\n", humanize.Bytes(uint64(h.keys.sum)))
}

func (h *histogram) PrintValues() {
	fmt.Println("VALUES", h.descr)
	h.values.Print()
	fmt.Printf("total size: %v\n", humanize.Bytes(uint64(h.values.sum)))
}
