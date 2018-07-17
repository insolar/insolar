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

package relay

// Proxy contains proxy addresses
type Proxy interface {
	// AddProxyNode add an address to proxy list.
	AddProxyNode(address string)
	// RemoveProxyNode removes proxy address from proxy list.
	RemoveProxyNode(address string)
	// GetNextProxyAddress returns a next address to send from proxy list.
	GetNextProxyAddress() string
	// Count return proxyList length.
	Count() int
}

// Note: thread unsafe!!!
type proxy struct {
	proxyList []string
	iterator  int
}

// CreateProxy instantiates proxy struct.
func CreateProxy() Proxy {
	return &proxy{proxyList: make([]string, 0), iterator: 0}
}

// AddProxyNode add an address to proxy list.
func (p *proxy) AddProxyNode(address string) {
	i := p.getProxyIndex(address)

	if i != -1 {
		return
	}
	p.proxyList = append(p.proxyList, address)
}

// RemoveProxyNode removes proxy address from proxy list.
func (p *proxy) RemoveProxyNode(address string) {
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

// Count return proxyList length.
func (p *proxy) Count() int {
	return len(p.proxyList)
}
