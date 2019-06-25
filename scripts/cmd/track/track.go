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
