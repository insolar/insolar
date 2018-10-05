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

package rpc

// Factory allows to create new RPC.
type Factory interface {
	Create() RPC
}

type rpcFactory struct {
	methods map[string]RemoteProcedure
}

// NewRPCFactory creates new RPC Factory.
func NewRPCFactory(methods map[string]RemoteProcedure) Factory {
	return &rpcFactory{
		methods: methods,
	}
}

// Create creates and registers new remote procedure.
func (rpcFactory *rpcFactory) Create() RPC {
	newRPC := NewRPC()
	for name, method := range rpcFactory.methods {
		newRPC.RegisterMethod(name, method)
	}
	return newRPC
}
