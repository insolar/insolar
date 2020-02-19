// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
