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

package object

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdLocker_Same(t *testing.T) {
	tl := newtestlocker()
	id := insolar.ID{0x0A}

	counter := 0
	numParallelAccessors := 800
	wg := sync.WaitGroup{}
	wg.Add(numParallelAccessors)
	for i := 0; i < numParallelAccessors; i++ {
		go func() {
			tl.Lock("", id)
			counter++
			tl.Unlock("", id)
			wg.Done()
		}()
	}

	wg.Wait()
	require.Equal(t, numParallelAccessors, counter)
}

func TestIdLocker_Lock_PulseDoesntMatter(t *testing.T) {
	tl := newtestlocker()
	id1 := *insolar.NewID(0, []byte{0x0A})
	id2 := *insolar.NewID(1, []byte{0x0A})
	end := make(chan bool)
	go func() {
		tl.Lock("lock1", id1)
		tl.Unlock("unlock1", id2)
		close(end)
	}()
	select {
	case <-end:
	case <-time.After(5 * time.Minute):
		// Different record.ID should not lock each other.
		// 5s should be enough for any slow test environment.
		t.Fatalf(
			"Probably got deadlock (id1=%v, id2=%v).",
			id1.String(), id2.String(),
		)
	}
}

func TestIdLocker_Different(t *testing.T) {
	tl := newtestlocker()
	id1 := *insolar.NewID(0, []byte{0x0A})
	id2 := *insolar.NewID(0, []byte{0x0B})
	end := make(chan bool)
	go func() {
		tl.Lock("lock1", id1)
		tl.Lock("lock2", id2)
		tl.Unlock("unlock1", id1)
		tl.Unlock("unlock2", id2)
		close(end)
	}()
	select {
	case <-end:
	case <-time.After(5 * time.Minute):
		// Different record.ID should not lock each other.
		// 5s should be enough for any slow test environment.
		t.Fatalf(
			"Probably got deadlock (id1=%v, id2=%v).",
			id1.String(), id2.String(),
		)
	}

	expectsteps := []string{
		"before-lock1",
		"before-lock2",
		"before-unlock1",
		"before-unlock2",
	}
	assert.Equal(t, expectsteps, tl.synclist.list, "steps in proper order")
}

// test helpers

type synclist struct {
	sync.Mutex
	list []string
}

type testlock struct {
	lock     IndexLocker
	synclist *synclist
}

func newtestlocker() *testlock {
	return &testlock{
		lock:     NewIndexLocker(),
		synclist: &synclist{list: []string{}},
	}
}

func (l *synclist) Add(name string) {
	l.Lock()
	l.list = append(l.list, name)
	l.Unlock()
}

func (l *synclist) String() string {
	s := []string{"Steps:"}
	for n, step := range l.list {
		s = append(s, fmt.Sprintf("  %v: %v", n, step))
	}
	return strings.Join(s, "\n")
}

func (tl *testlock) Lock(name string, id insolar.ID) {
	tl.synclist.Add("before-" + name)
	tl.lock.Lock(id)
}

func (tl *testlock) Unlock(name string, id insolar.ID) {
	tl.synclist.Add("before-" + name)
	tl.lock.Unlock(id)
}
