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

package api

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
)

// ContractService is a service that provides ability to add custom contracts
type ContractService struct {
	runner *Runner
}

// NewContractService creates new Contract service instance.
func NewContractService(runner *Runner) *ContractService {
	return &ContractService{runner: runner}
}

// UploadArgs is arguments that Contract.Upload accepts.
type UploadArgs struct {
	Code string
	Name string
}

// UploadReply is reply that Contract.Upload returns
type UploadReply struct {
	PrototypeRef insolar.Reference `json:"PrototypeRef"`
}

// Upload builds code and return prototype ref
func (s *ContractService) Upload(r *http.Request, args *UploadArgs, reply *UploadReply) error {
	_, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ ContractService.Upload ] Incoming request: %s", r.RequestURI)

	if len(args.Name) == 0 {
		return errors.New("params.name is missing")
	}

	if len(args.Code) == 0 {
		return errors.New("params.code is missing")
	}

	insgocc, err := goplugintestutils.BuildPreprocessor()
	if err != nil {
		return errors.Wrap(err, "can't build preprocessor")
	}
	cb := goplugintestutils.NewContractBuilder(s.runner.ArtifactManager, insgocc)

	contractMap := make(map[string]string)
	contractMap[args.Name] = args.Code

	err = cb.Build(contractMap)
	if err != nil {
		return errors.Wrap(err, "can't build contract")
	}

	reply.PrototypeRef = *cb.Prototypes[args.Name]
	return nil
}
