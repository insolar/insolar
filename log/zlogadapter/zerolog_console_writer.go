// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package zlogadapter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/insolar/insolar/log/logadapter"
	"github.com/rs/zerolog"
)

var fieldsOrder = []string{
	zerolog.TimestampFieldName,
	zerolog.LevelFieldName,
	zerolog.MessageFieldName,
	zerolog.CallerFieldName,
}

var _ io.WriteCloser = &closableConsoleWriter{}

type closableConsoleWriter struct {
	zerolog.ConsoleWriter
}

func (p *closableConsoleWriter) Close() error {
	if c, ok := p.Out.(io.Closer); ok {
		return c.Close()
	}
	return errors.New("unsupported: Close")
}

func (p *closableConsoleWriter) Sync() error {
	if c, ok := p.Out.(*os.File); ok {
		return c.Sync()
	}
	return errors.New("unsupported: Sync")
}

func newDefaultTextOutput(out io.Writer) io.WriteCloser {
	return &closableConsoleWriter{zerolog.ConsoleWriter{
		Out:          out,
		NoColor:      true,
		TimeFormat:   logadapter.TimestampFormat,
		PartsOrder:   fieldsOrder,
		FormatCaller: formatCaller(),
	}}
}

func formatCaller() zerolog.Formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if len(cwd) > 0 {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
			c = "file=" + c
		}
		return c
	}
}

var cwd string

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		cwd = ""
		fmt.Println("couldn't get current working directory: ", err.Error())
	}
}
