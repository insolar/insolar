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

package requestsqueue

import (
	"context"

	"github.com/insolar/insolar/logicrunner/common"
)

// SplitNFromMany utility function that extracts from multiple queues slits to no more than N elements and rest.
func SplitNFromMany(ctx context.Context, n int, qs ...RequestsQueue) ([]*common.Transcript, []*common.Transcript) {
	res := make([]*common.Transcript, 0, n)
	rest := []*common.Transcript(nil)
	for i := 0; i < int(numberOfSources); i++ {
		for _, q := range qs {
			if len(res) < n {
				res = append(res, q.TakeAllOriginatedFrom(ctx, RequestSource(i))...)
			} else {
				rest = append(rest, q.TakeAllOriginatedFrom(ctx, RequestSource(i))...)
			}
		}
	}

	return res, rest
}
