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

// +build functest

package functest

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"
)

func TestPressureOnSystem(t *testing.T) {
	var contractCode = `
package main

import "github.com/insolar/insolar/logicrunner/builtin/foundation"

type One struct {
	foundation.BaseContract
	Number int
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Inc_API = true
func (c *One) Inc() (int, error) {
	c.Number++
	return c.Number, nil
}

var INSATTR_Get_API = true
func (c *One) Get() (int, error) {
	return c.Number, nil
}

var INSATTR_Dec_API = true
func (c *One) Dec() (int, error) {
	c.Number--
	return c.Number, nil
}
`
	protoRef := uploadContractOnce(t, "testPressure", contractCode)

	t.Run("one object, sequential calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		objectRef := callConstructor(syncT, protoRef, "New")

		for i := 0; i < 100; i++ {
			result := callMethod(syncT, objectRef, "Inc")
			require.Empty(syncT, result.Error)
			result = callMethod(syncT, objectRef, "Dec")
			require.Empty(syncT, result.Error)
		}
	})

	t.Run("one object, parallel calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		objectRef := callConstructor(syncT, protoRef, "New")

		wg := sync.WaitGroup{}
		wg.Add(10)
		for g := 0; g < 10; g++ {
			go func() {
				defer wg.Done()
				result := callMethod(syncT, objectRef, "Inc")
				require.Empty(syncT, result.Error)
				result = callMethod(syncT, objectRef, "Dec")
				require.Empty(syncT, result.Error)
			}()
		}
		wg.Wait()
	})

	t.Run("ten objects, sequential calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		wg := sync.WaitGroup{}
		wg.Add(10)
		for g := 0; g < 10; g++ {
			objectRef := callConstructor(syncT, protoRef, "New")
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					result := callMethod(syncT, objectRef, "Inc")
					require.Empty(syncT, result.Error)
					result = callMethod(syncT, objectRef, "Dec")
					require.Empty(syncT, result.Error)
				}
			}()
		}
		wg.Wait()
	})

	t.Run("ten objects, parallel calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		wg := sync.WaitGroup{}
		wg.Add(100)
		for g := 0; g < 10; g++ {
			objectRef := callConstructor(syncT, protoRef, "New")
			for c := 0; c < 10; c++ {
				go func() {
					defer wg.Done()
					for i := 0; i < 2; i++ {
						result := callMethod(syncT, objectRef, "Inc")
						require.Empty(syncT, result.Error)
						result = callMethod(syncT, objectRef, "Dec")
						require.Empty(syncT, result.Error)
					}
				}()
			}
		}
		wg.Wait()
	})
}
