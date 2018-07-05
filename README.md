Abstract networking layer
-------------------------

_This project is still in early development state.
It is not recommended to use it in production environment._

This project is an implementation of [Kademlia DHT](https://en.wikipedia.org/wiki/Kademlia).
It is mostly based on original specification but has multiple backward-incompatible changes.


Installation
============

    go get github.com/insolar/network


Usage
=====

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
=========

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.
