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

package handle

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	KeyProcName  = insmetrics.MustTagKey("proc_name")
	KeyErrorCode = insmetrics.MustTagKey("error_code")
)

var (
	statProcError = stats.Int64(
		"proc_errors",
		"How many procedures return errors",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statProcError.Name(),
			Description: statProcError.Description(),
			Measure:     statProcError,
			TagKeys:     []tag.Key{KeyProcName, KeyErrorCode},
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
