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
