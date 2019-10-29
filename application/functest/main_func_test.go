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

// +build functest

package functest

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestMain(m *testing.M) {
	os.Exit(launchnet.Run(func() int {
		err := setMigrationDaemonsRef()
		if err != nil {
			fmt.Println(errors.Wrap(err, "[ setup ] get reference daemons by public key failed ").Error())
		}

		return m.Run()
	}))
}
