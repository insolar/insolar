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

package logadapter

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/inssyslog"
	"github.com/pkg/errors"
	"os"
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
		w, err := inssyslog.ConnectSyslogByParam(param, "insolar")
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
