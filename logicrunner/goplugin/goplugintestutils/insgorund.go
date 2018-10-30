package goplugintestutils

import (
	"os"
	"os/exec"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

// StartInsgorund starts `insgorund` process
func StartInsgorund(cmdPath, lProto, listen, upstreamProto, upstreamAddr string) (func(), error) {
	id := testutils.RandomString()
	log.Debug("Starting 'insgorund' ", id)

	stackTrace := (string)(debug.Stack())
	cancelWarning := make(chan error, 1)
	go func() {
		select {
		case <-time.After(60 * time.Second):
			log.Error("'insgorund' is running for a minute, was started by:\n", stackTrace)
		case <-cancelWarning:
		}
	}()

	var args []string
	if listen != "" {
		args = append(args, "-l", listen)
	} else {
		return nil, errors.New("listen is required to start `insgorund`")
	}
	if lProto != "" {
		args = append(args, "--proto", lProto)
	}

	if upstreamAddr != "" {
		args = append(args, "--rpc", upstreamAddr)
	} else {
		return nil, errors.New("address of the upstream is required to start `insgorund`")
	}
	if upstreamProto != "" {
		args = append(args, "--rpc-proto", upstreamProto)
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
		log.Debug("stopping 'insgorund' ", id)

		close(cancelWarning)

		p := runner.Process
		err := p.Kill()
		if err != nil {
			log.Error("couldn't kill process: ", err)
		}

		// Wait for the process to finish or kill it after a timeout:
		done := make(chan error, 1)
		go func() {
			done <- runner.Wait()
		}()

		select {
		case <-time.After(3 * time.Second):
			log.Debug("waited for insgorund to finish and got tired")
			if err := p.Signal(syscall.SIGTERM); err != nil {
				log.Fatal("failed to terminate process: ", err)
			}
		case err := <-done:
			if err != nil {
				log.Debug("process finished, error: ", err)
			} else {
				log.Debug("process finished successfully")
			}
		}
	}, nil
}
