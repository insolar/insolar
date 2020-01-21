// Copyright 2020 Insolar Network Ltd.
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

//usr/bin/env go run "$0" "$@"; exit "$?"
// realpath.go - because we don't wand depend on coreutils on MacOS X for building binaries

// +build tools

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func main() {
	gotPath := os.Args[1]
	absPath, err := filepath.Abs(gotPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get absolute path for '%v': %v\n", gotPath, err)
		os.Exit(1)
	}
	fmt.Println(path.Clean(absPath))
}
