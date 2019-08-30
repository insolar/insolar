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

package log

import (
	"sort"
	"strings"
	"sync"

	"github.com/insolar/insolar/insolar"
)

type Controller struct {
	mu sync.RWMutex
	fc *filterChecker
}

func NewController() *Controller {
	return &Controller{
		fc: newFilterChecker(),
	}
}

func (c *Controller) checker() *filterChecker {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.fc
}

func (c *Controller) Set(s string, level insolar.LogLevel) {
	c.mu.Lock()
	c.fc = c.fc.WithSet(s, level)
	c.mu.Unlock()
}

func (c *Controller) Del(s string) bool {
	c.mu.Lock()
	var ok bool
	c.fc, ok = c.fc.WithRemove(s)
	c.mu.Unlock()
	return ok
}

func (c *Controller) List() []insolar.LogControllerItem {
	c.mu.RLock()
	fc := c.fc
	c.mu.RUnlock()
	return fc.list()
}

type filterChecker struct {
	filters map[string]insolar.LogLevel
	ordered []string
}

func newFilterChecker() *filterChecker {
	return &filterChecker{
		filters: make(map[string]insolar.LogLevel),
	}
}

func (fc *filterChecker) copy() *filterChecker {
	fcNew := newFilterChecker()
	for k, v := range fc.filters {
		fcNew.filters[k] = v
	}
	for _, v := range fc.ordered {
		fcNew.ordered = append(fcNew.ordered, v)
	}
	return fcNew
}

func (fc *filterChecker) WithSet(s string, level insolar.LogLevel) *filterChecker {
	fcNew := fc.copy()
	fcNew.set(s, level)
	return fcNew
}

func (fc *filterChecker) set(s string, level insolar.LogLevel) {
	fc.filters[s] = level
	fc.updateOrdered()
}

func (fc *filterChecker) list() []insolar.LogControllerItem {
	items := make([]insolar.LogControllerItem, 0, len(fc.ordered))
	for _, pat := range fc.ordered {
		items = append(items, insolar.LogControllerItem{
			Level:  fc.filters[pat],
			Prefix: pat,
		})
	}
	return items
}

func (fc *filterChecker) updateOrdered() {
	fc.ordered = make([]string, 0, len(fc.filters))
	for k := range fc.filters {
		fc.ordered = append(fc.ordered, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fc.ordered)))
}

func (fc *filterChecker) WithRemove(s string) (*filterChecker, bool) {
	if _, ok := fc.filters[s]; !ok {
		return nil, false
	}
	fcNew := fc.copy()
	fcNew.remove(s)
	return fcNew, true
}

func (fc *filterChecker) remove(s string) {
	delete(fc.filters, s)
	// fc.rt.Delete(s)
	fc.updateOrdered()
}

func (c *Controller) Check(path string, level insolar.LogLevel) insolar.LogLevel {
	return c.checker().check(path, level)
}

func (fc *filterChecker) check(path string, level insolar.LogLevel) insolar.LogLevel {
	if len(fc.filters) > 0 {
		if level, ok := fc.search(path); ok {
			// fmt.Println("filterChecker:", path, "->", level.String())
			return level
		}
	}
	return level
}

func (fc *filterChecker) search(s string) (insolar.LogLevel, bool) {
	// fmt.Println("search:", s)
	for _, f := range fc.ordered {
		// fmt.Printf("   '%v' <- ('%v' contains '%v')\n", strings.Contains(s, f), s, f)
		// if strings.Contains(s, f) {
		if strings.HasPrefix(s, f) {
			lvl := fc.filters[f]
			// fmt.Println("....", lvl.String(), fc.filters[f], true)
			return lvl, true
		}
	}
	// fmt.Println(0, "false")
	return 0, false
}

//
// func (fc *filterChecker) search(s string) (zerolog.Level, bool) {
// 	val, ok := fc.rt.Get(s)
// 	if !ok {
// 		return 0, ok
// 	}
// 	return val.(zerolog.Level), ok
// }
