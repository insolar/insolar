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
	"fmt"
	"time"
)

type SignalStatistic struct {
	from int
	to   int
	id   int
	time time.Time
}

type Signal struct {
	from int
	id   int
}

func calculateLoad() {
	percentiles := make(map[int]map[int]time.Time)

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			fmt.Println()
			fmt.Println(" ----- start -------")
			for id, sygnalStat := range percentiles {
				percent := float64(len(sygnalStat)*100) / float64(len(network))
				var maxTime time.Time
				var minTime time.Time
				for _, t := range sygnalStat {
					if maxTime.IsZero() {
						maxTime = t
					}
					if minTime.IsZero() {
						minTime = t
					}

					if t.After(maxTime) {
						maxTime = t
					}
					if t.Before(minTime) {
						minTime = t
					}

				}
				fmt.Println(fmt.Sprintf("signal - %v, percent - %v, time - %v", id, percent, maxTime.Sub(minTime).Seconds()))
			}
			fmt.Println("-------------------")
		}
	}()

	for event := range commonBus {
		value, ok := percentiles[event.id]
		if !ok {
			percentiles[event.id] = make(map[int]time.Time)
			value = percentiles[event.id]
		}

		_, ok = value[event.to]
		if !ok {
			value[event.to] = event.time
		}
	}
}
