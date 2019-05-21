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

// Package builtin is implementation of builtin contracts engine
package builtin

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/preprocessor"
)

// BuiltIn is a contract runner engine
type BuiltIn struct {
	RefRegistry map[insolar.Reference]string
	// Prototype -> Code + Versions
	PrototypeRegistry    map[string]preprocessor.ContractWrapper
	PrototypeRefRegistry map[insolar.Reference]string
	// Code ->
	CodeRegistry    map[string]preprocessor.ContractWrapper
	CodeRefRegistry map[insolar.Reference]string
}

// NewBuiltIn is an constructor
func NewBuiltIn(eb insolar.MessageBus, am artifacts.Client) *BuiltIn {
	return &BuiltIn{}
}

func (bi *BuiltIn) CallConstructor(ctx context.Context, callCtx *insolar.LogicCallContext, code insolar.Reference, name string, args insolar.Arguments) (objectState []byte, err error) {
	panic("implement me")
}

func (bi *BuiltIn) Stop() error {
	return nil
}

// CallMethod runs a method on contract
func (bi *BuiltIn) CallMethod(ctx context.Context, callCtx *insolar.LogicCallContext, codeRef insolar.Reference, data []byte, method string, args insolar.Arguments) (newObjectState []byte, methodResults insolar.Arguments, err error) {
	panic("implement me")
}
