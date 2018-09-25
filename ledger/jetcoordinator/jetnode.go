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
	"bytes"

	"github.com/insolar/insolar/core"
)

type JetNode struct {
	ref   core.RecordRef
	left  *JetNode
	right *JetNode
}

func (jn *JetNode) GetContaining(objRef *core.RecordRef) *core.RecordRef {
	if jn.left == nil || jn.right == nil {
		return &jn.ref
	}

	// Ignore pulse number when selecting jet affinity. Object reference can be generated without knowing its pulse.
	if bytes.Compare(objRef[core.PulseNumberSize:], jn.ref[core.PulseNumberSize:]) < 0 {
		return jn.left.GetContaining(objRef)
	}
	return jn.right.GetContaining(objRef)
}
