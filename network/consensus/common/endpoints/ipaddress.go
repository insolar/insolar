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
