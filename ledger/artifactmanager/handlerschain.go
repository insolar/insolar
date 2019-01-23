/*
 *    Copyright 2019 Insolar
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

package artifactmanager

import (
	"github.com/insolar/insolar/core"
)

type Handler func(core.MessageHandler) core.MessageHandler

type HandlersChain struct {
	handlers []Handler
}

func NewChain(handlers ...Handler) *HandlersChain {
	return &HandlersChain{handlers: append([]Handler{}, handlers...)}
}

func (hc *HandlersChain) Extend(chain *HandlersChain) *HandlersChain {
	return hc.Append(chain.handlers...)
}

func (hc *HandlersChain) Append(handlers ...Handler) *HandlersChain {
	newCons := make([]Handler, 0, len(hc.handlers)+len(handlers))
	newCons = append(newCons, hc.handlers...)
	newCons = append(newCons, handlers...)

	return &HandlersChain{handlers: newCons}
}

func (hc *HandlersChain) PrependAndCopy(handler Handler) *HandlersChain {
	result := make([]Handler, len(hc.handlers)+1)
	copy(result[1:], hc.handlers)
	result[0] = handler

	return &HandlersChain{handlers: result}
}

func (hc *HandlersChain) Then(mh core.MessageHandler) core.MessageHandler {
	if mh == nil {
		panic("MessageHandler in Then function-call can't be nil")
	}

	for i := range hc.handlers {
		mh = hc.handlers[len(hc.handlers)-1-i](mh)
	}

	return mh
}
