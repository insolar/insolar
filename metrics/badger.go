// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func badgerCollector(namespace string) prometheus.Collector {
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
		exportname := name
		if exportname != "" {
			exportname = namespace + "_" + exportname
		}
		exports[name] = prometheus.NewDesc(
			exportname,
			"badger db metric "+name,
			nil, nil,
		)
	}
	return prometheus.NewExpvarCollector(exports)
}
