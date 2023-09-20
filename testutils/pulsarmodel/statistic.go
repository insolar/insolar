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
