// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build windows

package inssyslog

import (
	"github.com/pkg/errors"
)

func ConnectDefaultSyslog(tag string) (LogLevelWriteCloser, error) {
	return nil, errors.New("not implemented for Windows")
}

func ConnectRemoteSyslog(network, raddr string, tag string) (LogLevelWriteCloser, error) {
	return nil, errors.New("not implemented for Windows")
}
