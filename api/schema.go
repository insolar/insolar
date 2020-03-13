// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/rpc/v2"

	"gopkg.in/yaml.v2"
)

type MAP map[interface{}]interface{}

type SchemaService struct {
	Data MAP
}

func NewSchemaService(ar *Runner) *SchemaService {
	ss := new(SchemaService)
	path := ar.cfg.SwaggerPath

	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Panicf("Can't read schema from '%s'", path)
	}
	err = yaml.NewDecoder(f).Decode(&ss.Data)
	if err != nil {
		log.Panicf("Can't parse schema from '%s' : %s", path, err)
	}

	return ss
}

func (s *SchemaService) Get(r *http.Request, args *SeedArgs, _ *rpc.RequestBody, reply *MAP) error {
	ctx, instr := instrumenter.NewMethodInstrument("SchemaService.get")
	defer instr.End()

	msg := fmt.Sprint("Incoming request: ", r.RequestURI)
	instr.Annotate(msg)

	logger := inslogger.FromContext(ctx)
	logger.Info("[ SchemaService.get ] ", msg)

	*reply = s.Data

	return nil
}
