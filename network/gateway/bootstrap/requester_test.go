//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package bootstrap

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/insolar/insolar/network/consensus/adapters"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/insolar"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	mock "github.com/insolar/insolar/testutils/network"
)

func TestRequester_Authorize(t *testing.T) {
	t.Skip("Until merge")
	cert := GetTestCertificate()

	options := network.ConfigureOptions(configuration.NewConfiguration())

	cs := testutils.NewCryptographyServiceMock(t)
	sig := insolar.SignatureFromBytes([]byte("lalal"))
	cs.SignMock.Return(&sig, nil)

	r := NewRequester(options)
	r.(*requester).CryptographyService = cs
	resp, err := r.Authorize(context.Background(), cert)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRequester_Bootstrap(t *testing.T) {
	options := network.ConfigureOptions(configuration.NewConfiguration())

	hn := mock.NewHostNetworkMock(t)
	hn.SendRequestToHostMock.Set(func(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) (r network.Future, r1 error) {
		return nil, errors.New("123")
	})

	p := &packet.Permit{}
	candidateProfile := adapters.Candidate{}
	r := NewRequester(options)
	// inject HostNetwork
	r.(*requester).HostNetwork = hn

	resp, err := r.Bootstrap(context.Background(), p, candidateProfile, insolar.GenesisPulse)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func GetTestCertificate() *certificate.Certificate {
	buff := bytes.NewBufferString(`
{
  "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+2RsLu5z3nFEimNiesiLDH2Kw1GM\nvgYylDXAmZxpbGjQZ5FqHuXF+DJrwKYzDyfBDEQz6Tu/aeA2CgRZvqbKug==\n-----END PUBLIC KEY-----\n",
  "reference": "1tJBuMQ1SW9Q3fUW8YoateDhfqKBP3GhFEpHH95R8E.11111111111111111111111111111111",
  "role": "heavy_material",
  "majority_rule": 0,
  "min_roles": {
    "virtual": 1,
    "heavy_material": 1,
    "light_material": 1
  },
  "pulsar_public_keys": null,
  "root_domain_ref": "1tJEDNVffdf4PQjxhKvQVc3D166RqhmFBS2gkBpGva.11111111111111111111111111111111",
  "bootstrap_nodes": [
    {
      "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+2RsLu5z3nFEimNiesiLDH2Kw1GM\nvgYylDXAmZxpbGjQZ5FqHuXF+DJrwKYzDyfBDEQz6Tu/aeA2CgRZvqbKug==\n-----END PUBLIC KEY-----\n",
      "host": "127.0.0.1:13831",
      "network_sign": "4pctXVOJNOBZO09Nbd8DNxXM5foSfeTec52DgTemDYVO5WddFCqjUdeKNtRdfNTmwYdyPBKtSjFm5x1TOAPLHA==",
      "node_sign": "0IldPk9aVLKPNF3vFVVsJx4o94DxGyUEK9GgvRguVMgVfzSh48k4ymBe4bzmEt/Zfw4LHHi+OYrVLUu+eTTSNQ==",
      "node_ref": "1tJBuMQ1SW9Q3fUW8YoateDhfqKBP3GhFEpHH95R8E.11111111111111111111111111111111"
    },
    {
      "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEH4Q3jPLcsuajEsHuJS3si8uvk2H\nAx/b9Mx2z0TlqX6ql1YkM0HYLJUAF6ftxUTEP5igrApFw8h3ypHLMOV3Wg==\n-----END PUBLIC KEY-----\n",
      "host": "127.0.0.1:23832",
      "network_sign": "J/XNGuQf8WPDfCP/60+zrjyWw7rIyDKGLokrPVxTwU7LFVo5NGuL7vBjpIG3P9JTL1ez01y7LI+TaPpZOUNukA==",
      "node_sign": "6hl/+blUtMA80coqgH6ThHFJxbZioCQCCZ4v+pZ10yheYwX7QS5ANmEhHboHdlt1R4QNRKPsWQr7q6hgaKvNaA==",
      "node_ref": "1tJDdakD4TeVHYzsiuYciE2eCm2N7uxMpN1iqanmiP.11111111111111111111111111111111"
    },
    {
      "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAES8wyC6bcD05FAdubdXPtKAxZalyJ\n4t8F0/64lBUNMXzM75l6iO7GnhgcpCWLZ71Wmd0fkWjlFQDUdzXrWcxH+Q==\n-----END PUBLIC KEY-----\n",
      "host": "127.0.0.1:33833",
      "network_sign": "J5ikEJUlssGzRpl9DlKpJy8IH7DwmdkzP8WPcdHBWV0o5UDzLdt0sL3jhBS3xZtWtUhLwUOsVgN0tcbkXaoEXw==",
      "node_sign": "5/UIiZ4pz1hqAQnCGwMxXpOWODGLZ+/ip8UR+gmkiYQmN8faglzMkRxyGs+OdaN39T72cuYdFO5KUUujKAaxNg==",
      "node_ref": "1tJDQSknYb9G6yvWTfeD3K379E6xHxV1RAHMvhWA7N.11111111111111111111111111111111"
    },
    {
      "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEiJ6U8z+0puX7djtRDPwlDtzPi4rv\nw05XEirPjzF2zxRSZea0XojBIBKhT9d4id+LR9tFG2Vt+jQJCjmF7nri9Q==\n-----END PUBLIC KEY-----\n",
      "host": "127.0.0.1:43834",
      "network_sign": "OGgm2sAyVVB7vtRs4dWeoGqtGl7qML5bfhnG40WA8Pb5ew4Sl+gGk7GMj+F0DKpycSV/riGpgVbI7hdOFpsBkQ==",
      "node_sign": "L0XpoNFMSAvb/zbz1zE7Dc5ParDkFDpkZWfTOXJLlRzhKVKfyzTRoLq3aPz6gINA/3f8wfIejggVmqSK5mStbA==",
      "node_ref": "1tJDMvuYDKsH4gpBBATR9HsPo2ZY5iKYL9iF271Mig.11111111111111111111111111111111"
    },
    {
      "public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE7SzA4AftGRUvsiUpKd+JcNo47EdE\nXGt6dXbE9AKLPNkBKgHscLNo0ZOjC1jCajX8teWpWAyNsVuyadoR5/Q/ig==\n-----END PUBLIC KEY-----\n",
      "host": "127.0.0.1:53835",
      "network_sign": "95dbB3OonTHi36y1zZ8u4jE0Nt+hP86P/pC3Z6d8ZXM0RIyMM7GJolzCKPyYE+IfkYUX0HgLXX/XoWp0RPyppw==",
      "node_sign": "tQ5UbbVHebPJYMs2oOaPzDDNwnYaDpZXGnizfbxwJk8RuSYwmr5/6X77nrBarJ1sgC/NjN1jmZ/aQBBy75b7Eg==",
      "node_ref": "1tJDuWkkAxbb2yMRTzb6Xun9rLtjkhtanQaXzEGYDU.11111111111111111111111111111111"
    }
  ]
}
`)
	kp := platformpolicy.NewKeyProcessor()
	publicKey, err := kp.ImportPublicKeyPEM([]byte("-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+2RsLu5z3nFEimNiesiLDH2Kw1GM\nvgYylDXAmZxpbGjQZ5FqHuXF+DJrwKYzDyfBDEQz6Tu/aeA2CgRZvqbKug==\n-----END PUBLIC KEY-----\n"))
	if err != nil {
		panic(err)
	}
	c, err := certificate.ReadCertificateFromReader(publicKey, platformpolicy.NewKeyProcessor(), buff)
	if err != nil {
		panic(err)
	}
	return c
}
