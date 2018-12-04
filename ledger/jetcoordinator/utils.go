/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package jetcoordinator

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/utils/entropy"
)

func selectByEntropy(
	scheme core.PlatformCryptographyScheme,
	e core.Entropy,
	values []core.RecordRef,
	count int,
) ([]core.RecordRef, error) { // nolint: megacheck
	in := make([][]byte, 0, len(values))
	for _, value := range values {
		in = append(in, value[:])
	}
	res, err := entropy.SelectByEntropy(scheme, e[:], in, count)
	if err != nil {
		return nil, err
	}
	out := make([]core.RecordRef, 0, len(res))
	for _, value := range res {
		var coreref core.RecordRef
		copy(coreref[:], value)
		out = append(out, coreref)
	}
	return out, nil
}
