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

package pulsar

import (
	"bytes"
	"os"

	"github.com/insolar/insolar/log"
)

func capture(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

//func TestNeighbour_CheckAndRefreshConnection_RefreshSuccess(t *testing.T) {
//	client := &pulsartestutil.MockRPCClientWrapper{}
//	client.On("CreateConnection", configuration.TCP, "expectedAddress").Return(nil)
//	client.On("Lock")
//	client.On("Unlock")
//	neighbour := &Neighbour{
//		ConnectionAddress: "expectedAddress",
//		ConnectionType:    configuration.TCP,
//		OutgoingClient:    client,
//	}
//
//	writtenLog := capture(func() { neighbour.HandleConnectionError(rpc.ErrShutdown) })
//
//	assert.Contains(t, writtenLog, "Restarting RPC Connection to expectedAddress due to error connection is shut down")
//	client.AssertCalled(t, "CreateConnection", configuration.TCP, "expectedAddress")
//	client.AssertCalled(t, "Lock")
//	client.AssertCalled(t, "Unlock")
//}
//
//func TestNeighbour_CheckAndRefreshConnection_RefreshFailed(t *testing.T) {
//	client := &pulsartestutil.MockRPCClientWrapper{}
//	client.On("CreateConnection", configuration.TCP, "expectedAddress").Return(errors.New("oops"))
//	client.On("Lock")
//	client.On("Unlock")
//	neighbour := &Neighbour{
//		ConnectionAddress: "expectedAddress",
//		ConnectionType:    configuration.TCP,
//		OutgoingClient:    client,
//	}
//
//	writtenLog := capture(func() { neighbour.HandleConnectionError(rpc.ErrShutdown) })
//
//	assert.Contains(t, writtenLog, "Restarting RPC Connection to expectedAddress due to error connection is shut down")
//	assert.Contains(t, writtenLog, "Refreshing connection to expectedAddress failed due to error oops")
//	client.AssertCalled(t, "CreateConnection", configuration.TCP, "expectedAddress")
//	client.AssertCalled(t, "Lock")
//	client.AssertCalled(t, "Unlock")
//}
