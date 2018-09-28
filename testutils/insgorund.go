package testutils

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// StartInsgorund starts `insgorund` process
func StartInsgorund(cmdPath string, listen string, upstreamAddr string) (func(), error) {
	var args []string
	if listen != "" {
		args = append(args, "-l", listen)
	} else {
		return nil, errors.New("listen is required to start `insgorund`")
	}
	if upstreamAddr != "" {
		args = append(args, "--rpc", upstreamAddr)
	} else {
		return nil, errors.New("address of the upstream is required to start `insgorund`")
	}

	if cmdPath == "" {
		return nil, errors.New("command's path is required to start `insgorund`")
	}

	runner := exec.Command(cmdPath, args...)
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err := runner.Start()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't start `insgorund`")
	}
	// XXX: dirty hack
	time.Sleep(200 * time.Millisecond)

	return func() {
		err := runner.Process.Signal(syscall.SIGINT)
		if err != nil {
			return
		}
		// XXX: dirty hack
		time.Sleep(200 * time.Millisecond)
	}, nil
}
