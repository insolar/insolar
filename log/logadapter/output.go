// Copyright 2020 Insolar Network Ltd.
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

package logadapter

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/inssyslog"
)

func OpenLogBareOutput(output insolar.LogOutput, param string) (BareOutput, error) {
	switch output {
	case insolar.StdErrOutput:
		w := os.Stderr
		return BareOutput{
			Writer:         w,
			FlushFn:        w.Sync,
			ProtectedClose: true,
		}, nil
	case insolar.SysLogOutput:
		executableName := filepath.Base(os.Args[0])
		w, err := inssyslog.ConnectSyslogByParam(param, executableName)
		if err != nil {
			return BareOutput{}, err
		}
		return BareOutput{
			Writer:         w,
			FlushFn:        w.Flush,
			ProtectedClose: false,
		}, nil
	default:
		return BareOutput{}, errors.New("unknown output " + output.String())
	}
}
