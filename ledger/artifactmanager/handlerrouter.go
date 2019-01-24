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

// Alias for wrapping func-s
type Handler func(core.MessageHandler) core.MessageHandler

// Build return wrapping core.MessageHandler
// If want call wrapHandler1(wrapHandler2(wrapHandler3(handler))),
// we should use Build(handler, wrapHandler1, wrapHandler2, wrapHandler3).
func Build(handler core.MessageHandler, wrapHandlers ...Handler) core.MessageHandler {
	result := handler

	for i := range wrapHandlers {
		result = wrapHandlers[len(wrapHandlers)-1-i](result)
	}

	return result
}
