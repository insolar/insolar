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
	s.RegisterService(api.NewStorageExporterService(apiRunner), "exporter")
	http.Handle("/rpc", s)
	http.ListenAndServe("localhost:8080", s)
}
