Insolar
===============
Blockchain platform

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2150/badge)](https://bestpractices.coreinfrastructure.org/projects/2150)

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
### [Network](network/dhtnetwork)
Kademlia DHT based blockchain network layer.
 - Support of heterogeneous network topology.
 - Network routing with a host or host group becoming relays for others hosts.
 - Ability to limit number of gateways to corporate host group via relays
   to keep the host group secure.

See [package readme](network/dhtnetwork) for more details.


### [Ledger](ledger)
Record storage engine backed by [BadgerDB](https://github.com/dgraph-io/badger).


### [Virtual machines](vm)
Various engines for smart contract execution:
 - [wasm](vm/wasm) - WebAssembly implementation of smart contracts


### [Application layer](application)
Application module describes interaction of system components with each other.
Every component of the system is a `SmartContract`. Members of the system are given the opportunity to build their own dApps by publishing smart contracts in `Domain` instances.
Domains define the visibility scope for the child contracts and their interaction policies. Actually, `Domain` is subclass of `SmartContract`.

See [package readme](application) for more details.


### [Configuration](configuration)

Provides configuration params for all Insolar components and helper for config resources management.


### [Metrics](metrics)

Using Prometheus monitoring system and time series database for collecting and store metrics


Installation
------------

    go get github.com/insolar/insolar


Generate default configuration file
------------

    go run cmd/insolar/* --cmd=default_config

Example
------------
    # Start node
    ./scripts/insolard/launch.sh

    # In other terminal:
    # Create user
    curl --data '{"query_type": "create_member", "name": "Peter"}' "localhost:19191/api/v1?"
    # Dump user info
    curl --data '{"query_type": "dump_all_users"}' "localhost:19191/api/v1?"

Docker container
------------

    docker pull insolar/insolar
    docker run -ti insolar/insolar


Contributing
------------
See [Contributing Guidelines](.github/CONTRIBUTING.md).


License
-------
This project is licensed under the terms of the Apache license 2.0.
Please see [LICENSE](LICENSE) for more information.
