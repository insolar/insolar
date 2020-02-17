// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	TagContractMethodName = insmetrics.MustTagKey("vm_execution_contract_method_name")
	TagContractPrototype  = insmetrics.MustTagKey("vm_execution_contract_prototype")
)

var (
	ContractExecutionTime = stats.Float64(
		"vm_contracts_timing",
		"time spent executing contract",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        ContractExecutionTime.Name(),
			Description: ContractExecutionTime.Description(),
			Measure:     ContractExecutionTime,
			TagKeys:     []tag.Key{TagContractMethodName, TagContractPrototype},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}
