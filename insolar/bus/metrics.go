// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bus

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	kb = 1 << 10
	mb = 1 << 20
)

var (
	tagMessageType = insmetrics.MustTagKey("message_type")
	tagMessageRole = insmetrics.MustTagKey("message_role")
)

var (
	statSentBytes = stats.Int64(
		"bus_sent",
		"sent messages stats",
		stats.UnitDimensionless,
	)
	statSentTime = stats.Float64(
		"bus_sent_latency",
		"time spent on sending parcels",
		stats.UnitMilliseconds,
	)
	statRetries = stats.Int64(
		"bus_sent_retries",
		"retries on send messages",
		stats.UnitDimensionless,
	)

	statReply = stats.Int64(
		"bus_reply",
		"reply messages stats",
		stats.UnitDimensionless,
	)
	statReplyTimeouts = stats.Int64(
		"bus_reply_timeouts",
		"reply messages stats",
		stats.UnitDimensionless,
	)
	statReplyTime = stats.Float64(
		"bus_reply_latency",
		"time spent on sending parcels",
		stats.UnitMilliseconds,
	)
	statReplyError = stats.Int64(
		"bus_reply_error",
		"reply error messages stats",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		// sent stats
		&view.View{
			Name:        "bus_sent_total",
			Description: "sent messages total count",
			Measure:     statSentBytes,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagMessageType, tagMessageRole},
		},
		&view.View{
			Name:        "bus_sent_bytes",
			Description: "sent messages payload size",
			Measure:     statSentBytes,
			Aggregation: view.Distribution(1*kb, 10*kb, 100*kb, 1*mb, 10*mb, 100*mb),
			TagKeys:     []tag.Key{tagMessageType, tagMessageRole},
		},
		&view.View{
			Name:        "bus_sent_milliseconds",
			Description: "sent messages latency",
			Measure:     statSentTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Name:        "bus_sent_retries",
			Description: "sent messages retries count",
			Measure:     statRetries,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{tagMessageType},
		},
		// reply stats
		&view.View{
			Name:        "bus_reply_total",
			Description: "reply messages total count",
			Measure:     statReply,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Name:        "bus_reply_milliseconds",
			Description: "reply messages latency",
			Measure:     statReplyTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Name:        "bus_reply_timeouts",
			Description: "reply messages total count",
			Measure:     statReplyTimeouts,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Name:        "bus_reply_error_total",
			Description: "reply error messages total count",
			Measure:     statReplyError,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
