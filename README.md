Insolar
===============
Enterprise-ready blockchain platform

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2150/badge)](https://bestpractices.coreinfrastructure.org/projects/2150)

[![Build Status](https://travis-ci.org/insolar/insolar.svg?branch=master)](https://travis-ci.org/insolar/insolar)
[![Go Report Card](https://goreportcard.com/badge/github.com/insolar/insolar)](https://goreportcard.com/report/github.com/insolar/insolar)
[![GoDoc](https://godoc.org/github.com/insolar/insolar?status.svg)](https://godoc.org/github.com/insolar/insolar)
[![codecov](https://codecov.io/gh/insolar/insolar/branch/master/graph/badge.svg)](https://codecov.io/gh/insolar/insolar)


Overview
--------
**Insolar** is building a 4th generation blockchain platform for business aimed to enable seamless interactions between companies and unlock new growth opportunities. In addition to the blockchain platform, Insolar will provide blockchain services and ecosystem support for companies that are looking to develop and deploy blockchain solutions. Insolar will feature most complete and secure set of production-ready business blockchain tools and services to quickly build or launch blockchain enterprise applications, accelerating the progression path from initial proof-of-concept to full-scale production.

The world’s most innovative companies in finance, logistics, consumer goods, energy, healthcare, transportation, manufacturing and others will be turning to Insolar to create applications and networks that deliver tangible business success. They recognise that even in today’s digital economy, vast amounts of value continue to be trapped inside processes and organisations that don’t connect. Insolar is their remedy, helping them discover and design business value in blockchain networks — starting, accelerating and innovating strategies that replace longstanding business friction with trust and transparency. Delegating trust to a blockchain means that businesses can pursue broader networks, onboard new partners, and enter new ecosystems with ease. Blockchain-based networks that support multiparty collaboration around shared, trusted data and process automation across organisational boundaries bring benefits at many levels, starting with efficiency gains and culminating in reinventing how entire industry ecosystems operate.

Insolar is a global team of 60+ people in North America and Europe, including a 35-strong engineering team with practical blockchain engineering know-how, and 10 leading blockchain academics from major institutions (York University, ETH Zurich, Princeton).


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
    # Start net of nodes
    ./scripts/insolard/launchnet.sh -g


    # In other terminal:
    
    # Build insolar
    make insolar
  
    # Send request example
    ./bin/insolar -c=send_request --config=./scripts/insolard/configs/root_member_keys.json --root_as_caller --params=params.json

   ##### See [insolar readme](cmd/insolar) for more details. 

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
