/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func badgerCollector() prometheus.Collector {
	exports := map[string]*prometheus.Desc{}
	metricnames := []string{
		"badger_disk_reads_total",
		"badger_disk_writes_total",
		"badger_read_bytes",
		"badger_written_bytes",
		"badger_lsm_level_gets_total",
		"badger_lsm_bloom_hits_total",
		"badger_gets_total",
		"badger_puts_total",
		"badger_blocked_puts_total",
		"badger_memtable_gets_total",
		"badger_lsm_size_bytes",
		"badger_vlog_size_bytes",
		"badger_pending_writes_total",
	}
	for _, name := range metricnames {
		exports[name] = prometheus.NewDesc(
			name,
			name,
			nil, nil,
		)
	}
	return prometheus.NewExpvarCollector(exports)
}
