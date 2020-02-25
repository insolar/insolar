//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"bytes"
	"os"
	"time"
)

// Capturer has flags whether capture stdout/stderr or not.
type Capturer struct {
	captureStdout bool
	captureStderr bool
}

// CaptureStdout captures stdout.
func CaptureStdout(f func()) (string, error) {
	capturer := &Capturer{captureStdout: true}
	return capturer.capture(f)
}

// CaptureStderr captures stderr.
func CaptureStderr(f func()) (string, error) {
	capturer := &Capturer{captureStderr: true}
	return capturer.capture(f)
}

// CaptureOutput captures stdout and stderr.
func CaptureOutput(f func()) (string, error) {
	capturer := &Capturer{captureStdout: true, captureStderr: true}
	return capturer.capture(f)
}

func (capturer *Capturer) capture(fn func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err.Error())
	}

	if capturer.captureStdout {
		stdout := os.Stdout
		os.Stdout = w
		defer func() {
			os.Stdout = stdout
		}()
	}

	if capturer.captureStderr {
		stderr := os.Stderr
		os.Stderr = w
		defer func() {
			os.Stderr = stderr
		}()
	}

	defer w.Close()

	var buf bytes.Buffer
	var retErr error
	go func() {
		defer r.Close()
		_, e := buf.ReadFrom(r)
		if e != nil {
			retErr = err
			return
		}
	}()

	fn()

	// here we need sleep because not of all pipes closed
	time.Sleep(time.Millisecond * 10)
	return buf.String(), retErr
}
