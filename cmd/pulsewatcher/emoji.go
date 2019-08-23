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

package main

import (
	"github.com/insolar/insolar/api/requester"

	"sync"
)

type Emoji struct {
	mu        sync.RWMutex
	registred map[uint32]string
	light     []string
	virtual   []string
}

func NewEmoji() *Emoji {
	return &Emoji{
		registred: make(map[uint32]string),
		light:     []string{"😀", "😆", "😎", "😭", "😴", "♈️", "♉️", "♊️️", "♋️", "♌️", "♍️", "♎️", "♏️", "♐️", "♑️", "️♒️", "♓️"},
		virtual:   []string{"⚽", "🏀", "🏈", "🏐", "🏉", "🚗", "🚕", "🚙", "🚌", "🚒", "🚛", "🚜", "🚑️", "🚎", "🏎", "🚐", "🚚"},
	}
}

//todo: one url has many shortISs if node restart
func (e *Emoji) RegisterNode(_ string, n requester.Node) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.registred[n.ID]; ok {
		return
	}

	var x string
	switch n.Role {
	case "heavy_material":
		e.registred[n.ID] = "😈"
	case "light_material":
		// pop front
		x, e.light = e.light[0], e.light[1:]
		e.registred[n.ID] = x
	case "virtual":
		x, e.virtual = e.virtual[0], e.virtual[1:]
		e.registred[n.ID] = x
	}
}

func (e *Emoji) GetEmoji(n requester.Node) string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if res, ok := e.registred[n.ID]; ok {
		return res
	}
	return "️⛔️"
}
