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

package hot

import (
	"github.com/insolar/insolar/insolar"
)

type WriteAccessor interface {
	// Begin requests writing access for pulse number. If requested pulse is closed, ErrClosed will be returned.
	// The caller must call returned "done" function when finished writing.
	Begin(insolar.PulseNumber) (done func(), err error)
}

type WriteManager interface {
	// Open marks pulse number as opened for writing. It can be used later by Begin from accessor.
	Open(insolar.PulseNumber) error
	// CloseAndWait immediately marks pulse number as closed for writing and blocks until all writes are done.
	CloseAndWait(insolar.PulseNumber) error
}
