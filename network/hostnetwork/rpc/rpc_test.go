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

import (
	"errors"
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/stretchr/testify/assert"
)

func TestNewRPC(t *testing.T) {
	r := NewRPC()

	assert.Equal(t, r, &rpc{
		methodTable: make(map[string]RemoteProcedure),
	})
}

func TestRPC_Invoke_ReturnsErrorForNonExistingMethod(t *testing.T) {
	r := NewRPC()
	_, err := r.Invoke(nil, "test_method", nil)

	assert.EqualError(t, err, "method does not exist")
}

func TestRPC_RegisterMethod(t *testing.T) {
	r := NewRPC()
	_, err := r.Invoke(nil, "test_method", nil)
	assert.Error(t, err)

	r.RegisterMethod("test_method", func(sender *host.Host, args [][]byte) ([]byte, error) {
		return []byte("hello world"), nil
	})

	res, err := r.Invoke(nil, "test_method", nil)
	assert.NoError(t, err)
	assert.Equal(t, res, []byte("hello world"))
}

func TestRPC_Invoke_RecoversFromPanic(t *testing.T) {
	r := NewRPC()
	r.RegisterMethod("panic_method", func(sender *host.Host, args [][]byte) ([]byte, error) {
		panic("test_panic")
	})

	res, err := r.Invoke(nil, "panic_method", nil)
	assert.Nil(t, res)
	assert.EqualError(t, err, "panic: test_panic")
}

func TestRPC_Invoke_ReturnsErrorFromMethod(t *testing.T) {
	r := NewRPC()
	r.RegisterMethod("error_method", func(sender *host.Host, args [][]byte) ([]byte, error) {
		return nil, errors.New("example error")
	})

	res, err := r.Invoke(nil, "error_method", nil)
	assert.Nil(t, res)
	assert.EqualError(t, err, "example error")
}
