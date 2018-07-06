Insolar – Network
===============
Abstract networking layer

[![Build Status](https://travis-ci.org/insolar/network.svg?branch=master)](https://travis-ci.org/insolar/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/network)](https://goreportcard.com/report/github.com/insolar/network)
[![GoDoc](https://godoc.org/github.com/insolar/network?status.svg)](https://godoc.org/github.com/insolar/network)

_This project is still in early development state.
It is not recommended to use it in production environment._

Overview
--------
**Insolar** is a blockchain platform developed by INS

This library is an implementation of [Kademlia DHT](https://en.wikipedia.org/wiki/Kademlia).
It is mostly based on original specification but has multiple backward-incompatible improvements.

The main feature of our implementation is the support of heterogeneous network topology,
meaning that different types of computers and devices can communicate with each other,
using various OS and/or protocols. 

In classical peer-to-peer networks, it is presumed that any node can communicate directly
with any other node on the network. But in a real corporate environment, this condition
is often unacceptable for a variety of reasons including security.

We have added routing awareness to the network, where individual nodes or group of nodes
can be relays for others. Thus, despite various network restrictions (firewalls, NATs etc.),
the network continues to function.

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

Installation
------------

    go get github.com/insolar/network


Usage
-----

```go
package main

import (
	"github.com/insolar/network"
	"github.com/insolar/network/connection"
	"github.com/insolar/network/node"
	"github.com/insolar/network/resolver"
	"github.com/insolar/network/rpc"
	"github.com/insolar/network/store"
	"github.com/insolar/network/transport"
)

func main() {
	configuration := network.NewNetworkConfiguration(
		resolver.NewStunResolver(""),
		connection.NewConnectionFactory(),
		transport.NewUTPTransportFactory(),
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}))

	dhtNetwork, err := configuration.CreateNetwork("0.0.0.0:31337", &network.Options{})
	if err != nil {
		panic(err)
	}
	defer configuration.CloseNetwork()

	dhtNetwork.Listen()
}
```

For more detailed usage example see [cmd/example/main.go](cmd/example/main.go)


Contributing
------------

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.
