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

package pool

import (
	"context"
	"io"
	"net"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

// Consuming 1 byte; only usable for outgoing connections.
func connectionClosedByPeer(ctx context.Context, conn net.Conn) bool {
	logger := inslogger.FromContext(ctx)

	n, err := conn.Read(make([]byte, 1))

	if err == io.EOF || n > 0 {
		if err != nil {
			logger.Errorln("[ connectionClosedByPeer ] Failed to close connection: ", err.Error())
		} else {
			logger.Debug("[ connectionClosedByPeer ] Close connection to %s", conn.RemoteAddr())
		}

		return true
	}

	return false
}
