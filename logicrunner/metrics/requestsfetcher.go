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
