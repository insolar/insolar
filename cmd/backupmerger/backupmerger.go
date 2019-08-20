///
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
///

package main

import (
	"os"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	targetDBPath    string
	backupFileName  string
	numberOfWorkers int
	help            bool
)

func usage() {
	pflag.Usage()
	os.Exit(0)
}

func parseInputParams() {

	pflag.StringVarP(
		&targetDBPath, "target-db", "t", "", "directory where backup will be roll to (required)")
	pflag.StringVarP(
		&backupFileName, "bkp-name", "n", "", "file name if incremental backup (required)")
	pflag.IntVarP(
		&numberOfWorkers, "workers-num", "w", 1, "number of workers to read backup file")
	pflag.BoolVarP(
		&help, "help", "h", false, "show this help")

	pflag.Parse()

	if help {
		usage()
	}

	if len(targetDBPath) == 0 || len(backupFileName) == 0 {
		println("bkp-name and target-db are required\n")
		usage()
	}
}

func printError(err error, message string) {
	println(errors.Wrap(err, "ERROR "+message).Error())
}

func main() {
	parseInputParams()

	var merged bool
	ops := badger.DefaultOptions(targetDBPath)
	bdb, err := badger.Open(ops)
	if err != nil {
		printError(err, "failed to open badger")
		return
	}
	defer func() {
		err = bdb.Close()
		if err != nil {
			printError(err, "failed to close db")
			return
		}
		if merged {
			println()
			println("successfully merged " + backupFileName + " to " + targetDBPath)
		}
	}()

	bkpFile, err := os.Open(backupFileName)
	if err != nil {
		printError(err, "failed to open backup file")
		return
	}

	err = bdb.Load(bkpFile, numberOfWorkers)
	if err != nil {
		printError(err, "failed to load backup")
		return
	}
	merged = true
}
