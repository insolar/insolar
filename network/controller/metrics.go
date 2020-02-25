// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package controller

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMessageType = insmetrics.MustTagKey("messageType")
	tagPacketType  = insmetrics.MustTagKey("packetType")
)

var (
	statParcelsSentSizeBytes = stats.Int64(
		"network_parcels_sent_size",
		"size of sent parcels",
		stats.UnitBytes,
	)
	statParcelsReplySizeBytes = stats.Int64(
		"network_parcels_reply_size",
		"size of replies to parcels",
		stats.UnitBytes,
	)
	statPacketsReceived = stats.Int64(
		"network_packets_received",
		"number of received packets",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Measure:     statParcelsSentSizeBytes,
			Aggregation: view.Distribution(16, 32, 64, 128, 256, 512, 1024, 16*1<<10, 512*1<<10, 1<<20),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Measure:     statParcelsReplySizeBytes,
			Aggregation: view.Distribution(16, 32, 64, 128, 256, 512, 1024, 16*1<<10, 512*1<<10, 1<<20),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Measure:     statPacketsReceived,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagPacketType},
		},
	)
	if err != nil {
		panic(err)
	}
}
