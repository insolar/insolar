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

package light

import (
	"context"
	"sync"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/require"
)

var testPassed = false
var testLock sync.Mutex

func TestComponents(t *testing.T) {
	// make sure test is not executed multiple times in parallel
	testLock.Lock()
	defer testLock.Unlock()
	if testPassed {
		// Dirty hack. This test doesn't work properly with -count 10
		// because it registers HTTP handlers. Sadly there is no way
		// to unregister HTTP handlers in net/http package.
		return
	}

	ctx := context.Background()
	cfg := configuration.NewConfiguration()
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.CertificatePath = "testdata/certificate.json"

	c, err := newComponents(ctx, cfg)
	require.NoError(t, err)
	err = c.Start(ctx)
	require.NoError(t, err)
	err = c.Stop(ctx)
	require.NoError(t, err)

	testPassed = true
}
