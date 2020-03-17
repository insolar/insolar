// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
)

func TestPressureOnSystem(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
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

func TestCoinPassing(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	var contractCode = `
package main

import "github.com/insolar/insolar/insolar"
import "github.com/insolar/insolar/logicrunner/builtin/foundation"
import "github.com/insolar/insolar/applicationbase/proxy/testCoinPassing"
import "errors"

type One struct {
	foundation.BaseContract
	Amount int
}

func New(n int) (*One, error) {
	return &One{Amount: n}, nil
}

var INSATTR_Balance_API = true
func (c *One) Balance() (int, error) {
	return c.Amount, nil
}

var INSATTR_Transfer_API = true
func (c *One) Transfer(ref insolar.Reference) (error) {
	if c.Amount <= 0 {
		return errors.New("Oops")
	}

	c.Amount -= 1

	w := testCoinPassing.GetObject(ref)
	err := w.Accept(1)
	if err != nil {
		return err
	}
	return nil
}

//ins:saga(Rollback)
func (c *One) Accept(amount int) error {
	c.Amount += amount
	return nil
}

func (c *One) Rollback(amount int) error {
	c.Amount -= amount
	return nil
}
`
	protoRef := uploadContractOnce(t, "testCoinPassing", contractCode)

	t.Run("pass one coin in parallel between two wallets", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		w1Ref := callConstructor(syncT, protoRef, "New", 1)
		w2Ref := callConstructor(syncT, protoRef, "New", 0)

		transfers := 30
		wg := sync.WaitGroup{}
		wg.Add(transfers * 2)

		var successes uint32
		var errors uint32

		f := func(from, to *insolar.Reference) {
			for i := 0; i < transfers; i++ {
				res := callMethod(syncT, from, "Transfer", to)
				if res.ExtractedError != "" {
					atomic.AddUint32(&errors, 1)
				} else if res.Error == nil {
					atomic.AddUint32(&successes, 1)
				}
				wg.Done()
				time.Sleep(100 * time.Millisecond)
			}
		}
		go f(w1Ref, w2Ref)
		go f(w2Ref, w1Ref)
		wg.Wait()

		require.Greater(syncT, successes, uint32(0))
		require.Greater(syncT, errors, uint32(0))

		getBalance := func() float64 {
			r1 := callMethod(syncT, w1Ref, "Balance")
			require.Empty(syncT, r1.Error)
			r2 := callMethod(syncT, w2Ref, "Balance")
			require.Empty(syncT, r2.Error)
			return r1.ExtractedReply.(float64) + r2.ExtractedReply.(float64)
		}

		pass := false
		for i := 0; i < 10; i++ {
			bal := getBalance()
			if bal != float64(1) {
				time.Sleep(1 * time.Second)
				continue
			}
			pass = true
		}
		if !pass {
			require.Fail(t, "balance missmatch")
		}

		for i := 0; i < 3; i++ {
			bal := getBalance()
			require.Equal(t, float64(1), bal)
			time.Sleep(1 * time.Second)
		}
	})
}
