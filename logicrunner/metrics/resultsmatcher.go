// Copyright 2020 Insolar Network Ltd.
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

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	ResultsMatcherAddStillExecution = stats.Int64(
		"vm_results_matcher_still_execution",
		"AddStillExecution calls",
		stats.UnitDimensionless,
	)
	ResultsMatcherAddUnwantedResponse = stats.Int64(
		"vm_results_matcher_unwanted_response",
		"AddUnwantedResponse calls",
		stats.UnitDimensionless,
	)

	ResultsMatcherSentResults = stats.Int64(
		"vm_results_matcher_sent_results",
		"sent results to executor",
		stats.UnitDimensionless,
	)
	ResultsMatcherDroppedResults = stats.Int64(
		"vm_results_matcher_dropped_results",
		"dropped results in pulse",
		stats.UnitDimensionless,
	)
	ResultMatcherLoopDetected = stats.Int64(
		"vm_results_matcher_loop_detected",
		"unwanted results that loops between node/nodes",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        ResultsMatcherAddStillExecution.Name(),
			Description: ResultsMatcherAddStillExecution.Description(),
			Measure:     ResultsMatcherAddStillExecution,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ResultsMatcherAddUnwantedResponse.Name(),
			Description: ResultsMatcherAddUnwantedResponse.Description(),
			Measure:     ResultsMatcherAddUnwantedResponse,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ResultsMatcherSentResults.Name(),
			Description: ResultsMatcherSentResults.Description(),
			Measure:     ResultsMatcherSentResults,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ResultsMatcherDroppedResults.Name(),
			Description: ResultsMatcherDroppedResults.Description(),
			Measure:     ResultsMatcherDroppedResults,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ResultMatcherLoopDetected.Name(),
			Description: ResultMatcherLoopDetected.Description(),
			Measure:     ResultMatcherLoopDetected,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
