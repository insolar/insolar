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
package storage

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func Test_IDLockTheSame(t *testing.T) {
	tl := newtestlocker()
	id1 := core.RecordID{0x0A}
	id2 := core.RecordID{0x0A}
	start1 := make(chan bool)
	start2 := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		<-start1
		tl.Lock("lock1", &id1)
		close(start2)

		time.Sleep(time.Millisecond * 50)
		tl.Unlock("unlock1", &id1)
		wg.Done()
	}()
	go func() {
		<-start2
		tl.Lock("lock2", &id2)
		tl.Unlock("unlock2", &id2)
		wg.Done()
	}()
	close(start1)
	wg.Wait()

	expectsteps := []string{
		"before-lock1",
		"before-lock2",
		"before-unlock1",
		"before-unlock2",
	}
	assert.Equal(t, expectsteps, tl.synclist.list, "steps in proper order")
}

func Test_IDLockDifferent(t *testing.T) {
	tl := newtestlocker()
	id1 := core.NewRecordID(0, []byte{0x0A})
	id2 := core.NewRecordID(1, []byte{0x0A})
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
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Got deadlock. Different record.ID should not lock each other.")
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
	lock     *IDLocker
	synclist *synclist
}

func newtestlocker() *testlock {
	return &testlock{
		lock:     NewIDLocker(),
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

func (tl *testlock) Lock(name string, id *core.RecordID) {
	tl.synclist.Add("before-" + name)
	tl.lock.Lock(id)
}

func (tl *testlock) Unlock(name string, id *core.RecordID) {
	tl.synclist.Add("before-" + name)
	tl.lock.Unlock(id)
}
