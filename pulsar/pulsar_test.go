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
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockListener struct {
	mock.Mock
}

func (mock *mockListener) Accept() (net.Conn, error) {
	panic("implement me")
}

func (mock *mockListener) Close() error {
	panic("implement me")
}

func (mock *mockListener) Addr() net.Addr {
	panic("implement me")
}

func TestNewPulsar_WithoutNeighbours(t *testing.T) {
	assertObj := assert.New(t)
	expectedPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBOAIBAAJAVGhjmmpIL8vDlqTpW25w+atXN9uW/hZvYPb/4ZlmOqZ5wWrDsTym
xunzzq3VDhBqQefMEwqAM2aTzKj4TBmKEwIDAQABAkA+zclGrMv3XDq0jRHg6QUA
kB9+PVJUzmajFEWCG7x36GijaMPS28lGr2uaQBcxaBvoqFfCjqmjg/nmjypF3YvB
AiEAoibr/sbsuzg5APwG5/9JWPu1JDMBB/e5LgNHO1emzKECIQCFQpLDyIVPjKaN
YDjUEigtmpKZtMv3XQLHjWl7iTMGMwIgMAghfd3FAAw+bnk5Pn2TZ4Vf+fIVyxtp
QiT8c6qaISECIGub+Nw0vsIgOBaODxXhm6RH3/5TKyoTZ70xCm8BubxVAiB7+fg/
vF+t7yqqR9T1g2Xv0KJpkquwBKNliiQnVwbuhA==
-----END RSA PRIVATE KEY-----`
	config := configuration.Pulsar{
		ConnectionType: "testType",
		ListenAddress:  "listedAddress",
		PrivateKey:     expectedPrivateKey,
	}
	actualConnectionType := ""
	actualAddress := ""

	result, err := NewPulsar(config, func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &mockListener{}, nil
	})

	if err != nil {
		t.Errorf("Error happened %v", err)
	}
	parsedKey, _ := ParseRsaPrivateKeyFromPemStr(expectedPrivateKey)
	assertObj.Equal(parsedKey, result.PrivateKey)
	assertObj.Equal("testType", actualConnectionType)
	assertObj.Equal("listedAddress", actualAddress)
	assertObj.IsType(result.Sock, &mockListener{})
	assertObj.NotNil(result.PrivateKey)
}

func TestNewPulsar_WithNeighbours(t *testing.T) {
	assertObj := assert.New(t)
	firstExpectedKey := `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALeFt8LnSBE13PHr5hI7L3JeHHg+CsAj
FoB1dP0Fq8BIRHbZIEawayjE1j1jvfpPNkVwqMEop+8utqy1XXJ1uL0CAwEAAQ==
-----END PUBLIC KEY-----`
	secondExpectedKey := `-----BEGIN PUBLIC KEY-----
MFswDQYJKoZIhvcNAQEBBQADSgAwRwJAZvrEQZj39XFoaQ+bho1J98yGXWyi729X
cYrmtcKWHcEvaSIFLUSC9Ec7VGeSS5H20r9YF/o5mo0SW6GJ8+Wg5QIDAQAB
-----END PUBLIC KEY-----`
	expectedPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBOAIBAAJAVGhjmmpIL8vDlqTpW25w+atXN9uW/hZvYPb/4ZlmOqZ5wWrDsTym
xunzzq3VDhBqQefMEwqAM2aTzKj4TBmKEwIDAQABAkA+zclGrMv3XDq0jRHg6QUA
kB9+PVJUzmajFEWCG7x36GijaMPS28lGr2uaQBcxaBvoqFfCjqmjg/nmjypF3YvB
AiEAoibr/sbsuzg5APwG5/9JWPu1JDMBB/e5LgNHO1emzKECIQCFQpLDyIVPjKaN
YDjUEigtmpKZtMv3XQLHjWl7iTMGMwIgMAghfd3FAAw+bnk5Pn2TZ4Vf+fIVyxtp
QiT8c6qaISECIGub+Nw0vsIgOBaODxXhm6RH3/5TKyoTZ70xCm8BubxVAiB7+fg/
vF+t7yqqR9T1g2Xv0KJpkquwBKNliiQnVwbuhA==
-----END RSA PRIVATE KEY-----`
	config := configuration.Pulsar{
		ConnectionType: "testType",
		ListenAddress:  "listedAddress",
		PrivateKey:     expectedPrivateKey,
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "first", PublicKey: firstExpectedKey},
			{ConnectionType: "pct", Address: "second", PublicKey: secondExpectedKey},
		},
	}

	result, err := NewPulsar(config, func(connectionType string, address string) (net.Listener, error) {
		return &mockListener{}, nil
	})

	if err != nil {
		t.Errorf("Error happened %v", err)
	}
	assertObj.Equal(2, len(result.Neighbours))
	assertObj.Equal("tcp", result.Neighbours[firstExpectedKey].ConnectionType.String())
	assertObj.Equal("pct", result.Neighbours[secondExpectedKey].ConnectionType.String())
}
