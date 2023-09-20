// +build bloattest

package functest

import (
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestPressureOnSystem(t *testing.T) {
	t.Run("one object, sequential calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		objectRef := callConstructor(syncT.T, "first", "New")
		exp := callMethod(syncT, objectRef, "first.Get").(float64)

		for i := 0; i < 100; i++ {
			result := callMethod(syncT, objectRef, "first.Inc")
			require.Equal(syncT, exp+1, result.(float64))
			result = callMethod(syncT, objectRef, "first.Dec")
			require.Empty(syncT, exp, result.(float64))
		}
	})

	t.Run("one object, parallel calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		objectRef := callConstructor(syncT.T, "first", "New")
		exp := callMethod(syncT, objectRef, "first.Get").(float64)

		wg := sync.WaitGroup{}
		wg.Add(10)
		for g := 0; g < 10; g++ {
			go func() {
				defer wg.Done()
				callMethod(syncT, objectRef, "first.Inc")
				callMethod(syncT, objectRef, "first.Dec")
			}()
		}
		wg.Wait()
		act := callMethod(syncT, objectRef, "first.Get").(float64)
		require.Equal(syncT.T, exp, act)
	})

	t.Run("ten objects, sequential calls", func(t *testing.T) {
		syncT := &testutils.SyncT{T: t}

		wg := sync.WaitGroup{}
		wg.Add(10)
		for g := 0; g < 10; g++ {
			objectRef := callConstructor(syncT.T, "first", "New")
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					callMethod(syncT, objectRef, "first.Inc")
					callMethod(syncT, objectRef, "first.Dec")
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
			objectRef := callConstructor(syncT, "first", "New")
			for c := 0; c < 10; c++ {
				go func() {
					defer wg.Done()
					for i := 0; i < 2; i++ {
						callMethod(syncT, objectRef, "first.Inc")
						callMethod(syncT, objectRef, "first.Dec")
					}
				}()
			}
		}
		wg.Wait()
	})
}

func TestCoinPassing(t *testing.T) {
	syncT := &testutils.SyncT{T: t}

	balance := 100
	member1 := callConstructorWithParameters(t, "first", "NewWithNumber", map[string]interface{}{"amount": balance})
	member2 := callConstructorWithParameters(t, "first", "NewWithNumber", map[string]interface{}{"amount": balance})

	transfers := 100
	amount := 1
	wg := sync.WaitGroup{}
	wg.Add(transfers * 2)

	f := func(from, to string) {
		for i := 0; i < transfers; i++ {
			_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "first.TransferTo",
				map[string]interface{}{"reference": from, "toRef": to, "amount": amount})
			require.NoError(t, err)
			wg.Done()
			time.Sleep(100 * time.Millisecond)
		}
	}
	go f(member1, member2)
	go f(member2, member1)
	wg.Wait()

	act := callMethod(syncT, member1, "first.Get").(float64)
	require.Equal(syncT.T, balance, int(act), "member1 balance not as expected")
	act = callMethod(syncT, member2, "first.Get").(float64)
	require.Equal(syncT.T, balance, int(act), "member2 balance not as expected")
}
