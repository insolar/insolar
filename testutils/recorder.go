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

package testutils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

type synclist struct {
	sync.Mutex
	items []string
}

func (l *synclist) add(name string) {
	l.Lock()
	l.items = append(l.items, name)
	l.Unlock()
}

// Recorder records input data line by line in underline list.
type Recorder struct {
	list *synclist
}

// NewRecoder produces new Recorder instance.
func NewRecoder() *Recorder {
	return &Recorder{
		list: &synclist{},
	}
}

// Write implements io.Write interface
func (rec *Recorder) Write(p []byte) (int, error) {
	pr, pw := io.Pipe()
	s := bufio.NewScanner(pr)

	ch := make(chan struct{})
	go func() {
		for s.Scan() {
			rec.Add(s.Text())
		}
		close(ch)
	}()

	n, err := pw.Write(p)
	if err != nil {
		<-ch
		// defer pw.Close() ?
		return n, err
	}
	err = pw.Close()
	<-ch
	return n, err
}

// Add appends string to Recorder.
func (rec *Recorder) Add(s string) {
	rec.list.add(s)
}

// String stringifies recorder content.
func (rec *Recorder) String() string {
	if len(rec.list.items) == 0 {
		return ""
	}
	s := []string{"Steps:"}
	for n, step := range rec.list.items {
		s = append(s, fmt.Sprintf("%v:%v", n, step))
	}
	return strings.Join(s, ", ")
}

// StringMultiline stringifies recorder content in multiple lines.
func (rec *Recorder) StringMultiline() string {
	if len(rec.list.items) == 0 {
		return ""
	}
	s := []string{"Steps:"}
	for n, step := range rec.list.items {
		s = append(s, fmt.Sprintf("  %v: %v", n, step))
	}
	return strings.Join(s, "\n")
}

// Items returns Recorder's underlying array.
func (rec *Recorder) Items() []string {
	return rec.list.items
}
