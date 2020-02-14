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
