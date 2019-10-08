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

package proc

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statHotsAbandoned = stats.Int64(
		"requests_abandoned",
		"How many abandoned requests in hot data",
		stats.UnitDimensionless,
	)
	statSetRequestTotal = stats.Int64(
		"proc_set_request_total",
		"How many requests have been set",
		stats.UnitDimensionless,
	)
	statSetRequestDuplicate = stats.Int64(
		"proc_set_request_duplicate",
		"How many requests have been duplicated",
		stats.UnitDimensionless,
	)
	statSetRequestSuccess = stats.Int64(
		"proc_set_request_success",
		"How many requests have been saved successfully",
		stats.UnitDimensionless,
	)
	statSetRequestError = stats.Int64(
		"proc_set_request_error",
		"How many requests have been saved successfully",
		stats.UnitDimensionless,
	)
	statSetResultTotal = stats.Int64(
		"proc_set_result_total",
		"How many results have been set",
		stats.UnitDimensionless,
	)
	statSetResultDuplicate = stats.Int64(
		"proc_set_result_duplicate",
		"How many results have been duplicated",
		stats.UnitDimensionless,
	)
	statSetResultError = stats.Int64(
		"proc_set_result_error",
		"How many results finished with errors",
		stats.UnitDimensionless,
	)
	statSetResultSuccess = stats.Int64(
		"proc_set_result_success",
		"How many results have been saved successfully",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statHotsAbandoned.Name(),
			Description: statHotsAbandoned.Description(),
			Measure:     statHotsAbandoned,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetRequestTotal.Name(),
			Description: statSetRequestTotal.Description(),
			Measure:     statSetRequestTotal,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetRequestSuccess.Name(),
			Description: statSetRequestSuccess.Description(),
			Measure:     statSetRequestSuccess,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetRequestError.Name(),
			Description: statSetRequestError.Description(),
			Measure:     statSetRequestError,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetRequestDuplicate.Name(),
			Description: statSetRequestDuplicate.Description(),
			Measure:     statSetRequestDuplicate,
			Aggregation: view.Count(),
		},

		&view.View{
			Name:        statSetResultTotal.Name(),
			Description: statSetResultTotal.Description(),
			Measure:     statSetResultTotal,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetResultSuccess.Name(),
			Description: statSetResultSuccess.Description(),
			Measure:     statSetResultSuccess,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetResultError.Name(),
			Description: statSetResultError.Description(),
			Measure:     statSetResultError,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statSetResultDuplicate.Name(),
			Description: statSetResultDuplicate.Description(),
			Measure:     statSetResultDuplicate,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
