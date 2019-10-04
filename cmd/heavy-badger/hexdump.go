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
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/spf13/cobra"
)

func (app *appCtx) valueHexDumpCommand() *cobra.Command {
	var key []byte
	var fixCmd = &cobra.Command{
		Use:   "dump-value",
		Short: "dump value by key",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires a argument with hex key value")
			}
			var err error
			key, err = hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("provided argument should be valid hex string: %v", err)
			}
			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			db, close := openDB(app.dataDir)
			defer close()

			var value []byte
			err := db.Backend().View(func(txn *badger.Txn) error {
				item, err := txn.Get(key)
				if err != nil {
					return err
				}
				value, err = item.ValueCopy(value)
				return err
			})
			if err != nil {
				fatalf("failed to get key from badger: %v", err)
			}
			// fmt.Printf("%v\n", value)
			io.Copy(os.Stdout, bytes.NewReader(value))
			// hex.EncodeToString(value)
		},
	}

	return fixCmd
}
