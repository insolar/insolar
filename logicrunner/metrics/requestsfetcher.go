// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	RequestFetcherFetchCall = stats.Int64(
		"vm_request_fetcher_fetch_call",
		"RequestFetcher fetch calls",
		stats.UnitDimensionless,
	)
	RequestFetcherFetchUnique = stats.Int64(
		"vm_request_fetcher_fetch_unique",
		"RequestFetcher fetch unique responses",
		stats.UnitDimensionless,
	)
	RequestFetcherFetchKnown = stats.Int64(
		"vm_request_fetcher_fetch_known",
		"RequestFetcher fetch known responses",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        RequestFetcherFetchCall.Name(),
			Description: RequestFetcherFetchCall.Description(),
			Measure:     RequestFetcherFetchCall,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        RequestFetcherFetchUnique.Name(),
			Description: RequestFetcherFetchUnique.Description(),
			Measure:     RequestFetcherFetchUnique,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        RequestFetcherFetchKnown.Name(),
			Description: RequestFetcherFetchKnown.Description(),
			Measure:     RequestFetcherFetchKnown,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
