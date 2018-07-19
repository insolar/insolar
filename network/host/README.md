Insolar – Host Network
===============
Abstract networking layer

Key components
--------------
### [Transport](https://godoc.org/github.com/insolar/network/transport)
Network transport interface. It allows to abstract our network from physical transport.
It can either be IP based network or any other kind of message courier (e.g. an industrial message bus). 

### [Node](https://godoc.org/github.com/insolar/network/node)
Node is a fundamental part of networking system. Each node has:
 - one real network address (IP or any other transport protocol address)
 - multiple abstract network IDs (either node's own or ones belonging to relayed nodes)

### [Routing](https://godoc.org/github.com/insolar/network/routing)
It is actually a Kademlia hash table used to store network nodes and calculate distances between them.
See [Kademlia whitepaper](https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf) and
[XLattice design specification](http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html) for details.


### [Message](https://godoc.org/github.com/insolar/network/message)
A set of data transferred by this module between nodes.
 - Request message
 - Response message
 
 Now messages are serialized simply with encoding/gob.
 In future there will be a powerful robust serialization system based on Google's Protocol Buffers.

### [RPC](https://godoc.org/github.com/insolar/network/rpc)
RPC module allows higher level components to register methods that can be called by other network nodes.

Usage
-----

```go
package main

import (
	"github.com/insolar/insolar/network/host"
	"github.com/insolar/insolar/network/host/connection"
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/insolar/insolar/network/host/resolver"
	"github.com/insolar/insolar/network/host/rpc"
	"github.com/insolar/insolar/network/host/store"
	"github.com/insolar/insolar/network/host/transport"
)

func main() {
	configuration := host.NewNetworkConfiguration(
		resolver.NewStunResolver(""),
		connection.NewConnectionFactory(),
		transport.NewUTPTransportFactory(),
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy())

	dhtNetwork, err := configuration.CreateNetwork("0.0.0.0:31337", &host.Options{})
	if err != nil {
		panic(err)
	}
	defer configuration.CloseNetwork()

	dhtNetwork.Listen()
}
```

For more detailed usage example see [cmd/example/network/host/main.go](../../cmd/example/network/host/main.go)
