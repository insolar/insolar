/*
 *    Copyright 2018 INS Ecosystem
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

package transport

import (
	"testing"

	"github.com/insolar/network/connection"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStoreFactory(t *testing.T) {
	expectedFactory := &utpTransportFactory{}
	actualFactory := NewUTPTransportFactory()

	assert.Equal(t, expectedFactory, actualFactory)
}

func TestMemoryStoreFactory_Create(t *testing.T) {
	conn, err := connection.NewConnectionFactory().Create("127.0.0.1:8080")
	assert.NoError(t, err)
	defer conn.Close()

	transport, err := NewUTPTransportFactory().Create(conn)

	assert.NoError(t, err)
	assert.Implements(t, (*Transport)(nil), transport)
}
