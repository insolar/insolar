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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/launchnet"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	os.Exit(launchnet.Run(func() int {
		err := setMigrationDaemonsRef()
		if err != nil {
			fmt.Println(errors.Wrap(err, "[ setup ] get reference daemons by public key failed ").Error())
		}

		var b bytes.Buffer
		buffer := bufio.NewWriter(&b)

		closeLogReader := testutils.NodesErrorLogReader(filepath.Join("..", defaults.LaunchnetDiscoveryNodesLogsDir()), buffer)
		defer close(closeLogReader)

		runResult := m.Run()
		if runResult > 0 {
			return runResult
		}

		// waiting few pulses for possible errors
		time.Sleep(30 * time.Second)

		// check logs for errors
		err = buffer.Flush()
		if err != nil {
			fmt.Println(errors.Wrap(
				err,
				"Can't flush buffer").Error())
		}

		// change Require.Error - make wrapper with save errors

		if b.Len() > 0 {
			_, err = b.WriteTo(os.Stdout)
			if err != nil {
				fmt.Println(errors.Wrap(
					err,
					"Discovery nodes contains errors, but there was an error while writing it to output ").Error())
			}

			return 1
		}
		return 0
	}))
}
