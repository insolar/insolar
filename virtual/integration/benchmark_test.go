package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/semaphore"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

// BenchmarkUserCreation benchmarks parallel user creation (with maximum 100 workers at a time)
func BenchmarkUserCreation(b *testing.B) {
	ctx := context.Background()
	cfg := DefaultVMConfig()

	s, err := NewVirtualServer(b, ctx, cfg).WithGenesis().PrepareAndStart()
	require.NoError(b, err)
	defer s.Stop(ctx)

	var (
		iterations = b.N
		helper     = ServerHelper{s}
		syncAssert = NewTSAssert(b)
		sema       = semaphore.NewWeighted(100)
	)

	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()

			ctx, _ := inslogger.WithTraceField(ctx, uuid.New().String())

			err := sema.Acquire(ctx, 1)
			if err == nil {
				defer sema.Release(1)
			} else {
				panic(fmt.Sprintf("unexpected: %s", err.Error()))
			}

			_, err = helper.createUser(ctx)
			syncAssert.NoError(err, "failed to create user")
		}()
	}
	wg.Wait()
}
