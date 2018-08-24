Insolar – Host Network
===============
Physical networking layer

[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar/network/hostnetwork)](https://goreportcard.com/report/github.com/insolar/insolar/network/hostnetwork)
[![GoDoc](https://godoc.org/github.com/insolar/insolar/network/hostnetwork?status.svg)](https://godoc.org/github.com/insolar/insolar/network/hostnetwork)


Overview
--------

We took [Kademlia DHT](https://en.wikipedia.org/wiki/Kademlia) original specifications and made significant improvements to make it ready
for real world application by enterprises.

#### Key features of our blockchain network layer:
 - **Support of heterogeneous network topology** with different types of hosts being able to communicate with each other.
   In classic peer-to-peer networks, any host can communicate directly with any other host on the network.
   In a real enterprise environment, this condition is often unacceptable for a variety of reasons including security.
 - **Network routing with a host or host group becoming relays** for others hosts.
   The network can continue to function despite various network restrictions such as firewalls, NATs, etc.
 - **Ability to limit number of gateways to corporate host group via relays** to keep the host group secure while being
   able to interact with the rest of the network through relays. This feature mitigates the risk of DDoS attacks.


Key components
--------------
### [Transport](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/transport)
Network transport interface. It allows to abstract our network from physical transport.
It can either be IP based network or any other kind of packet courier (e.g. an industrial packet bus). 

### [Host](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/host)
Host is a fundamental part of networking system. Each host has:
 - one real network address (IP or any other transport protocol address)
 - multiple abstract network IDs (either host's own or ones belonging to relayed hosts)

### [Routing](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/routing)
It is actually a Kademlia hash table used to store network hosts and calculate distances between them.
See [Kademlia whitepaper](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf) and
[XLattice design specification](http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html) for details.


### [Packet](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/packet)
A set of data transferred by this module between hosts.
 - Request packet
 - Response packet
 
 Now packets are serialized simply with encoding/gob.
 In future there will be a powerful robust serialization system based on Google's Protocol Buffers.

### [RPC](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/rpc)
RPC module allows higher level components to register methods that can be called by other network hosts.

Usage
-----

```go
package main

import (
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/configuration"
)

func main() {
	cfg := configuration.NewConfiguration().Host
	cfg.Address = "0.0.0.0:31337"

	network, err := hostnetwork.NewHostNetwork(cfg)
	if err != nil {
		panic(err)
	}
	defer network.Disconnect()

	network.Listen()
}
```

For more detailed usage example see [cmd/example/network/hostnetwork/main.go](../../cmd/example/network/hostnetwork/main.go)
