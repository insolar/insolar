Insolar – Host Network
===============
Physical networking layer

[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar/network/hostnetwork)](https://goreportcard.com/report/github.com/insolar/insolar/network/host)
[![GoDoc](https://godoc.org/github.com/insolar/insolar/network/hostnetwork?status.svg)](https://godoc.org/github.com/insolar/insolar/network/host)


Overview
--------

We took [Kademlia DHT](https://en.wikipedia.org/wiki/Kademlia) original specifications and made significant improvements to make it ready
for real world application by enterprises.

#### Key features of our blockchain network layer:
 - **Support of heterogeneous network topology** with different types of nodes being able to communicate with each other.
   In classic peer-to-peer networks, any node can communicate directly with any other node on the network.
   In a real enterprise environment, this condition is often unacceptable for a variety of reasons including security.
 - **Network routing with a node or node group becoming relays** for others nodes.
   The network can continue to function despite various network restrictions such as firewalls, NATs, etc.
 - **Ability to limit number of gateways to corporate node group via relays** to keep the node group secure while being
   able to interact with the rest of the network through relays. This feature mitigates the risk of DDoS attacks.


Key components
--------------
### [Transport](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/transport)
Network transport interface. It allows to abstract our network from physical transport.
It can either be IP based network or any other kind of packet courier (e.g. an industrial packet bus). 

### [Node](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/node)
Node is a fundamental part of networking system. Each node has:
 - one real network address (IP or any other transport protocol address)
 - multiple abstract network IDs (either node's own or ones belonging to relayed nodes)

### [Routing](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/routing)
It is actually a Kademlia hash table used to store network nodes and calculate distances between them.
See [Kademlia whitepaper](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf) and
[XLattice design specification](http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html) for details.


### [Packet](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/packet)
A set of data transferred by this module between nodes.
 - Request packet
 - Response packet
 
 Now packets are serialized simply with encoding/gob.
 In future there will be a powerful robust serialization system based on Google's Protocol Buffers.

### [RPC](https://godoc.org/github.com/insolar/insolar/network/hostnetwork/rpc)
RPC module allows higher level components to register methods that can be called by other network nodes.

Usage
-----

```go
package main

import (
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/connection"
	"github.com/insolar/insolar/network/hostnetwork/node"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"
)

func main() {
	configuration := hostnetwork.NewNetworkConfiguration(
		resolver.NewStunResolver(""),
		connection.NewConnectionFactory(),
		transport.NewUTPTransportFactory(),
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy())

	dhtNetwork, err := configuration.CreateNetwork("0.0.0.0:31337", &hostnetwork.Options{})
	if err != nil {
		panic(err)
	}
	defer configuration.CloseNetwork()

	dhtNetwork.Listen()
}
```

For more detailed usage example see [cmd/example/network/hostnetwork/main.go](../../cmd/example/network/hostnetwork/main.go)
