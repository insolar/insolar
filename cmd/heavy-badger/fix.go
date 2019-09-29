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
	"github.com/dgraph-io/badger"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar/store"
)

func (app *appCtx) addFixCommand(parent *cobra.Command) {
	var fixCmd = &cobra.Command{
		Use:   "fix",
		Short: "opens and closes badger database. Could fix 'Database was not properly closed' error.",
		Run: func(_ *cobra.Command, _ []string) {
			if err := checkDirectory(app.dataDir); err != nil {
				fatalf("Database directory '%v' open failed. Error: \"%v\"", app.dataDir, err)
			}

			ops := badger.DefaultOptions(app.dataDir)
			dbWrapped, err := store.NewBadgerDB(ops)
			if err != nil {
				fatalf("failed open database directory %v: %v", app.dataDir, err)
			}
			err = dbWrapped.Backend().Close()
			if err != nil {
				fatalf("failed close database directory %v: %v", app.dataDir, err)
			}
		},
	}

	parent.AddCommand(fixCmd)
}
