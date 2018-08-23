Insolar
===============
Blockchain platform

[![Build Status](https://travis-ci.org/insolar/insolar.svg?branch=master)](https://travis-ci.org/insolar/insolar)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar)](https://goreportcard.com/report/github.com/insolar/insolar)
[![GoDoc](https://godoc.org/github.com/insolar/insolar?status.svg)](https://godoc.org/github.com/insolar/insolar)
[![codecov](https://codecov.io/gh/insolar/insolar/branch/master/graph/badge.svg)](https://codecov.io/gh/insolar/insolar)

_This project is still in early development state.
It is not recommended to use it in production environment._

Overview
--------
**Insolar** is the next generation high-performance scalable blockchain platform
designed with the express purpose to meet an immense business scope.
The enterprise-grade distributed ledger cloud platform will help to increase
business velocity, create new revenue streams, and reduce cost and risk
by securely extending enterprise SaaS and on-premises applications
to drive tamper-resistant transactions on a trusted business network.

Insolar supports public and private blockchains and is able to customize
different blockchains for different applications. Insolar team will
constantly provide common modules on the underlying infrastructure
for different kinds of distributed scenarios.

We value the expansion of the ecosystem which operates across chains,
systems, industries and applications. With a range of protocols and modules,
data and information will be connected to support various business scenarios.
Our goal is to build the underlying blockchain infrastructure to bridge
the real world and the distributed digital world. With this, companies
from different industries will be able to develop applications
for a range of scenarios and collaborate with other entities on the platform.


Components
----------
### [Network](network/hostnetwork)
Kademlia DHT based blockchain network layer.
 - Support of heterogeneous network topology.
 - Network routing with a node or node group becoming relays for others nodes.
 - Ability to limit number of gateways to corporate node group via relays
   to keep the node group secure.

See [package readme](network/hostnetwork) for more details.


### [Ledger](ledger)
Record storage engine backed by [LevelDB](https://github.com/syndtr/goleveldb).



### [Virtual machines](vm)
Various engines for smart contract execution:
 - [wasm](vm/wasm) - WebAssembly implementation of smart contracts


### [Genesis](genesis)
Genesis module describes interaction of system components with each other.
Every component of the system is a `SmartContract`. Members of the system are given the opportunity to build their own dApps by publishing smart contracts in `Domain` instances.
Domains define the visibility scope for the child contracts and their interaction policies. Actually, `Domain` is subclass of `SmartContract`.

See [package readme](genesis) for more details.


Installation
------------

    go get github.com/insolar/insolar


Contributing
------------
Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.

License
-------
This project is licensed under the terms of the Apache license 2.0.
Please see [LICENSE](LICENSE) for more information.
