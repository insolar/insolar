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

package relay

// Proxy contains proxy addresses.
type Proxy interface {
	// AddProxyHost add an address to proxy list.
	AddProxyHost(address string)
	// RemoveProxyHost removes proxy address from proxy list.
	RemoveProxyHost(address string)
	// GetNextProxyAddress returns a next address to send from proxy list.
	GetNextProxyAddress() string
	// ProxyHostsCount return added proxy count.
	ProxyHostsCount() int
}

// Note: thread unsafe!!!
type proxy struct {
	proxyList []string
	iterator  int
}

// NewProxy instantiates proxy struct.
func NewProxy() Proxy {
	return &proxy{
		proxyList: make([]string, 0),
		iterator:  0,
	}
}

// AddProxyHost add an address to proxy list.
func (p *proxy) AddProxyHost(address string) {
	i := p.getProxyIndex(address)

	if i != -1 {
		return
	}
	p.proxyList = append(p.proxyList, address)
}

// RemoveProxyHost removes proxy address from proxy list.
func (p *proxy) RemoveProxyHost(address string) {
	i := p.getProxyIndex(address)

	if i == -1 {
		return
	}

	p.proxyList = append(p.proxyList[:i], p.proxyList[i+1:]...)
}

// getProxyIndex returns index of address in proxy list.
func (p *proxy) getProxyIndex(address string) int {
	for i := 0; i < len(p.proxyList); i++ {
		if p.proxyList[i] == address {
			return i
		}
	}

	return -1
}

// GetNextProxyAddress returns a next address to send from proxy list.
func (p *proxy) GetNextProxyAddress() string {
	if len(p.proxyList) == 0 {
		return ""
	}
	if p.iterator >= len(p.proxyList) {
		p.iterator = 0
	}
	p.iterator++
	return p.proxyList[p.iterator-1]
}

// ProxyHostsCount return added proxy count.
func (p *proxy) ProxyHostsCount() int {
	return len(p.proxyList)
}
