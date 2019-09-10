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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const logDir = ".artifacts/launchnet/logs/"
const statLogMessage = "stat_log_message"
const typeSent = "sent"
const typeReply = "reply"

var pattern = regexp.MustCompile(".*output.log$")

var ignoreTracePrefixes = []string{
	"main",
	"pulse",
}

type StatLog struct {
	StatType    string  `json:"stat_type"`
	TraceID     string  `json:"traceid"`
	Message     string  `json:"message"`
	MessageType string  `json:"message_type"`
	ReplyTimeMS float32 `json:"reply_time_ms"`
}

type Stats struct {
	lock sync.RWMutex
	// Trace id to stat.
	stats map[string]*TraceStats
}

func NewStats() *Stats {
	return &Stats{stats: map[string]*TraceStats{}}
}

func (s *Stats) GetOrCreate(trace string) *TraceStats {
	s.lock.Lock()
	defer s.lock.Unlock()

	if stat, ok := s.stats[trace]; ok {
		return stat
	}

	stat := NewTraceStats()
	s.stats[trace] = stat
	return stat
}

type TraceStats struct {
	sync.RWMutex
	First, Last time.Time
	// Message type to reply times.
	ReplyTimings map[string][]float32
	// Message type to sent count.
	SentCounts map[string]uint64
}

func NewTraceStats() *TraceStats {
	return &TraceStats{
		ReplyTimings: map[string][]float32{},
		SentCounts:   map[string]uint64{},
	}
}

func main() {
	shouldParse := func(log StatLog) bool {
		if log.Message != statLogMessage {
			return false
		}
		if log.TraceID == "" {
			return false
		}
		for _, i := range ignoreTracePrefixes {
			if strings.HasPrefix(log.TraceID, i) {
				return false
			}
		}

		return true
	}

	parseFile := func(stats *Stats, filename string) {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log := StatLog{}
			err = json.Unmarshal(scanner.Bytes(), &log)
			if err != nil {
				continue
			}

			if !shouldParse(log) {
				continue
			}

			parseLog(log, stats.GetOrCreate(log.TraceID))
		}
	}

	var files []string
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !pattern.MatchString(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	checkError(err)

	stats := NewStats()
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, path := range files {
		path := path
		go func() {
			parseFile(stats, path)
			wg.Done()
		}()
	}
	wg.Wait()

	aggregate := Aggregate{
		Sent:  aggregateSent(stats),
		Reply: aggregateReply(stats),
	}
	// aggregateJSON, err := json.Marshal(aggregate)
	// checkError(err)
	out := bufio.NewWriter(os.Stdout)
	_, err = out.Write([]byte(aggregate.String()))
	checkError(err)
	err = out.Flush()
	checkError(err)
}

func parseLog(log StatLog, stat *TraceStats) {
	stat.Lock()
	defer stat.Unlock()

	switch log.StatType {
	case typeSent:
		stat.SentCounts[log.MessageType] += 1
	case typeReply:
		stat.ReplyTimings[log.MessageType] = append(stat.ReplyTimings[log.MessageType], log.ReplyTimeMS)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Aggregate struct {
	Sent  AggSent
	Reply AggReply
}

func (a *Aggregate) String() string {
	b := strings.Builder{}
	b.WriteString("[Sent]\n")
	b.WriteString(a.Sent.String())
	b.WriteString("\n")
	b.WriteString("[Reply]\n")
	b.WriteString(a.Reply.String())
	b.WriteString("\n")

	return b.String()
}

type AggSent struct {
	// Message type to count.
	AVGCount map[string]float64 `json:"avg_count"`
}

func (a *AggSent) String() string {
	b := strings.Builder{}

	b.WriteString("[Average count per trace]\n")
	for msgType, count := range a.AVGCount {
		b.WriteString("    " + msgType + ": " + fmt.Sprintf("%f", count) + "\n")
	}

	return b.String()
}

func aggregateSent(logStats *Stats) AggSent {
	logStats.lock.RLock()
	defer logStats.lock.RUnlock()

	var aggs []AggSent
	for _, stats := range logStats.stats {
		agg := NewAggSent()
		for msgType, count := range stats.SentCounts {
			agg.AVGCount[msgType] = float64(count)
		}
		aggs = append(aggs, agg)
	}

	totals := NewAggSent()
	for _, agg := range aggs {
		for msgType, count := range agg.AVGCount {
			totals.AVGCount[msgType] += count
		}
	}

	avgDivider := len(aggs)
	for msgType, _ := range totals.AVGCount {
		totals.AVGCount[msgType] /= float64(avgDivider)
	}

	return totals
}

func NewAggSent() AggSent {
	return AggSent{AVGCount: map[string]float64{}}
}

type AggReply struct {
	// Message type to reply time.
	AVGReplyTime map[string]float64 `json:"avg_reply_time"`
}

func NewAggReply() AggReply {
	return AggReply{AVGReplyTime: map[string]float64{}}
}

func (a *AggReply) String() string {
	b := strings.Builder{}

	b.WriteString("[Average reply times per trace, ms]\n")
	for msgType, replyTime := range a.AVGReplyTime {
		b.WriteString("    " + msgType + ": " + fmt.Sprintf("%f", replyTime) + "\n")
	}

	return b.String()
}

func aggregateReply(logStats *Stats) AggReply {
	logStats.lock.RLock()
	defer logStats.lock.RUnlock()

	var aggs []AggReply
	for _, stats := range logStats.stats {
		agg := NewAggReply()
		for msgType, timings := range stats.ReplyTimings {
			var summ float64
			for _, t := range timings {
				summ += float64(t)
			}
			agg.AVGReplyTime[msgType] = summ
		}
		aggs = append(aggs, agg)
	}

	totals := NewAggReply()
	for _, agg := range aggs {
		for msgType, replyTime := range agg.AVGReplyTime {
			totals.AVGReplyTime[msgType] += replyTime
		}
	}

	avgDivider := len(aggs)
	for msgType, _ := range totals.AVGReplyTime {
		totals.AVGReplyTime[msgType] /= float64(avgDivider)
	}

	return totals
}
