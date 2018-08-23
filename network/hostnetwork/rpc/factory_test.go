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
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/node"

	"github.com/stretchr/testify/assert"
)

func rpcTestMethod1(sender *node.Node, args [][]byte) ([]byte, error) {
	return []byte("testMethod1"), nil
}

func rpcTestMethod2(sender *node.Node, args [][]byte) ([]byte, error) {
	return []byte("testMethod2"), nil
}

var rpcTestMethods = map[string]RemoteProcedure{
	"testMethod1": rpcTestMethod1,
	"testMethod2": rpcTestMethod2,
}

func TestNewRPCFactory(t *testing.T) {
	actualFactory := NewRPCFactory(rpcTestMethods)
	expectedFactory := &rpcFactory{
		methods: rpcTestMethods,
	}

	assert.Equal(t, expectedFactory, actualFactory)
}

func TestRpcFactory_Create(t *testing.T) {
	actualRPC := NewRPCFactory(rpcTestMethods).Create()

	assert.Implements(t, (*RPC)(nil), actualRPC)

	address, _ := node.NewAddress("127.0.0.1:31337")
	tests := []struct {
		name   string
		result []byte
		err    error
	}{
		{"testMethod1", []byte("testMethod1"), nil},
		{"testMethod2", []byte("testMethod2"), nil},
		{"testMethodNotExist", nil, errors.New("method does not exist")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := actualRPC.Invoke(node.NewNode(address), test.name, [][]byte{})
			assert.Equal(t, test.result, res)
			assert.Equal(t, test.err, err)
		})
	}
}
