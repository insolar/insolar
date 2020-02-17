// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package endpoints

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"

	"github.com/pkg/errors"
)

const (
	ipSize        = net.IPv6len
	portSize      = 2
	ipAddressSize = ipSize + portSize

	maxPortNumber = ^uint16(0)
)

var defaultByteOrder = binary.BigEndian

type IPAddress [ipAddressSize]byte

func NewIPAddress(address string) (IPAddress, error) {
	var addr IPAddress

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return addr, errors.Errorf("invalid address: %s", address)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return addr, errors.Errorf("invalid ip: %s", host)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return addr, errors.Errorf("invalid port number: %s", port)
	}

	return addr, newIPAddress(ip, portNumber, &addr)
}

func newIPAddress(ip net.IP, portNumber int, addr *IPAddress) error {
	switch ipSize {
	case net.IPv6len:
		ip = ip.To16()
	case net.IPv4len:
		ip = ip.To4()
	default:
		panic("not implemented")
	}

	if portNumber > int(maxPortNumber) || portNumber <= 0 {
		return errors.Errorf("invalid port number: %d", portNumber)
	}

	portBytes := make([]byte, portSize)
	defaultByteOrder.PutUint16(portBytes, uint16(portNumber))

	copy(addr[:], ip)
	copy(addr[ipSize:], portBytes)

	return nil
}

func (a IPAddress) String() string {
	r := bytes.NewReader(a[:])

	ipBytes := make([]byte, ipSize)
	_, _ = r.Read(ipBytes)

	portBytes := make([]byte, portSize)
	_, _ = r.Read(portBytes)

	ip := net.IP(ipBytes)
	portNumber := defaultByteOrder.Uint16(portBytes)

	host := ip.String()
	port := strconv.Itoa(int(portNumber))

	return net.JoinHostPort(host, port)
}
