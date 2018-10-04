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

package reply

import (
	"github.com/insolar/insolar/core"
)

// Common - the most common reply
type Common struct {
	Data   []byte
	Result []byte
}

// Type returns type of the reply
func (r *Common) Type() core.ReplyType {
	return TypeCommon
}

type CallConstructor struct {
	Ref     *core.RecordRef
}

// Type returns type of the reply
func (r *CallConstructor) Type() core.ReplyType {
	return TypeCommon
}
