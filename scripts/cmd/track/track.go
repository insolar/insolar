// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

var lineRegex = regexp.MustCompile(`^([^ :]+):\d+:([^ ]*) `)
var ErrNotMatch = errors.New("line does no match")

type line struct {
	file   string
	time   string
	source string
}

func main() {
	bufIN := bufio.NewReader(os.Stdin)
	bufOUT := bufio.NewWriter(os.Stdout)

	var lines []line
	for {
		text, err := bufIN.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		ln, err := parse(text)
		if err != nil {
			panic(err)
		}
		lines = append(lines, ln)
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].time < lines[j].time
	})

	var currentFile string
	for _, ln := range lines {
		if ln.file != currentFile {
			_, err := bufOUT.WriteString("\n" + ln.file + "\n")
			if err != nil {
				panic(err)
			}
			currentFile = ln.file
		}
		_, err := bufOUT.WriteString(ln.source)
		if err != nil {
			panic(err)
		}
		bufOUT.Flush()
	}
}

func parse(ln string) (line, error) {
	g := lineRegex.FindStringSubmatch(ln)

	if len(g) < 3 {
		return line{}, ErrNotMatch
	}

	return line{file: g[1], time: g[2], source: ln}, nil
}
