/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

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
		"network/parcels/sent/size",
		"size of sent parcels",
		stats.UnitBytes,
	)
	statParcelsReplySizeBytes = stats.Int64(
		"network/parcels/reply/size",
		"size of replies to parcels",
		stats.UnitBytes,
	)
	statPacketsReceived = stats.Int64(
		"network/packets/received",
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
