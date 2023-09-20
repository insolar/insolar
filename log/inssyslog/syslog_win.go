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
