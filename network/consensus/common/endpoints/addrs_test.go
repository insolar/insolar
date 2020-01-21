// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package endpoints

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/insolar"

	"github.com/stretchr/testify/require"
)

func TestIsLocalHost(t *testing.T) {
	var ha *Name
	//noinspection ALL
	require.False(t, ha.IsLocalHost())

	h := Name("")
	require.True(t, h.IsLocalHost())

	h = "addr"
	require.False(t, h.IsLocalHost())
}

func TestEquals(t *testing.T) {
	var h1 *Name
	h2 := Name("")
	//noinspection ALL
	require.False(t, h1.Equals(h2))

	h3 := h2
	require.True(t, h2.Equals(h3))

	h2 = "addr"
	h3 = "addr"
	require.True(t, h2.Equals(h3))

	h3 = "addr1"
	require.False(t, h2.Equals(h3))
}

func TestEqualsToString(t *testing.T) {
	var h1 *Name
	//noinspection ALL
	require.False(t, h1.EqualsToString(""))

	h2 := Name("")
	require.True(t, h2.EqualsToString(""))

	h2 = "addr"
	require.True(t, h2.EqualsToString("addr"))

	h2 = "addr"
	require.False(t, h2.EqualsToString("addr1"))
}

func TestString(t *testing.T) {
	s := ""
	h1 := Name(s)
	require.Equal(t, s, h1.String())

	s = "addr"
	h1 = Name(s)
	require.Equal(t, s, h1.String())
}

func TestNewHostIdentityFromHolder(t *testing.T) {
	in := NewInboundMock(t)
	addr := Name("addr")
	in.GetNameAddressMock.Set(func() Name { return addr })
	skh := cryptkit.NewSignatureKeyHolderMock(t)
	in.GetTransportKeyMock.Set(func() cryptkit.SignatureKeyHolder { return skh })
	ch := cryptkit.NewCertificateHolderMock(t)
	in.GetTransportCertMock.Set(func() cryptkit.CertificateHolder { return ch })
	inc := NewHostIdentityFromHolder(in)
	require.Equal(t, addr, inc.Addr)

	require.Equal(t, skh, inc.Key)

	require.Equal(t, ch, inc.Cert)
}

func TestShortNodeIDAsByteString(t *testing.T) {
	require.NotEmpty(t, ShortNodeIDAsByteString(insolar.ShortNodeID(123)))
}

func TestAsByteString(t *testing.T) {
	inc := InboundConnection{Addr: "test"}
	require.NotEmpty(t, inc.AsByteString())
}

func TestGetNameAddress(t *testing.T) {
	addr := Name("test")
	inc := InboundConnection{Addr: addr}
	require.Equal(t, addr, inc.GetNameAddress())
}

func TestGetTransportKey(t *testing.T) {
	skh := cryptkit.NewSignatureKeyHolderMock(t)
	inc := InboundConnection{Key: skh}
	require.Equal(t, skh, inc.GetTransportKey())
}

func TestGetTransportCert(t *testing.T) {
	ch := cryptkit.NewCertificateHolderMock(t)
	inc := InboundConnection{Cert: ch}
	require.Equal(t, ch, inc.GetTransportCert())
}

func TestEqualOutboundEndpoints(t *testing.T) {
	ob1 := NewOutboundMock(t)
	require.False(t, EqualOutboundEndpoints(nil, ob1))

	require.False(t, EqualOutboundEndpoints(ob1, nil))

	require.True(t, EqualOutboundEndpoints(ob1, ob1))

	et1 := NameEndpoint
	ob1.GetEndpointTypeMock.Set(func() NodeEndpointType { return *(&et1) })
	ob2 := NewOutboundMock(t)
	et2 := RelayEndpoint
	ob2.GetEndpointTypeMock.Set(func() NodeEndpointType { return *(&et2) })
	require.False(t, EqualOutboundEndpoints(ob1, ob2))

	et2 = et1
	addr1 := Name("addr")
	addr2 := Name("addr2")
	ob1.GetNameAddressMock.Set(func() Name { return *(&addr1) })
	ob2.GetNameAddressMock.Set(func() Name { return *(&addr2) })
	require.False(t, EqualOutboundEndpoints(ob1, ob2))

	addr2 = addr1
	require.True(t, EqualOutboundEndpoints(ob1, ob2))

	et1 = IPEndpoint
	et2 = et1
	ip1 := IPAddress{}
	ip2 := IPAddress{1}
	ob1.GetIPAddressMock.Set(func() IPAddress { return *(&ip1) })
	ob2.GetIPAddressMock.Set(func() IPAddress { return *(&ip2) })
	require.False(t, EqualOutboundEndpoints(ob1, ob2))

	ip2 = ip1
	require.True(t, EqualOutboundEndpoints(ob1, ob2))

	et1 = RelayEndpoint
	et2 = et1
	rID1 := insolar.ShortNodeID(1)
	rID2 := insolar.ShortNodeID(2)
	ob1.GetRelayIDMock.Set(func() insolar.ShortNodeID { return *(&rID1) })
	ob2.GetRelayIDMock.Set(func() insolar.ShortNodeID { return *(&rID2) })
	require.False(t, EqualOutboundEndpoints(ob1, ob2))

	rID2 = rID1
	require.True(t, EqualOutboundEndpoints(ob1, ob2))

	et1 = NodeEndpointType(4)
	et2 = et1
	require.Panics(t, func() { EqualOutboundEndpoints(ob1, ob2) })
}

func TestEqualListOfOutboundEndpoints(t *testing.T) {
	require.True(t, EqualListOfOutboundEndpoints(nil, nil))

	var o, p []Outbound
	require.True(t, EqualListOfOutboundEndpoints(p, nil))

	require.True(t, EqualListOfOutboundEndpoints(nil, o))

	ob1 := NewOutboundMock(t)
	p = append(p, ob1)
	require.False(t, EqualListOfOutboundEndpoints(p, o))

	o = append(o, ob1)
	require.True(t, EqualListOfOutboundEndpoints(p, o))

	ob2 := NewOutboundMock(t)
	ob2.GetEndpointTypeMock.Set(func() NodeEndpointType { return IPEndpoint })
	p = append(p, ob2)
	require.False(t, EqualListOfOutboundEndpoints(p, o))

	ob3 := NewOutboundMock(t)
	ob3.GetEndpointTypeMock.Set(func() NodeEndpointType { return NameEndpoint })
	o = append(o, ob3)
	require.False(t, EqualListOfOutboundEndpoints(p, o))

	ob3.GetEndpointTypeMock.Set(func() NodeEndpointType { return IPEndpoint })
	ob2.GetIPAddressMock.Set(func() IPAddress { return IPAddress{1} })
	ob3.GetIPAddressMock.Set(func() IPAddress { return IPAddress{1} })
	require.True(t, EqualListOfOutboundEndpoints(p, o))
}
