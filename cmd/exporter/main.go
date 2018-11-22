/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/gorilla/rpc/v2"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/exporter"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/platformpolicy"
)

func main() {
	// FIXME: this is a temporary implementation. Make me pretty!
	ctx := context.Background()

	// Ledger
	ledgerConf := configuration.NewLedger()
	ledgerConf.Storage.DataDirectory = flag.Arg(1)
	db, err := storage.NewDB(ledgerConf, nil)
	if err != nil {
		panic(err)
	}
	db.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()
	exp := exporter.NewExporter(db)
	err = db.Init(ctx)
	if err != nil {
		panic(err)
	}

	// API
	apiConf := configuration.NewAPIRunner()
	apiRunner, err := api.NewRunner(&apiConf)
	if err != nil {
		panic(err)
	}
	apiRunner.StorageExporter = exp

	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	err = s.RegisterService(api.NewStorageExporterService(apiRunner), "exporter")
	if err != nil {
		panic(err)
	}
	http.Handle("/rpc", s)
	err = http.ListenAndServe("localhost:8080", s)
	if err != nil {
		panic(err)
	}
}
