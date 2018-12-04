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

package contractrequester

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
)

func mockMessageBus(t *testing.T, result core.Reply) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		return result, nil
	}
	return mbMock
}
