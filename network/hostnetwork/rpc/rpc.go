/*
 *    Copyright 2018 INS Ecosystem
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

package rpc

import (
	"errors"
	"fmt"

	"github.com/insolar/insolar/network/hostnetwork/node"
)

// RemoteProcedure is remote procedure call function.
type RemoteProcedure func(sender *node.Node, args [][]byte) ([]byte, error)

// RPC is remote procedure call module.
type RPC interface {
	// Invoke is used to actually call remote procedure.
	Invoke(sender *node.Node, method string, args [][]byte) ([]byte, error)
	// RegisterMethod allows to register new function in RPC module.
	RegisterMethod(name string, method RemoteProcedure)
}

type rpc struct {
	methodTable map[string]RemoteProcedure
}

// NewRPC creates new RPC module.
func NewRPC() RPC {
	return &rpc{
		methodTable: make(map[string]RemoteProcedure),
	}
}

// Invoke calls registered function or returns error.
func (rpc *rpc) Invoke(sender *node.Node, methodName string, args [][]byte) (result []byte, err error) {
	method, exist := rpc.methodTable[methodName]
	if !exist {
		return nil, errors.New("method does not exist")
	}

	defer func() {
		if r := recover(); r != nil {
			result, err = nil, fmt.Errorf("panic: %s", r)
		}
	}()

	result, err = method(sender, args)
	return
}

// RegisterMethod registers new function in RPC module.
func (rpc *rpc) RegisterMethod(name string, method RemoteProcedure) {
	rpc.methodTable[name] = method
}
