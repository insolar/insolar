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
