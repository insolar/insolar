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
