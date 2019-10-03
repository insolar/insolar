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
	"fmt"
	"os"
	"strings"
)

func (dbs *dbScanner) failIfStrictf(format string, args ...interface{}) {
	if dbs.nonStrict {
		format = "WARNING: " + format
	} else {
		format = "ERROR: " + format
	}
	_, _ = fmt.Fprintf(os.Stderr, "\n"+format+"\n\n", args...)
	if !dbs.nonStrict {
		os.Exit(1)
	}
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func printLine(s string) {
	fmt.Println(strings.Repeat(s, 50))
}
