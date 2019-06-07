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

package executor

import (
	"testing"

	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/stretchr/testify/require"
)

func TestJetSplitter_New(t *testing.T) {
	jc := jet.NewCoordinatorMock(t)
	js := jet.NewStorageMock(t)
	da := drop.NewAccessorMock(t)
	rsp := recentstorage.NewProviderMock(t)
	splitter := NewJetSplitter(
		jc,
		js,
		js,
		da,
		rsp,
	)
	require.NotNil(t, splitter, "jet splitter created")
}
